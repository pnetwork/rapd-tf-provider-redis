# TF Provider Redis (Terraform Plugin Framework)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Quick start
1. please install first
```sh
make build VERSION=0.1.0
make apply-example EXAMPLE_NAME=resource TF_LOG_LEVEL=INFO
```
