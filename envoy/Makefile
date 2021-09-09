IMAGE=local-envoy
IMAGE_TAG=latest
NAMESPACE=ingress
CLUSTER=cluster01

build:
	docker build -t ${IMAGE}:${IMAGE_TAG} ./

kindly-push:
	kind load docker-image ${IMAGE}:${IMAGE_TAG} --name=${CLUSTER}

kindly-deploy:
	kubectl create namespace ${NAMESPACE} || true
	kubectl delete configmap envoy-config -n ${NAMESPACE} || true
	kubectl create configmap envoy-config --from-file=envoy.yaml=envoy-ext-authz-config.yaml -n ${NAMESPACE}
	kubectl apply -f deployment.yaml -n ${NAMESPACE}

all: build kindly-push kindly-deploy
