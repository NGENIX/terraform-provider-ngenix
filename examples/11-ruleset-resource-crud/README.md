# Ruleset examples

## Example 1 - Ruleset with ONE rule

```
resource "ngenix_ruleset" "ngenixrulesetexampleone" {
  name = "ngenixrulesetexampleone"
  enabled = true
  rules = [
    {
      name = "rule_one"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
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
      ],
      actions = [
        {
          action = "deny"
          params = [ "403", "wontbe" ]
        }
      ]
    }
  ]
}
```

## Example 2 - Ruleset with TWO rules (two conditions and one action)

```
resource "ngenix_ruleset" "ngenixrulesetterrulestest" {
  name    = "RULESETTWO"
  enabled = true
  rules   = [
    {
      name = "HTTP_MET"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
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
      actions = [
        {
          action = "deny"
          params = [ "403", "denied" ]
        }
      ]
    },
    {
      name = "ASN"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
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
          params   = [
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
          params = [ "403", "no_access" ]
        }
      ]
      
    }
  ]
}
```

## Example 3 - Ruleset with SIX rules as example 'TERRRULES'

```
resource "ngenix_ruleset" "ngenixrulesetterrulestest" {
  name    = "TERRRULESTEST"
  enabled = true
  rules   = [
    {
      name = "HTTP_MET"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
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
      actions = [
        {
          action = "deny"
          params = [ "403", "access_denied" ]
        }
      ]
    },
    {
      name = "ASN"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
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
          params   = [
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
          params = [ "403", "no_access" ]
        }
      ]
    },
    {
      name       = "Addr_no"
      enabled    = true
      conditions = [
        {
          function = "in"
          negation = false
          params   = [
            {
              variable           = "ip.addr"
            },
            {
              trafficpattern_ref = {
                id = 1003850
              }
            },
          ]
        },
      ]
      actions    = [
        {
          action = "deny"
          params = [ "403", "access_denied" ]
        },
      ]
    },
    {
      name       = "HASH_ddd"
      enabled    = true
      conditions = [
        {
          function = "in"
          negation = true
          params   = [
            {
              variable           = "tls.ngx_fingerprint_hash"
            },
            {
              trafficpattern_ref = {
                id = 1003855
              }
            },
          ]
        },
      ]
      actions    = [
        {
          action = "allow"
          params = [ "", "" ]
        },
      ]
    },
    {
      name       = "strings-no"
      enabled    = true
      conditions = [
        {
          function = "in"
          negation = false
          params   = [
            {
              variable           = "http.request.path"
            },
            {
              trafficpattern_ref = {
                id = 1003852
              }
            },
          ]
        },
      ]
      actions    = [
        {
          action = "jsChallenge"
          params = [ "3600" ]
        },
      ]
    },
    {
      name       = "CountryNO"
      enabled    = true
      conditions = [
        {
          function = "in"
          negation = false
          params   = [
            {
              variable           = "ip.geoip.country"
            },
            {
              trafficpattern_ref = {
                id = 1003853
              }
            },
          ]
        },
      ]
      actions    = [
        {
          action = "deny"
          params = [ "403", "nonono" ]
        },
      ]
    },
  ]
}
```