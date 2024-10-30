// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	email = os.Getenv("EMAIL")
	token = os.Getenv("TOKEN")
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Ngenix client is properly configured.
	// It is also possible to use the NGENIX_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfigString = `
provider "ngenix" {
  host     = "https://api.ngenix.net/api/v3/"
  username = "%s"
  password = "%s"
}
`
)

var providerConfig = fmt.Sprintf(providerConfigString, email, token)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ngenix": providerserver.NewProtocol6WithError(New("test")()),
}
