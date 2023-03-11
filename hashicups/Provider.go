package hashicups

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/resources"
	"terraform-provider-hashicups-pf/services/baseservice/models"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &hashicupsProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &hashicupsProvider{}
}

// hashicupsProvider is the provider implementation.
type hashicupsProvider struct{}

// Metadata returns the provider type name.
func (p *hashicupsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hashicups"
}

// Schema defines the provider-level schema for configuration data.
func (p *hashicupsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    false,
				Required:    true,
				Description: "The domain name of your Jira Server e.g. jira.app-dev.company.com",
			},
			"authorizationMethod": schema.StringAttribute{
				Optional:    false,
				Required:    true,
				Description: "Valid values are: Basic | Bearer",
			},
			"token": schema.StringAttribute{
				Optional:    false,
				Required:    true,
				Description: "token value that follows the authorization method",
			},
		},
	}
}

type jiraServerProviderModel struct {
	Host                types.String `tfsdk:"host"`
	AuthorizationMethod types.String `tfsdk:"authorizationMethod"`
	Token               types.String `tfsdk:"token"`
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *hashicupsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config jiraServerProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsNull() || config.Host.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Host",
			"The provider cannot create an http request using an empty host",
		)
	}

	if config.Token.IsNull() || config.Token.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown token",
			"The provider cannot create an http request using an empty token",
		)
	}

	if config.AuthorizationMethod.IsNull() || (config.AuthorizationMethod.ValueString() != "Bearer" && config.AuthorizationMethod.ValueString() != "Basic") {
		resp.Diagnostics.AddAttributeError(
			path.Root("authorizationMethod"),
			"Unknown AuthorizationMethod",
			"The provider cannot create an http request using an unknown AuthorizationMethod host, e.g. Basic or Bearer",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := models.JiraServerBase{
		Token:               config.Token.ValueString(),
		AuthorizationMethod: config.AuthorizationMethod.ValueString(),
		Domain:              config.Host.ValueString(),
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *hashicupsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *hashicupsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewProjectRoleResource,
	}
}
