default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	ginkgo -r -timeout=120s ./...

testacc:
	TF_ACC=1 ginkgo -r -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
