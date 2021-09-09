# envoy-ext-authz-server

Sample implementation for Envoy ExtAuthz Server

## Deployment steps

### Deploy KIND Cluster
Follow instructions here: [kind-clusters](https://github.com/abhide/kind-clusters#create-new-cluster)
```
➜  kind-clusters git:(master) make init 
kind create cluster --config=cluster01.yaml --name=cluster01 --wait=300s
Creating cluster "cluster01" ...
 ✓ Ensuring node image (kindest/node:v1.19.1) 🖼
 ✓ Preparing nodes 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
 ✓ Waiting ≤ 5m0s for control-plane = Ready ⏳ 
 • Ready after 29s 💚
Set kubectl context to "kind-cluster01"
You can now use your cluster with:

kubectl cluster-info --context kind-cluster01

Thanks for using kind! 😊
```

### Deploy simple-httpapp
Follow instructions here: [simple-httpapp](https://github.com/abhide/simple-httpapp)
```
➜  simple-httpapp git:(master) ✗ make all 
go fmt ./
docker build -t simple-httpapp:latest ./
Sending build context to Docker daemon  81.92kB
Step 1/8 : FROM golang:alpine3.12
 ---> baed0e68a17f
Step 2/8 : WORKDIR /go/src/github.com/abhide/simple-httpapp/
 ---> Running in 5667c56ac6c0
Removing intermediate container 5667c56ac6c0
 ---> 0894b2afded0
Step 3/8 : COPY main.go .
 ---> 78a8f329facc
Step 4/8 : RUN go build -o simple-httpapp ./main.go
 ---> Running in a2cbc2f34918
Removing intermediate container a2cbc2f34918
 ---> ef4d11816016
Step 5/8 : FROM alpine:3.12
 ---> 48b8ec4ed9eb
Step 6/8 : WORKDIR /root/
 ---> Using cache
 ---> 8fa8bc160dc3
Step 7/8 : COPY --from=0 /go/src/github.com/abhide/simple-httpapp/simple-httpapp .
 ---> 331606d5aa23
Step 8/8 : CMD ["./simple-httpapp"]
 ---> Running in c0af969c8ac5
Removing intermediate container c0af969c8ac5
 ---> 4feeec8c0ed8
Successfully built 4feeec8c0ed8
Successfully tagged simple-httpapp:latest
kind load docker-image simple-httpapp:latest --name=cluster01
Image: "simple-httpapp:latest" with ID "sha256:4feeec8c0ed8aef3f732d81d25b987b745ad7af1124624b7cc71708493345892" not yet present on node "cluster01-control-plane", loading...
kubectl create namespace v1 || true
namespace/v1 created
kubectl apply -f k8s/simpleapp-v1.yaml -n v1
configmap/simple-httpapp-v1-config created
deployment.apps/simple-httpapp-v1 created
service/simple-httpapp-v1-svc created
kubectl create namespace v2 || true
namespace/v2 created
kubectl apply -f k8s/simpleapp-v2.yaml -n v2
configmap/simple-httpapp-v2-config created
deployment.apps/simple-httpapp-v2 created
service/simple-httpapp-v2-svc created
```

### Build and Deploy Envoy ExtAuthz Server
```
➜  envoy-ext-authz-server git:(main) ✗ make all 
docker build -t envoy-ext-authz-server:latest ./
Sending build context to Docker daemon  14.42MB
Step 1/10 : FROM golang:alpine3.12
 ---> baed0e68a17f
Step 2/10 : WORKDIR /go/src/github.com/abhide/envoy-ext-authz-server/
 ---> Using cache
 ---> ec125cf211b9
Step 3/10 : COPY main.go .
 ---> Using cache
 ---> 35e482cff2fd
Step 4/10 : COPY go.mod .
 ---> Using cache
 ---> 4ec45d132a5a
Step 5/10 : COPY go.sum .
 ---> Using cache
 ---> 89f853627f4f
Step 6/10 : RUN go build -o envoy-ext-authz-server ./main.go
 ---> Using cache
 ---> afe5ba61280f
Step 7/10 : FROM alpine:3.12
 ---> 48b8ec4ed9eb
Step 8/10 : WORKDIR /root/
 ---> Using cache
 ---> 8fa8bc160dc3
Step 9/10 : COPY --from=0 /go/src/github.com/abhide/envoy-ext-authz-server/envoy-ext-authz-server .
 ---> Using cache
 ---> 76eec2ab5ef0
Step 10/10 : CMD ["./envoy-ext-authz-server"]
 ---> Using cache
 ---> 3ca810873dfd
Successfully built 3ca810873dfd
Successfully tagged envoy-ext-authz-server:latest
kind load docker-image envoy-ext-authz-server:latest --name=cluster01
Image: "envoy-ext-authz-server:latest" with ID "sha256:3ca810873dfd12f573c5e6799cce9f418c2dee258110d2029b125d5d4333aacd" not yet present on node "cluster01-control-plane", loading...
kubectl create namespace ext-authz-server || true
namespace/ext-authz-server created
kubectl apply -f k8s/deploy.yaml -n ext-authz-server
deployment.apps/envoy-ext-authz-server created
service/envoy-ext-authz-server-svc created
```

### Deploy Envoy with ExtAuthz Configuration
```
```