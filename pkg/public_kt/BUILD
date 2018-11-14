package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_kotlin//kotlin:kotlin.bzl", "kt_jvm_library")

kt_jvm_library(
    name = "public_kt_lib",
    srcs = glob(["main.kt"]),
    deps = ["//pkg/common_kt:common_kt_lib"],
)

java_binary(
    name = "java_bin",
    main_class = "com.github.soushin.multirpojects.publickt.MainKt",
    runtime_deps = [":public_kt_lib"],
)

# --

load("@io_bazel_rules_docker//container:container.bzl", "container_image")

container_image(
    name = "public_kt_image",
    base = "@java_base_image//image",
    files = [":java_bin_deploy.jar"],
    cmd = [
        "java",
        "-jar",
        "java_bin_deploy.jar",
    ],
)
