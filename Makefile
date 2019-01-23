PACKAGE = public_go
NAMESPACE = default

.PHONY: build
build: compile

.PHONY: clean
clean: bazel clean

# dependencies

.PHONY: install
install:
	brew install bazel
	brew install skaffold

.PHONY: dep
dep:
	dep ensure

.PHONY: dep-update
dep-update:
	rm -rf ./vendor
	dep ensure -update

.PHONY: dep-wire
dep-wire:
	go get github.com/google/wire/cmd/wire
	cd pkg/${PACKAGE} && wire

# build

.PHONY: gazelle
gazelle:
	bazel run gazelle
	bazel run gazelle -- update-repos -from_file ./Gopkg.lock

.PHONY: compile
compile: gazelle gen-proto
	bazel build //pkg/${PACKAGE}:${PACKAGE}

.PHONY: run
run: gazelle gen-proto
	bazel run //pkg/${PACKAGE}:${PACKAGE}

# proto

.PHONY: gen-proto
gen-proto:
	$(eval protos = $(shell find ./proto -type d -d 1 | sed 's/\.\/proto\///g' | xargs))
	@for p in $(protos); do bazel build //proto/$$p:proto_buf && mv -f bazel-genfiles/proto/$$p/proto_buf/proto/$$p/*.pb.go proto/$$p; done

# test

.PHONY: test-go
test-go: gazelle dep-wire
	bazel test //pkg/${PACKAGE}:go_default_test

.PHONY: test-go-all
test-go-all: gazelle dep-wire
	bazel query //... | grep "go_default_test" | xargs bazel test

# container

.PHONY: container-build
container-build:
	bazel query //... | grep "//pkg/${PACKAGE}:container_image" | xargs bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64

.PHONY: container-push
container-push:
	$(eval image = $(shell git rev-parse --abbrev-ref @ | sed 's/\//_/g'))
	bazel query //... | grep "//pkg/${PACKAGE}:container_push" | xargs bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 --define IMAGE_TAG=$(image)

# deploy
.PHONY: deploy
deploy: gazelle
	skaffold deploy -f pkg/${PACKAGE}/k8s/overlays/${NAMESPACE}/skaffold.yaml
