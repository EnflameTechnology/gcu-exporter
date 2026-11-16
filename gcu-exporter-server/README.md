# GCU Exporter Server

GCU Exporter Server 是一个 HTTP 服务器，当收到请求时会执行 `gcu-exporter` 命令并将结果返回给客户端。

## 功能特性

- 接收 HTTP 请求
- 自动执行 `gcu-exporter` 命令生成 Prometheus 格式的指标
- 将生成的指标数据返回给请求端
- 支持配置超时时间

## 使用方法

### 编译

```bash
go build -o gcu-exporter-server
```

### 运行

```bash
./gcu-exporter-server
```

### 命令行参数

- `--web.listen-address`: 监听地址（默认: `:9400`）
- `--web.telemetry-path`: 指标路径（默认: `/metrics`）
- `--web-config-file`: Web 配置文件路径，可启用 HTTP Basic Auth（默认: 空，不启用鉴权）
- `--exporter.path`: gcu-exporter 二进制文件路径（默认: `./gcu-exporter`）
- `--data.file`: 临时数据文件路径（默认: `/tmp/gcu_monitor_data`）
- `--exporter.timeout`: 执行 exporter 命令的超时时间，单位：秒（默认: `30`）

### HTTP Basic Auth

架构

```
客户端 → :9400 (server，无鉴权) → 转发 Authorization → :9410 (daemon，web-config-file)
```

对外 HTTP 入口（默认 `:9400`）支持 `--web-config-file` 参数， 它会将 `--web-config-file` 转发到 `daemon` 进程(默认 `:9410`) 启用 Basic Auth，配置格式与 [exporter-toolkit](https://github.com/prometheus/exporter-toolkit) / node_exporter 一致。未指定该参数时行为与旧版本完全一致。

配置文件示例（密码须为 bcrypt hash）：

```yaml
basic_auth_users:
  prometheus: $2y$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p61LGi3E17VWK18hQQ2aHy
```

生成 bcrypt hash：

```bash
python3 -c 'import bcrypt; print(bcrypt.hashpw(b"my-secret-password", bcrypt.gensalt()).decode())'
# 或: htpasswd -nBC 10 "" | tr -d ':\n'
```

启动示例：

```bash
./gcu-exporter-server --web-config-file=/etc/gcu-exporter/web-config.yml
```

验证：

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:9400/metrics          # 期望 401
curl -s -u prometheus:my-secret-password http://localhost:9400/metrics | head   # 期望 200
```

Prometheus scrape 侧配置明文密码（与 hash 对应）：

```yaml
scrape_configs:
  - job_name: gcu-exporter
    basic_auth:
      username: prometheus
      password_file: /etc/prometheus/secrets/gcu-exporter-password
      # OR
      # password: my-secret-password
    static_configs:
      - targets: ['localhost:9400']
```

### 未开启验证 示例

```bash
# 使用默认配置
./gcu-exporter-server

# 指定监听地址和 exporter 路径
./gcu-exporter-server \
  --web.listen-address=:8080 \
  --exporter.path=/usr/local/bin/gcu-exporter \
  --data.file=/tmp/gcu_metrics.txt

# 访问指标
curl http://localhost:9400/metrics
```

## 工作原理

1. 服务器启动并监听指定的地址和路径
2. 当收到 HTTP 请求时：
   - 执行 `gcu-exporter --output.format prometheus --output.pretty --output.file /tmp/gcu_monitor_data`
   - 等待命令执行完成
   - 读取生成的数据文件
   - 将文件内容作为 HTTP 响应返回给客户端

## 注意事项

- 确保 `gcu-exporter` 二进制文件存在且可执行
- 确保有权限写入临时数据文件
- 确保有权限读取临时数据文件

