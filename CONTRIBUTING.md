# Contributing

The HashiCorp documentation is a great place to start for learning how to contribute to a Terraform
provider. Some things you may find helpful include:

1. [Documentation writing and formatting guide](https://www.terraform.io/registry/providers/docs)
2. [Documentation preview tool](https://registry.terraform.io/tools/doc-preview)
3. [Provider development with SDKv2](https://www.terraform.io/plugin/sdkv2)
4. [Terraform provider Discuss forum](https://discuss.hashicorp.com/c/terraform-providers)

When contributing to the FireHydrant Terraform provider, adding or changing an 
attribute in a resource or data source generally involves:

1. Adding/changing it in the API client and tests located in the [/firehydrant directory](./firehydrant).
2. Adding/changing in the provider and tests located in the [/provider directory](./provider).
3. Adding/changing it in the documentation in the [/docs directory](./docs).
4. Adding an entry to the [changelog](./CHANGELOG.md).

## Adding a new resource

### API client

If support for the endpoints you need to interact with is not already part of the API client,
you must add it.

1. Add the necessary create/get/update/delete functions to a new client in the API client.
   All new clients added to the API client should get their own file. See [service dependencies](./firehydrant/service_dependencies.go)
   for an example of what this should look like. 
2. Add your new client to the API client by adding it to the Client interface and adding a function that returns an
   interface for interaction with your new client in [client.go](./firehydrant/client.go)

### Provider

1. Create a new resource file that follows the naming convention `<resource_name>_resource.go` in the [/provider directory](./provider)
   <resource_name> should be singular. 
2. Copy the [new resource template](./contributing/new_resource_templates.md#resource-template), paste it
   into your file and modify it to fit your use case. The template has comments to help explain the overall flow and what
   various parts do. For more complex data, you will find it helpful to look at other resources or other providers like 
   the TFE provider, the Google provider, or the AWS provider. 
3. If there are any format restrictions for any optional or required attributes, try to add validation functions. 
4. Create a new resource test file that follows the naming convention `<resource_name>_resource_test.go` in the [/provider directory](./provider)
   <resource_name> should be singular and this file should be named the same as your resource file but with `_test` at the end. 
5. Copy the [new resource test template](./contributing/new_resource_templates.md#resource-test-template), paste it
   into your file and modify it to fit your use case. The template has comments to help explain the overall flow and what
   various parts do. For more complex data, you will find it helpful to look at other resources or other providers like
   the TFE provider, the Google provider, or the AWS provider.
6. Create a new resource markdown file for documentation that follows the naming convention `<resource_name>.md` in the [/docs/resources directory](./docs/resources)
   <resource_name> should be singular.
7. Copy the existing documentation for another resource (task lists and services are good examples), paste it into your file 
   and modify it to fit your use case. 

## Adding a new data source

### API client

If support for the endpoints you need to interact with is not already part of the API client,
you must add it.

1. Add the necessary create/get/update/delete functions to a new client in the API client.
   All new clients added to the API client should get their own file. See [service dependencies](./firehydrant/service_dependencies.go)
   for an example of what this should look like. 
2. Add your new client to the API client by adding it to the Client interface and adding a function that returns an
   interface for interaction with your new client in [client.go](./firehydrant/client.go)

### Provider

1. Create a new data source file that follows the naming convention `<data_source_name>_data.go` in the [/provider directory](./provider)
   <data_source_name> should be singular unless you are creating a data source that retrieves multiple instances of something, like the
   services data source.
2. Copy the [new data source template](./contributing/new_data_source_templates.md#data-source-template), paste it
   into your file and modify it to fit your use case. The template has comments to help explain the overall flow and what
   various parts do. For more complex data, you will find it helpful to look at other data sources or other providers like
   the TFE provider, the Google provider, or the AWS provider.
3. Create a new data source test file that follows the naming convention `<data_source_name>_data_test.go` in the [/provider directory](./provider)
   <data_source_name> should be singular unless you are creating a data source that retrieves multiple instances of something, like the
   services data source. This test file should be named the same as your data source file but with `_test` at the end.
4. Copy the [new data source test template](./contributing/new_data_source_templates.md#data-source-test-template), paste it
   into your test file and modify it to fit your use case. For more complex data, you will find it helpful to look at other data sources or other providers like
   the TFE provider, the Google provider, or the AWS provider.
5. Create a new data source markdown file for documentation that follows the naming convention `<data_source_name>.md` in the [/docs/data-sources directory](./docs/data-sources)
   <data_source_name> should be singular unless you are creating a data source that retrieves multiple instances of something, like the
   services data source.
6. Copy the existing documentation for another data source (task lists and services are good examples), paste it into your file
   and modify it to fit your use case.

## Writing documentation

When writing documentation for resources and data sources, the patterns you should follow are:

1. Every new doc should include at least one example usage config that uses all the possible attributes.
2. Whenever you add a new attribute, you should make sure to add it to the documentation, including 
   the example usage configs.
3. Example usage configs should be run through the formatter by running `terraform fmt`.
4. The arguments reference section should have required attributes first, in alphabetical order, followed 
   by optional attributes in alphabetical order. 
5. Each attribute should have a description.
6. The attributes reference section should have `id` first, followed by the rest of the exported attributes,
   in alphabetical order. 

## Building the provider

Building the provider requires you to have Go 1.16 installed.

1. Clone this repository
2. Navigate into the repository's directory
3. Run `make build`. This will build the provider binary and store it in the
   repository's directory.

You can use this local build of the provider binary in many ways, which you can read more about in
the [provider installation documentation](https://www.terraform.io/cli/config/config-file#provider-installation).

For provider development, you will want to use [Terraform's development overrides](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers).
You can set the development overrides up in one of two ways:

### 1. (Preferred) Setting up development overrides in your Terraform CLI config file

1. Create a new directory called `providers` somewhere you can easily access (but not inside the same folder as your Terraform config).
2. Move the executable created to the `providers` directory you created earlier.
   ```
   mv terraform-provider-firehydrant /PATH/TO/YOUR/DIRECTORY/providers
   ```
3. Add the following to your ~/.terraformrc (or your preferred Terraform CLI config) file:
   ```
   provider_installation {

     # Use /PATH/TO/YOUR/DIRECTORY/providers/terraform-provider-firehydrant
     # as an overridden package directory for the firehydrant/firehydrant provider.
     # This disables the version and checksum verifications for this provider and
     # forces Terraform to look for the null provider plugin in the given directory.
     dev_overrides {
       "firehydrant/firehydrant" = "/PATH/TO/YOUR/DIRECTORY/providers"
     }

     # For all other providers, install them directly from their origin provider
     # registries as normal. If you omit this, Terraform will _only_ use
     # the dev_overrides block, and so no other providers will be available.
     direct {}
   }
   ```

### 2. Setting up development overrides in the provider repository directory

1. First, setup a local dev_terraformrc file by copying `local/examples/dev_terraformrc.example` to a local file
   ```shell
   cp local/examples/dev_terraformrc.example local/examples/dev_terraformrc
   ```
2. Next, modify the `dev_terraformrc` file to point to the FULL PATH to the `local/bin` subdirectory of the provider on your machine.
3. Then use the `local` target in the `Makefile` to build the provider binary and copy it to the `local/bin` directory.
   ```shell
   make local
   ```
4. Finally, set the environment variable, either in your shell or on the CLI to tell Terraform to use your `dev_terraformrc`
   ```shell
   export TF_CLI_CONFIG_FILE=/home/developer/terraform-provider-firehydrant/local/examples/dev_terraformrc
   ```

## Running tests

For information on running tests, see the [Tests guide](./TESTS.md).

## Updating the Changelog

In general, follow
[the Terraform provider best practices for versioning and changelog updates](https://www.terraform.io/plugin/sdkv2/best-practices/versioning).

Only update the `Unreleased` section, changing the unreleased version placeholder to an appropriate version if necessary,
using [Semantic Versioning](https://semver.org/) as a guideline.

Please use the template below when updating the changelog:
```
<change category>:

* **New Resource:** `name_of_new_resource` ([#10](link-to-PR))
* resource/resource_name: description of change or bug fix ([#10](link-to-PR))

<change category>:

* resource/resource_name: description of change or bug fix ([#10](link-to-PR))

...
```

### Change categories

- BREAKING CHANGES: This section documents in brief any incompatible changes and how to handle them.
- NOTES: Additional information for potentially unexpected upgrade behavior, upcoming deprecations, or to highlight very important crash fixes (e.g. due to upstream API changes)
- FEATURES: These are major new improvements that deserve a special highlight, such as a new resource or data source.
- ENHANCEMENTS: Smaller features added to the project such as a new attribute for a resource.
- BUG FIXES: Any bugs that were fixed.