.PHONY: build
build: compile

# dependencies

.PHONY: dep
dep:
	go get github.com/google/go-cloud/wire/cmd/wire

.PHONY: dep-go
go-dep:
	cd pkg/public_go && wire

# build

.PHONY: gazelle
gazelle:
	bazel run gazelle

.PHONY: compile
compile: dep-go gazelle
	bazel build //pkg/public_go:public_go

# test

.PHONY: test-go
test-go: dep-go gazelle
	bazel run //pkg/common_go/util:go_default_test
	bazel run //pkg/public_go:go_default_test
