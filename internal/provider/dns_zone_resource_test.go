package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDnzZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_dnszone" "test" {
  name = "ngenixterraformacctest.ru"
  dns_records = [
    {
      name = "vm-a-record"
      type = "A"
      data = "23.12.76.128"
    },
    {
      name = "config-ref-88903"
      type = "A"
      config_ref = {
        id = 88903
      }
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name.
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "name", "ngenixterraformacctest.ru"),
					// Verify number of items.
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.#", "2"),
					// Verify first dns record item.
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.0.name", "vm-a-record"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.0.type", "A"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.0.data", "23.12.76.128"),
					// Verify second dns record item.
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.1.name", "config-ref-88903"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.1.type", "A"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.1.config_ref.id", "88903"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ngenix_dnszone.test", "id"),
					resource.TestCheckResourceAttrSet("ngenix_dnszone.test", "last_updated"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "ngenix_dnszone.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the Ngenix
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_dnszone" "test" {
  name = "ngenixterraformacctest.ru"
  dns_records = [
    {
      name = "vm-a-record"
      type = "A"
      data = "23.12.76.128"
    },
    {
      name = "config-ref-88903"
      type = "A"
      config_ref = {
        id = 88903
      }
    },
	{
      name = "test-cname-record",
      type = "CNAME",
      data = "terraform-internal.express42.com."
    },
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated.
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.0.name", "vm-a-record"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.0.type", "A"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.0.data", "23.12.76.128"),
					// Verify second dns record item.
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.1.name", "config-ref-88903"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.1.type", "A"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.1.config_ref.id", "88903"),
					// Verify third dns record item updated.
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.2.name", "test-cname-record"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.2.type", "CNAME"),
					resource.TestCheckResourceAttr("ngenix_dnszone.test", "dns_records.2.data", "terraform-internal.express42.com."),
				),
			},
			// Delete testing automatically occurs in TestCase.
		},
	})
}
