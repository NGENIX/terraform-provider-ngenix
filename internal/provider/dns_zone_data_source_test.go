package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDnsZoneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing.
			{
				Config: providerConfig + `data "ngenix_dns_zones" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of dns zones returned.
					resource.TestCheckResourceAttr("data.ngenix_dns_zones.test", "dns_zones.#", "17"),
				),
			},
		},
	})
}
