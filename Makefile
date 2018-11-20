.PHONY: build
build: compile

# dependencies

.PHONY: dep
dep:
	go get github.com/google/go-cloud/wire/cmd/wire

.PHONY: dep-go
dep-go:
	cd pkg/public_go && wire

# build

.PHONY: gazelle
gazelle:
	bazel run gazelle

.PHONY: compile
compile: dep-go gazelle gen-proto
	bazel build //pkg/public_go:public_go

# proto

.PHONY: gen-proto
gen-proto:
	$(eval protos = $(shell find ./proto -type d -d 1 | sed 's/\.\/proto\///g' | xargs))
	@for p in $(protos); do bazel build //proto/$$p:proto_buf && mv -f bazel-genfiles/proto/$$p/proto_buf/proto/$$p/$$p.pb.go proto/$$p; done

# test

.PHONY: test-go
test-go: dep-go gazelle
	bazel run //pkg/common_go/util:go_default_test
	bazel run //pkg/public_go:go_default_test


