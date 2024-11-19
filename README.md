# Terraform Provider Ngenix

This is a Terraform provider for Ngenix Planform based on a [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework). It uses [Ngenix plaform API](https://developer.ngenix.net/platformApi) for managing objects.

This provider supports only these Ngenix platform objects for now:

1. [DNS](https://docs.ngenix.net/dns) zones with record sets.
2. [Traffic pattern](https://docs.ngenix.net/upravlenie-pravilami-obrabotki-zaprosov/kak-sozdat-spisok-znachenii).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21
- [Ngenix API core library](https://github.com/NGENIX/terraform-ngenix-api-core)

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
git clone git@github.com:NGENIX/terraform-provider-ngenix.git
cd terraform-provider-ngenix
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

At the moment, the provider allows you to manage three objects of Ngenix platform: 
1. DNS zones
2. Traffic patterxns

Each object supports basic CRUD operations (create/read/update/delete) and import command also. Two Terraform entities are implemented for each object - `data source` & `resource`.

You can find documentation for these in the `docs` directory. It could be generated by this command - `go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name ngenix`

All examples are described in the `examples` directory.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

Please don't forget to add new functionality to Ngenix API core library before adding it support to Terraform provider.

### How to build the Provider

1. Clone NGENIX/terraform-ngenix-provider - `git clone git@github.com:NGENIX/terraform-provider-ngenix.git`
2. Clone NGENIX/terraform-ngenix-api-core repository as a module for Terraform provider - `git clone git@github.com:NGENIX/terraform-ngenix-api-core.git terraform-provider-ngenix/ngenix-restapi`
3. Go to terraform-ngenix-provider folder - `cd terraform-ngenix-provider`
4. Download dependencies - `go mod tidy` and `go mod download`
5. Build and install provider - `go build .` and `go install .`

This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate ./...` or a run `go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name ngenix`

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Prepare Terraform for local development

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

#### Instruction

1. Find the `GOBIN` path where Go installs your binaries

```
$ go env GOBIN
/Users/<Username>/go/bin
```

2. Create a new file called `.terraformrc` in your home directory (~), then add the `dev_overrides` block below. Change the `<PATH>` to the value returned from the `go env GOBIN` command above.

```
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/hashicups" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

3. Use the `go install` command from the example repository's root directory to compile the provider into a binary and install it in your GOBIN path.

```
$ go install .
```

4. Go to `examples/1-provider-install-verification` folder and run a `terraform plan`.

Running a Terraform plan will report the provider override, as well as an error about the missing data source. Even though there was an error, this verifies that Terraform was able to successfully start the locally installed provider and interact with it in your development environment.

## Acceptance Testing 

All tests are located in the `internal/provider` directory with code and have prefix `_test.go`

Before running tests, setup a provider variables - `host`, `username` and `password` in a `provider_test.go` file.

```
	providerConfig = `
provider "ngenix" {
  host     = "https://api.ngenix.net/api/v3/"
  username = "EMAIL/token"
  password = "TOKEN"
}
`
)
```

For running all tests, execute the command `go test -v -cover ./internal/provider/`

To run a single test, you need to pass the test name in the launch command `TF_ACC=1 go test -count=1 -run='TestDnzZoneResource' -v`

## Publishing the Provider

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing) so that others can use it.
