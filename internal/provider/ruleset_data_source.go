package provider

import (
	"context"
	"fmt"

	"ngenix/restapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ruleSetDataSource{}
	_ datasource.DataSourceWithConfigure = &ruleSetDataSource{}
)

// RulesetDataSource is a helper function to simplify the provider implementation.
func RulesetDataSource() datasource.DataSource {
	return &ruleSetDataSource{}
}

// ruleSetDataSource is the data source implementation.
type ruleSetDataSource struct {
	client *restapi.Client
}

// dnsRecordSetDataSourceModel maps the data source schema data.
type ruleSetDataSourceModel struct {
	Rulesets []ruleSetModel `tfsdk:"rulesets"`
}

// ruleSetModel maps schema data.
type ruleSetModel struct {
	Name    types.String         `tfsdk:"name"`
	Enabled types.Bool           `tfsdk:"enabled"`
	Rules   []rulesetsRulesModel `tfsdk:"rules"`
}

// dnsRecordSetModel maps schema data.
type rulesetsRulesModel struct {
	//ID         types.Int64              `tfsdk:"id"`
	Name       types.String             `tfsdk:"name"`
	Enabled    types.Bool               `tfsdk:"enabled"`
	Conditions []rulesetConditionsModel `tfsdk:"conditions"`
	Actions    []rulesetActionsModel    `tfsdk:"actions"`
}

type rulesetConditionsModel struct {
	Function types.String                  `tfsdk:"function"`
	Negation types.Bool                    `tfsdk:"negation"`
	Params   []rulesetConditionParamsModel `tfsdk:"params"`
}

type rulesetConditionParamsModel struct {
	Variable          types.String               `tfsdk:"variable"`
	Value             types.String               `tfsdk:"value"`
	TrafficPatternRef *trafficPatternRefModel    `tfsdk:"trafficpattern_ref"`
	Option            *conditionParamOptionModel `tfsdk:"option"`
}

type conditionParamOptionModel struct {
	Name  types.String `tfsdk:"name"` // ignore_case
	Value types.Bool   `tfsdk:"bool"`
}

type trafficPatternRefModel struct {
	ID types.Int64 `tfsdk:"id"`
}

type rulesetActionsModel struct {
	Action types.String   `tfsdk:"action"`
	Params []types.String `tfsdk:"params"`
}

// Metadata returns the data source type name.
func (d *ruleSetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rulesets"
}

// Configure adds the provider configured client to the data source.
func (d *ruleSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ruleSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rulesets": schema.ListNestedAttribute{
				Description: "A list of Rulesets.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Ruleset name.",
						},
						"enabled": schema.BoolAttribute{
							Required:    true,
							Description: "Ruleset enabled flag.",
						},
						"rules": schema.ListNestedAttribute{
							Description: "A list of Ruleset rules.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:    true,
										Description: "Ruleset name",
									},
									"enabled": schema.BoolAttribute{
										Required:    true,
										Description: "Ruleset enabled flag.",
									},
									"conditions": schema.ListNestedAttribute{
										Description: "Ruleset rules conditions list.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"function": schema.StringAttribute{
													Required:    true,
													Description: "Ruleset condition function type",
													Validators: []validator.String{
														stringvalidator.OneOf([]string{"in", "equals", "contains", "begins", "ends", "true", "exists"}...),
													},
												},
												"negation": schema.BoolAttribute{
													Required:    true,
													Description: "Ruleset negation (default = false)",
												},
												"params": schema.ListNestedAttribute{
													Description: "Ruleset rules conditions.",
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"variable": schema.StringAttribute{
																Optional:    true,
																Description: "Ruleset variable",
															},
															"value": schema.StringAttribute{
																Optional:    true,
																Description: "Ruleset value",
															},
															"trafficpattern_ref": schema.ObjectAttribute{
																AttributeTypes: map[string]attr.Type{
																	"id": types.Int64Type,
																},
																Optional:    true,
																Description: "Ruleset condition traffic pattern reference id.",
															},
															"option": schema.ObjectAttribute{
																AttributeTypes: map[string]attr.Type{
																	"name":  types.StringType,
																	"value": types.BoolType,
																},
																Optional:    true,
																Description: "Ruleset condition parameter option.",
															},
														},
													},
												},
											},
										},
									},
									"actions": schema.ListNestedAttribute{
										Description: "Ruleset rules actions list.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"action": schema.StringAttribute{
													Required:    true,
													Description: "Ruleset action name",
													Validators: []validator.String{
														stringvalidator.OneOf([]string{"allow", "deny", "redirect", "setHeader", "delHeader", "jsChallenge"}...),
													},
												},
												"params": schema.ListAttribute{
													ElementType: types.StringType,
													Required:    true,
													Description: "Ruleset action parameter",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Blocks:      map[string]schema.Block{},
		Description: "Interface with the Ngenix Platform REST API",
	}
}

// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-data-source-read.
// Read refreshes the Terraform state with the latest data.
func (d *ruleSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ruleSetDataSourceModel

	// Getting all Rulesets for provided username.
	rulesetsList := d.client.GetAllRulesetsList()
	// Map response body to model.
	for _, ruleset := range rulesetsList {
		rulesetState := ruleSetModel{
			// Name.
			Name:    types.StringValue(ruleset.Name),
			Enabled: types.BoolValue(ruleset.Enabled),
		}
		// Ruleset Rules.
		for _, rule := range ruleset.Rules {
			rulesConditions := []rulesetConditionsModel{}
			rulesActions := []rulesetActionsModel{}
			if rule.Conditions != nil && rule.Actions != nil {
				// Conditions.
				for _, condition := range rule.Conditions {
					conditionsParams := []rulesetConditionParamsModel{}
					for _, param := range condition.Params {
						if param.TrafficPatternRef == nil && param.Option == nil {
							conditionsParams = append(conditionsParams, rulesetConditionParamsModel{
								Variable: types.StringPointerValue(param.Variable),
								Value:    types.StringPointerValue(param.Value),
							})
						} else if param.TrafficPatternRef != nil && param.Option == nil {
							conditionsParams = append(conditionsParams, rulesetConditionParamsModel{
								Variable: types.StringPointerValue(param.Variable),
								Value:    types.StringPointerValue(param.Value),
								TrafficPatternRef: &trafficPatternRefModel{
									ID: types.Int64Value(param.TrafficPatternRef.ID),
								},
							})
						} else if param.TrafficPatternRef != nil && param.Option != nil {
							conditionsParams = append(conditionsParams, rulesetConditionParamsModel{
								Variable: types.StringPointerValue(param.Variable),
								Value:    types.StringPointerValue(param.Value),
								TrafficPatternRef: &trafficPatternRefModel{
									ID: types.Int64Value(param.TrafficPatternRef.ID),
								},
								Option: &conditionParamOptionModel{
									Name:  types.StringValue(param.Option.Name),
									Value: types.BoolValue(param.Option.Value),
								},
							})
						}
					}
					rulesConditions = append(rulesConditions, rulesetConditionsModel{
						Function: types.StringValue(condition.Function),
						Negation: types.BoolValue(condition.Negation),
						Params:   conditionsParams,
					})
				}
				// Actions.
				for _, action := range rule.Actions {
					var actionParams []types.String
					for _, aParam := range action.Params {
						actionParams = append(actionParams, types.StringValue(aParam))
					}
					rulesActions = append(rulesActions, rulesetActionsModel{
						Action: types.StringValue(action.Action),
						Params: actionParams,
					})
				}
			}
			// Ruleset Model.
			rulesetState.Rules = append(rulesetState.Rules, rulesetsRulesModel{
				Name:       types.StringValue(rule.Name),
				Enabled:    types.BoolValue(rule.Enabled),
				Conditions: rulesConditions,
				Actions:    rulesActions,
			})
		}

		state.Rulesets = append(state.Rulesets, rulesetState)
	}

	// Set state.
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
