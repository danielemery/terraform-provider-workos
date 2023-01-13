package workos

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	workosOrganizations "github.com/workos/workos-go/pkg/organizations"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &organizationsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationsDataSource{}
)

// NewOrganizationsDataSource is a helper function to simplify the provider implementation.
func NewOrganizationsDataSource() datasource.DataSource {
	return &organizationsDataSource{}
}

// organizationsDataSource is the data source implementation.
type organizationsDataSource struct {
	client *workosOrganizations.Client
}

type organizationsDataSourceModel struct {
	ID            types.String        `tfsdk:"id"`
	Organizations []organizationModel `tfsdk:"organizations"`
}

type organizationModel struct {
	ID                               types.String   `tfsdk:"id"`
	Name                             types.String   `tfsdk:"name"`
	AllowProfilesOutsideOrganization types.Bool     `tfsdk:"allow_profiles_outside_organization"`
	Domains                          []types.String `tfsdk:"domains"`
	CreatedAt                        types.String   `tfsdk:"created_at"`
	UpdatedAt                        types.String   `tfsdk:"updated_at"`
}

func (d *organizationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*workosClient).Organizations
}

// Metadata returns the data source type name.
func (d *organizationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

// Schema defines the schema for the data source.
func (d *organizationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"organizations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                                  schema.StringAttribute{Computed: true},
						"name":                                schema.StringAttribute{Computed: true},
						"allow_profiles_outside_organization": schema.BoolAttribute{Computed: true},
						"domains": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"created_at": schema.StringAttribute{Computed: true},
						"updated_at": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *organizationsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organizationsDataSourceModel

	organizations, err := d.client.ListOrganizations(ctx, workosOrganizations.ListOrganizationsOpts{})
	if err != nil {
		resp.Diagnostics.AddError("Unable to read WorkOS Organizations", err.Error())
		return
	}

	for _, organization := range organizations.Data {
		organizationState := buildOrganizationState(organization)
		state.Organizations = append(state.Organizations, organizationState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func buildOrganizationState(organization workosOrganizations.Organization) organizationModel {
	organizationState := organizationModel{
		ID:                               types.StringValue(organization.ID),
		Name:                             types.StringValue(organization.Name),
		AllowProfilesOutsideOrganization: types.BoolValue(organization.AllowProfilesOutsideOrganization),
		CreatedAt:                        types.StringValue(organization.CreatedAt),
		UpdatedAt:                        types.StringValue(organization.UpdatedAt),
	}
	for _, domain := range organization.Domains {
		organizationState.Domains = append(organizationState.Domains, types.StringValue(domain.Domain))
	}
	return organizationState
}
