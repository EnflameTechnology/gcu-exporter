package collector

import (
	"testing"
	"github.com/prometheus/client_golang/prometheus"
)

func TestConstants(t *testing.T) {
	expectedNamespace := "enflame_gcu"
	if namespace != expectedNamespace {
		t.Errorf("Expected namespace to be %s, but got %s", expectedNamespace, namespace)
	}

	expectedVgcuNamespace := "enflame_vgcu"
	if vgcu_namespace != expectedVgcuNamespace {
		t.Errorf("Expected vgcu_namespace to be %s, but got %s", expectedVgcuNamespace, vgcu_namespace)
	}
}

func TestScrapeDurationDesc(t *testing.T) {
	expectedScrapeDurationDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"gcu-exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	if scrapeDurationDesc.String() != expectedScrapeDurationDesc.String() {
		t.Errorf("Expected scrapeDurationDesc to be %s, but got %s", expectedScrapeDurationDesc, scrapeDurationDesc)
	}
}

func TestScrapeSuccessDesc(t *testing.T) {
	expectedScrapeSuccessDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"gcu-exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
	if scrapeSuccessDesc.String() != expectedScrapeSuccessDesc.String() {
		t.Errorf("Expected scrapeSuccessDesc to be %s, but got %s", expectedScrapeSuccessDesc, scrapeSuccessDesc)
	}
}