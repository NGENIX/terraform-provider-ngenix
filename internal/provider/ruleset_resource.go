package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"ngenix/restapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ruleSetResource{}
	_ resource.ResourceWithConfigure   = &ruleSetResource{}
	_ resource.ResourceWithImportState = &ruleSetResource{}
)

// RulesetDataSource is a helper function to simplify the provider implementation.
func RulesetResource() resource.Resource {
	return &ruleSetResource{}
}

// ruleSetDataSource is the data source implementation.
type ruleSetResource struct {
	client *restapi.Client
}

// ruleSetModel maps schema data.
type ruleSetResourceModel struct {
	ID          types.String                 `tfsdk:"id"`
	Name        types.String                 `tfsdk:"name"`
	Enabled     types.Bool                   `tfsdk:"enabled"`
	Rules       []rulesetsRulesResourceModel `tfsdk:"rules"`
	LastUpdated types.String                 `tfsdk:"last_updated"`
}

// dnsRecordSetModel maps schema data.
type rulesetsRulesResourceModel struct {
	Name       types.String                 `tfsdk:"name"`
	Enabled    types.Bool                   `tfsdk:"enabled"`
	Conditions []rulesetConditionsModelItem `tfsdk:"conditions"`
	Actions    []rulesetActionsModelItem    `tfsdk:"actions"`
}

type rulesetConditionsModelItem struct {
	Function types.String                      `tfsdk:"function"`
	Negation types.Bool                        `tfsdk:"negation"`
	Params   []rulesetConditionParamsModelItem `tfsdk:"params"`
}

type rulesetConditionParamsModelItem struct {
	Variable          types.String                   `tfsdk:"variable"`
	Value             types.String                   `tfsdk:"value"`
	TrafficPatternRef *trafficPatternRef             `tfsdk:"trafficpattern_ref"`
	Option            *conditionParamOptionModelItem `tfsdk:"option"`
}

type conditionParamOptionModelItem struct {
	Name  types.String `tfsdk:"name"`
	Value types.Bool   `tfsdk:"bool"`
}

type trafficPatternRef struct {
	ID types.Int64 `tfsdk:"id"`
}

type rulesetActionsModelItem struct {
	Action types.String   `tfsdk:"action"`
	Params []types.String `tfsdk:"params"`
}

// Metadata returns the data source type name.
func (r *ruleSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ruleset"
}

// Configure adds the provider configured client to the data source.
func (r *ruleSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*restapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the data source.
func (r *ruleSetResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
				Description: "Ruleset name.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Ruleset enabled flag.",
			},
			"rules": schema.ListNestedAttribute{
				Description: "A list of Ruleset rules. Maximum items - 40",
				Required:    true,
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
							Description: "Ruleset rules conditions list. Max conditions - 5",
							Required:    true,
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
										Optional:    true,
										//MaxItems:    3,
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
							Description: "Ruleset rules actions list. Max actions - 5",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"action": schema.StringAttribute{
										Optional:    true,
										Description: "Ruleset action name",
										Validators: []validator.String{
											stringvalidator.OneOf([]string{"allow", "deny", "redirect", "setHeader", "delHeader", "jsChallenge"}...),
										},
									},
									"params": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "Ruleset action parameters",
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

// Ruleset to Rules Item model.
func (r *ruleSetResource) RulesetToRulesItemTransformation(rules []restapi.Rules) ([]rulesetsRulesResourceModel, error) {
	// Rules.
	rulesetRules := []rulesetsRulesResourceModel{}
	for _, rule := range rules {
		rulesConditions := []rulesetConditionsModelItem{}
		rulesActions := []rulesetActionsModelItem{}
		// Conditions.
		for _, condition := range rule.Conditions {
			conditionsParams := []rulesetConditionParamsModelItem{}
			for _, param := range condition.Params {
				if param.TrafficPatternRef != nil && param.Option == nil {
					conditionsParams = append(conditionsParams, rulesetConditionParamsModelItem{
						Variable: types.StringPointerValue(param.Variable),
						Value:    types.StringPointerValue(param.Value),
						TrafficPatternRef: &trafficPatternRef{
							ID: types.Int64Value(param.TrafficPatternRef.ID),
						},
					})
				} else if param.TrafficPatternRef == nil && param.Option == nil {
					conditionsParams = append(conditionsParams, rulesetConditionParamsModelItem{
						Variable: types.StringPointerValue(param.Variable),
						Value:    types.StringPointerValue(param.Value),
					})
				} else if param.TrafficPatternRef != nil && param.Option != nil {
					conditionsParams = append(conditionsParams, rulesetConditionParamsModelItem{
						Variable: types.StringPointerValue(param.Variable),
						Value:    types.StringPointerValue(param.Value),
						TrafficPatternRef: &trafficPatternRef{
							ID: types.Int64Value(param.TrafficPatternRef.ID),
						},
						Option: &conditionParamOptionModelItem{
							Name:  types.StringValue(param.Option.Name),
							Value: types.BoolValue(param.Option.Value),
						},
					})
				}
			}
			rulesConditions = append(rulesConditions, rulesetConditionsModelItem{
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
			rulesActions = append(rulesActions, rulesetActionsModelItem{
				Action: types.StringValue(action.Action),
				Params: actionParams,
			})
		}
		// Ruleset Model.
		rulesetRules = append(rulesetRules, rulesetsRulesResourceModel{
			Name:       types.StringValue(rule.Name),
			Enabled:    types.BoolValue(rule.Enabled),
			Conditions: rulesConditions,
			Actions:    rulesActions,
		})
	}
	return rulesetRules, nil
}

// Ruleset to Rules Item model.
func (r *ruleSetResource) RulesModelItemTransformation(rules []rulesetsRulesResourceModel) ([]restapi.Rules, error) {
	// Rules.
	rulesetRules := []restapi.Rules{}
	for _, rule := range rules {
		rulesConditions := []restapi.Conditions{}
		rulesActions := []restapi.Actions{}
		// Conditions.
		for _, condition := range rule.Conditions {
			conditionsParams := []restapi.ConditionParams{}
			for _, param := range condition.Params {
				if param.TrafficPatternRef != nil && param.Option == nil {
					conditionsParams = append(conditionsParams, restapi.ConditionParams{
						Variable: param.Variable.ValueStringPointer(),
						Value:    param.Value.ValueStringPointer(),
						TrafficPatternRef: &restapi.TrafficPatternRef{
							ID: param.TrafficPatternRef.ID.ValueInt64(),
						},
					})
				} else if param.TrafficPatternRef == nil && param.Option == nil {
					conditionsParams = append(conditionsParams, restapi.ConditionParams{
						Variable: param.Variable.ValueStringPointer(),
						Value:    param.Value.ValueStringPointer(),
					})
				} else if param.TrafficPatternRef != nil && param.Option != nil {
					conditionsParams = append(conditionsParams, restapi.ConditionParams{
						Variable: param.Variable.ValueStringPointer(),
						Value:    param.Value.ValueStringPointer(),
						TrafficPatternRef: &restapi.TrafficPatternRef{
							ID: param.TrafficPatternRef.ID.ValueInt64(),
						},
						Option: &restapi.ParamOption{
							Name:  param.Option.Name.ValueString(),
							Value: param.Option.Value.ValueBool(),
						},
					})
				}
			}
			rulesConditions = append(rulesConditions, restapi.Conditions{
				Function: condition.Function.ValueString(),
				Negation: condition.Negation.ValueBool(),
				Params:   conditionsParams,
			})
		}
		// Actions.
		for _, action := range rule.Actions {
			var actionParams []string
			for _, aParam := range action.Params {
				actionParams = append(actionParams, aParam.ValueString())
			}
			rulesActions = append(rulesActions, restapi.Actions{
				Action: action.Action.ValueString(),
				Params: actionParams,
			})
		}
		// Ruleset Rules struct.
		rulesetRules = append(rulesetRules, restapi.Rules{
			Name:       rule.Name.ValueString(),
			Enabled:    rule.Enabled.ValueBool(),
			Conditions: rulesConditions,
			Actions:    rulesActions,
		})
	}
	return rulesetRules, nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *ruleSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var plan ruleSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ruleset rules.
	rules, err := r.RulesModelItemTransformation(plan.Rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Ruleset rules creation and validation process",
			fmt.Sprintf("Could not create Ruleset Rules, please check a page with errors, error: %s", err.Error()),
		)
		return
	}

	ruleset := restapi.Ruleset{
		Name:    plan.Name.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
		CustomerRef: &restapi.RulesetCustomerRef{
			ID: r.client.CustomerId(),
		},
		Rules: rules,
	}

	// Create new Ruleset.
	createdRs, err := r.client.CreateRuleset(ruleset)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating ruleset",
			fmt.Sprintf("Could not create Ruleset, unexpected error: %s", err.Error()),
		)
		return
	}

	// From Ruleset to RS model.
	rulesModel, err := r.RulesetToRulesItemTransformation(createdRs.Rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Ruleset rules creation and validation process",
			fmt.Sprintf("Could not create Ruleset, error: %s", err.Error()),
		)
		return
	}

	rsId := r.client.GetRulesetIDByName(createdRs.Name)
	plan.ID = types.StringValue(strconv.Itoa(rsId))
	plan.Name = types.StringValue(createdRs.Name)
	plan.Enabled = types.BoolValue(createdRs.Enabled)
	plan.Rules = rulesModel
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "Ruleset was created successfully!")

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ruleSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state.
	var state ruleSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Getting Ruleset by ID.
	rsId := r.client.GetRulesetIDByName(state.Name.ValueString())
	ruleset, err := r.client.GetRulesetById(rsId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ngenix Ruleset",
			fmt.Sprintf("Could not read Ngenix Ruleset by ID = %d, error: %s", rsId, err.Error()),
		)
		return
	}

	// Rules.
	rulesetRules, _ := r.RulesetToRulesItemTransformation(ruleset.Rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Ruleset creation and validation process",
			fmt.Sprintf("Could not read Ruleset, error: %s", err.Error()),
		)
		return
	}

	state.Name = types.StringValue(ruleset.Name)
	state.Enabled = types.BoolValue(ruleset.Enabled)
	state.Rules = rulesetRules

	tflog.Trace(ctx, "Ruleset was read successfully!")

	// Set refreshed state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// // Update updates the resource and sets the updated Terraform state on success.
func (r *ruleSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan ruleSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ruleset rules.
	rules, err := r.RulesModelItemTransformation(plan.Rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Ruleset rules creation and validation process",
			fmt.Sprintf("Could not create Ruleset Rules, please check a page with errors, error: %s", err.Error()),
		)
		return
	}
	//
	ruleset := restapi.Ruleset{
		Name:    plan.Name.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
		Rules:   rules,
	}

	// Update existing Ruleset.
	rsId := r.client.GetRulesetIDByName(plan.Name.ValueString())
	_, err = r.client.UpdateRulesetById(ruleset, rsId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Ngenix Ruleset",
			fmt.Sprintf("Could not update Ruleset (PATCH), unexpected error: %s", err.Error()),
		)
		return
	}

	// getting updated ruleset.
	updatedRuleset, err := r.client.GetRulesetById(rsId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ngenix Ruleset by ID",
			fmt.Sprintf("Could not read Ngenix Ruleset by ID %d, error: %s", rsId, err.Error()),
		)
		return
	}

	// Update Ruleset resource state with updated items and timestamp.
	rulesModel, err := r.RulesetToRulesItemTransformation(updatedRuleset.Rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Traffic pattern creation and validation process",
			fmt.Sprintf("Could not create Traffic pattern, error: %s", err.Error()),
		)
		return
	}
	//
	rsId = r.client.GetRulesetIDByName(updatedRuleset.Name)
	plan.ID = types.StringValue(strconv.Itoa(rsId))
	plan.Name = types.StringValue(updatedRuleset.Name)
	plan.Enabled = types.BoolValue(updatedRuleset.Enabled)
	plan.Rules = rulesModel
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	tflog.Trace(ctx, "Traffic pattern was updated successfully!")

	// Set state to fully populated data.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete the resource and removes the Terraform state on success.
func (r *ruleSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state.
	var state ruleSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Ruleset.
	rsId := r.client.GetRulesetIDByName(state.Name.ValueString())
	err := r.client.DeleteRulesetById(rsId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Ngenix ruleset",
			fmt.Sprintf("Could not delete Ruleset (DELETE), unexpected error: : %s", err.Error()),
		)
		return
	}

	tflog.Trace(ctx, "Ruleset was deleted successfully!")
}

// Import Traffic pattern state by TP ID.
func (r *ruleSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the ID from the import ID string. This usually contains the unique identifier for the DNS Zone.
	resourceID := req.ID

	// Set the ID field in the state.
	resp.Diagnostics.Append(
		resp.State.Set(ctx, &ruleSetResourceModel{
			ID: types.StringValue(resourceID),
		})...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the Traffic patterns by ID using the client.
	rsIdInt, _ := strconv.Atoi(resourceID)
	ruleset, err := r.client.GetRulesetById(rsIdInt)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching resource", fmt.Sprintf("Could not fetch DNS Zone with ID %s: %s", resourceID, err))
		return
	}

	// From Ruleset to Ruleset Item Model.

	// Convert Traffic patterns data to the resource model.
	rulesetRules, _ := r.RulesetToRulesItemTransformation(ruleset.Rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while Ruleset creation and validation process",
			fmt.Sprintf("Could not read Ruleset, error: %s", err.Error()),
		)
		return
	}

	// Writing an updated / imported traffic patterns to the state.
	state := ruleSetResourceModel{
		ID:          types.StringValue(resourceID),
		Name:        types.StringValue(ruleset.Name),
		Enabled:     types.BoolValue(ruleset.Enabled),
		Rules:       rulesetRules,
		LastUpdated: types.StringValue(time.Now().Format(time.RFC850)),
	}

	tflog.Trace(ctx, "Traffic Pattern was imported successfully!")

	// Set the state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
