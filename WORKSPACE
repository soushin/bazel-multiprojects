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

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# ================================================================
# Kotlin support requires rules_kotlin
# ================================================================

rules_kotlin_version = "87bd13f91d166a8070e9bbfbb0861f6f76435e7a"

http_archive(
    name = "io_bazel_rules_kotlin",
    urls = ["https://github.com/bazelbuild/rules_kotlin/archive/%s.zip" % rules_kotlin_version],
    type = "zip",
    strip_prefix = "rules_kotlin-%s" % rules_kotlin_version,
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
    digest = "sha256:628939ac8bf3f49571d05c6c76b8688cb4a851af6c7088e599388259875bde20",
)

container_pull(
    name = "ops_base_image",
    registry = "index.docker.io",
    repository = "soushin/bazel-multiprojects-ops-base",
    digest = "sha256:5c272657409aa2f1b85f56003a34d685e7f80bbbf8c19eb001f9b9b1f34644fa",
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
    urls = ["https://github.com/stackb/rules_proto/archive/d86ca6bc56b1589677ec59abfa0bed784d6b7767.tar.gz"],
    sha256 = "36f11f56f6eb48a81eb6850f4fb6c3b4680e3fc2d3ceb9240430e28d32c47009",
    strip_prefix = "rules_proto-d86ca6bc56b1589677ec59abfa0bed784d6b7767",
)

# ================================================================
# Go projects dependencies
# ================================================================

go_repository(
    name = "com_github_golang_protobuf",
    commit = "aa810b61a9c79d51363740d207bb46cf8e620ed5",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "org_golang_google_genproto",
    commit = "bd91e49a0898e27abb88c339b432fa53d7497ac0",
    importpath = "google.golang.org/genproto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "610586996380ceef02dd726cc09df7e00a3f8e56",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "2a47403f2ae58167c4d75973960ccc62b28cb0d8",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "org_golang_x_text",
    commit = "f21a4dfb5e38f5895301dc265a8def02365cc3d0",
    importpath = "golang.org/x/text",
)

go_repository(
    name = "com_github_google_go_cloud",
    commit = "1929e0c4fa0bd3defff57736e5f821d9983aad91",
    importpath = "github.com/google/go-cloud",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "df014850f6dee74ba2fc94874043a9f3f75fbfd8",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "com_github_ghodss_yaml",
    commit = "0ca9ea5df5451ffdf184b4428c902747c2c11cd7",
    importpath = "github.com/ghodss/yaml",
)

go_repository(
    name = "com_github_gogo_protobuf",
    commit = "636bf0302bc95575d69441b25a2603156ffdddf1",
    importpath = "github.com/gogo/protobuf",
)

go_repository(
    name = "com_github_google_gofuzz",
    commit = "24818f796faf91cd76ec7bddd72458fbced7a6c1",
    importpath = "github.com/google/gofuzz",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    commit = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75",
    importpath = "github.com/inconshreveable/mousetrap",
)

go_repository(
    name = "com_github_spf13_cobra",
    commit = "ef82de70bb3f60c65fb8eebacbb2d122ef517385",
    importpath = "github.com/spf13/cobra",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "298182f68c66c05229eb03ac171abe6e309ee79a",
    importpath = "github.com/spf13/pflag",
)

go_repository(
    name = "in_gopkg_inf_v0",
    commit = "d2d2541c53f18d2a059457998ce2876cc8e67cbf",
    importpath = "gopkg.in/inf.v0",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "51d6538a90f86fe93ac480b35f37b2be17fef232",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "io_k8s_api",
    build_file_proto_mode = "disable",
    commit = "8b7507fac302640dd5f1efbf9643199952cc58db",
    importpath = "k8s.io/api",
)

go_repository(
    name = "io_k8s_apimachinery",
    build_file_proto_mode = "disable",
    commit = "af2f90f9922d5bb462e28b799bd27342aeb794b3",
    importpath = "k8s.io/apimachinery",
)

go_repository(
    name = "io_k8s_klog",
    commit = "a5bc97fbc634d635061f3146511332c7e313a55a",
    importpath = "k8s.io/klog",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "66b9c49e59c6c48f0ffce28c2d8b8a5678502c6d",
    importpath = "github.com/gorilla/websocket",
)

go_repository(
    name = "com_github_nlopes_slack",
    commit = "38ac5e8c495bc1d4d299c6a14930376f32145eb8",
    importpath = "github.com/nlopes/slack",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "645ef00459ed84a119197bfb8d8205042c6df63d",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "com_github_kelseyhightower_envconfig",
    commit = "f611eb38b3875cc3bd991ca91c51d06446afa14c",
    importpath = "github.com/kelseyhightower/envconfig",
)

go_repository(
    name = "org_uber_go_atomic",
    commit = "1ea20fb1cbb1cc08cbd0d913a96dead89aa18289",
    importpath = "go.uber.org/atomic",
)

go_repository(
    name = "org_uber_go_multierr",
    commit = "3c4937480c32f4c13a875a1829af76c98ca3d40a",
    importpath = "go.uber.org/multierr",
)

go_repository(
    name = "org_uber_go_zap",
    commit = "ff33455a0e382e8a81d14dd7c922020b6b5e7982",
    importpath = "go.uber.org/zap",
)

go_repository(
    name = "com_github_go_redis_redis",
    commit = "b3d9bf10f6666b2ee5100a6f3f84f4caf3b4e37d",
    importpath = "github.com/go-redis/redis",
)

go_repository(
    name = "com_github_google_go_github",
    commit = "35781f7f4db7b3d7fc3359527472838da65023c6",
    importpath = "github.com/google/go-github",
)

go_repository(
    name = "com_github_google_go_querystring",
    commit = "44c6ddd0a2342c386950e880b658017258da92fc",
    importpath = "github.com/google/go-querystring",
)

go_repository(
    name = "org_golang_google_appengine",
    commit = "4a4468ece617fc8205e99368fa2200e9d1fad421",
    importpath = "google.golang.org/appengine",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "d668ce993890a79bda886613ee587a69dd5da7a6",
    importpath = "golang.org/x/oauth2",
)
