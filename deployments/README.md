
## 镜像构建

执行当前目录下的 `build-image.sh` 可以自动构建好镜像。

## 自定义镜像名称

build-image.sh 里默认的镜像路径与名称为: `artifact.enflame.cn/enflame_docker_images/enflame/gcu-exporter:latest`，如下：

```conf
# Currently supports ubuntu, tlinux, openeuler
OS="ubuntu"

# Currently supports docker, ctr, podman, nerdctl
CLI_NAME="docker"

# The repository name
REPO_NAME="artifact.enflame.cn/enflame_docker_images/enflame"

# The image name
IMAGE_NAME="gcu-exporter"

# The image tag
TAG="latest"

# The namespace used by nerdctl, ctr
NAMESPACE="k8s.io"

```

可以根据自己的需要自定义这个镜像路径与名称。


## 部署示例

```bash
# x86
kubectl apply -f yaml/gcu-exporter.yaml
# ARM
kubectl apply -f yaml/gcu-exporter-for-arm.yaml
```

## 使用示例

参考 examples/README.md