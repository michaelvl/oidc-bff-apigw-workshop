.PHONY: deploy-all
deploy-all: create-cluster deploy-metallb deploy-istio-base

.PHONY: deploy-istio-base
deploy-istio-base: deploy-gateway-api deploy-istio-ambient

.PHONY: create-cluster
create-cluster:
	kind create cluster --config test/kind-config

.PHONY: delete-cluster
delete-cluster:
	kind delete cluster -n kind

#################
.PHONY: deploy-gateway-api
deploy-gateway-api:
	kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd/experimental?ref=v1.0.0" | kubectl apply -f -

.PHONY: deploy-metallb
deploy-metallb:
	kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.12/config/manifests/metallb-native.yaml
	kubectl wait --namespace metallb-system --for=condition=ready pod --selector=app=metallb --timeout=90s
	scripts/kind-metallb-configure.sh

ISTIO_VERSION := 1.20.0

.PHONY: deploy-istio-ambient
deploy-istio-ambient:
	kubectl create namespace istio-system
	helm upgrade -i istio-base --repo https://istio-release.storage.googleapis.com/charts --version $(ISTIO_VERSION) base    -n istio-system --set defaultRevision=default
	helm upgrade -i istiod     --repo https://istio-release.storage.googleapis.com/charts --version $(ISTIO_VERSION) istiod  -n istio-system --values test/istiod-values.yaml --wait
	helm upgrade -i istio-cni  --repo https://istio-release.storage.googleapis.com/charts --version $(ISTIO_VERSION) cni --values test/istio-cni-values.yaml
	helm upgrade -i ztunnel    --repo https://istio-release.storage.googleapis.com/charts  --version $(ISTIO_VERSION) ztunnel -n istio-system

#################
IMAGE ?= ghcr.io/michaelvl/oidc-bff-apigw-workshop:latest

.PHONY: container
container:
	docker build -t $(IMAGE) .

.PHONY: kind-load-image
kind-load-image:
	kind load docker-image $(IMAGE) --name kind

#################
.PHONY: run-cdn-local
run-cdn-local:
	docker run --rm -e STATIC_FILES_PATH=/apps/spa -p 5030:5030 $(IMAGE)

.PHONY: deploy-spa
deploy-spa:
	kubectl apply -f kubernetes/spa-cdn.yaml
	kubectl apply -f kubernetes/spa-redis-session-store.yaml
	kubectl apply -f kubernetes/spa-login-bff.yaml
	kubectl apply -f kubernetes/gateway-httproutes.yaml
