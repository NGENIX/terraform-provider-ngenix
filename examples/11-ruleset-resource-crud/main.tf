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

// Example 1
# resource "ngenix_ruleset" "ngenixrulesetex1" {
#   name    = "ngenixrulesetexample"
#   enabled = true
#   rules = [
#     {
#       name    = "rule_one"
#       enabled = true
#       conditions = [
#         {
#           function = "in"
#           negation = true
#           params = [
#             {
#               variable = "http.request.host" // ip.geoip.country
#             },
#             {
#               trafficpattern_ref = {
#                 id = 1003852
#               }
#             }
#           ]
#         }
#       ],
#       actions = [
#         {
#           action = "allow"
#           params = ["", ""]
#         }
#       ]
#     }
#   ]
# }

// Example 2
# resource "ngenix_ruleset" "ngenixrulesetterrulestest" {
#   enabled = true
#   name    = "TERRRULESTEST"
#   rules = [
#     {
#       actions = [
#         {
#           action = "deny"
#           params = ["403", "access_denied"]
#         }
#       ]
#       conditions = [
#         {
#           function = "in"
#           negation = true
#           params = [
#             {
#               variable = "http.request.method"
#             },
#             {
#               trafficpattern_ref = {
#                 id = 1003851
#               }
#             }
#           ]
#         }
#       ]
#       enabled = true
#       name    = "HTTP_MET"
#     },
#     {
#       actions = [
#         {
#           action = "deny"
#           params = ["403", "access_denied"]
#         }
#       ]
#       conditions = [
#         {
#           function = "in"
#           negation = true
#           params = [
#             {
#               variable = "http.request.method"
#             },
#             {
#               trafficpattern_ref = {
#                 id = 1003851
#               }
#             }
#           ]
#         },
#         {
#           function = "in"
#           negation = false
#           params = [
#             {
#               variable = "ip.asn"
#             },
#             {
#               trafficpattern_ref = {
#                 id = 1003854
#               }
#             }
#           ]
#         }
#       ]
#       enabled = true
#       name    = "ASN"
#     },
#     {
#       actions = [
#         {
#           action = "deny"
#           params = ["403", "access_denied"]
#         },
#       ]
#       conditions = [
#         {
#           function = "in"
#           negation = false
#           params = [
#             {
#               variable = "ip.addr"
#             },
#             {
#               trafficpattern_ref = {
#                 id = 1003850
#               }
#             },
#           ]
#         },
#       ]
#       enabled = true
#       name    = "Addr_no"
#     }
#   ]
# }

// Example 3

resource "ngenix_ruleset" "ngenixrulesetterrulestest" {
  enabled = true
  name    = "TERRRULESTESTSIX"
  rules = [
    {
      name    = "HTTP_MET"
      enabled = true
      actions = [
        {
          action = "deny"
          params = ["403", "access_denied"]
        }
      ]
      conditions = [
        {
          function = "in"
          negation = true
          params = [
            {
              variable = "http.request.method"
            },
            {
              trafficpattern_ref = {
                id = 1003851
              }
            }
          ]
        }
      ]
    },
    {
      name    = "ASN"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params = [
            {
              variable = "http.request.method"
            },
            {
              trafficpattern_ref = {
                id = 1003851
              }
            }
          ]
        },
        {
          function = "in"
          negation = false
          params = [
            {
              variable = "ip.asn"
            },
            {
              trafficpattern_ref = {
                id = 1003854
              }
            }
          ]
        }
      ]
      actions = [
        {
          action = "deny"
          params = ["403", "no_access"]
        }
      ]
    },
    {
      name    = "Addr_no"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = false
          params = [
            {
              variable = "ip.addr"
            },
            {
              trafficpattern_ref = {
                id = 1003850
              }
            }
          ]
        },
      ]
      actions = [
        {
          action = "deny"
          params = ["403", "access_denied"]
        }
      ]
    },
    {
      name    = "strings-no-no"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = false
          params = [
            {
              variable = "http.request.path"
            },
            {
              trafficpattern_ref = {
                id = 1004269
              }
            }
          ]
        }
      ]
      actions = [
        {
          action = "jsChallenge"
          params = ["3600"]
        }
      ]
    },
    {
      name    = "HASH_ddd"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params = [
            {
              variable = "tls.ngx_fingerprint_hash"
            },
            {
              trafficpattern_ref = {
                id = 1003855
              }
            }
          ]
        }
      ]
      actions = [
        {
          action = "allow"
          params = ["", ""]
        },
      ]
    },
    {
      name    = "CountryNO"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = false
          params = [
            {
              variable = "ip.geoip.country"
            },
            {
              trafficpattern_ref = {
                id = 1003853
              }
            }
          ]
        }
      ]
      actions = [
        {
          action = "deny"
          params = ["403", "nonono"]
        }
      ]
    },
    {
      name    = "strings-no"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = false
          params = [
            {
              variable = "http.request.path"
            },
            {
              trafficpattern_ref = {
                id = 1003852
              }
            }
          ]
        }
      ]
      actions = [
        {
          action = "jsChallenge"
          params = ["3600"]
        }
      ]
    }
  ]
}

# output "ngenixrulesetex1_ruleset" {
#   value = ngenix_ruleset.ngenixrulesetex1
# }
