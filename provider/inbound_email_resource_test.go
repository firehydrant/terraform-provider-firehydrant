package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccInboundEmailResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckInboundEmailResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInboundEmailResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInboundEmailResourceExists("firehydrant_inbound_email.test"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "name", fmt.Sprintf("test-inbound-email-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "slug", fmt.Sprintf("test-inbound-email-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "description", "Test inbound email description"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "status_cel", "email.body.contains('has recovered') ? 'CLOSED' : 'OPEN'"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "level_cel", "email.body.contains('panic') ? 'ERROR' : 'INFO'"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "allowed_senders.#", "1"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "allowed_senders.0", "@firehydrant.com"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "target.0.type", "Team"),
					resource.TestCheckResourceAttrSet("firehydrant_inbound_email.test", "target.0.id"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "rules.0", "email.body.contains(\"hello\")"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "rule_matching_strategy", "all"),
					testAccCheckInboundEmailResourceEmailAddressFormat("firehydrant_inbound_email.test"),
				),
			},
		},
	})
}

func TestAccInboundEmailResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckInboundEmailResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInboundEmailResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInboundEmailResourceExists("firehydrant_inbound_email.test"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "name", fmt.Sprintf("test-inbound-email-%s", rName)),
					testAccCheckInboundEmailResourceEmailAddressFormat("firehydrant_inbound_email.test"),
				),
			},
			{
				Config: testAccInboundEmailResourceConfig_update(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInboundEmailResourceExists("firehydrant_inbound_email.test"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "name", fmt.Sprintf("updated-inbound-email-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "description", "Updated test inbound email description"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "allowed_senders.#", "2"),
					resource.TestCheckResourceAttr("firehydrant_inbound_email.test", "rules.#", "2"),
				),
			},
		},
	})
}

func testAccCheckInboundEmailResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Inbound Email ID is set")
		}

		client, err := getTestClient()
		if err != nil {
			return fmt.Errorf("Error getting client: %s", err)
		}

		_, err = client.InboundEmails().Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error fetching inbound email with ID %s: %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckInboundEmailResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getTestClient()
		if err != nil {
			return fmt.Errorf("Error getting client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "firehydrant_inbound_email" {
				continue
			}

			_, err := client.InboundEmails().Get(context.Background(), rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Inbound Email still exists")
			}
		}

		return nil
	}
}

func testAccCheckInboundEmailResourceEmailAddressFormat(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		emailAddress := rs.Primary.Attributes["email"]
		if emailAddress == "" {
			return fmt.Errorf("No email address set")
		}

		// Check if the email address matches the expected format
		// This is a basic check; adjust the regex as needed to match the actual format
		match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@signals\.email$`, emailAddress)
		if !match {
			return fmt.Errorf("Email address %s does not match the expected format", emailAddress)
		}

		return nil
	}
}

func getTestClient() (firehydrant.Client, error) {
	apiKey := os.Getenv("FIREHYDRANT_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FIREHYDRANT_API_KEY must be set for acceptance tests")
	}

	client, err := firehydrant.NewRestClient(apiKey)
	if err != nil {
		return nil, fmt.Errorf("Error creating client: %s", err)
	}

	return client, nil
}

func testAccInboundEmailResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test" {
  name = "test-team-%s"
}

resource "firehydrant_inbound_email" "test" {
  name                   = "test-inbound-email-%s"
  slug                   = "test-inbound-email-%s"
  description            = "Test inbound email description"
  status_cel             = "email.body.contains('has recovered') ? 'CLOSED' : 'OPEN'"
  level_cel              = "email.body.contains('panic') ? 'ERROR' : 'INFO'"
  allowed_senders        = ["@firehydrant.com"]
  target {
    type = "Team"
    id   = firehydrant_team.test.id
  }
  rules                  = ["email.body.contains(\"hello\")"]
  rule_matching_strategy = "all"
}
`, rName, rName, rName)
}

func testAccInboundEmailResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test" {
  name = "test-team-%s"
}

resource "firehydrant_inbound_email" "test" {
  name                   = "updated-inbound-email-%s"
  slug                   = "test-inbound-email-%s"
  description            = "Updated test inbound email description"
  status_cel             = "email.body.contains('resolved') ? 'CLOSED' : 'OPEN'"
  level_cel              = "email.body.contains('critical') ? 'ERROR' : 'INFO'"
  allowed_senders        = ["@firehydrant.com", "@example.com"]
  target {
    type = "Team"
    id   = firehydrant_team.test.id
  }
  rules                  = ["email.body.contains(\"hello\")", "email.body.contains(\"urgent\")"]
  rule_matching_strategy = "any"
}
`, rName, rName, rName)
}
