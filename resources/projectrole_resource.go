package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/projectroleservice"
	models2 "terraform-provider-hashicups-pf/services/projectroleservice/models"
)

var (
	_ resource.Resource              = &projectRoleResource{}
	_ resource.ResourceWithConfigure = &projectRoleResource{}
)

type projectRoleModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type projectRoleResource struct {
	client models.JiraServerBase
}

func (p projectRoleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	p.client = request.ProviderData.(models.JiraServerBase)
}

func (p projectRoleResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_projectrole"
}

func (p projectRoleResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p projectRoleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	projectRoleService := projectroleservice.ProjectRoleService{
		JiraServerBase: p.client,
	}

	var plan projectRoleModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	createdProjectRole, err := projectRoleService.CreateRole(models2.ProjectRoleCreateRequestModel{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating project role",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, &createdProjectRole)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (p projectRoleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	projectRoleService := projectroleservice.ProjectRoleService{
		JiraServerBase: p.client,
	}

	var state projectRoleModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	projectRole, err := projectRoleService.GetRole(models2.ProjectRoleGetRequestModel{
		Name: state.Name.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Project Role",
			"Could not read Project Role w. name "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(projectRole.Name)
	state.Description = types.StringValue(projectRole.Description)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (p projectRoleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (p projectRoleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func NewProjectRoleResource() resource.Resource {
	return &projectRoleResource{}
}
