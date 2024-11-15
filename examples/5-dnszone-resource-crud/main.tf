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

resource "ngenix_dnszone" "ngenixterraformex4" {
  name = "ngenixterraformex12.ru"
  dns_records = [
    {
      name = "canonical-name-internal-server"
      type = "CNAME"
      data = "terraform-internal-ex12.express42.com."
    }
  ]
  //comment = "Terraform change"
}

# resource "ngenix_dnszone" "ngenixterraformex2" {
#   name = "ngenixterraformex1.ru"
#   dns_records = [
#     {
#       name = "config-ref-88903"
#       type = "A"
#       config_ref = {
#         id = 88903
#       }
#     }
#   ]
# }

# output "ngenixterraform_dnszone" {
#   value = ngenix_dnszone.ngenixterraformex4
# }
