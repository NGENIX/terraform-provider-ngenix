package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTrafficPatternAddrResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_traffic_pattern" "test" {
  name = "tst-addr-tp"
  type = "commonlist"
  content_type = "addr"
  patterns = [
    {
      addr = "98.164.15.2/32"
      expires = 1924166191
    },
    {
      addr = "23.56.67.89/32"
      expires = 1924165191
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name / type and content_type,
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "name", "tst-addr-tp"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "type", "commonlist"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "content_type", "addr"),
					// Verify number of items,
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.#", "2"),
					// Verify first traffic pattern item,
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.addr", "98.164.15.2/32"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.expires", "1924166191"),
					// Verify second traffic pattern item,
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.addr", "23.56.67.89/32"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.expires", "1924165191"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ngenix_traffic_pattern.test", "id"),
					resource.TestCheckResourceAttrSet("ngenix_traffic_pattern.test", "last_updated"),
				),
			},
			// ImportState testing,
			{
				ResourceName:      "ngenix_traffic_pattern.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups.
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_traffic_pattern" "test" {
  name = "tst-addr-tp"
  type = "commonlist"
  content_type = "addr"
  patterns = [
    {
      addr = "98.163.15.2/32"
      expires = 1924166191
    },
    {
      addr = "23.56.67.89/32"
      expires = 1924165191
    },
	{
      addr = "12.12.34.45/32"
      expires = 1924164191
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first traffic pattern item updated.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.addr", "98.163.15.2/32"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.expires", "1924166191"),
					// Verify second traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.addr", "23.56.67.89/32"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.expires", "1924165191"),
					// Verify third traffic pattern item updated.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.2.addr", "12.12.34.45/32"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.2.expires", "1924164191"),
				),
			},
			// Delete testing automatically occurs in TestCase.
		},
	})
}

func TestTrafficPatternHttpMethodResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_traffic_pattern" "test" {
  name = "tst-http-method-tp"
  type = "commonlist"
  content_type = "httpMethod"
  patterns = [
    {
      http_method = "OPTIONS"
    },
    {
      http_method = "GET"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name / type and content_type.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "name", "tst-http-method-tp"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "type", "commonlist"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "content_type", "httpMethod"),
					// Verify number of items.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.#", "2"),
					// Verify first traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.http_method", "OPTIONS"),
					// Verify second traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.http_method", "GET"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ngenix_traffic_pattern.test", "id"),
					resource.TestCheckResourceAttrSet("ngenix_traffic_pattern.test", "last_updated"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "ngenix_traffic_pattern.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the Ngenix API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_traffic_pattern" "test" {
  name = "tst-http-method-tp"
  type = "commonlist"
  content_type = "httpMethod"
  patterns = [
    {
      http_method = "OPTIONS"
    },
    {
      http_method = "GET"
    },
	{
      http_method = "POST"
    },
	{
      http_method = "DELETE"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.http_method", "OPTIONS"),
					// Verify second traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.http_method", "GET"),
					// Verify third traffic pattern item updated.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.2.http_method", "POST"),
					// Verify fourth traffic pattern item updated.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.3.http_method", "DELETE"),
				),
			},
			// Delete testing automatically occurs in TestCase.
		},
	})
}

func TestTrafficPatternAsnResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_traffic_pattern" "test" {
  name = "tst-asn-tp"
  type = "commonlist"
  content_type = "asn"
  patterns = [
    {
      asn = 3456
	  comment = "ASN3456"
    },
    {
      asn = 9065
	  comment = "ASN9065"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name / type and content_type.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "name", "tst-asn-tp"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "type", "commonlist"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "content_type", "asn"),
					// Verify number of items.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.#", "2"),
					// Verify first traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.asn", "3456"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.comment", "ASN3456"),
					// Verify second traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.asn", "9065"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.comment", "ASN9065"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ngenix_traffic_pattern.test", "id"),
					resource.TestCheckResourceAttrSet("ngenix_traffic_pattern.test", "last_updated"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "ngenix_traffic_pattern.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the Ngenix API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing.
			{
				Config: providerConfig + `
resource "ngenix_traffic_pattern" "test" {
  name = "tst-asn-tp"
  type = "commonlist"
  content_type = "asn"
  patterns = [
    {
      asn = 3456
	  comment = "ASN3456"
    },
    {
      asn = 9065
	  comment = "ASN9065"
    },
    {
      asn = 123
	  comment = "ASN123"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.asn", "3456"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.0.comment", "ASN3456"),
					// Verify second traffic pattern item.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.asn", "9065"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.1.comment", "ASN9065"),
					// Verify third traffic pattern item updated.
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.2.asn", "123"),
					resource.TestCheckResourceAttr("ngenix_traffic_pattern.test", "patterns.2.comment", "ASN123"),
				),
			},
			// Delete testing automatically occurs in TestCase.
		},
	})
}
