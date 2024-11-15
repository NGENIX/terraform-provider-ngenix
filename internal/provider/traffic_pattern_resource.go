package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"ngenix/restapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                = &trafficPatternResource{}
	_ resource.ResourceWithConfigure   = &trafficPatternResource{}
	_ resource.ResourceWithImportState = &trafficPatternResource{}
)

// TrafficPatternResource is a helper function to simplify the provider implementation.
func TrafficPatternResource() resource.Resource {
	return &trafficPatternResource{}
}

// TrafficPatternDataSource is the data source implementation.
type trafficPatternResource struct {
	client *restapi.Client
}

// TrafficPatternResourceModel maps schema data.
type TrafficPatternResourceModel struct {
	ID          types.String     `tfsdk:"id"`
	Name        types.String     `tfsdk:"name"`
	Type        types.String     `tfsdk:"type"`
	ContentType types.String     `tfsdk:"content_type"`
	Patterns    []*PatternsModel `tfsdk:"patterns"`
	LastUpdated types.String     `tfsdk:"last_updated"`
}

// TrafficPatternsModel maps schema data.
type PatternsModel struct {
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

var (
	// Traffic pattern requirements.
	TrafficPatternTypes   = []string{"blacklist", "whitelist", "filterlist", "commonlist"}
	TrafficPatternROTypes = []string{"blacklist", "filterlist"}

	// Параметры addr, commonString, countryCode, httpMethod, asn и md5HashString взаимоисключающиеs.
	TrafficPatternContentTypes = []string{"addr", "commonString", "countryCode", "httpMethod", "asn", "md5HashString"}

	// Type whitelist is not compatible with following Content types.
	WhitelistIncompatibleTypes = []string{"commonString", "countryCode", "httpMethod", "asn", "md5HashString"}

	// ttl / expires are not compatible with the following Content types.
	TtlExpiresIncompatibleTypes = []string{"commonString", "countryCode", "httpMethod", "asn", "md5HashString"}

	// Поля ttl / expires взаимоисключающие.
	//	- поле ttl может быть null.
	//  - поле expires нет.

	// Supported HTTP methods.
	PatternHttpMethods = []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE", "HEAD", "PURGE"}
)

// Configure adds the provider configured client to the data source.
func (r *trafficPatternResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Metadata returns the resource type name.
func (r *trafficPatternResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_traffic_pattern"
}

// Schema defines the schema for the data source.
func (r *trafficPatternResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Traffic Pattern.",
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
				Description: "Traffic pattern name",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(42),
				},
			},
			"type": schema.StringAttribute{ // [ blacklist, whitelist, filterlist, commonlist ]
				Required:    true,
				Description: "Traffic pattern type",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"blacklist", "whitelist", "filterlist", "commonlist"}...),
				},
			},
			"content_type": schema.StringAttribute{ // [ addr, commonString, countryCode, httpMethod, asn, md5HashString ]
				Required:    true,
				Description: "Traffic pattern content type",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"addr", "commonString", "countryCode", "httpMethod", "asn", "md5HashString"}...),
				},
			},
			"patterns": schema.ListNestedAttribute{
				Description: "A list of Traffic patterns.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"addr": schema.StringAttribute{
							Optional:    true,
							Description: "Traffic pattern CIDR block - network CIDR block after validation",
						},
						"common_string": schema.StringAttribute{
							Optional:    true,
							Description: "Traffic pattern common string",
						},
						"country_code": schema.StringAttribute{
							Optional:    true,
							Description: "Traffic pattern country code",
						},
						"http_method": schema.StringAttribute{
							Optional:    true,
							Description: "Traffic pattern HTTP method",
						},
						"asn": schema.Int64Attribute{
							Optional:    true,
							Description: "Traffic pattern ASN",
						},
						"md5hash_string": schema.StringAttribute{
							Optional:    true,
							Description: "Traffic pattern MD5 hash string",
						},
						"ttl": schema.Int64Attribute{
							Optional:    true,
							Description: "Traffic pattern TTL",
						},
						"expires": schema.Int64Attribute{
							Optional:    true,
							Description: "Traffic pattern expiration date",
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Traffic pattern comment",
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{},
	}
}

func (r *trafficPatternResource) TrafficPatternModelTransformation(trafficPatterns []*PatternsModel, contentType string) ([]*restapi.Patterns, error) {
	patterns := []*restapi.Patterns{}
	if trafficPatterns == nil {
		return patterns, nil
	}
	for _, pattern := range trafficPatterns {
		switch {
		case contentType == "addr":
			if pattern.Ttl.IsNull() && pattern.Expires.IsNull() {
				patterns = append(patterns, &restapi.Patterns{
					Addr:    pattern.Addr.ValueStringPointer(),
					Comment: pattern.Comment.ValueString(),
				})
			} else if !pattern.Ttl.IsNull() && pattern.Expires.IsNull() {
				patterns = append(patterns, &restapi.Patterns{
					Addr:    pattern.Addr.ValueStringPointer(),
					Ttl:     pattern.Ttl.ValueInt64Pointer(),
					Comment: pattern.Comment.ValueString(),
				})
			} else if pattern.Ttl.IsNull() && !pattern.Expires.IsNull() {
				patterns = append(patterns, &restapi.Patterns{
					Addr:    pattern.Addr.ValueStringPointer(),
					Expires: pattern.Expires.ValueInt64Pointer(),
					Comment: pattern.Comment.ValueString(),
				})
			}
		case contentType == "commonString":
			patterns = append(patterns, &restapi.Patterns{
				CommonString: pattern.CommonString.ValueStringPointer(),
				Ttl:          pattern.Ttl.ValueInt64Pointer(),
				Expires:      pattern.Expires.ValueInt64Pointer(),
				Comment:      pattern.Comment.ValueString(),
			})
		case contentType == "countryCode":
			patterns = append(patterns, &restapi.Patterns{
				CountryCode: pattern.CountryCode.ValueStringPointer(),
				Ttl:         pattern.Ttl.ValueInt64Pointer(),
				Expires:     pattern.Expires.ValueInt64Pointer(),
				Comment:     pattern.Comment.ValueString(),
			})
		case contentType == "httpMethod":
			patterns = append(patterns, &restapi.Patterns{
				HttpMethod: pattern.HttpMethod.ValueStringPointer(),
				Ttl:        pattern.Ttl.ValueInt64Pointer(),
				Expires:    pattern.Expires.ValueInt64Pointer(),
				Comment:    pattern.Comment.ValueString(),
			})
		case contentType == "asn":
			patterns = append(patterns, &restapi.Patterns{
				Asn:     pattern.Asn.ValueInt64Pointer(),
				Ttl:     pattern.Ttl.ValueInt64Pointer(),
				Expires: pattern.Expires.ValueInt64Pointer(),
				Comment: pattern.Comment.ValueString(),
			})
		case contentType == "md5HashString":
			patterns = append(patterns, &restapi.Patterns{
				Md5HashString: pattern.Md5HashString.ValueStringPointer(),
				Ttl:           pattern.Ttl.ValueInt64Pointer(),
				Expires:       pattern.Expires.ValueInt64Pointer(),
				Comment:       pattern.Comment.ValueString(),
			})
		}
	}
	return patterns, nil
}

func (r *trafficPatternResource) TrafficPatternItemModelTransformation(trafficPatterns []*restapi.Patterns, contentType *string) ([]*PatternsModel, error) {
	patternsModel := []*PatternsModel{}
	if trafficPatterns == nil {
		return patternsModel, nil
	}
	for _, pattern := range trafficPatterns {
		switch {
		case *contentType == "addr":
			if pattern.Ttl == nil && pattern.Expires == nil {
				patternsModel = append(patternsModel, &PatternsModel{
					Addr:    types.StringPointerValue(pattern.Addr),
					Comment: types.StringValue(pattern.Comment),
				})
			} else if pattern.Ttl != nil && pattern.Expires == nil {
				patternsModel = append(patternsModel, &PatternsModel{
					Addr:    types.StringPointerValue(pattern.Addr),
					Ttl:     types.Int64Value(*pattern.Ttl),
					Comment: types.StringValue(pattern.Comment),
				})
			} else if pattern.Ttl == nil && pattern.Expires != nil {
				patternsModel = append(patternsModel, &PatternsModel{
					Addr:    types.StringPointerValue(pattern.Addr),
					Expires: types.Int64Value(*pattern.Expires),
					Comment: types.StringValue(pattern.Comment),
				})
			}
		case *contentType == "commonString":
			patternsModel = append(patternsModel, &PatternsModel{
				CommonString: types.StringPointerValue(pattern.CommonString),
				Ttl:          types.Int64PointerValue(pattern.Ttl),
				Expires:      types.Int64PointerValue(pattern.Expires),
				Comment:      types.StringValue(pattern.Comment),
			})
		case *contentType == "countryCode":
			patternsModel = append(patternsModel, &PatternsModel{
				CountryCode: types.StringPointerValue(pattern.CountryCode),
				Ttl:         types.Int64PointerValue(pattern.Ttl),
				Expires:     types.Int64PointerValue(pattern.Expires),
				Comment:     types.StringValue(pattern.Comment),
			})
		case *contentType == "httpMethod":
			patternsModel = append(patternsModel, &PatternsModel{
				HttpMethod: types.StringPointerValue(pattern.HttpMethod),
				Ttl:        types.Int64PointerValue(pattern.Ttl),
				Expires:    types.Int64PointerValue(pattern.Expires),
				Comment:    types.StringValue(pattern.Comment),
			})
		case *contentType == "asn":
			patternsModel = append(patternsModel, &PatternsModel{
				Asn:     types.Int64PointerValue(pattern.Asn),
				Ttl:     types.Int64PointerValue(pattern.Ttl),
				Expires: types.Int64PointerValue(pattern.Expires),
				Comment: types.StringValue(pattern.Comment),
			})
		case *contentType == "md5HashString":
			patternsModel = append(patternsModel, &PatternsModel{
				Md5HashString: types.StringPointerValue(pattern.Md5HashString),
				Ttl:           types.Int64PointerValue(pattern.Ttl),
				Expires:       types.Int64PointerValue(pattern.Expires),
				Comment:       types.StringValue(pattern.Comment),
			})
		}
	}
	return patternsModel, nil
}

// This function validates:
// 1. pattern field based on provided content type.
// 2. pattern fields base validattion for ip addr / max characters / match in a range and others ...
func (r *trafficPatternResource) validateContentTypeWithPattern(contentType *string, patterns []*PatternsModel) error {
	switch {
	case *contentType == "addr":
		for _, pattern := range patterns {
			if pattern.Addr.IsNull() || pattern.Addr.IsUnknown() {
				return fmt.Errorf("pattern must contain addr field for content type 'addr'")
			}
			ipRegex := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2]?\d))?$`)
			if !ipRegex.MatchString(pattern.Addr.ValueString()) {
				return fmt.Errorf("addr value is not a valid IP address")
			}
		}
	case *contentType == "commonString":
		for _, pattern := range patterns {
			if pattern.CommonString.IsNull() || pattern.CommonString.IsUnknown() {
				return fmt.Errorf("pattern must contain common_string field for content type 'commonString'")
			}
			if len(pattern.CommonString.ValueString()) > 255 {
				return fmt.Errorf("common_string value exceeds 255 characters")
			}
		}
	case *contentType == "countryCode":
		for _, pattern := range patterns {
			if pattern.CountryCode.IsNull() || pattern.CountryCode.IsUnknown() {
				return fmt.Errorf("pattern must contain country_code field for content type 'countryCode'")
			}
			ccRegex := regexp.MustCompile(`^[A-Z]{2}$`)
			if !ccRegex.MatchString(pattern.CountryCode.ValueString()) {
				return fmt.Errorf("country_code value is not valid")
			}
		}
	case *contentType == "httpMethod":
		validMethods := map[string]bool{
			"GET": true, "POST": true, "OPTIONS": true, "PUT": true,
			"PATCH": true, "DELETE": true, "HEAD": true, "PURGE": true,
		}
		for _, pattern := range patterns {
			if pattern.HttpMethod.IsNull() || pattern.HttpMethod.IsUnknown() {
				return fmt.Errorf("pattern must contain http_method field for content type 'httpMethod'")
			}
			if !validMethods[pattern.HttpMethod.ValueString()] {
				return fmt.Errorf("http_method value is not valid")
			}
		}
	case *contentType == "asn":
		for _, pattern := range patterns {
			if pattern.Asn.IsNull() || pattern.Asn.IsUnknown() {
				return fmt.Errorf("pattern must contain asn field for content type 'asn'")
			}
			if pattern.Asn.ValueInt64() < 0 || pattern.Asn.ValueInt64() > 4294967295 {
				return fmt.Errorf("asn value is not valid")
			}
		}
	case *contentType == "md5HashString":
		for _, pattern := range patterns {
			if pattern.Md5HashString.IsNull() || pattern.Md5HashString.IsUnknown() {
				return fmt.Errorf("pattern must contain md5hash_string field for content type 'md5HashString'")
			}
			md5Regex := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
			if !md5Regex.MatchString(pattern.Md5HashString.ValueString()) {
				return fmt.Errorf("md5hash_string value is not valid")
			}
		}
	default:
		return fmt.Errorf("content_type value is not valid")
	}
	return nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *trafficPatternResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan TrafficPatternResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Requirement - Checking content type values in patterns.
	err := r.validateContentTypeWithPattern(plan.ContentType.ValueStringPointer(), plan.Patterns)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic Pattern validation process with Content type field",
			fmt.Sprintf("Could not create Traffic Pattern, error: %s", err.Error()),
		)
		return
	}

	// Requirement - Checking type value in a range of list.
	if !restapi.IsValueInRange(plan.Type.ValueString(), TrafficPatternTypes) {
		resp.Diagnostics.AddError(
			"Traffic pattern type is out of range!",
			"Traffic pattern type should be in range of [ blacklist, whitelist, filterlist, commonlist ]!",
		)
		return
	}

	// Requirement: Types [filterlist, blacklist] are Read-Only.
	if restapi.IsValueInRange(plan.Type.ValueString(), TrafficPatternROTypes) {
		resp.Diagnostics.AddError(
			"Traffic pattern types [blacklist, filterlist] are Read-Only",
			"Traffic pattern types [blacklist, filterlist] are Read-Only! It is not possible to manage them by Terraform",
		)
		return
	}

	// Requirement - Checking content type value in a range of list.
	if !restapi.IsValueInRange(plan.ContentType.ValueString(), TrafficPatternContentTypes) {
		resp.Diagnostics.AddError(
			"Traffic pattern Content Type is out of range!",
			"Traffic pattern Content Type should be in range of [ addr, commonString, countryCode, httpMethod, asn, md5HashString ]!",
		)
		return
	}

	// Requirement - Type whitelist is not compatible with content_type httpMethod / countryCode / ...
	if plan.Type.ValueString() == "whitelist" {
		if restapi.IsValueInRange(plan.ContentType.ValueString(), WhitelistIncompatibleTypes) {
			resp.Diagnostics.AddError(
				"Type whitelist is not compatible with content types - [httpMethod, countryCode]",
				"Failed - Type whitelist is not compatible with next content types - commonString, countryCode, httpMethod, asn, md5HashString",
			)
			return
		}
	}

	// Requirement: TTL and Expires are not compatible with next Content Types.
	for _, pattern := range plan.Patterns {
		if restapi.IsValueInRange(plan.ContentType.ValueString(), TtlExpiresIncompatibleTypes) {
			if !pattern.Ttl.IsNull() || !pattern.Expires.IsNull() {
				resp.Diagnostics.AddError(
					"TTL & Expires are not compatible with this content type",
					"Failed - TTL & Expires are not compatible with next content Types - commonString, countryCode, httpMethod, asn, md5HashString",
				)
				return
			}
		}
	}

	// Requirements.
	// - Both fields TTL & Expires are not supported.
	// - Expires must not be null.
	networkWithMask := ""
	if plan.ContentType.ValueString() == "addr" {
		for _, pattern := range plan.Patterns {
			// Validating IP address.
			networkWithMask, err = restapi.GetNetworkAddressWithMask(pattern.Addr.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"IP address is not suitable, ip address should be in format: ip address / netmask",
					fmt.Sprintf("Failed - Could no get network address from pattern.addr %s field, error - %s, please add netmask!", pattern.Addr.ValueString(), err),
				)
				return
			}
			// Compare IP address with network address after validation.
			if pattern.Addr.ValueString() == networkWithMask {
				fmt.Printf("IP address %s is equal the network address %s", pattern.Addr.ValueString(), networkWithMask)
			} else {
				resp.Diagnostics.AddError(
					"Invalid IP address",
					fmt.Sprintf("Failed - IP address %s is NOT equal the network IP address %s, please enter the network IP address instead", pattern.Addr.ValueString(), networkWithMask),
				)
				return
			}
			// Both TTL and Expires are not supported.
			if !pattern.Ttl.IsNull() && !pattern.Expires.IsNull() {
				resp.Diagnostics.AddError(
					"Both fields TTL & Expires are not supported",
					"Failed - Both fields TTL & Expires are not supported - Assign only one field TTL or Expires",
				)
				return
			}
		}
	}

	// Traffic patterns.
	patterns, err := r.TrafficPatternModelTransformation(plan.Patterns, plan.ContentType.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic Patterns creation and validation process",
			fmt.Sprintf("Could not create Traffic Patterns, please check a page with errors, error: %s", err.Error()),
		)
		return
	}

	trafficPattern := restapi.TrafficPattern{
		Name:        plan.Name.ValueStringPointer(),
		Type:        plan.Type.ValueStringPointer(),
		ContentType: plan.ContentType.ValueStringPointer(),
		CustomerRef: &restapi.TPCustomerRef{
			ID: r.client.CustomerId(),
		},
		Patterns: patterns,
	}

	// Create new Traffic pattern.
	createdTP, err := r.client.CreateNewTrafficPatternForCustomer(trafficPattern, r.client.CustomerId())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Traffic pattern",
			fmt.Sprintf("Could not create Traffic pattern, unexpected error: %s", err.Error()),
		)
		return
	}

	// From TrafficPattern to TPModel
	patternsModel, err := r.TrafficPatternItemModelTransformation(createdTP.Patterns, createdTP.ContentType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic pattern creation and validation process",
			fmt.Sprintf("Could not create Traffic pattern, error: %s", err.Error()),
		)
		return
	}
	tpId := r.client.GetTPIDByName(*createdTP.Name)
	plan.ID = types.StringValue(strconv.Itoa(tpId))
	plan.Name = types.StringValue(*createdTP.Name)
	plan.Type = types.StringValue(*createdTP.Type)
	plan.ContentType = types.StringValue(*createdTP.ContentType)
	plan.Patterns = patternsModel
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "Traffic pattern was created successfully!")

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *trafficPatternResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state TrafficPatternResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Getting Traffic Pattern by ID
	tpId, _ := strconv.Atoi(state.ID.ValueString())
	//tpId := r.client.GetTPIDByName(state.Name.ValueString())
	trafficPattern, err := r.client.GetTrafficPatternById(tpId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ngenix Traffic Pattern",
			fmt.Sprintf("Could not read Ngenix Traffic Pattern by ID = %d, error: %s", tpId, err.Error()),
		)
		return
	}

	// Map response body to model.
	patterns, err := r.TrafficPatternItemModelTransformation(trafficPattern.Patterns, trafficPattern.ContentType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic pattern creation and validation process",
			fmt.Sprintf("Could not create Traffic pattern, error: %s", err.Error()),
		)
		return
	}
	state.Name = types.StringValue(*trafficPattern.Name)
	state.Type = types.StringValue(*trafficPattern.Type)
	state.ContentType = types.StringValue(*trafficPattern.ContentType)
	state.Patterns = patterns

	tflog.Trace(ctx, "Traffic Pattern was read successfully!")

	// Set refreshed state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *trafficPatternResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan TrafficPatternResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Requirement - Checking content type values in patterns.
	err := r.validateContentTypeWithPattern(plan.ContentType.ValueStringPointer(), plan.Patterns)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic Pattern validation process with content type field",
			fmt.Sprintf("Could not create Traffic Pattern, error: %s", err.Error()),
		)
		return
	}

	// Requirement - Checking type value in a range of list.
	if !restapi.IsValueInRange(plan.Type.ValueString(), TrafficPatternTypes) {
		resp.Diagnostics.AddError(
			"Traffic pattern type is out of range!",
			"Traffic pattern type should be in range of [ blacklist, whitelist, filterlist, commonlist ]!",
		)
		return
	}

	// Requirement: Types [filterlist, blacklist] are Read-Only.
	if restapi.IsValueInRange(plan.Type.ValueString(), TrafficPatternROTypes) {
		resp.Diagnostics.AddError(
			"Traffic pattern types [blacklist, filterlist] are Read-Only",
			"Traffic pattern types [blacklist, filterlist] are Read-Only! It is not possible to manage them by Terraform",
		)
		return
	}

	// Requirement - Checking content type value in a range of list.
	if !restapi.IsValueInRange(plan.ContentType.ValueString(), TrafficPatternContentTypes) {
		resp.Diagnostics.AddError(
			"Traffic pattern Content Type is out of range!",
			"Traffic pattern Content Type should be in range of [ addr, commonString, countryCode, httpMethod, asn, md5HashString ]!",
		)
		return
	}

	// Requirements.
	// - Both fields TTL & Expires are not supported.
	// - Expires must not be null.
	networkWithMask := ""
	if plan.ContentType.ValueString() == "addr" {
		for _, pattern := range plan.Patterns {
			// Validating IP address.
			networkWithMask, err = restapi.GetNetworkAddressWithMask(pattern.Addr.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"IP address is not suitable, ip address should be in format: ip address / netmask",
					fmt.Sprintf("Failed - Could no get network address from pattern.addr %s field, error - %s, please add netmask!", pattern.Addr.ValueString(), err),
				)
				return
			}
			// Compare IP address with network address after validation.
			if pattern.Addr.ValueString() == networkWithMask {
				fmt.Printf("IP address %s is equal the network address %s", pattern.Addr.ValueString(), networkWithMask)
			} else {
				resp.Diagnostics.AddError(
					"Invalid IP address",
					fmt.Sprintf("Failed - IP address %s is NOT equal the network IP address %s, please enter the network IP address instead", pattern.Addr.ValueString(), networkWithMask),
				)
				return
			}
			// Both TTL and Expires are not supported.
			if !pattern.Ttl.IsNull() && !pattern.Expires.IsNull() {
				resp.Diagnostics.AddError(
					"Both fields TTL & Expires are not supported",
					"Failed - Both fields TTL & Expires are not supported - Assign only one field TTL or Expires",
				)
				return
			}
		}
	}

	// Traffic patterns.
	patterns, er := r.TrafficPatternModelTransformation(plan.Patterns, plan.ContentType.ValueString())
	if er != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic Patterns creation and validation process",
			fmt.Sprintf("Could not create Traffic Patterns, please check a page with errors, error: %s", er.Error()),
		)
		return
	}
	//
	trafficPattern := restapi.TrafficPattern{
		Name:     plan.Name.ValueStringPointer(),
		Patterns: patterns,
	}

	// Update existing Traffic pattern.
	tpId, _ := strconv.Atoi(plan.ID.ValueString())
	updatedTrafficPattern, errTp := r.client.UpdateTrafficPatternById(trafficPattern, tpId)
	if errTp != nil {
		resp.Diagnostics.AddError(
			"Error Updating Ngenix Traffic Pattern",
			fmt.Sprintf("Could not update Traffic Pattern (PATCH), unexpected error: %s", errTp.Error()),
		)
		return
	}

	// Update Traffic pattern resource state with updated items and timestamp.
	patternsModel, err := r.TrafficPatternItemModelTransformation(updatedTrafficPattern.Patterns, updatedTrafficPattern.ContentType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic pattern creation and validation process",
			fmt.Sprintf("Could not create Traffic pattern, error: %s", err.Error()),
		)
		return
	}
	tpId = r.client.GetTPIDByName(*updatedTrafficPattern.Name)
	plan.ID = types.StringValue(strconv.Itoa(tpId))
	plan.Name = types.StringValue(*updatedTrafficPattern.Name)
	plan.Type = types.StringValue(*updatedTrafficPattern.Type)
	plan.ContentType = types.StringValue(*updatedTrafficPattern.ContentType)
	plan.Patterns = patternsModel
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "Traffic pattern was updated successfully!")

	// Set state to fully populated data.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete the resource and removes the Terraform state on success.
func (r *trafficPatternResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state.
	var state TrafficPatternResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Traffic Pattern.
	tpId := r.client.GetTPIDByName(state.Name.ValueString())
	err := r.client.DeleteTrafficPatternById(tpId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Ngenix traffic pattern",
			fmt.Sprintf("Could not delete Traffic Pattern (DELETE), unexpected error: : %s", err.Error()),
		)
		return
	}

	tflog.Trace(ctx, "Traffic pattern was deleted successfully!")
}

// Import Traffic pattern state by TP ID.
func (r *trafficPatternResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the ID from the import ID string. This usually contains the unique identifier for the DNS Zone.
	resourceID := req.ID

	// Set the ID field in the state.
	resp.Diagnostics.Append(
		resp.State.Set(ctx, &TrafficPatternResourceModel{
			ID: types.StringValue(resourceID),
		})...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the Traffic patterns by ID using the client.
	tpIdInt, _ := strconv.Atoi(resourceID)
	trafficPattern, err := r.client.GetTrafficPatternById(tpIdInt)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching resource", fmt.Sprintf("Could not fetch DNS Zone with ID %s: %s", resourceID, err))
		return
	}

	// Convert Traffic patterns data to the resource model.
	patterns, err := r.TrafficPatternItemModelTransformation(trafficPattern.Patterns, trafficPattern.ContentType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic pattern creation and validation process",
			fmt.Sprintf("Could not create Traffic pattern, error: %s", err.Error()),
		)
		return
	}
	// Writing an updated / imported traffic patterns to the state.
	state := TrafficPatternResourceModel{
		ID:          types.StringValue(resourceID),
		Name:        types.StringValue(*trafficPattern.Name),
		Type:        types.StringValue(*trafficPattern.Type),
		ContentType: types.StringValue(*trafficPattern.ContentType),
		Patterns:    patterns,
		LastUpdated: types.StringValue(time.Now().Format(time.RFC850)),
	}

	tflog.Trace(ctx, "Traffic Pattern was imported successfully!")

	// Set the state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
