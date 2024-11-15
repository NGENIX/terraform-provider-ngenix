package provider

import (
	"context"
	"fmt"

	"ngenix/restapi"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsZoneDataSource{}
)

func DnsZoneDataSource() datasource.DataSource {
	return &dnsZoneDataSource{}
}

type dnsZoneDataSource struct {
	client *restapi.Client
}

// Data model.
// dnsZoneDataSourceModel maps the data source schema data.
type dnsZoneDataSourceModel struct {
	DnsZones []dnsZoneModel `tfsdk:"dns_zones"`
}

// dnsZoneModel maps schema data.
type dnsZoneModel struct {
	Id      types.Int64       `tfsdk:"id"`
	Name    types.String      `tfsdk:"name"`
	Records []dnsRecordsModel `tfsdk:"dns_records"`
	Comment types.String      `tfsdk:"comment"`
}

type configRefModel struct {
	ID types.Int64 `tfsdk:"id"`
}

type targetGroupRefModel struct {
	ID types.Int64 `tfsdk:"id"`
}

type dnsRecordsModel struct {
	Name           types.String         `tfsdk:"name"`
	Type           types.String         `tfsdk:"type"`
	Data           types.String         `tfsdk:"data"`
	ConfigRef      *configRefModel      `tfsdk:"config_ref"`
	TargetGroupRef *targetGroupRefModel `tfsdk:"targetgroup_ref"`
}

// Configure adds the provider configured client to the data source.
func (d *dnsZoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*restapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *restapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *dnsZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zones"
}

func (d *dnsZoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dns_zones": schema.ListNestedAttribute{
				Description: "List of DNS zones records.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "DNS zone ID.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "DNS zone name.",
						},
						"dns_records": schema.ListNestedAttribute{
							Description: "List of DNS zone record sets.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:    true,
										Description: "DNS recordset name.",
									},
									"type": schema.StringAttribute{
										Required:    true,
										Description: "DNS recordset type.",
									},
									"data": schema.StringAttribute{
										Required:    true,
										Description: "DNS recordset data.",
									},
									"config_ref": schema.ObjectAttribute{
										AttributeTypes: map[string]attr.Type{
											"id": types.Int64Type,
										},
										Optional:    true,
										Description: "DNS recordset config reference id.",
									},
									"targetgroup_ref": schema.ObjectAttribute{
										AttributeTypes: map[string]attr.Type{
											"id": types.Int64Type,
										},
										Optional:    true,
										Description: "DNS recordset targetgroup reference id.",
									},
								},
							},
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "A comment string.",
						},
					},
				},
			},
		},
	}
}

func (d *dnsZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsZoneDataSourceModel

	// Getting all DNS zones for provided username.
	dnsZonesList := d.client.GetAllDnsZonesList()
	// Map response body to model.
	for _, dnszone := range dnsZonesList {
		dnsZoneState := dnsZoneModel{
			Id:   types.Int64Value(int64(dnszone.ID)),
			Name: types.StringValue(dnszone.Name),
		}
		// DNS record set.
		recordConfigRef := &configRefModel{}
		recordTargetGroupRef := &targetGroupRefModel{}
		for _, record := range dnszone.Records {
			if record.ConfigRef != nil && record.TargetGroupRef != nil {
				recordConfigRef = &configRefModel{
					ID: types.Int64Value(record.ConfigRef.ID),
				}
				recordTargetGroupRef = &targetGroupRefModel{
					ID: types.Int64Value(record.TargetGroupRef.ID),
				}
			}
			dnsZoneState.Records = append(dnsZoneState.Records, dnsRecordsModel{
				Name:           types.StringValue(record.Name),
				Type:           types.StringValue(record.Type),
				Data:           types.StringValue(record.Data),
				ConfigRef:      recordConfigRef,
				TargetGroupRef: recordTargetGroupRef,
			})
		}
		dnsZoneState.Comment = types.StringValue(dnszone.Comment)

		state.DnsZones = append(state.DnsZones, dnsZoneState)
	}

	// Set state.
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
