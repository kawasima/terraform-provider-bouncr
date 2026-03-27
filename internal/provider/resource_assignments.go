package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kawasima/bouncr-client-go"
	"github.com/rs/xid"
)

var _ resource.Resource = &assignmentsResource{}

type assignmentsResource struct {
	client *bouncr.Client
}

type assignmentsResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Assignments types.List   `tfsdk:"assignment"`
}

type assignmentModel struct {
	Group types.String `tfsdk:"group"`
	Role  types.String `tfsdk:"role"`
	Realm types.String `tfsdk:"realm"`
}

func NewAssignmentsResource() resource.Resource {
	return &assignmentsResource{}
}

func (r *assignmentsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assignments"
}

func (r *assignmentsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"assignment": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Required: true,
						},
						"role": schema.StringAttribute{
							Required: true,
						},
						"realm": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *assignmentsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *assignmentsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan assignmentsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []assignmentModel
	resp.Diagnostics.Append(plan.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignmentsRequest := make([]bouncr.AssignmentRequest, 0, len(assignments))
	for _, a := range assignments {
		assignmentsRequest = append(assignmentsRequest, bouncr.AssignmentRequest{
			Group: bouncr.IDObject{Name: a.Group.ValueString()},
			Role:  bouncr.IDObject{Name: a.Role.ValueString()},
			Realm: bouncr.IDObject{Name: a.Realm.ValueString()},
		})
	}

	_, err := r.client.CreateAssignments(ctx, &assignmentsRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating assignments", err.Error())
		return
	}

	plan.ID = types.StringValue(xid.New().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *assignmentsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state assignmentsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []assignmentModel
	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resultAssignments := make([]assignmentModel, 0, len(assignments))
	for _, a := range assignments {
		assignmentReq := &bouncr.AssignmentRequest{
			Group: bouncr.IDObject{Name: a.Group.ValueString()},
			Role:  bouncr.IDObject{Name: a.Role.ValueString()},
			Realm: bouncr.IDObject{Name: a.Realm.ValueString()},
		}
		assignment, err := r.client.FindAssignment(ctx, assignmentReq)
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		if err != nil {
			resp.Diagnostics.AddError("Error reading assignment", err.Error())
			return
		}
		resultAssignments = append(resultAssignments, assignmentModel{
			Group: types.StringValue(assignment.Group.Name),
			Role:  types.StringValue(assignment.Role.Name),
			Realm: types.StringValue(assignment.Realm.Name),
		})
	}

	listVal, diags := types.ListValueFrom(ctx, state.Assignments.ElementType(ctx), resultAssignments)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Assignments = listVal
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *assignmentsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan assignmentsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state assignmentsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *assignmentsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state assignmentsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var assignments []assignmentModel
	resp.Diagnostics.Append(state.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignmentsRequest := make([]bouncr.AssignmentRequest, 0, len(assignments))
	for _, a := range assignments {
		assignmentsRequest = append(assignmentsRequest, bouncr.AssignmentRequest{
			Group: bouncr.IDObject{Name: a.Group.ValueString()},
			Role:  bouncr.IDObject{Name: a.Role.ValueString()},
			Realm: bouncr.IDObject{Name: a.Realm.ValueString()},
		})
	}

	err := r.client.DeleteAssignments(ctx, &assignmentsRequest)
	if isNotFound(err) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error deleting assignments", err.Error())
	}
}
