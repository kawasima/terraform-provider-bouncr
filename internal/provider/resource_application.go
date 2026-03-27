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
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

type applicationResource struct {
	client *bouncr.Client
}

type applicationResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	PassTo      types.String `tfsdk:"pass_to"`
	VirtualPath types.String `tfsdk:"virtual_path"`
	TopPage     types.String `tfsdk:"top_page"`
	Realms      types.List   `tfsdk:"realm"`
}

type realmModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	URL         types.String `tfsdk:"url"`
}

func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"pass_to": schema.StringAttribute{
				Required: true,
			},
			"virtual_path": schema.StringAttribute{
				Required: true,
			},
			"top_page": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"realm": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"description": schema.StringAttribute{
							Required: true,
						},
						"url": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CreateApplication(ctx, &bouncr.ApplicationCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		PassTo:      plan.PassTo.ValueString(),
		VirtualPath: plan.VirtualPath.ValueString(),
		TopPage:     plan.TopPage.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating application", err.Error())
		return
	}

	if !plan.Realms.IsNull() && len(plan.Realms.Elements()) > 0 {
		var realms []realmModel
		resp.Diagnostics.Append(plan.Realms.ElementsAs(ctx, &realms, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, realm := range realms {
			_, err := r.client.CreateRealm(ctx, plan.Name.ValueString(), &bouncr.RealmCreateRequest{
				Name:        realm.Name.ValueString(),
				Description: realm.Description.ValueString(),
				URL:         realm.URL.ValueString(),
			})
			if err != nil {
				resp.Diagnostics.AddError("Error creating realm", err.Error())
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state applicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := r.client.FindApplication(ctx, state.Name.ValueString())
	if isNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading application", err.Error())
		return
	}

	state.Name = types.StringValue(application.Name)
	state.Description = types.StringValue(application.Description)
	state.PassTo = types.StringValue(application.PassTo)
	state.VirtualPath = types.StringValue(application.VirtualPath)
	state.TopPage = types.StringValue(application.TopPage)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan applicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state applicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateApplication(ctx, state.Name.ValueString(), &bouncr.ApplicationUpdateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		PassTo:      plan.PassTo.ValueString(),
		VirtualPath: plan.VirtualPath.ValueString(),
		TopPage:     plan.TopPage.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", err.Error())
		return
	}

	var oldRealms, newRealms []realmModel
	if !state.Realms.IsNull() {
		resp.Diagnostics.Append(state.Realms.ElementsAs(ctx, &oldRealms, false)...)
	}
	if !plan.Realms.IsNull() {
		resp.Diagnostics.Append(plan.Realms.ElementsAs(ctx, &newRealms, false)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	oldRealmNames := make(map[string]bool)
	for _, r := range oldRealms {
		oldRealmNames[r.Name.ValueString()] = true
	}
	newRealmNames := make(map[string]bool)
	for _, r := range newRealms {
		newRealmNames[r.Name.ValueString()] = true
	}

	appName := plan.Name.ValueString()
	for _, realm := range newRealms {
		if !oldRealmNames[realm.Name.ValueString()] {
			_, err := r.client.CreateRealm(ctx, appName, &bouncr.RealmCreateRequest{
				Name:        realm.Name.ValueString(),
				Description: realm.Description.ValueString(),
				URL:         realm.URL.ValueString(),
			})
			if err != nil {
				resp.Diagnostics.AddError("Error creating realm", err.Error())
				return
			}
		}
	}
	for _, realm := range oldRealms {
		if !newRealmNames[realm.Name.ValueString()] {
			err := r.client.DeleteRealm(ctx, appName, realm.Name.ValueString())
			if err != nil && !isNotFound(err) {
				resp.Diagnostics.AddError("Error deleting realm", err.Error())
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state applicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApplication(ctx, state.Name.ValueString())
	if isNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error deleting application", err.Error())
	}
}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
