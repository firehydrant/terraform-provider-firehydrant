resource "firehydrant_team" "firefighters" {
  name = "Firefighters"
}

data "firehydrant_team" "firefighters" {
  id = "857a83c4-17d1-4362-a4e8-42d1a3d19ed1"
}

data "firehydrant_teams" "all_teams" {
}

data "firehydrant_teams" "test_teams" {
  query = "Test"
}
