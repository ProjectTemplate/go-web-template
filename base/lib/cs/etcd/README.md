## 部署安装

```shell
ETCD_VER=v3.5.17

rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
  docker rmi gcr.io/etcd-development/etcd:${ETCD_VER} || true && \
  docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  --name etcd-gcr-${ETCD_VER} \
  gcr.io/etcd-development/etcd:${ETCD_VER} \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --log-level info \
  --logger zap \
  --log-outputs stderr
  
docker exec etcd-gcr-${ETCD_VER}/usr/local/bin/etcd --version
docker exec etcd-gcr-${ETCD_VER}/usr/local/bin/etcdctl version
docker exec etcd-gcr-${ETCD_VER}/usr/local/bin/etcdutl version
docker exec etcd-gcr-${ETCD_VER}/usr/local/bin/etcdctl endpoint health
docker exec etcd-gcr-${ETCD_VER}/usr/local/bin/etcdctl put foo bar
docker exec etcd-gcr-${ETCD_VER}/usr/local/bin/etcdctl get foo
docker exec etcd-gcr-${ETCD_VER}/usr/local/bin/etcdctl get hello

```