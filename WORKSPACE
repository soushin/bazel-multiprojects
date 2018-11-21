workspace(name = "com_github_soushin_bazelmultiprojects")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# ================================================================
# Go support requires rules_go, bazel_gazelle
# ================================================================

http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.16.2/rules_go-0.16.2.tar.gz",
    sha256 = "f87fa87475ea107b3c69196f39c82b7bbf58fe27c62a338684c20ca17d1d8613",
)
http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.15.0/bazel-gazelle-0.15.0.tar.gz"],
    sha256 = "6e875ab4b6bf64a38c352887760f21203ab054676d9c1b274963907e0768740d",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

# ================================================================
# Go projects dependencies
# ================================================================

load("@bazel_gazelle//:deps.bzl", "go_repository")

go_repository(
    name = "com_github_google_cloud",
    commit = "2152f209f3c907645f7ebdcacdb2c18cd89e6fa8",
    importpath = "github.com/google/go-cloud",
)

go_repository(
    name = "org_golang_google_grpc",
    importpath = "github.com/grpc/grpc-go",
    tag = "v1.16.0",
)

# ================================================================
# Kotlin support requires rules_kotlin
# ================================================================

rules_kotlin_version = "87bd13f91d166a8070e9bbfbb0861f6f76435e7a"
http_archive(
    name = "io_bazel_rules_kotlin",
    urls = ["https://github.com/bazelbuild/rules_kotlin/archive/%s.zip" % rules_kotlin_version],
    type = "zip",
    strip_prefix = "rules_kotlin-%s" % rules_kotlin_version
)
load("@io_bazel_rules_kotlin//kotlin:kotlin.bzl", "kotlin_repositories", "kt_register_toolchains")
kotlin_repositories()
kt_register_toolchains()

# ================================================================
# Docker support requires rules_docker
# ================================================================

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "29d109605e0d6f9c892584f07275b8c9260803bf0c6fcb7de2623b2bedc910bd",
    strip_prefix = "rules_docker-0.5.1",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.5.1.tar.gz"],
)

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

container_pull(
    name = "java_base_image",
    registry = "index.docker.io",
    repository = "library/openjdk",
    tag = "12-ea-18-jdk-alpine3.8",
)

container_pull(
    name = "distroless_base_image",
    registry = "gcr.io",
    repository = "distroless/base",
    digest = "sha256:628939ac8bf3f49571d05c6c76b8688cb4a851af6c7088e599388259875bde20"
)

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)
_go_image_repos()

# ================================================================
# Protobuf support requires rules_proto
# ================================================================

http_archive(
	name = "build_stack_rules_proto",
	urls = ["https://github.com/stackb/rules_proto/archive/1d6550fc2e62.tar.gz"],
	sha256 = "113e6792f5b20679285c86d91c163cc8c4d2b4d24d7a087ae4f233b5d9311012",
	strip_prefix = "rules_proto-1d6550fc2e625d47dc4faadac92d7cb20e3ba5c5",
)
