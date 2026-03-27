package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kawasima/bouncr-client-go"
)

var (
	_ resource.Resource                = &oidcApplicationResource{}
	_ resource.ResourceWithImportState = &oidcApplicationResource{}
)

type oidcApplicationResource struct {
	client *bouncr.Client
}

type oidcApplicationResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	HomeURI     types.String `tfsdk:"home_uri"`
	CallbackURI types.String `tfsdk:"callback_uri"`
	GrantTypes  types.Set    `tfsdk:"grant_types"`
	Permissions types.Set    `tfsdk:"permissions"`
}

func NewOidcApplicationResource() resource.Resource {
	return &oidcApplicationResource{}
}

func (r *oidcApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oidc_application"
}

func (r *oidcApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"home_uri": schema.StringAttribute{
				Optional: true,
			},
			"callback_uri": schema.StringAttribute{
				Optional: true,
			},
			"grant_types": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"permissions": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *oidcApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*bouncr.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *bouncr.Client")
		return
	}
	r.client = client
}

func (r *oidcApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oidcApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var grantTypes []string
	resp.Diagnostics.Append(plan.GrantTypes.ElementsAs(ctx, &grantTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	if !plan.Permissions.IsNull() {
		resp.Diagnostics.Append(plan.Permissions.ElementsAs(ctx, &permissions, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	_, err := r.client.CreateOidcApplication(ctx, &bouncr.OidcApplicationCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		HomeURI:     plan.HomeURI.ValueString(),
		CallbackURI: plan.CallbackURI.ValueString(),
		GrantTypes:  grantTypes,
		Permissions: permissions,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating OIDC application", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oidcApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oidcApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oidcApp, err := r.client.FindOidcApplication(ctx, state.Name.ValueString())
	if isNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading OIDC application", err.Error())
		return
	}

	state.Name = types.StringValue(oidcApp.Name)
	state.Description = types.StringValue(oidcApp.Description)
	state.HomeURI = types.StringValue(oidcApp.HomeURI)
	state.CallbackURI = types.StringValue(oidcApp.CallbackURI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *oidcApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan oidcApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state oidcApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var grantTypes []string
	resp.Diagnostics.Append(plan.GrantTypes.ElementsAs(ctx, &grantTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	if !plan.Permissions.IsNull() {
		resp.Diagnostics.Append(plan.Permissions.ElementsAs(ctx, &permissions, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	_, err := r.client.UpdateOidcApplication(ctx, state.Name.ValueString(), &bouncr.OidcApplicationUpdateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		HomeURI:     plan.HomeURI.ValueString(),
		CallbackURI: plan.CallbackURI.ValueString(),
		GrantTypes:  grantTypes,
		Permissions: permissions,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating OIDC application", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oidcApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state oidcApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOidcApplication(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting OIDC application", err.Error())
	}
}

func (r *oidcApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
