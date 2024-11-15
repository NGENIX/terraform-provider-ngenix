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

data "ngenix_dns_zones" "example" {}

# locals {
#   dns_zones_records = flatten([
#     for zone in data.ngenix_dns_zones.example.dns_zones : [
#       for record in zone.dns_records : {
#         zone_name = zone.name
#         record_name = record.name
#         record_type = record.type
#       } 
#     ]
#   ])
# }

output "example_dnszone" {
  value = data.ngenix_dns_zones.example
}
# output "dnszone_id" {
#   value = [for zone in data.ngenix_dns_zones.example.dns_zones : zone.id]
# }
# output "dnszone_name" {
#   value = [for zone in data.ngenix_dns_zones.example.dns_zones : zone.name]
# }
# output "dnszone_comment" {
#   value = [for zone in data.ngenix_dns_zones.example.dns_zones : zone.comment]
# }
# output "dnszone_a_records" {
#   value = [for record in local.dns_zones_records : record.record_name if record.record_type == "A"]
# }

