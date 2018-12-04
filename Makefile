PACKAGE = public_go

.PHONY: build
build: compile

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
	go get github.com/google/go-cloud/wire/cmd/wire
	cd pkg/${PACKAGE} && wire

# build

.PHONY: gazelle
gazelle:
	bazel run gazelle
	bazel run gazelle -- update-repos -from_file ./Gopkg.lock

.PHONY: compile
compile: gazelle dep-wire gen-proto
	bazel query //... | grep "//pkg/${PACKAGE}" | xargs bazel build --define IMAGE_TAG=test

.PHONY: run
run: gazelle dep-wire gen-proto
	bazel query //... | grep "//pkg/${PACKAGE}" | xargs bazel run --define IMAGE_TAG=test

# proto

.PHONY: gen-proto
gen-proto:
	$(eval protos = $(shell find ./proto -type d -d 1 | sed 's/\.\/proto\///g' | xargs))
	@for p in $(protos); do bazel build //proto/$$p:proto_buf && mv -f bazel-genfiles/proto/$$p/proto_buf/proto/$$p/$$p.pb.go proto/$$p; done

# test

.PHONY: test-go
test-go: gazelle dep-wire
	bazel query //... | grep "//pkg/${PACKAGE}" | xargs bazel test --define IMAGE_TAG=latest

.PHONY: test-go-all
test-go-all: gazelle dep-wire
	bazel query //... | grep "//pkg" | xargs bazel test --define IMAGE_TAG=latest

# container

.PHONY: container-build
container-build:
	bazel query //... | grep "//pkg/${PACKAGE}:container_image" | xargs bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64

.PHONY: container-push
container-push:
	$(eval image = $(shell git rev-parse --abbrev-ref @ | sed 's/\//_/g'))
	bazel query //... | grep "//pkg/${PACKAGE}:container_push" | xargs bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 --define IMAGE_TAG=$(image)
