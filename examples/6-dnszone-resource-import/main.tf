terraform {
  required_providers {
    ngenix = {
      source = "ngenix.net/api/ngenix"
    }
  }
}

provider "ngenix" {
  host     = var.ngenix_api_url
  username = var.ngenix_email
  password = var.ngenix_token
}

resource "ngenix_dnszone" "terraformimporttestex1" {
  name = "terraformimporttestex1.ru"
  dns_records = [
    {
      name = "terraformimporttestex1"
      type = "CNAME"
      data = "tf-import-example1.ngenix.net."
    },
  ]
}

# output "ngenixterraform_dnszone" {
#   value = ngenix_dnszone.terraformimporttestex1
# }
