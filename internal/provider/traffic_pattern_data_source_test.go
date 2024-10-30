package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTrafficPatternDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: providerConfig + `data "ngenix_traffic_patterns" "example" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of traffic patterns returned.
					resource.TestCheckResourceAttr("data.ngenix_traffic_patterns.example", "traffic_patterns.#", "20"),
				),
			},
		},
	})
}
