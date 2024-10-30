package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRulesetOneRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_ruleset" "test" {
  name    = "TERRRULESACCTEST"
  enabled = true
  rules   = [
    {
      name = "HTTP_MET"
	  enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
            {
              variable = "http.request.method"
            },
            {
              trafficpattern_ref = {
                id = 1003851
              }
            }
          ]
        }
      ]
      actions = [
        {
          action = "deny"
          params = [ "403", "access_denied" ]
        }
      ]
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name / enabled.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "name", "TERRRULESACCTEST"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "enabled", "true"),
					// Verify number of rules.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.#", "1"),
					// Verify rules items.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.name", "HTTP_MET"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.enabled", "true"),
					// Verify rules conditions item.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.#", "1"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.function", "in"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.negation", "true"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.params.#", "2"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.params.0.variable", "http.request.method"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.params.1.trafficpattern_ref.id", "1003851"),
					// Verify rules action item
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.action", "deny"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.params.#", "2"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.params.0", "403"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.params.1", "access_denied"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ngenix_ruleset.test", "id"),
					resource.TestCheckResourceAttrSet("ngenix_ruleset.test", "last_updated"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "ngenix_ruleset.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the Ngenix
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_ruleset" "test" {
  name    = "TERRRULESACCTEST"
  enabled = true
  rules   = [
    {
      name = "HTTP_MET"
	  enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
            {
              variable = "http.request.method"
            },
            {
              trafficpattern_ref = {
                id = 1003851
              }
            }
          ]
        }
      ]
      actions = [
        {
          action = "deny"
          params = [ "403", "access_denied" ]
        }
      ]
    },
    {
      name = "ASN"
	  enabled = true
      conditions = [
        {
          function = "in"
          negation = false
          params   = [
            {
              variable = "ip.asn"
            },
            {
              trafficpattern_ref = {
                id = 1003854
              }
            }
          ]
        }
      ]
	  actions = [
        {
          action = "deny"
          params = [ "403", "access_denied" ]
        }
      ]
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name / enabled.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "name", "TERRRULESACCTEST"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "enabled", "true"),
					// Verify number of rules.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.#", "2"),
					// Verify rules items.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.name", "HTTP_MET"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.enabled", "true"),
					// Verify rules conditions item.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.#", "1"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.function", "in"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.negation", "true"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.params.#", "2"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.params.0.variable", "http.request.method"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.conditions.0.params.1.trafficpattern_ref.id", "1003851"),
					// Verify rules action item.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.action", "deny"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.params.#", "2"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.params.0", "403"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.0.actions.0.params.1", "access_denied"),
					// Verify rules items.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.name", "ASN"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.enabled", "true"),
					// Verify rules conditions item.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.conditions.#", "1"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.conditions.0.function", "in"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.conditions.0.negation", "false"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.conditions.0.params.#", "2"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.conditions.0.params.0.variable", "ip.asn"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.conditions.0.params.1.trafficpattern_ref.id", "1003854"),
					// Verify rules action item.
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.actions.#", "1"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.actions.0.action", "deny"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.actions.0.params.#", "2"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.actions.0.params.0", "403"),
					resource.TestCheckResourceAttr("ngenix_ruleset.test", "rules.1.actions.0.params.1", "access_denied"),
				),
			},
			// Delete testing automatically occurs in TestCase.
		},
	})
}
