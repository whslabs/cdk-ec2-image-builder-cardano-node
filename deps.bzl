load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "com_github_aws_aws_cdk_go_awscdk_v2",
        importpath = "github.com/aws/aws-cdk-go/awscdk/v2",
        sum = "h1:14tUA2rKbUgXkn8+vfGhxUg5uGtdIXbjIgUuOTyaftw=",
        version = "v2.39.0",
    )
    go_repository(
        name = "com_github_aws_constructs_go_constructs_v10",
        importpath = "github.com/aws/constructs-go/constructs/v10",
        sum = "h1:2dkR80x3ptvKnYpQ7BkbfE6dwYAQ2o2g/79K8N8W8jA=",
        version = "v10.1.85",
    )
    go_repository(
        name = "com_github_aws_jsii_runtime_go",
        importpath = "github.com/aws/jsii-runtime-go",
        sum = "h1:A6o9DpZD0+IeFrXJ/qBPX7VJne5Vuk2KSfrG5Ez2dz8=",
        version = "v1.65.0",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_masterminds_semver_v3",
        importpath = "github.com/Masterminds/semver/v3",
        sum = "h1:hLg3sBzpNErnxhQtUy/mmLR2I9foDujNK030IGemrRc=",
        version = "v3.1.1",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:M2gUjqZET1qApGOWNSnZ49BAIMX4F/1plDv3+l31EJ4=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:pSgiaMZlXftHpm5L7V1+rVB+AZJydKsMxsQBIJw4PKk=",
        version = "v1.8.0",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
        version = "v0.0.0-20161208181325-20d25e280405",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
