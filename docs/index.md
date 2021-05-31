---
page_title: "firehydrant Provider"
subcategory: ""
description: |-

---

# FireHydrant Provider

Welcome to the FireHydrant Terraform provider! With this provider you can create and manage resources on your [FireHydrant](https://www.firehydrant.io) organization such as incident runbooks, services, teams, and more!



## Schema

### Required

- **api_key** (String, Optional) This is your API key (typically a bot token in FireHydrant) that is used to manage resources in FireHydrant. If set, the environment variable `FIREHYDRANT_API_KEY` will be used.

### Optional
- **firehydrant_base_url** (String, Optional)
