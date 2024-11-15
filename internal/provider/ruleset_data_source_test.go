package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRulesetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: providerConfig + `data "ngenix_rulesets" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of traffic patterns returned.
					resource.TestCheckResourceAttr("data.ngenix_rulesets.example", "rulesets.#", "2"),
				),
			},
		},
	})
}
