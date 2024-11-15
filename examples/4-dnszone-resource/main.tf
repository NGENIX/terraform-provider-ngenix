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

# Example 1
# resource "ngenix_dnszone" "ngenix_tf_ex1" {
#   dns_zone = {
#     name = "ngenix-tf-ex1.ru",
#     customer_ref = {
#       id = 21046
#     }
#   }
# }

#output "ngenix_tf_ex1_dnszone" {
#  value = ngenix_dnszone.ngenix_tf_ex1
#}

# Example 2
# resource "ngenix_dnszone" "ngenix_tf_ex2" {
#   dns_zone = {
#     name = "ngenix-tf-ex2.ru",
#     customer_ref = {
#       id = 21046
#     },
#     dns_records = [
#       {
#         name = "internal-server",
#         type = "A",
#         data = "109.24.56.23"
#       },
#       {
#         name = "canonical-name-internal-server",
#         type = "CNAME",
#         data = "internal-server.express42.com."
#       }
#     ],
#     comment = "Created by Terraform"
#   }
# }

#output "ngenixterraform1_dnszone" {
#  value = ngenix_dnszone.ngenixterraform1
#}

resource "ngenix_dnszone" "ngenix_tf_ex3" {
  name = "ngenix-tf-ex3.ru"
  dns_records = [
    {
      name = "terraform-ex1",
      type = "A",
      data = "45.87.34.23"
    },
    {
      name = "cname-terraform-ex3",
      type = "CNAME",
      data = "ngenix-tf-ex3.express42.com."
    }
  ]
}

output "ngenix_tf_ex3_dnszone" {
  value = ngenix_dnszone.ngenix_tf_ex3
}
