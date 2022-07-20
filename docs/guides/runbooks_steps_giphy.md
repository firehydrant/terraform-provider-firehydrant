---
page_title: "Step Configuration - Giphy"
subcategory: "Runbooks"
---

# Giphy

~> **Note** You must have the Giphy integration installed in FireHydrant
for any Giphy runbook steps to work properly.

The FireHydrant Giphy integration allows GIFs to automatically be added to lift your 
team’s spirits and celebrate as an incident progresses.

### Available Steps

* [Incident Channel GIF](#incident-channel-gif)

## Incident Channel GIF

The Giphy **Incident Channel GIF** step allows GIFs to automatically 
be added as an incident progresses.

### Incident Channel GIF - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "giphy_incident_channel_gif" {
  integration_slug = "giphy"
  slug             = "incident_channel_gif"
}

resource "firehydrant_runbook" "giphy_incident_channel_gif_runbook" {
  name = "giphy-incident-channel-gif-runbook"

  steps {
    name      = "Post A Gif from Giphy"
    action_id = data.firehydrant_runbook_action.giphy_incident_channel_gif.id

    config = jsonencode({
      phrases = "untitled goose game\nmagical"
    })
  }
}
```

### Incident Channel GIF - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `repeats` - (Optional) Whether this step should repeat. Defaults to `false`.
  When this value is `true`, `repeats_duration` _must_ be provided.
* `repeats_duration` - (Optional) How often this step should repeat in ISO8601.
  Example: PT10M [Format Spec](https://www.digi.com/resources/documentation/digidocs/90001437-13/reference/r_iso_8601_duration_format.htm)
  This value _must_ be provided if `repeats` is `true`. This value _must not_ be provided if `repeats` is `false`.

The `config` block supports:

* `keywords` - (Optional) A list of keywords, separated by newlines
* `phrases` - (Optional) A list of random phrases, separated by newlines
