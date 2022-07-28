---
page_title: "Conditional Logic"
subcategory: "Runbooks"
---

# Conditional Logic

The JSON encoded strings used in the runbook `attachment rule` and runbook steps `rule` attributes allows for
configuring when a runbook is attached to an incident or when a runbook step is executed based off of conditions
met within the incident. The values within the `attachment_rule` and steps `rule` attributes consist of a tuple of
keys for `logic` and `user_data`.

## Types of Conditions - Attachment Rule

We derive the types of conditions that can be set in the runbook `attachment_rule` attribute from our API's runbook
attachment attributes endpoint.

```
https://api.firehydrant.io/v1/fh-attributes/data_bags/system-runbook-attachment-attributes
```

In the example payload below you can see that we have a list of attributes that can be selected and we can see
the `opcode` or operators that can be used for these attributes. Given these operators, we can infer that we can
either have a condition where the Slack channel does or does not exist.

```json
{
  "name": "system-runbook-attachment-attributes",
  "attributes": [
    {
      "key": "incident_slack_channel",
      "type": "SlackChannel",
      "human_name": "Incident Slack channel"
    }
  ],
  "types": {
    "SlackChannel": {
      "operators": [
        {
          "human_name": "exists",
          "opcode": "exists"
        },
        {
          "human_name": "does not exist",
          "opcode": "does_not_exist"
        }
      ]
    }
  }
}
```

## Types of Conditions - Steps Rule

We derive the types of conditions that can be set in the runbook steps `rule` attribute from our API's step
execution attributes endpoint.

```
https://api.firehydrant.io/v1/fh-attributes/data_bags/step-execution-attributes
```

In the example payload below you can see that we have a list of attributes that can be selected and we can see
the `opcode` or operators that can be used for these attributes. Given these operators, we can infer that we can
either have a condition where the previous runbook step has completed, errored, or started.

```json
{
  "name": "step-execution-attributes",
  "attributes": [
    {
      "key": "previous_runbook_step",
      "type": "RunbookStep",
      "human_name": "Previous Runbook step"
    }
  ],
  "types": {
    "RunbookStep": {
      "operators": [
        {
          "human_name": "has completed",
          "opcode": "runbook_step_completed"
        },
        {
          "human_name": "has errored",
          "opcode": "runbook_step_errored"
        },
        {
          "human_name": "has started",
          "opcode": "runbook_step_started"
        }
      ],
      "values": {
        "async": "/fh-attributes/values/RunbookStep"
      }
    }
  }
}
```

## Value Based Conditions

Operators of conditions can be assigned values that come from data saved in different parts of our system.
For example, the below data we get back from the runbook attachment attributes API shows that we can select
multiple operators that take in an array of `IncidentRole` data values. Looking at the `IncidentRole` type,
we also see values that specify an `async` url that we can use to get the data necessary to fill in the `type`,
`value`, and `label` fields required for the `user_data` attribute.

```json
{
  "name": "system-runbook-attachment-attributes",
  "attributes": [
    {
      "key": "incident_assigned_roles",
      "type": "Array[IncidentRole]",
      "human_name": "Incident assigned roles"
    }
  ],
  "types": {
    "Array[IncidentRole]": {
      "operators": [
        {
          "human_name": "includes any of",
          "opcode": "includes_any",
          "arguments": ["Array[IncidentRole]"]
        },
        {
          "human_name": "includes all of",
          "opcode": "includes_all",
          "arguments": ["Array[IncidentRole]"]
        },
        {
          "human_name": "does not include",
          "opcode": "includes_none_of",
          "arguments": ["Array[IncidentRole]"]
        },
        {
          "human_name": "is empty",
          "opcode": "is_empty"
        }
      ],
      "hints": {
        "related_types": ["IncidentRole"],
        "collection": true
      }
    },
    "IncidentRole": {
      "operators": [
        {
          "human_name": "is",
          "opcode": "eq",
          "arguments": ["IncidentRole"]
        }
      ],
      "values": {
        "async": "/fh-attributes/values/IncidentRole"
      }
    }
  }
}
```

Using the above data, we can create something like the rule below which specifies that the runbook should be attached
when `incident_assigned_roles` has the role assigned of `Commander` or `Communication`.

```hcl
attachment_rule = jsonencode({
  logic = {
    includes_any = [
      {
        var = "incident_assigned_roles"
      },
      {
        var = "usr.1"
      }
    ]
  }
  user_data = {
    "1" = {
      type = "Array[IncidentRole]"
      value = [
        {
          type = "IncidentRole"
          value = data.firehydrant_incident_role.commander.id
          label = data.firehydrant_incident_role.commander.name
        },
        {
          type = "IncidentRole"
          value = data.firehydrant_incident_role.communications.id
          label = data.firehydrant_incident_role.communications.name
        }
      ]
      label = null
    }
  }
})
```

## Logic

The `logic` object consists of one or more sets conditions and varies based on the types of conditions added.

```hcl
attachment_rule = jsonencode({
  logic = {
    eq = [
      {
        var = "incident_current_milestone"
      },
      {
        var = "usr.1"
      }
    ]
  }
  user_data = {
    "1" = {
      type  = "Milestone"
      value = "resolved"
      label = "Resolved"
    }
  }
})
```

### Logic - Argument Reference

The `logic` block supports:

* `<OPERATOR>` - (Required) A block that represents the conditions to evaluate.
  Only one top level operator can be specified.
  Valid values for `<OPERATOR>` include:
    - `and` - Runs more than one logical operator check and requires all to return true. 
      Will take in an array of objects of any operator.
    - `or` - Runs more than one logical operator check and requires at least one of them to return true. 
      Will take in an array of objects of any operator.
    - `eq` - Requires a given argument to equal a given value.
    - `exists` - Requires a given argument to exist.
    - `does_not_exist` - Requires a given argument to not exist.
    - `is_one_of` - Requires a given argument to match at least one of a given set of values.
    - `includes_any` - Requires a given argument to include at least one of a given set of values.
    - `includes_all` - Requires a given argument to include all of a given set of values.
    - `includes_none_of` - Requires a given argument to include none of a given set of values.
    - `is_empty` - Requires a given argument to be empty.
    - `>` - Requires a given argument to be greater than a given value.
    - `<=` - Requires a given argument to be less than or equal to a given value.
    - `manually` - Requires a given argument to include at least one of a given set of values.

## Logic - Single Condition Example

The runbook will attach if a Slack incident channel exists.

```hcl
attachment_rule = jsonencode({
  logic = {
    exists = [
      {
        var = "incident_slack_channel"
      }
    ]
  }
  user_data = {}
})
```

## Logic - Multi Condition Or Example

The runbook step will execute if a Slack incident channel exists or if a Microsoft Teams channel exists.

```hcl
rule = jsonencode({
  logic = {
    or = [
      {
        exists = [
          {
            var = "incident_microsoft_teams_channel"
          }
        ]
      },
      {
        exists = [
          {
            var = "incident_slack_channel"
          }
        ]
      }
    ]
  }
  user_data = {}
})
```

## Logic - Multi Condition And Example

The runbook will attach if a Slack incident channel exists and a Microsoft Teams channel exists.

```hcl
attachment_rule = jsonencode({
  logic = {
    and = [
      {
        exists = [
          {
            var = "incident_microsoft_teams_channel"
          }
        ]
      },
      {
        exists = [
          {
            var = "incident_slack_channel"
          }
        ]
      }
    ]
  }
  user_data = {}
})
```

## User Data

The `user_data` attribute is used to map directly to variables in the `logic` attribute. In the example below,
you can see `var = "usr.1"` in the `logic` attribute, which maps to `user_data["1"]`. The key value in this case is "1"
but it can be anything that you would like it to be, it just needs to be a unique value within the `user_data` attribute that
maps to a `var` in the `logic` attribute.

```hcl
attachment_rule = jsonencode({
  logic = {
    eq = [
      {
        var = "incident_current_milestone"
      },
      {
        var = "usr.1"
      }
    ]
  }
  user_data = {
    "1" = {
      type  = "Milestone"
      value = "resolved"
      label = "Resolved"
    }
  }
})
```

### User Data - Argument Reference

The `user_data` block supports:

* `<KEY>` - (Required) A block that represents data associated with a `var` value in the `logic` attribute.
  `<KEY>` can be whatever you want, as long as it is a unique value within the `user_data` attribute. If 
  your `logic` attribute doesn't use `user_data`, you should set the value to `{}`.

The `<KEY>` block supports:

Unless `user_data` is set to `{}`, all three of these values _must_ be provided. 

* `type` - (Optional) The type of the data represented by this variable block.
* `label` - (Optional) The label of the data represented by this variable block.
* `value` - (Optional) The value of the data represented by this variable block.

## Arguments & Operators

Below is the current list of shared arguments and their operators that are available to be 
used in the runbook `attachment_rule` or the runbook steps `rule` attributes.

```
- incident_slack_channel
    - "exists"
    - "does_not_exist"
- incident_microsoft_teams_channel
    - "exists"
    - "does_not_exist"
- incident_current_milestone
    - "eq"
      - arg: Milestone
    - "is_one_of"
      - arg: Array[Milestone]
- incident_current_severity
    - "eq"
      - arg: Severity
    - "is_one_of"
      - arg: Array[Severity]
- incident_current_priority
    - "eq"
      - arg: Priority
    - "is_one_of"
      - arg: Array[Priority]
- incident_tags
    - "includes_any"
      - arg: Array[IncidentTag]
    - "includes_all"
      - arg: Array[IncidentTag]
    - "is_empty"
- incident_ticket
    - "exists"
    - "does_not_exist"
- incident_time_since_opened
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_last_note
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_started
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_detected
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_acknowledged
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_investigating
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_identified
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_mitigated
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_resolved
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_postmortem_started
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_time_since_milestone_postmortem_completed
    - ">"
      - arg: Duration
    - "<="
      - arg: Duration
- incident_assigned_roles
    - "includes_any"
      - arg: Array[IncidentRole]
    - "includes_all"
      - arg: Array[IncidentRole]
    - "includes_none_of"
      - arg: Array[IncidentRole]
    - "is_empty"
- incident_impacted_infrastructure
    - "includes_any"
      - arg: Array[Infrastructure]
    - "includes_all"
      - arg: Array[Infrastructure]
    - "includes_none_of"
      - arg: Array[Infrastructure]
    - "is_empty"
- incident_impacted_service_tiers
    - "includes_any"
      - arg: Array[ServiceTier]
    - "includes_all"
      - arg: Array[ServiceTier]
    - "includes_none_of"
      - arg: Array[ServiceTier]
    - "is_empty"
- incident_attached_runbooks
    - "includes_any"
      - arg: Array[Runbook]
    - "includes_all"
      - arg: Array[Runbook]
    - "includes_none_of"
      - arg: Array[Runbook]
    - "is_empty"
```

## Arguments & Operators - Attachment Rule

In addition to the shared list of arguments and their operators, below is the current list of
arguments and their operators that are only available to be used in the runbook `attachment_rule` 
attribute. 

```
- when_invoked
    - "manually"
```

## Arguments & Operators - Steps Rule

In addition to the shared list of arguments and their operators, below is the current list of 
arguments and their operators that are only available to be used in the runbook steps `rule` 
attribute.

```
- previous_runbook_step
    - "runbook_step_completed"
    - "runbook_step_errored"
    - "runbook_step_started"
```