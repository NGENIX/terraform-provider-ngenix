// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"ngenix/restapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &ngenixProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ngenixProvider{
			version: version,
		}
	}
}

// ngenixProvider is the provider implementation.
type ngenixProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ngenixProviderModel maps provider schema data to a Go type.
type ngenixProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *ngenixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ngenix"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *ngenixProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "URI for Ngenix API. May also be provided via NGENIX_HOST environment variable.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username email for Ngenix API in format email/token. May also be provided via NGENIX_USERNAME environment variable.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "User token for Ngenix API. May also be provided via NGENIX_PASSWORD environment variable.",
			},
		},
	}
}

// Configure prepares a Ngenix API client for data sources and resources.
func (p *ngenixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Ngenix client")

	// Retrieve provider data from configuration
	var config ngenixProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the attributes, it must be a known value.
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Ngenix API Host",
			"The provider cannot create the Ngenix API client as there is an unknown configuration value for the Ngenix API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NGENIX_OST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Ngenix API Username",
			"The provider cannot create the Ngenix API client as there is an unknown configuration value for the Ngenix API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NGENIX_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Ngenix API Password",
			"The provider cannot create the Ngenix API client as there is an unknown configuration value for the Ngenix API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NGENIX_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("NGENIX_HOST")
	username := os.Getenv("NGENIX_USERNAME")
	password := os.Getenv("NGENIX_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Ngenix API Host",
			"The provider cannot create the Ngenix API client as there is a missing or empty value for the Ngenix API host. "+
				"Set the host value in the configuration or use the NGENIX_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Ngenix API Username",
			"The provider cannot create the Ngenix API client as there is a missing or empty value for the Ngenix API username. "+
				"Set the username value in the configuration or use the NGENIX_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Ngenix API Password",
			"The provider cannot create the Ngenix API client as there is a missing or empty value for the Ngenix API password. "+
				"Set the password value in the configuration or use the NGENIX_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ngenix_host", host)
	ctx = tflog.SetField(ctx, "ngenix_username", username)
	ctx = tflog.SetField(ctx, "ngenix_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ngenix_password")

	tflog.Debug(ctx, "Creating Ngenix client")

	// Create a new Ngenix client using the configuration values.
	client, err := restapi.NewClient(host, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Ngenix API Client",
			"An unexpected error occurred when creating the Ngenix API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Ngenix Client Error: "+err.Error(),
		)
		return
	}

	// Make the Ngenix client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Ngenix client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *ngenixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		DnsZoneDataSource,
		TrafficPatternDataSource,
		RulesetDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *ngenixProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		DnsZoneResource,
		TrafficPatternResource,
		RulesetResource,
	}
}
