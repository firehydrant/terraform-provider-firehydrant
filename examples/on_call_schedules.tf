data "firehydrant_on_call_schedule" "schedule" {
  id      = "my-schedule-id"
  team_id = "my-team-id"
}

data "firehydrant_on_call_schedules" "all_schedules" {
  team_id = "my-team-id"
}

data "firehydrant_on_call_schedules" "primary_schedules" {
  team_id = "my-team-id"
  query   = "primary"
}
