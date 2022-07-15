---
page_title: "Step Configuration - Confluence Cloud"
subcategory: "Runbooks"
---

# Confluence Cloud

~> **Note** You must have the Confluence Cloud integration installed in FireHydrant
for any Confluence Cloud runbook steps to work properly.

The FireHydrant Confluence Cloud integration allows FireHydrant users to export retrospective
reports to Confluence Cloud spaces after they are published.

### Available Steps

* [Export Retrospective](#export-retrospective)

## Export Retrospective

The [Confluence Cloud **Export Retrospective** step](https://support.firehydrant.com/hc/en-us/articles/4416252162196-Exporting-retrospectives-to-Confluence-)
allows FireHydrant users to export retrospective reports to Confluence Cloud spaces after they are published.

### Export Retrospective - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "confluence_cloud_export_retro" {
  integration_slug = "confluence_cloud"
  slug             = "export_retrospective"
}

resource "firehydrant_runbook" "confluence_cloud_runbook" {
  name        = "confluence-cloud-runbook"
  description = "This is an example configuration for a runbook that uses a Confluence Cloud step"

  steps {
    name      = "Export Retrospective to Confluence"
    action_id = data.firehydrant_runbook_action.confluence_cloud_export_retro.id

    config = jsonencode({
      title_template = "Retrospective: {{ retro.name }}"
      body_template  = "# Retrospective: {{ retro.name }}\n\n**Incident:** [{{ incident.severity }} - {{ incident.name}}]({{ incident.incident_url }}) ({{ retro.incident_active_duration }})\n\n{{ retro.summary }}\n\n{%- if incident.customer_impact_summary != blank or retro.impacts != empty -%}\n### Impact\n{% if incident.customer_impact_summary != blank -%}\n  {{ incident.customer_impact_summary }}\n{% endif %}\n\n\n{% if retro.impacts != empty %}\n| Type | Name | Condition |\n|------|------|-----------|\n{%- for impact in retro.impacts %}\n| {{ impact.type }} | {{ impact.name }} | {{ impact.condition }} |\n{%- endfor %}\n{% endif %}\n{% endif %}\n\n\n### Milestones\n\n| Milestone | Started | Duration |\n|-----------|---------|----------|\n{%- for milestone in retro.milestones %}\n| {{ milestone.type }} | {{ milestone.occurred_at }} | {{ milestone.duration }} |\n{%- endfor %}\n\n\n{% if retro.incident_roles != empty %}\n### Responders\n\n| Role | User |\n|------|------|\n{%- for role in retro.incident_roles %}\n| {{ role.name }} | {{ role.user }} |\n{%- endfor %}\n{% endif %}\n\n\n## Analysis\n\n{% if retro.contributing_factors != empty -%}\n### Contributing Factors\n{% for factor in retro.contributing_factors %}\n- {{ factor.summary }}\n{%- endfor %}\n{%- endif %}\n\n\n{% if retro.questions != empty -%}\n### Questions\n\n{% for question in retro.questions %}\n**{{ question.title }}**\n\n{{ question.body }}\n{% endfor %}\n{% endif %}\n\n\n---\n\n## Important events\n\n{% if retro.starred_events != empty %}\n{% for event in retro.starred_events %}\n_{{ event.occurred_at}}_ ({{event.created_by}})\n\n{{event.body}}\n{% endfor %}\n{% endif %}\n"

      parent_page = {
        label = "Dalmation Home"
        value = "xxxxxx"
      }
    })

    automatic = true
  }
}
```

### Export Retrospective - Steps Argument Reference

~> **Note** For this step to successfully execute, the step must run after the retrospective has started.
You can automatically execute the step when the current incident milestone is "Retrospective Started"
(aka "Postmortem Started") or "Retrospective Completed" (aka "Postmortem Completed") by setting the step `rule` attribute
to check for the correct milestones or by setting the runbook `attachment_rule` to check for the correct milestones.
You can also execute the step manually by setting the step `automatic` attribute to `false` or by setting the runbook
`attachment_rule` to execute the runbook manually.

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `body_template` - (Required) A body for the Confluence Cloud page that is created.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `title_template` - (Required) A name for the Confluence Cloud page that is created. 
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `parent_page` - (Optional) A parent page to nest the new Confluence Cloud page beneath.

The `parent_page` block supports:

* `label` - (Required) The title of the parent Confluence Cloud page.
* `value` - (Required) The ID of the parent Confluence Cloud page.
