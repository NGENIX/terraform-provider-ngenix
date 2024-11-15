package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"ngenix/restapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dnsZoneResource{}
	_ resource.ResourceWithConfigure   = &dnsZoneResource{}
	_ resource.ResourceWithImportState = &dnsZoneResource{}
)

// NewDnsZoneResource is a helper function to simplify the provider implementation.
func DnsZoneResource() resource.Resource {
	return &dnsZoneResource{}
}

// dnsZoneResource is the resource implementation.
type dnsZoneResource struct {
	client *restapi.Client
}

// Data model
// dnsZoneDataSourceModel maps the data source schema data.
type dnsZoneResourceModel struct {
	ID          types.String          `tfsdk:"id"`
	Name        types.String          `tfsdk:"name"`
	Records     []dnsRecordsItemModel `tfsdk:"dns_records"`
	Comment     types.String          `tfsdk:"comment"`
	LastUpdated types.String          `tfsdk:"last_updated"`
}

type configRefItemModel struct {
	ID types.Int64 `tfsdk:"id"`
}

type targetGroupRefItemModel struct {
	ID types.Int64 `tfsdk:"id"`
}

type dnsRecordsItemModel struct {
	Name           types.String             `tfsdk:"name"`
	Type           types.String             `tfsdk:"type"`
	Data           types.String             `tfsdk:"data"`
	ConfigRef      *configRefItemModel      `tfsdk:"config_ref"`
	TargetGroupRef *targetGroupRefItemModel `tfsdk:"targetgroup_ref"`
}

var DnsRecordTypes = []string{"A", "CNAME", "MX", "AAAA", "SRV", "NS", "TXT", "CAA"}

func isDnsTypeInRange(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// Configure adds the provider configured client to the data source.
func (d *dnsZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the resource type name.
func (r *dnsZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnszone"
}

// Schema defines the schema for the resource.
func (r *dnsZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DNS zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Placeholder identifier attribute.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last Terraform update of the DNS zone.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "DNS zone name",
			},
			"dns_records": schema.ListNestedAttribute{
				Optional:    true,
				Description: "DNS zone records",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "DNS record name",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "DNS record type",
						},
						"data": schema.StringAttribute{
							Optional:    true,
							Description: "DNS record data",
						},
						"config_ref": schema.ObjectAttribute{
							AttributeTypes: map[string]attr.Type{
								"id": types.Int64Type,
							},
							Optional:    true,
							Description: "DNS record config reference id",
						},
						"targetgroup_ref": schema.ObjectAttribute{
							AttributeTypes: map[string]attr.Type{
								"id": types.Int64Type,
							},
							Optional:    true,
							Description: "DNS record target group reference id",
						},
					},
				},
			},
			"comment": schema.StringAttribute{
				Description: "DNS zone resource comment",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Changed by Terraform"),
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
				},
			},
		},
	}
}

func (r *dnsZoneResource) DNSRecordsModelTransformation(records []dnsRecordsItemModel) ([]restapi.Records, error) {
	recordConfigRef := restapi.ConfigRef{}
	recordTargetGroupRef := restapi.TargetGroupRef{}
	dnsRecordSet := []restapi.Records{}
	for _, record := range records {
		// Requirement: Records.Type is in range onf values.
		if isDnsTypeInRange(record.Type.ValueString(), DnsRecordTypes) {
			// Requirement: DNS record Type = A.
			if record.Type.ValueString() == "A" {
				if record.ConfigRef != nil && len(record.Data.ValueString()) != 0 {
					return nil, errors.New("both fields <data> and <config_ref> are not supported")
				} else if record.TargetGroupRef != nil && len(record.Data.ValueString()) != 0 {
					return nil, errors.New("both fields <data> and <targetgroup_ref> are not supported")
				} else if record.ConfigRef == nil && record.TargetGroupRef == nil {
					dnsRecordSet = append(dnsRecordSet, restapi.Records{
						Name: record.Name.ValueString(),
						Type: record.Type.ValueString(),
						Data: record.Data.ValueString(),
					})
				} else if record.ConfigRef != nil {
					recordConfigRef = restapi.ConfigRef{
						ID: record.ConfigRef.ID.ValueInt64(),
					}
					dnsRecordSet = append(dnsRecordSet, restapi.Records{
						Name:      record.Name.ValueString(),
						Type:      record.Type.ValueString(),
						ConfigRef: &recordConfigRef,
					})
				} else if record.TargetGroupRef != nil {
					recordTargetGroupRef = restapi.TargetGroupRef{
						ID: record.TargetGroupRef.ID.ValueInt64(),
					}
					dnsRecordSet = append(dnsRecordSet, restapi.Records{
						Name:           record.Name.ValueString(),
						Type:           record.Type.ValueString(),
						TargetGroupRef: &recordTargetGroupRef,
					})
				} else {
					return nil, errors.New("misunderstanding with data, config_ref and targetgroup_ref fields")
				}
			} else {
				dnsRecordSet = append(dnsRecordSet, restapi.Records{
					Name: record.Name.ValueString(),
					Type: record.Type.ValueString(),
					Data: record.Data.ValueString(),
				})
			}
		} else {
			return nil, errors.New("dns record type value is not in range [A, CNAME, MX, AAAA, SRV, NS, TXT, CAA]")
		}
	}
	return dnsRecordSet, nil
}

func (r *dnsZoneResource) DNSRecordsItemModelTransformation(records []restapi.Records) ([]dnsRecordsItemModel, error) {
	recordConfigRefItem := configRefItemModel{}
	recordTargetGroupRefItem := targetGroupRefItemModel{}
	dnsRecordsItems := []dnsRecordsItemModel{}
	for _, record := range records {
		// Requirement: Records.Type is in range onf values.
		if isDnsTypeInRange(record.Type, DnsRecordTypes) {
			// Requirement: DNS record Type = A.
			if record.Type == "A" {
				if record.ConfigRef != nil && len(record.Data) != 0 {
					return nil, errors.New("both fields <data> and <config_ref> are not supported")
				} else if record.TargetGroupRef != nil && len(record.Data) != 0 {
					return nil, errors.New("both fields <data> and <targetgroup_ref> are not supported")
				} else if record.ConfigRef == nil && record.TargetGroupRef == nil {
					dnsRecordsItems = append(dnsRecordsItems, dnsRecordsItemModel{
						Name: types.StringValue(record.Name),
						Type: types.StringValue(record.Type),
						Data: types.StringValue(record.Data),
					})
				} else if record.ConfigRef != nil {
					recordConfigRefItem = configRefItemModel{
						ID: types.Int64Value(record.ConfigRef.ID),
					}
					dnsRecordsItems = append(dnsRecordsItems, dnsRecordsItemModel{
						Name:      types.StringValue(record.Name),
						Type:      types.StringValue(record.Type),
						ConfigRef: &recordConfigRefItem,
					})
				} else if record.TargetGroupRef != nil {
					recordTargetGroupRefItem = targetGroupRefItemModel{
						ID: types.Int64Value(record.TargetGroupRef.ID),
					}
					dnsRecordsItems = append(dnsRecordsItems, dnsRecordsItemModel{
						Name:           types.StringValue(record.Name),
						Type:           types.StringValue(record.Type),
						TargetGroupRef: &recordTargetGroupRefItem,
					})
				} else {
					return nil, errors.New("misunderstanding with data, config_ref and targetgroup_ref fields, you could use ONLY ONE field at a moment")
				}
			} else {
				dnsRecordsItems = append(dnsRecordsItems, dnsRecordsItemModel{
					Name: types.StringValue(record.Name),
					Type: types.StringValue(record.Type),
					Data: types.StringValue(record.Data),
				})
			}
		} else {
			return nil, errors.New("dns record type should be in range [A, CNAME, MX, AAAA, SRV, NS, TXT, CAA]")
		}
	}
	return dnsRecordsItems, nil
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *dnsZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan dnsZoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Check DNS zone existence.
	if r.client.DnsZoneExist(plan.Name.ValueString()) {
		resp.Diagnostics.AddError(
			"DNS zone has already exist!",
			"Could not create DNS zone - DNS zone has already exist!",
		)
		return
	}

	// Generate API request body from PLAN.
	dnsRecords, err := r.DNSRecordsModelTransformation(plan.Records)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while DNS zone Records creation and validation process",
			fmt.Sprintf("Could not create DNS zone Records, please check a page with errors, error: %s", err.Error()),
		)
		return
	}
	// Default comment value.
	comment := "Created by Terraform"
	if !plan.Comment.IsNull() {
		comment = plan.Comment.ValueString()
	}
	var dnsZone = restapi.DnsZone{
		Name: plan.Name.ValueString(),
		CustomerRef: &restapi.CustomerRef{
			ID: int64(r.client.CustomerId()),
		},
		Records: dnsRecords,
		Comment: comment,
	}

	// Create new DNS zone.
	createdDnszone, err := r.client.CreateDnsZone(dnsZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating DNS zone",
			fmt.Sprintf("Could not create DNS Zone, unexpected error: %s", err.Error()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values.
	dnsRecordsItems, err := r.DNSRecordsItemModelTransformation(createdDnszone.Records)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while DNS zone Records creation and validation process",
			fmt.Sprintf("Could not create DNS zone Records, error: %s", err.Error()),
		)
		return
	}
	// Default comment value
	createdZoneComment := "Created by Terraform"
	if createdDnszone.Comment != "" {
		createdZoneComment = createdDnszone.Comment
	}
	// Update state model from newly created DNS zone.
	zoneId := r.client.GetDnsZoneIDByName(plan.Name.ValueString())
	plan.ID = types.StringValue(strconv.Itoa(zoneId))
	plan.Name = types.StringValue(createdDnszone.Name)
	plan.Records = dnsRecordsItems
	plan.Comment = types.StringValue(createdZoneComment)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "DNS zone was created successfully!")

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *dnsZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state.
	var state dnsZoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed DNS zone value from Ngenix.
	zoneId := r.client.GetDnsZoneIDByName(state.Name.ValueString())
	fromZone, err := r.client.GetDnsZoneById(zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ngenix DNS zone",
			fmt.Sprintf("Could not read Ngenix DNS Zone by ID = %d, error: %s", zoneId, err.Error()),
		)
		return
	}

	// Overwrite items with refreshed state.
	dnsRecordsItems, err := r.DNSRecordsItemModelTransformation(fromZone.Records)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while DNS zone Records creation and validation process",
			fmt.Sprintf("Could not create DNS zone Records, error: %s", err.Error()),
		)
		return
	}
	state.Name = types.StringValue(fromZone.Name)
	state.Records = dnsRecordsItems
	state.Comment = types.StringValue(fromZone.Comment)

	tflog.Trace(ctx, "DNS zone was read successfully!")

	// Set refreshed state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan dnsZoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan.
	dnsRecords, er := r.DNSRecordsModelTransformation(plan.Records)
	if er != nil {
		resp.Diagnostics.AddError(
			"Error while DNS zone Records creation and validation process",
			fmt.Sprintf("Could not update DNS zone Records, please check a page with errors, error: %s", er.Error()),
		)
		return
	}

	var dnsZone = restapi.DnsZone{
		Records: dnsRecords,
		Comment: plan.Comment.ValueString(),
	}

	// Update existing DNS zone.
	zoneId := r.client.GetDnsZoneIDByName(plan.Name.ValueString())
	_, err := r.client.UpdateDnsZone(zoneId, dnsZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Ngenix DNS zones",
			fmt.Sprintf("Could not update DNS Zone (PATCH), unexpected error: %s", err.Error()),
		)
		return
	}

	// Fetch updated traffic pattern.
	updatedDnsZone, err := r.client.GetDnsZoneById(zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Ngenix DNS zone by ID",
			fmt.Sprintf("Could not update Ngenix DNS zone by ID %d, error: %s", zoneId, err.Error()),
		)
		return
	}

	// Update resource state with updated items and timestamp.
	dnsRecordsItems, err := r.DNSRecordsItemModelTransformation(updatedDnsZone.Records)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while DNS zone Records creation and validation process",
			fmt.Sprintf("Could not update DNS zone, error: %s", err.Error()),
		)
		return
	}
	updatedZoneId := r.client.GetDnsZoneIDByName(updatedDnsZone.Name)
	plan.ID = types.StringValue(strconv.Itoa(updatedZoneId))
	plan.Name = types.StringValue(updatedDnsZone.Name)
	plan.Records = dnsRecordsItems
	plan.Comment = types.StringValue(updatedDnsZone.Comment)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "DNS zone was updated successfully!")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state.
	var state dnsZoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing DNS zone
	zoneId := r.client.GetDnsZoneIDByName(state.Name.ValueString())
	err := r.client.DeleteDnsZone(zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Ngenix DNS zone",
			fmt.Sprintf("Could not delete DNS zone (DELETE), unexpected error: : %s", err.Error()),
		)
		return
	}

	tflog.Trace(ctx, "DNS zone was deleted successfully!")
}

func (r *dnsZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the ID from the import ID string. This usually contains the unique identifier for the DNS Zone.
	resourceID := req.ID

	// Set the ID field in the state.
	resp.Diagnostics.Append(
		resp.State.Set(ctx, &dnsZoneResourceModel{
			ID: types.StringValue(resourceID),
		})...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the DNS Zone by ID using the client.
	dnsZoneInt, _ := strconv.Atoi(resourceID)
	dnsZone, err := r.client.GetDnsZoneById(dnsZoneInt)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching resource", fmt.Sprintf("Could not fetch DNS Zone with ID %s: %s", resourceID, err))
		return
	}

	// Convert DNS Zone data to the resource model.
	dnsRecordsItems, err := r.DNSRecordsItemModelTransformation(dnsZone.Records)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while DNS zone Records creation and validation process",
			fmt.Sprintf("Could not create DNS zone Records, error: %s", err.Error()),
		)
		return
	}
	// Writing an updated / imported DNS model to the state.
	state := dnsZoneResourceModel{
		ID:          types.StringValue(resourceID),
		Name:        types.StringValue(dnsZone.Name),
		Records:     dnsRecordsItems,
		Comment:     types.StringValue(dnsZone.Comment),
		LastUpdated: types.StringValue(time.Now().Format(time.RFC850)),
	}

	// Set the state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
