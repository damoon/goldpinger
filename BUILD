# gazelle
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/damoon/goldpinger
gazelle(name = "gazelle")

# assets
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "assets",
    srcs = [":public"],
    mode = "0o644",
    visibility = ["//visibility:public"],
)

# kubernetes deployment
load("@k8s_deploy//:defaults.bzl", "k8s_deploy")

k8s_deploy(
  name = "dev",
  template = ":kubernetes.yaml",
  images = {
    "registry.31j.de/goldpinger/goldpinger": "//cmd/goldpinger:image"
  },
)
