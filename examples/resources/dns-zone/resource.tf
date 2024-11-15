# Manage DNS zone example
resource "ngenix_dnszone" "example" {
  dns_zone = {
    name = "example.ru",
    customer_ref = {
      id = 21046
    },
    dns_records = [
      {
        name = "example",
        type = "A",
        config_ref = {
          id = 52835
        }
      },
      {
        name = "example",
        type = "CNAME",
        data = "example.ngenix.net."
      },
      {
        name = "_internal._protocol.name",
        type = "SRV",
        data = "10 5 2000 test.srv.test."
      },
      {
        name = "@",
        type = "TXT",
        data = "ngenix-validation-id=dfd32323c3h266a49gba0erb7763d4cfb2"
      }
    ]
  }
}