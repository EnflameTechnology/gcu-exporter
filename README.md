## GCU-Exporter
This repository contains Prometheus GCU-Exporter, which exposes Enflame GCU metrics.

### Building from Source

```
git clone https://github.com/EnflameTechnology/gcu-exporter.git
cd gcu-exporter
./build.sh all

# This step will generate a package under dist folder:

dist/gcu-exporter_{VERSION}/
├── build-image.conf
├── build-image.sh
├── dockerfiles
│   ├── Dockerfile.openeuler
│   ├── Dockerfile.tlinux
│   └── Dockerfile.ubuntu
├── examples
│   ├── docker
│   ├── k8s
│   └── README.md
├── gcu-exporter
├── LICENSE.md
├── README.md
└── yaml
    ├── gcu-exporter-for-arm-non-privileged.yaml
    ├── gcu-exporter-for-arm.yaml
    ├── gcu-exporter-non-privileged.yaml
    └── gcu-exporter.yaml
```

### Building the container image

```
cd dist/gcu-exporter_{VERSION}/
# config build-image.conf
vim build-image.conf
build-image.sh
```


#### Collect metrics

To gather GCU metrics, to start the gcu-exporter.


## License

gcu-exporter is licensed under the Apache-2.0 license.

gcu-exporter was developed with reference to node_exporter (licensed under Apache-2.0, https://github.com/prometheus/node_exporter/blob/master/LICENSE). Many thanks to node_exporter.