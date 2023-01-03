package workos

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/pkg/organizations"
)

var (
	_ resource.Resource              = &organizationResource{}
	_ resource.ResourceWithConfigure = &organizationResource{}
)

func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

type organizationResource struct {
	client *organizations.Client
}

func (r *organizationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*workosClient).Organizations
}

func (r *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{Required: true},
			"allow_profiles_outside_organization": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"domains": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organizationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var domains []string
	for _, domain := range plan.Domains {
		domains = append(domains, domain.ValueString())
	}
	allowProfilesOutsideOrganization := false
	if !plan.AllowProfilesOutsideOrganization.IsNull() {
		allowProfilesOutsideOrganization = plan.AllowProfilesOutsideOrganization.ValueBool()
	}
	organization, err := r.client.CreateOrganization(ctx, organizations.CreateOrganizationOpts{
		Name:                             plan.Name.ValueString(),
		AllowProfilesOutsideOrganization: allowProfilesOutsideOrganization,
		Domains:                          domains,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating organization",
			"Could not create organization, unexpected error: "+err.Error(),
		)
		return
	}

	state := buildOrganizationState(organization)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organizationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organization, err := r.client.GetOrganization(ctx, organizations.GetOrganizationOpts{
		Organization: state.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization",
			"Could not read organization, unexpected error: "+err.Error(),
		)
		return
	}

	state = buildOrganizationState(organization)

}

func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organizationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var domains []string
	for _, domain := range plan.Domains {
		domains = append(domains, domain.ValueString())
	}
	allowProfilesOutsideOrganization := false
	if !plan.AllowProfilesOutsideOrganization.IsNull() {
		allowProfilesOutsideOrganization = plan.AllowProfilesOutsideOrganization.ValueBool()
	}
	organization, err := r.client.UpdateOrganization(ctx, organizations.UpdateOrganizationOpts{
		Organization:                     plan.ID.ValueString(),
		Name:                             plan.Name.ValueString(),
		AllowProfilesOutsideOrganization: allowProfilesOutsideOrganization,
		Domains:                          domains,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating organization",
			"Could not update organization, unexpected error: "+err.Error(),
		)
		return
	}

	state := buildOrganizationState(organization)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organizationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrganization(ctx, organizations.DeleteOrganizationOpts{
		Organization: state.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting organization",
			"Could not delete organization, unexpected error: "+err.Error(),
		)
		return
	}
}
