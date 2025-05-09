
## 版本申明

| 版本 | 修改内容                                   | 修改时间   |
| ---- | ------------------------------------------ | ---------- |
| v1.0 | 将之前的版本更新为Prometheus + K8S示例说明   | 10/10/2022 |
| v1.1 | 简化内容，更新说明                          | 6/30/2023 |
| v1.2 | 更新了一些格式与内容                        | 4/8/2024 |
| v1.3 | 更新一些格式                               | 4/9/2024 |
| v1.4 | 更新一些内容                               | 1/7/2025 |


## 简介

Gcu-exporter/examples是一个基于prometheus + grafana + gcu-exporter的可观测方案应用简单示例。从gcu-exporter采集的指标通过 Prometheus展示到Grafana以便用户获取或设置GCU设备的运行指标与告警信息。 相应的运行指标说明参考《GCU Exporter用户使用手册》，用户可以根据自己的具体需要进行可观测方案的定制化二次开发。

## 应用示例

### gcu-exporter 镜像构建

cd gcu-exporter_{VERSION}
执行，镜像构建脚本：

```
./build-image.sh
```


### k8s部署示例

在K8s集群已经部署好的前提下，采用`k8s/k8s-deploy.sh`脚本，如 ：

```bash
cd k8s
./k8s-deploy.sh --help

Usage: k8s-deploy.sh [command]

Commands:
    apply    Apply k8s yaml
    delete   Delete k8s yaml

Examples:
    k8s-deploy.sh apply
    k8s-deploy.sh delete
```

拉起运行指标观测示例，执行：

```bash
./k8s-deploy.sh apply
```
下线运行指标观测示例，执行：

```bash
./k8s-deploy.sh delete
```

以上步骤需要注意先配置`yaml/gcu-exporter.yaml`里的`gcu-exporter`镜像路径，镜像路径需要根据本地的实际情况修改，如下：

```yaml
      containers:
        - name: gcu-exporter
          image: artifact.enflame.cn/enflame_docker_images/enflame/gcu-exporter:latest
          imagePullPolicy: IfNotPresent #Always
          securityContext:
            privileged: true
```

另外用户还可以按需自我定制yaml目录下的yaml文件：

```yaml
yaml/gcu-exporter.yaml
yaml/grafana.yaml
yaml/namespace.yaml
yaml/prometheus.yaml
```

以上文件中prometheus 与 grafana的镜像如果连不了外网需要先下载后再导入。


### docker部署示例

docker部署示例，采用`docker/docker-compose.sh`，如 ：

```bash
cd docker
./docker-compose.sh --help

Usage: docker-compose.sh [command]

Commands:
    init        init docker compose
    up          docker compose up -d
    down        docker compose down

Examples:
    docker-compose.sh up
    docker-compose.sh down
```


拉起运行指标观测示例，执行：

```bash
./docker-compose.sh up
```

这一步如果出现 `docker-compose: command not found` 这样的log，如下：

```
[INFO] Action start : up
[INFO] docker-compose up -d
./docker-compose.sh: line 60: docker-compose: command not found
[ERROR] Action is failed : up

```

则需要先安装docker-compose 命令：

```bash
cd docker
./docker-compose.sh init
```

关闭运行指标观测示例，执行：

```bash
./docker-compose.sh down
```

以上示例需要先根据本地实际情况配置docker-compose.yaml里 `gcu-exporter`的镜像路径，例如`image: artifact.enflame.cn/enflame_docker_images/enflame/gcu-exporter:latest`。
prometheus 与 grafana的镜像如果连不了外网需要先下载后再导入，docker-compose.yaml 内容如下：

```yaml
version: '2.0'

services:
    prometheus:
        container_name: prometheus
        image: prom/prometheus:latest
        volumes:
            - ./prom/prometheus.yml:/etc/prometheus/prometheus.yml:ro
        ports:
            - 9090:9090
        network_mode: host

    grafana:
        container_name: grafana
        image: grafana/grafana:latest
        volumes:
            - /var/lib/grafana:/var/lib/grafana
        ports:
            - 3000:3000
        network_mode: host

    gcu-exporter:
        container_name: gcu-exporter
        image: artifact.enflame.cn/enflame_docker_images/enflame/gcu-exporter:latest
        privileged: true
        volumes:
            - /usr/lib/libefml.so:/usr/lib/libefml.so
            - /usr/local/efsmi:/usr/local/efsmi
        ports:
            - 9400:9400
        network_mode: host

```


### 通过Prometheus 查看运行指标

通过浏览器访问prometheus服务，访问http://\<NodeIP\>:9090, prometheus默认端口9090（注意配置K8S的端口映射）， 依次选择status-\> target 查看endpoint status，如果每个服务的status 为 UP代表节点运行程序正常启动，如果为DOWN 则代表节点运行程序异常。

### 通过Grafana查看运行指标

注：如果Grafana版本不一致，以下步骤与过程也可能会不一致，需要根据具体情况进行调整。

#### 登录Grafana web界面

在浏览器地址栏输入grafana服务的IP和端口，**http://\<NodeIp\>:3000** , grafana的默认端口是3000， 默认 **Username: admin**， 默认：**Password: admin**。


#### 添加Grafana数据源

在Grafana的首页里点击 `Add your first data source`选择`Prometheus` 作为数据源，再根据Prometheus的配置选项提示配置相应信息，比如在 Prometheus 的配置选项 URL 里 填写 `http://localhost:9090` , 然后再点 左下角的 `Save & Test` ，即可完成Prometheus的简单配置。


## 注意事项

- 本示例仅提供了gcu-exporter 的Prometheus + Grafana 使用简单示例。如果要在生产上使用，建议用户可以根据自己的具体要求参考运行指标说明文档 《GCU Exporter用户使用手册》，进行合理的二次开发；

- **注：本应用示例仅供参考，而非一键开箱即用方案。**



