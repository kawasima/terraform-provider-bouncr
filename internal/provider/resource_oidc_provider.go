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
	_ resource.Resource                = &oidcProviderResource{}
	_ resource.ResourceWithImportState = &oidcProviderResource{}
)

type oidcProviderResource struct {
	client *bouncr.Client
}

type oidcProviderResourceModel struct {
	Name                    types.String `tfsdk:"name"`
	ClientID                types.String `tfsdk:"client_id"`
	ClientSecret            types.String `tfsdk:"client_secret"`
	Scope                   types.String `tfsdk:"scope"`
	ResponseType            types.String `tfsdk:"response_type"`
	AuthorizationEndpoint   types.String `tfsdk:"authorization_endpoint"`
	TokenEndpoint           types.String `tfsdk:"token_endpoint"`
	TokenEndpointAuthMethod types.String `tfsdk:"token_endpoint_auth_method"`
	RedirectURI             types.String `tfsdk:"redirect_uri"`
	PkceEnabled             types.Bool   `tfsdk:"pkce_enabled"`
	JwksURI                 types.String `tfsdk:"jwks_uri"`
	Issuer                  types.String `tfsdk:"issuer"`
}

func NewOidcProviderResource() resource.Resource {
	return &oidcProviderResource{}
}

func (r *oidcProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oidc_provider"
}

func (r *oidcProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"client_id": schema.StringAttribute{
				Required: true,
			},
			"client_secret": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"scope": schema.StringAttribute{
				Required: true,
			},
			"response_type": schema.StringAttribute{
				Required: true,
			},
			"authorization_endpoint": schema.StringAttribute{
				Required: true,
			},
			"token_endpoint": schema.StringAttribute{
				Optional: true,
			},
			"token_endpoint_auth_method": schema.StringAttribute{
				Required: true,
			},
			"redirect_uri": schema.StringAttribute{
				Required: true,
			},
			"pkce_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"jwks_uri": schema.StringAttribute{
				Optional: true,
			},
			"issuer": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *oidcProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *oidcProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oidcProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CreateOidcProvider(ctx, &bouncr.OidcProviderCreateRequest{
		Name:                    plan.Name.ValueString(),
		ClientID:                plan.ClientID.ValueString(),
		ClientSecret:            plan.ClientSecret.ValueString(),
		Scope:                   plan.Scope.ValueString(),
		ResponseType:            plan.ResponseType.ValueString(),
		AuthorizationEndpoint:   plan.AuthorizationEndpoint.ValueString(),
		TokenEndpoint:           plan.TokenEndpoint.ValueString(),
		TokenEndpointAuthMethod: plan.TokenEndpointAuthMethod.ValueString(),
		RedirectURI:             plan.RedirectURI.ValueString(),
		PkceEnabled:             plan.PkceEnabled.ValueBool(),
		JwksURI:                 plan.JwksURI.ValueString(),
		Issuer:                  plan.Issuer.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating OIDC provider", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oidcProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oidcProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oidcProvider, err := r.client.FindOidcProvider(ctx, state.Name.ValueString())
	if isNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading OIDC provider", err.Error())
		return
	}

	state.Name = types.StringValue(oidcProvider.Name)
	state.ClientID = types.StringValue(oidcProvider.ClientID)
	state.ClientSecret = types.StringValue(oidcProvider.ClientSecret)
	state.Scope = types.StringValue(oidcProvider.Scope)
	state.ResponseType = types.StringValue(oidcProvider.ResponseType)
	state.AuthorizationEndpoint = types.StringValue(oidcProvider.AuthorizationEndpoint)
	state.TokenEndpoint = types.StringValue(oidcProvider.TokenEndpoint)
	state.TokenEndpointAuthMethod = types.StringValue(oidcProvider.TokenEndpointAuthMethod)
	state.RedirectURI = types.StringValue(oidcProvider.RedirectURI)
	state.PkceEnabled = types.BoolValue(oidcProvider.PkceEnabled)
	state.JwksURI = types.StringValue(oidcProvider.JwksURI)
	state.Issuer = types.StringValue(oidcProvider.Issuer)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *oidcProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan oidcProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state oidcProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateOidcProvider(ctx, state.Name.ValueString(), &bouncr.OidcProviderUpdateRequest{
		Name:                    plan.Name.ValueString(),
		ClientID:                plan.ClientID.ValueString(),
		ClientSecret:            plan.ClientSecret.ValueString(),
		Scope:                   plan.Scope.ValueString(),
		ResponseType:            plan.ResponseType.ValueString(),
		AuthorizationEndpoint:   plan.AuthorizationEndpoint.ValueString(),
		TokenEndpoint:           plan.TokenEndpoint.ValueString(),
		TokenEndpointAuthMethod: plan.TokenEndpointAuthMethod.ValueString(),
		RedirectURI:             plan.RedirectURI.ValueString(),
		PkceEnabled:             plan.PkceEnabled.ValueBool(),
		JwksURI:                 plan.JwksURI.ValueString(),
		Issuer:                  plan.Issuer.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating OIDC provider", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oidcProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state oidcProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOidcProvider(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting OIDC provider", err.Error())
	}
}

func (r *oidcProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

