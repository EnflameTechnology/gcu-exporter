// Copyright (c) 2024 Enflame. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/exporter-toolkit/web"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	Version string
)

const (
	defaultListenAddress              = ":9400"
	defaultMetricsPath                = "/metrics"
	defaultDataFile                   = "/tmp/gcu_monitor_data"
	defaultExporterPath               = "/usr/bin/gcu-exporter-cmdline"
	defaultMode                       = "daemon"
	defaultDaemonBin                  = "/usr/bin/gcu-exporter-daemon"
	defaultDaemonAddress              = "127.0.0.1:9410"
	defaultDaemonStartupTimeoutSec    = 10        // seconds
	defaultDaemonRestartEverySec      = 12 * 3600 // 12h in seconds
	defaultDaemonRetryIntervalSec     = 5         // seconds; short so transient driver flakes recover quickly
	defaultDaemonPostRecoveryGraceSec = 5 * 60    // 5m in seconds
	daemonMaxRetryAttempts            = 3
)

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(defaultListenAddress).String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default(defaultMetricsPath).String()
		mode = kingpin.Flag(
			"mode",
			"Data source mode: daemon (long-lived, default) | cmdline (one-shot exec) | smi (efsmi json primary, eflib fallback).",
		).Default(defaultMode).Envar("EXPORTER_MODE").Enum("daemon", "cmdline", "smi")
		exporterPath = kingpin.Flag(
			"exporter.path",
			"Path to the gcu-exporter-cmdline binary (used when --mode=cmdline).",
		).Default(defaultExporterPath).String()
		dataFile = kingpin.Flag(
			"data.file",
			"Path to the temporary file where cmdline mode metrics data will be written.",
		).Default(defaultDataFile).String()
		timeout = kingpin.Flag(
			"exporter.timeout",
			"Timeout for executing the exporter command in cmdline mode (seconds).",
		).Default("30").Int()
		daemonBin = kingpin.Flag(
			"daemon.bin",
			"Path to gcu-exporter-daemon binary (used when --mode=daemon).",
		).Default(defaultDaemonBin).Envar("DAEMON_BIN").String()
		daemonAddress = kingpin.Flag(
			"daemon.address",
			"Address daemon listens on; server fetches /metrics here.",
		).Default(defaultDaemonAddress).Envar("DAEMON_ADDRESS").String()
		daemonStartupTimeoutSec = kingpin.Flag(
			"daemon.startup-timeout",
			"Max seconds (integer) to wait for daemon /metrics to become ready at startup.",
		).Default(strconv.Itoa(defaultDaemonStartupTimeoutSec)).Envar("DAEMON_STARTUP_TIMEOUT").Int()
		daemonRestartIntervalSec = kingpin.Flag(
			"daemon.restart-interval",
			"Periodically restart the daemon child every N seconds (integer). 0 disables.",
		).Default(strconv.Itoa(defaultDaemonRestartEverySec)).Envar("DAEMON_RESTART_INTERVAL").Int()
		daemonRetryIntervalSec = kingpin.Flag(
			"daemon.retry-interval",
			"Wait N seconds (integer) before each restart attempt after an abnormal daemon exit. Up to 3 attempts at this fixed interval; if all fail, server exits and lets K8s rebuild the pod.",
		).Default(strconv.Itoa(defaultDaemonRetryIntervalSec)).Envar("DAEMON_RETRY_INTERVAL").Int()
		daemonPostRecoveryGraceSec = kingpin.Flag(
			"daemon.post-recovery-grace",
			"After a successful recovery from abnormal exit, wait N seconds (integer) and then exit the server so K8s rebuilds the pod with a fresh process state. 0 disables (server keeps running indefinitely after recovery).",
		).Default(strconv.Itoa(defaultDaemonPostRecoveryGraceSec)).Envar("DAEMON_POST_RECOVERY_GRACE").Int()
		webConfigFile = kingpin.Flag(
			"web-config-file",
			"Path to YAML web config file that can enable authentication.",
		).Default("").String()
		enablePprof = kingpin.Flag(
			"web.enable-pprof",
			"Expose /debug/pprof endpoints on the daemon (not recommended in production).",
		).Default("false").Bool()
	)

	kingpin.Version(Version)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Range validation. All time configs are integer seconds. Invalid values
	// fall back to defaults with WARN log rather than failing fast — daemon
	// config is non-critical at startup, silent rejection would surprise
	// operators more than a noisy fallback.
	*timeout = ensurePositiveInt(*timeout, 30, "exporter.timeout")
	*daemonStartupTimeoutSec = ensurePositiveInt(*daemonStartupTimeoutSec, defaultDaemonStartupTimeoutSec, "daemon.startup-timeout")
	*daemonRetryIntervalSec = ensureNonNegativeInt(*daemonRetryIntervalSec, defaultDaemonRetryIntervalSec, "daemon.retry-interval")
	*daemonRestartIntervalSec = ensureNonNegativeInt(*daemonRestartIntervalSec, defaultDaemonRestartEverySec, "daemon.restart-interval")
	*daemonPostRecoveryGraceSec = ensureNonNegativeInt(*daemonPostRecoveryGraceSec, defaultDaemonPostRecoveryGraceSec, "daemon.post-recovery-grace")

	log.Printf("Starting gcu-exporter-server version %s", Version)
	log.Printf("Listening on %s", *listenAddress)
	log.Printf("Metrics path: %s", *metricsPath)
	log.Printf("Mode: %s", *mode)

	handler := &metricsHandler{
		mode:          *mode,
		exporterPath:  *exporterPath,
		dataFile:      *dataFile,
		timeout:       time.Duration(*timeout) * time.Second,
		daemonAddress: *daemonAddress,
	}

	// daemon mode: spawn daemon child, supervise, wait ready
	var daemon *daemonSupervisor
	if *mode == "daemon" || *mode == "smi" {
		log.Printf("Daemon binary: %s", *daemonBin)
		log.Printf("Daemon address: %s", *daemonAddress)
		daemon = &daemonSupervisor{
			bin:               *daemonBin,
			address:           *daemonAddress,
			webConfigFile:     *webConfigFile,
			enablePprof:       *enablePprof,
			smiMode:           *mode == "smi",
			startupTimeout:    time.Duration(*daemonStartupTimeoutSec) * time.Second,
			restartInterval:   time.Duration(*daemonRestartIntervalSec) * time.Second,
			retryInterval:     time.Duration(*daemonRetryIntervalSec) * time.Second,
			postRecoveryGrace: time.Duration(*daemonPostRecoveryGraceSec) * time.Second,
			stopCh:            make(chan struct{}),
		}

		if err := logWebAuthStatus(*webConfigFile); err != nil {
			log.Fatalf("Web config error: %v", err)
		}
		if *enablePprof {
			log.Printf("WARN: pprof enabled on daemon; endpoints at http://%s/debug/pprof/ (not recommended in production)", *daemonAddress)
		}

		if err := daemon.start(); err != nil {
			log.Fatalf("Failed to start daemon: %v", err)
		}
		if err := daemon.waitReady(); err != nil {
			log.Fatalf("Daemon failed to become ready: %v", err)
		}
		log.Printf("Daemon ready, pid=%d", daemon.pid())
		log.Printf("Daemon periodic restart interval: %ds (0 = disabled)", *daemonRestartIntervalSec)
		log.Printf("Daemon abnormal-exit retry interval: %ds (max %d attempts)", *daemonRetryIntervalSec, daemonMaxRetryAttempts)
		log.Printf("Daemon post-recovery grace: %ds (0 = disabled; otherwise server exits this long after a successful recovery so K8s rebuilds pod)", *daemonPostRecoveryGraceSec)
		go daemon.supervise()
		go daemon.periodicRestart()
		handler.daemon = daemon
	} else {
		log.Printf("Exporter path: %s", *exporterPath)
		log.Printf("Data file: %s", *dataFile)
	}

	// SIGTERM / SIGINT graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		log.Printf("Signal received, shutting down")
		if daemon != nil {
			daemon.stop()
		}
		os.Exit(0)
	}()

	// Routes
	http.Handle(*metricsPath, handler)
	http.HandleFunc("/healthz", handler.healthz)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>Enflame GCU Exporter Server</title></head>
				<body>
				<h1>Enflame GCU Exporter Server</h1>
				<p><a href="` + *metricsPath + `">Metrics</a></p>
				<p><a href="/healthz">Healthz</a></p>
				</body>
				</html>`))
	})

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

type webAuthConfig struct {
	Users map[string]string `yaml:"basic_auth_users"`
}

func logWebAuthStatus(configPath string) error {
	if configPath == "" {
		log.Printf("Daemon web-config-file not set; daemon :9410 has no Basic Auth from config")
		return nil
	}
	if err := web.Validate(configPath); err != nil {
		return fmt.Errorf("invalid web config file %s: %w", configPath, err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read web config file %s: %w", configPath, err)
	}
	var cfg webAuthConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse web config file %s: %w", configPath, err)
	}
	if len(cfg.Users) > 0 {
		log.Printf(
			"Daemon Basic Auth enabled (web-config-file=%s, users=%d); "+
				"server :9400 does not authenticate, forwards client Authorization to daemon",
			configPath, len(cfg.Users),
		)
	} else {
		log.Printf(
			"Web config loaded for daemon but basic_auth_users empty; "+
				"daemon Basic Auth disabled (web-config-file=%s)",
			configPath,
		)
	}
	return nil
}

// ---------- metrics handler ----------

type metricsHandler struct {
	mode          string
	exporterPath  string
	dataFile      string
	timeout       time.Duration
	daemonAddress string
	daemon        *daemonSupervisor // nil unless mode=daemon
}

func (h *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch h.mode {
	case "daemon", "smi":
		h.serveFromDaemon(w, r)
	case "cmdline":
		h.serveFromExec(w, r)
	default:
		// Unreachable: kingpin Enum validates --mode at parse time.
		http.Error(w, "unknown mode: "+h.mode, http.StatusInternalServerError)
	}
}

func (h *metricsHandler) serveFromDaemon(w http.ResponseWriter, r *http.Request) {
	url := "http://" + h.daemonAddress + "/metrics"
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("build daemon request: %v", err), http.StatusInternalServerError)
		return
	}
	if auth := r.Header.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("daemon GET %s failed: %v", url, err)
		http.Error(w, fmt.Sprintf("daemon unreachable: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	} else {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("copy daemon response: %v", err)
	}
}

func (h *metricsHandler) serveFromExec(w http.ResponseWriter, r *http.Request) {
	if err := h.executeExporter(r.Context()); err != nil {
		log.Printf("Error executing exporter: %v", err)
		http.Error(w, fmt.Sprintf("Error executing exporter: %v", err), http.StatusInternalServerError)
		return
	}
	data, err := os.ReadFile(h.dataFile)
	if err != nil {
		log.Printf("Error reading data file %s: %v", h.dataFile, err)
		http.Error(w, fmt.Sprintf("Error reading data file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *metricsHandler) healthz(w http.ResponseWriter, r *http.Request) {
	if h.mode == "cmdline" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok (mode=%s)\n", h.mode)
		return
	}
	if h.daemon == nil || !h.daemon.alive() {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "daemon not alive")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok (mode=daemon, pid=%d, restarts=%d)\n", h.daemon.pid(), h.daemon.restartCount())
}

func (h *metricsHandler) executeExporter(reqCtx context.Context) error {
	exporterPath, err := filepath.Abs(h.exporterPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for exporter: %w", err)
	}
	if _, err := os.Stat(exporterPath); os.IsNotExist(err) {
		return fmt.Errorf("exporter binary not found at %s", exporterPath)
	}
	ctx, cancel := context.WithTimeout(reqCtx, h.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, exporterPath,
		"--output.format", "prometheus",
		"--output.pretty",
		"--output.file", h.dataFile,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	defer stderr.Close()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start exporter: %w", err)
	}
	go func() {
		stderrData, _ := io.ReadAll(stderr)
		if len(stderrData) > 0 {
			log.Printf("Exporter stderr: %s", string(stderrData))
		}
	}()
	go func() {
		stdoutData, _ := io.ReadAll(stdout)
		if len(stdoutData) > 0 {
			log.Printf("Exporter stdout: %s", string(stdoutData))
		}
	}()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("exporter command failed: %w", err)
	}
	return nil
}

// ---------- config validation helpers ----------

func ensurePositiveInt(v, def int, name string) int {
	if v <= 0 {
		log.Printf("WARN: --%s=%d is invalid (must be a positive integer), falling back to default %d", name, v, def)
		return def
	}
	return v
}

// ensureNonNegativeInt accepts 0 (which disables the corresponding behavior,
// e.g., periodic restart) but rejects negative values.
func ensureNonNegativeInt(v, def int, name string) int {
	if v < 0 {
		log.Printf("WARN: --%s=%d is invalid (must be a non-negative integer; 0 disables), falling back to default %d", name, v, def)
		return def
	}
	return v
}

// ---------- daemon supervisor ----------

type daemonSupervisor struct {
	bin               string
	address           string
	webConfigFile     string
	enablePprof       bool
	smiMode           bool
	startupTimeout    time.Duration
	restartInterval   time.Duration // 0 disables periodic restart
	retryInterval     time.Duration // wait between restart attempts after abnormal exit
	postRecoveryGrace time.Duration // 0 disables; otherwise server exits this long after a successful recovery

	mu                sync.Mutex
	cmd               *exec.Cmd
	postRecoveryTimer *time.Timer
	stopping          atomic.Bool
	restarts          atomic.Int32

	intentionalRestart atomic.Bool
	stopCh             chan struct{}
}

// schedulePostRecoveryExit arms a one-shot timer that exits the server after
// d.postRecoveryGrace once the daemon has been successfully restarted from an
// abnormal exit. eflib does not re-enumerate hot-added PCI devices at runtime,
// so a recovered daemon may keep an outdated device set; rebuilding the pod
// (via K8s) is the safest way to get a fresh enumerate. The grace window lets
// operators investigate before the rebuild kicks in.
func (d *daemonSupervisor) schedulePostRecoveryExit() {
	if d.postRecoveryGrace <= 0 {
		return
	}
	d.mu.Lock()
	if d.postRecoveryTimer != nil {
		d.postRecoveryTimer.Stop()
	}
	grace := d.postRecoveryGrace
	d.postRecoveryTimer = time.AfterFunc(grace, func() {
		log.Fatalf("daemon post-recovery grace (%v) expired; exiting to let K8s rebuild pod and re-enumerate devices", grace)
	})
	d.mu.Unlock()
	log.Printf("daemon recovered from abnormal exit; scheduled server exit in %v to let K8s rebuild pod", grace)
}

// periodicRestart fires SIGTERM on the daemon child every restart-interval,
// letting supervise() handle the actual respawn (with intentionalRestart
// flag set so it doesn't count toward the 5-in-60s failure ceiling).
// Zero or negative interval disables the goroutine entirely.
func (d *daemonSupervisor) periodicRestart() {
	if d.restartInterval <= 0 {
		return
	}
	ticker := time.NewTicker(d.restartInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("Periodic daemon restart fired (interval=%v)", d.restartInterval)
			d.triggerIntentionalRestart()
			// Give supervise() time to actually restart before next cycle.
			time.Sleep(2 * time.Second)
		case <-d.stopCh:
			return
		}
	}
}

func (d *daemonSupervisor) triggerIntentionalRestart() {
	d.mu.Lock()
	cmd := d.cmd
	d.mu.Unlock()
	if cmd == nil || cmd.Process == nil {
		return
	}
	d.intentionalRestart.Store(true)
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Printf("triggerIntentionalRestart SIGTERM failed: %v", err)
		d.intentionalRestart.Store(false)
	}
}

func (d *daemonSupervisor) start() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	args := []string{}
	if d.webConfigFile != "" {
		args = append(args, "--web-config-file="+d.webConfigFile)
	}
	if d.enablePprof {
		args = append(args, "--web.enable-pprof")
	}
	if d.smiMode {
		args = append(args, "--smi-mode")
	}
	cmd := exec.Command(d.bin, args...)
	if d.smiMode {
		// The daemon decides Init tolerance from the SMI_MODE env before flags
		// are parsed, so pass it as an env var (the --smi-mode flag also reads it).
		cmd.Env = append(os.Environ(), "SMI_MODE=1")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start daemon binary %s: %w", d.bin, err)
	}
	d.cmd = cmd
	return nil
}

func (d *daemonSupervisor) waitReady() error {
	deadline := time.Now().Add(d.startupTimeout)
	url := "http://" + d.address + "/metrics"
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			code := resp.StatusCode
			resp.Body.Close()
			if code == http.StatusOK || code == http.StatusUnauthorized {
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("daemon not ready at %s within %v", url, d.startupTimeout)
}

// supervise blocks on cmd.Wait(); on intentional restart it goes through a
// clean respawn; on abnormal exit it tries up to daemonMaxRetryAttempts
// restarts at the fixed retryInterval. If all attempts fail the whole server
// process exits so K8s rebuilds the pod.
func (d *daemonSupervisor) supervise() {
	for {
		d.mu.Lock()
		cmd := d.cmd
		d.mu.Unlock()
		err := cmd.Wait()
		if d.stopping.Load() {
			return
		}
		if d.intentionalRestart.Swap(false) {
			log.Printf("daemon intentional restart (previous wait err=%v)", err)
			if err := d.start(); err != nil {
				log.Fatalf("daemon intentional restart failed: %v", err)
			}
			if err := d.waitReady(); err != nil {
				log.Fatalf("daemon intentional restart not ready: %v", err)
			}
			d.restarts.Add(1)
			log.Printf("daemon intentional restart done, pid=%d (total restarts=%d)", d.pid(), d.restarts.Load())
			continue
		}
		log.Printf("daemon exited unexpectedly: %v", err)
		if d.retryInterval <= 0 {
			log.Fatalf("daemon retry disabled (--daemon.retry-interval=0); exiting to let K8s rebuild pod")
		}
		recovered := false
		for attempt := 1; attempt <= daemonMaxRetryAttempts; attempt++ {
			log.Printf("daemon restart attempt %d/%d in %v", attempt, daemonMaxRetryAttempts, d.retryInterval)
			if d.sleepOrStop(d.retryInterval) {
				return
			}
			if err := d.start(); err != nil {
				log.Printf("daemon restart attempt %d/%d: start failed: %v", attempt, daemonMaxRetryAttempts, err)
				continue
			}
			if err := d.waitReady(); err != nil {
				log.Printf("daemon restart attempt %d/%d: not ready: %v", attempt, daemonMaxRetryAttempts, err)
				continue
			}
			d.restarts.Add(1)
			log.Printf("daemon restarted on attempt %d/%d, pid=%d (total restarts=%d)", attempt, daemonMaxRetryAttempts, d.pid(), d.restarts.Load())
			recovered = true
			break
		}
		if !recovered {
			log.Fatalf("daemon failed all %d restart attempts at %v interval; exiting to let K8s rebuild pod", daemonMaxRetryAttempts, d.retryInterval)
		}
		d.schedulePostRecoveryExit()
	}
}

// sleepOrStop sleeps for d, returning true if interrupted by stop.
func (d *daemonSupervisor) sleepOrStop(dur time.Duration) bool {
	select {
	case <-time.After(dur):
		return false
	case <-d.stopCh:
		return true
	}
}

func (d *daemonSupervisor) stop() {
	d.stopping.Store(true)
	// wake periodicRestart() so it can return
	select {
	case <-d.stopCh:
	default:
		close(d.stopCh)
	}
	d.mu.Lock()
	if d.postRecoveryTimer != nil {
		d.postRecoveryTimer.Stop()
		d.postRecoveryTimer = nil
	}
	cmd := d.cmd
	d.mu.Unlock()
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Signal(syscall.SIGTERM)
		done := make(chan struct{})
		go func() { _ = cmd.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
		}
	}
}

func (d *daemonSupervisor) alive() bool {
	d.mu.Lock()
	cmd := d.cmd
	d.mu.Unlock()
	if cmd == nil || cmd.Process == nil {
		return false
	}
	// Signal 0: no-op signal that errors if process is gone.
	return cmd.Process.Signal(syscall.Signal(0)) == nil
}

func (d *daemonSupervisor) pid() int {
	d.mu.Lock()
	cmd := d.cmd
	d.mu.Unlock()
	if cmd == nil || cmd.Process == nil {
		return -1
	}
	return cmd.Process.Pid
}

func (d *daemonSupervisor) restartCount() int32 {
	return d.restarts.Load()
}
