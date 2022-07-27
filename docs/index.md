---
page_title: "Provider: FireHydrant"
---

# FireHydrant Provider

Welcome to the FireHydrant Terraform provider! With this provider you can create
and manage resources on your [FireHydrant](https://www.firehydrant.io) organization
such as incident runbooks, services, teams, and more!

Use the navigation on the left to read about the available resources and data sources.

~> **Note** This provider is still at version zero. To make sure that new versions with 
potentially breaking changes will not be automatically installed, you should 
[constrain the acceptable provider versions](https://www.terraform.io/language/providers/requirements#version-constraints) 
on the minor version.

## Example Usage

```hcl
terraform {
  required_providers {
    firehydrant = {
      source  = "firehydrant/firehydrant"
      version = "~> 0.3.0"
    }
  }
}

# Configure the FireHydrant Provider
provider "firehydrant" {
  api_key              = var.firehydrant_api_key
}

# Configure a FireHydrant priority
resource "firehydrant_priority" "example-priority" {
  slug        = "MYEXAMPLEPRIORITY"
  description = "This is an example priority"
}
```

## Argument Reference

The following arguments are supported:

* `api_key` - (Required) This is your API key that is used to manage resources in 
  FireHydrant. This value should be a bot token generated in FireHydrant.
  If set, the environment variable `FIREHYDRANT_API_KEY` will be used.
* `firehydrant_base_url` - (Optional) The FireHydrant API URL to connect to.
  Defaults to `https://api.firehydrant.io/v1/`. If set, the environment variable 
  `FIREHYDRANT_BASE_URL` will be used.
