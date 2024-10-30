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

// Addr
# resource "ngenix_traffic_pattern" "ngenixtpex1" {
#   name = "test-addr-ex3"
#   type = "commonlist"
#   content_type = "addr"
#   patterns = [
#     {
#       addr = "98.164.15.2/32" // 23.76.67.0/24
#       expires = 1724166191
#       //ttl = 18000
#       //comment = "Created by Terraform"
#     },
#     {
#       addr = "23.56.67.89/32" // 23.56.67.0/24
#       expires = 1724165191
#       //ttl = 4800
#       //comment = "Created by Terraform"
#     },
#     {
#       addr = "12.12.34.45/32"
#       expires = 1724165191
#       //ttl = 4800
#       //comment = "Created by Terraform"
#     }
#   ]
# }

// httpMethod
# resource "ngenix_traffic_pattern" "ngenixtpex2" {
#   name = "test-http-ex1"
#   type = "commonlist"
#   content_type = "httpMethod"
#   patterns = [
#     {
#       http_method = "OPTIONS"
#     },
#     {
#       http_method = "POST"
#     },
#     {
#       http_method = "DELETE"
#     }
#   ]
# }

// countryCode
# resource "ngenix_traffic_pattern" "ngenixtpex2" {
#   name = "test-http-ex2"
#   type = "commonlist"
#   content_type = "countryCode"
#   patterns = [
#     {
#       country_code = "RU"
#     },
#     {
#       country_code = "US"
#     },
#     {
#       country_code = "CH"
#     }
#   ]
# }

// commonString
# resource "ngenix_traffic_pattern" "ngenixtpex2" {
#   name = "test-string-ex"
#   type = "commonlist"
#   content_type = "commonString"
#   patterns = [
#     {
#       common_string = "Common string 1"
#     },
#     {
#       common_string = "Common string 2"
#     },
#     {
#       common_string = "Common string 3"
#     }
#   ]
# }

// asn
# resource "ngenix_traffic_pattern" "ngenixtpex2" {
#   name = "test-asn-ex"
#   type = "commonlist"
#   content_type = "asn"
#   patterns = [
#     {
#       asn = 3456
#       comment = "ASN3456"
#     },
#     {
#       asn = 429496
#       comment = "ASN429496"
#     },
#     {
#       asn = 9065
#       comment = "ASN9065"
#     }
#   ]
# }

// md5HashString
# resource "ngenix_traffic_pattern" "ngenixtpex2" {
#   name = "test-md5hashstring-ex"
#   type = "commonlist"
#   content_type = "md5HashString"
#   patterns = [
#     {
#       md5hash_string = "3eac2d002ae18a6e2f565d63cf955725"
#       comment = "md5 hash - WelcomeToTheJungle"
#     },
#     {
#       md5hash_string = "8b1a9953c4611296a827abf8c47804d7"
#       comment = "md5 hash - Hello"
#     },
#     {
#       md5hash_string = "f5a7924e621e84c9280a9a27e1bcb7f6"
#       comment = "md5 hash - World"
#     }
#   ]
# }

# output "ngenixtpex1_traffic_pattern" {
#   value = ngenix_traffic_pattern.ngenixtpex1
# }

# output "ngenixtpex2_traffic_pattern" {
#   value = ngenix_traffic_pattern.ngenixtpex2
# }