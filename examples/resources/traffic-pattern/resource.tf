# Manage Traffic pattern example
resource "ngenix_traffic_pattern" "testaddrex" {
  name         = "test-addr-ex"
  type         = "commonlist"
  content_type = "addr"
  patterns = [
    {
      addr    = "98.164.15.2"
      expires = 1924166191
    },
    {
      addr    = "23.56.67.89"
      expires = 1924165191
    }
  ]
}

resource "ngenix_traffic_pattern" "testhttpex" {
  name         = "test-http-ex"
  type         = "commonlist"
  content_type = "httpMethod"
  patterns = [
    {
      http_method = "OPTIONS"
    },
    {
      http_method = "POST"
    },
    {
      http_method = "DELETE"
    }
  ]
}