.PHONY: all
all: add
	@kubectl apply -k .

.PHONY: clean
clean:
	@kubectl delete -k .

.PHONY: build
build: add
	@kustomize build .

.PHONY: add
add:
	@kustomize edit add resource resources/*.yaml
