# Quick Start Guides

# Table of Contents
* [ClickHouse Operator Installation](#clickhouse-operator-installation)
* [Building ClickHouse Operator from Sources](#building-clickhouse-operator-from-sources)
* [Examples](#examples)
  * [Trivial Example](#trivial-example)
  * [Connect to ClickHouse Database](#connect-to-clickhouse-database)
  * [Simple Persistent Volume Example](#simple-persistent-volume-example)
  * [Custom Deployment with Pod and VolumeClaim Templates](#custom-deployment-with-pod-and-volumeclaim-templates)
  * [Custom Deployment with Specific ClickHouse Configuration](#custom-deployment-with-specific-clickhouse-configuration)

# Prerequisites
1. Operational Kubernetes instance
1. Properly configured `kubectl`
1. `curl`

# ClickHouse Operator Installation

Apply `clickhouse-operator` installation manifest. The simplest way - directly from `github`.

In case you are convenient to install operator into `kube-system` namespace just run:
```bash
kubectl apply -f https://raw.githubusercontent.com/Altinity/clickhouse-operator/master/manifests/operator/clickhouse-operator-install.yaml
``` 

In case you'd like to customize namespace where to install operator, or customize any operator's or ClickHouse's configuration, feel free to use installer script.
Please, `cd` into writable folder, because install script would download config files to build `.yaml` manifests from into current dir. 
```bash
cd ~
curl -s https://raw.githubusercontent.com/Altinity/clickhouse-operator/master/manifests/operator-installer/clickhouse-operator-install.sh | CHOPERATOR_NAMESPACE=test-clickhouse-operator bash
```
Take into account explicitly specified namespace
```bash
CHOPERATOR_NAMESPACE=test-clickhouse-operator
```
This namespace would be created and used to install `clickhouse-operator` into.
Install script would download some `.yaml` and `.xml` files and install `clickhouse-operator` into specified namespace.

In case no `CHOPERATOR_NAMESPACE` specified, as:
```bash
cd ~
curl -s https://raw.githubusercontent.com/Altinity/clickhouse-operator/master/manifests/operator-installer/clickhouse-operator-install.sh | bash
```
installer will create namespace `clickhouse-operator` and install **clickhouse-operator** into it.

Operator installation process
```text
Setup ClickHouse Operator into test-clickhouse-operator namespace
namespace/test-clickhouse-operator created
customresourcedefinition.apiextensions.k8s.io/clickhouseinstallations.clickhouse.altinity.com configured
serviceaccount/clickhouse-operator created
clusterrolebinding.rbac.authorization.k8s.io/clickhouse-operator configured
service/clickhouse-operator-metrics created
configmap/etc-clickhouse-operator-files created
configmap/etc-clickhouse-operator-confd-files created
configmap/etc-clickhouse-operator-configd-files created
configmap/etc-clickhouse-operator-templatesd-files created
configmap/etc-clickhouse-operator-usersd-files created
deployment.apps/clickhouse-operator created
```

Check `clickhouse-operator` is running:
```bash
kubectl get pods -n test-clickhouse-operator
```
```text
NAME                                 READY   STATUS    RESTARTS   AGE
clickhouse-operator-5ddc6d858f-drppt 1/1     Running   0          1m
```

## Building ClickHouse Operator from Sources

Complete instructions on how to build ClickHouse operator from sources as well as how to build a docker image and use it inside `kubernetes` described [here][build_from_sources].

# Examples

There are several ready-to-use [examples](./examples/). Below are few ones to start with.

## Create Custom Namespace
It is a good practice to have all components run in dedicated namespaces. Let's run examples in `test` namespace
```bash
kubectl create namespace test
```
```text
namespace/test created
```

## Trivial example

This is the trivial [1 shard 1 replica](./examples/01-simple-layout-01-1shard-1repl.yaml) example.

**WARNING**: Do not use it for anything other than 'Hello, world!', it does not have persistent storage!
 
```bash
kubectl apply -n test-clickhouse-operator -f https://raw.githubusercontent.com/Altinity/clickhouse-operator/master/docs/examples/01-simple-layout-01-1shard-1repl.yaml
```
```text
clickhouseinstallation.clickhouse.altinity.com/simple-01 created
```

Installation specification is straightforward and defines 1-replica cluster:
```yaml
apiVersion: "clickhouse.altinity.com/v1"
kind: "ClickHouseInstallation"
metadata:
  name: "simple-01"
```

Once cluster is created, there are two checks to be made.

```bash
kubectl get pods -n test-clickhouse-operator
```
```text
NAME                    READY   STATUS    RESTARTS   AGE
chi-b3d29f-a242-0-0-0   1/1     Running   0          10m
```

Watch out for 'Running' status. Also check services created by an operator:

```bash
kubectl get service -n test-clickhouse-operator
```
```text
NAME                    TYPE           CLUSTER-IP       EXTERNAL-IP                          PORT(S)                         AGE
chi-b3d29f-a242-0-0     ClusterIP      None             <none>                               8123/TCP,9000/TCP,9009/TCP      11m
clickhouse-example-01   LoadBalancer   100.64.167.170   abc-123.us-east-1.elb.amazonaws.com   8123:30954/TCP,9000:32697/TCP   11m
```

ClickHouse is up and running!

## Connect to ClickHouse Database

There are two ways to connect to ClickHouse database

1. In case previous command `kubectl get service -n test` reported **EXTERNAL-IP** (abc-123.us-east-1.elb.amazonaws.com in our case) we can directly access ClickHouse with:
```bash
clickhouse-client -h abc-123.us-east-1.elb.amazonaws.com
```
```text
ClickHouse client version 18.14.12.
Connecting to abc-123.us-east-1.elb.amazonaws.com:9000.
Connected to ClickHouse server version 19.4.3 revision 54416.
``` 
1. In case there is not **EXTERNAL-IP** available, we can access ClickHouse from inside Kubernetes cluster
```bash
kubectl -n test-clickhouse-operator exec -it chi-b3d29f-a242-0-0-0 -- clickhouse-client
```
```text
ClickHouse client version 19.4.3.11.
Connecting to localhost:9000 as user default.
Connected to ClickHouse server version 19.4.3 revision 54416.
```

## Simple Persistent Volume Example

In case of having Dynamic Volume Provisioning available - ex.: running on AWS - we are able to use PersistentVolumeClaims
Manifest is [available in examples](./examples/02-simple-layout-01-1shard-1repl-simple-persistent-volume.yaml)

```yaml
apiVersion: "clickhouse.altinity.com/v1"
kind: "ClickHouseInstallation"
metadata:
  name: "simple-02-pv1"
spec:
  defaults:
    templates:
      volumeClaimTemplate: volumeclaim-template
  configuration:
    clusters:
      - name: "simple-pv"
        layout:
          shardsCount: 1
          replicasCount: 1
  templates:
    volumeClaimTemplates:
      - name: volumeclaim-template
#        reclaimPolicy: Retain
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 500Mi
```

## Custom Deployment with Pod and VolumeClaim Templates

Let's install more complex example with:
1. Deployment specified
1. Pod template
1. VolumeClaim template

Manifest is [available in examples](./examples/02-simple-layout-03-1shard-1repl-deployment-persistent-volume.yaml)

```yaml
apiVersion: "clickhouse.altinity.com/v1"
kind: "ClickHouseInstallation"
metadata:
  name: "simple-02-pv2"
spec:
  configuration:
    clusters:
      - name: "deployment-pv"
        # Templates are specified for this cluster explicitly
        templates:
          podTemplate: pod-template-with-volume
          volumeClaimTemplate: storage-vc-template
        layout:
          shardsCount: 1
          replicasCount: 1

  templates:
    podTemplates:
      - name: pod-template-with-volume
        spec:
          containers:
            - name: clickhouse
              image: yandex/clickhouse-server:19.3.7
              ports:
                - name: http
                  containerPort: 8123
                - name: client
                  containerPort: 9000
                - name: interserver
                  containerPort: 9009
              volumeMounts:
                - name: storage-vc-template
                  mountPath: /var/lib/clickhouse

    volumeClaimTemplates:
      - name: storage-vc-template
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 1Gi
```

## Custom Deployment with Specific ClickHouse Configuration

You can tell operator to configure your ClickHouse, as shown in the example below ([link to the manifest](./examples/03-settings-01-overview.yaml)):

```yaml
apiVersion: "clickhouse.altinity.com/v1"
kind: "ClickHouseInstallation"
metadata:
  name: "settings-01"
spec:
  configuration:
    users:
      # test user has 'password' specified, while admin user has 'password_sha256_hex' specified
      test/password: qwerty
      test/networks/ip:
        - "127.0.0.1/32"
        - "192.168.74.1/24"
      test/profile: test_profile
      test/quota: test_quota
      test/allow_databases/database:
        - "dbname1"
        - "dbname2"
        - "dbname3"
      # admin use has 'password_sha256_hex' so actual password value is not published
      admin/password_sha256_hex: 8bd66e4932b4968ec111da24d7e42d399a05cb90bf96f587c3fa191c56c401f8
      admin/networks/ip: "127.0.0.1/32"
      admin/profile: default
      admin/quota: default
      # readonly user has 'password' field specified, not 'password_sha256_hex' as admin user above
      readonly/password: readonly_password
      readonly/profile: readonly
      readonly/quota: default
    profiles:
      test_profile/max_memory_usage: "1000000000"
      test_profile/readonly: "1"
      readonly/readonly: "1"
    quotas:
      test_quota/interval/duration: "3600"
    settings:
      compression/case/method: zstd
      disable_internal_dns_cache: 1
    files:
      dict1.xml: |
        <yandex>
            <!-- ref to file /etc/clickhouse-data/config.d/source1.csv -->
        </yandex>
      source1.csv: |
        a1,b1,c1,d1
        a2,b2,c2,d2
    clusters:
      - name: "standard"
        layout:
          shardsCount: 1
          replicasCount: 1
```

[build_from_sources]: ./operator_build_from_sources.md
