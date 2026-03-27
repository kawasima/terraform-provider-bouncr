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
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

type userResource struct {
	client *bouncr.Client
}

type userResourceModel struct {
	Account      types.String `tfsdk:"account"`
	Password     types.String `tfsdk:"password"`
	UserProfiles types.Map    `tfsdk:"user_profiles"`
}

func NewUserResource() resource.Resource {
	return &userResource{}
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"user_profiles": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &bouncr.UserCreateRequest{
		"account": plan.Account.ValueString(),
	}

	if !plan.UserProfiles.IsNull() {
		var profiles map[string]string
		resp.Diagnostics.Append(plan.UserProfiles.ElementsAs(ctx, &profiles, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for k, v := range profiles {
			(*input)[k] = v
		}
	}

	_, err := r.client.CreateUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	if !plan.Password.IsNull() && plan.Password.ValueString() != "" {
		_, err := r.client.CreatePasswordCredential(ctx, &bouncr.PasswordCredentialCreateRequest{
			Account:  plan.Account.ValueString(),
			Password: plan.Password.ValueString(),
			Initial:  false,
		})
		if err != nil {
			resp.Diagnostics.AddError("Error creating password credential", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.FindUser(ctx, state.Account.ValueString())
	if isNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	state.Account = types.StringValue(user.Account)
	if user.UserProfiles != nil {
		profiles := make(map[string]string)
		for k, v := range user.UserProfiles {
			if s, ok := v.(string); ok {
				profiles[k] = s
			}
		}
		mapVal, diags := types.MapValueFrom(ctx, types.StringType, profiles)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.UserProfiles = mapVal
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &bouncr.UserUpdateRequest{
		"account": plan.Account.ValueString(),
	}

	if !plan.UserProfiles.IsNull() {
		var profiles map[string]string
		resp.Diagnostics.Append(plan.UserProfiles.ElementsAs(ctx, &profiles, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for k, v := range profiles {
			(*input)[k] = v
		}
	}

	_, err := r.client.UpdateUser(ctx, state.Account.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, state.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("account"), req, resp)
}
