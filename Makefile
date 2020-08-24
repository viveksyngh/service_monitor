TAG?=latest
NAMESPACE?=viveksyngh

.PHONY: build
build:
	docker build -t ${NAMESPACE}/service_monitor:${TAG} .

.PHONY: push
push:
	docker push ${NAMESPACE}/service_monitor:${TAG}
