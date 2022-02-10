# FireHydrant Terraform Provider

![Banner](images/terraform-firehydrant.png)

![](https://github.com/firehydrant/terraform-provider-firehydrant/actions/workflows/ci.yml/badge.svg)

Welcome to the FireHydrant Terraform provider! With this provider you can create and manage resources on your [FireHydrant](https://www.firehydrant.io) organization such as incident runbooks, services, teams, and more!

To view the full documentation of this provider, we recommend reading the documentation on the [Terraform Registry](https://registry.terraform.io/providers/firehydrant/firehydrant/latest)

## Developing the provider

There are a few conveniences for developing the Firehydrant provider to ease the complexity of building the Terraform
providers and running examples.

### Setup

First, setup a local dev_terraformrc file by copying `local/examples/dev_terraformrc.example` to a local file

```shell
cp local/examples/dev_terraformrc.example local/examples/dev_terraformrc
```

Next, modify the `devrc` file to point to the FULL PATH to the `local/bin` subdirectory of the provider on your machine.

Then use the `local` target in the `Makefile` to build the provider binary and copy it to the `local/bin` directory.

```shell
make local
```

Finally, set the environment variable, either in your shell or on the CLI to tell Terraform to use your `dev_terraformrc`


```shell
export TF_CLI_CONFIG_FILE=/home/developer/terraform-provider-firehydrant/local/examples/dev_terraformrc
```

Now you can skip `terraform init` because the override will use the file in devrc instead of the normal provider. From
here the workflow for testing examples becomes creating and applying a plan using the new provider and the normal
Terraform workflow.
