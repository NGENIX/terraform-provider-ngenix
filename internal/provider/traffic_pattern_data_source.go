package provider

import (
	"context"
	"fmt"

	"ngenix/restapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &trafficPatternDataSource{}
	_ datasource.DataSourceWithConfigure = &trafficPatternDataSource{}
)

// TrafficPatternDataSource is a helper function to simplify the provider implementation.
func TrafficPatternDataSource() datasource.DataSource {
	return &trafficPatternDataSource{}
}

// TrafficPatternDataSource is the data source implementation.
type trafficPatternDataSource struct {
	client *restapi.Client
}

// TrafficPatternSourceModel maps schema data.
type TrafficPatternSourceModel struct {
	TrafficPatterns []TrafficPatternModel `tfsdk:"traffic_patterns"`
}

type TrafficPatternModel struct {
	Name        types.String        `tfsdk:"name"`
	Type        types.String        `tfsdk:"type"`
	ContentType types.String        `tfsdk:"content_type"`
	Patterns    []PatternsDataModel `tfsdk:"patterns"`
}

// TrafficPatternsModel maps schema data.
type PatternsDataModel struct {
	Addr          types.String `tfsdk:"addr"`
	CommonString  types.String `tfsdk:"common_string"`
	CountryCode   types.String `tfsdk:"country_code"`
	HttpMethod    types.String `tfsdk:"http_method"`
	Asn           types.Int64  `tfsdk:"asn"`
	Md5HashString types.String `tfsdk:"md5hash_string"`
	Ttl           types.Int64  `tfsdk:"ttl"`
	Expires       types.Int64  `tfsdk:"expires"`
	Comment       types.String `tfsdk:"comment"`
}

// Metadata returns the data source type name.
func (d *trafficPatternDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_traffic_patterns"
}

// Configure adds the provider configured client to the data source.
func (d *trafficPatternDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*restapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ngenix.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Schema defines the schema for the data source.
func (d *trafficPatternDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Traffic Pattern.",
		Attributes: map[string]schema.Attribute{
			"traffic_patterns": schema.ListNestedAttribute{
				Description: "List of Traffic Patterns.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Traffic pattern name",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Traffic pattern type",
						},
						"content_type": schema.StringAttribute{
							Required:    true,
							Description: "Traffic pattern type",
						},
						"patterns": schema.ListNestedAttribute{
							Description: "A list of Traffic patterns.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"addr": schema.StringAttribute{
										Required:    true,
										Description: "Traffic pattern CIDR block - network CIDR block after validation",
									},
									"common_string": schema.StringAttribute{
										Required:    true,
										Description: "Traffic pattern common string",
									},
									"country_code": schema.StringAttribute{
										Required:    true,
										Description: "Traffic pattern country code",
									},
									"http_method": schema.StringAttribute{
										Required:    true,
										Description: "Traffic pattern HTTP method",
									},
									"asn": schema.Int64Attribute{
										Required:    true,
										Description: "Traffic pattern ASN",
									},
									"md5hash_string": schema.StringAttribute{
										Required:    true,
										Description: "Traffic pattern MD5 hash string",
									},
									"ttl": schema.Int64Attribute{
										Required:    true,
										Description: "Traffic pattern TTL",
									},
									"expires": schema.Int64Attribute{
										Required:    true,
										Description: "Traffic pattern expiration date",
									},
									"comment": schema.StringAttribute{
										Required:    true,
										Description: "Traffic pattern comment",
									},
								},
							},
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{},
	}
}

// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-data-source-read.
// Read refreshes the Terraform state with the latest data.
func (d *trafficPatternDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state.
	var state TrafficPatternSourceModel
	// Getting all DNS zones for provided username.
	trafficPatternsList := d.client.GetAllTrafficPatternsList()
	// Map response body to model.
	for _, trafficPattern := range trafficPatternsList {
		tpState := TrafficPatternModel{
			Name:        types.StringPointerValue(trafficPattern.Name),
			Type:        types.StringPointerValue(trafficPattern.Type),
			ContentType: types.StringPointerValue(trafficPattern.ContentType),
		}
		// Traffic patterns.
		for _, pattern := range trafficPattern.Patterns {
			tpState.Patterns = append(tpState.Patterns, PatternsDataModel{
				Addr:          types.StringPointerValue(pattern.Addr),
				CommonString:  types.StringPointerValue(pattern.CommonString),
				CountryCode:   types.StringPointerValue(pattern.CountryCode),
				HttpMethod:    types.StringPointerValue(pattern.HttpMethod),
				Asn:           types.Int64PointerValue(pattern.Asn),
				Md5HashString: types.StringPointerValue(pattern.Md5HashString),
				Ttl:           types.Int64PointerValue(pattern.Ttl),
				Expires:       types.Int64PointerValue(pattern.Expires),
				Comment:       types.StringValue(pattern.Comment),
			})
		}
		state.TrafficPatterns = append(state.TrafficPatterns, tpState)
	}

	// Set state.
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
