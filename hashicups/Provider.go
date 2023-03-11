package hashicups

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-hashicups-pf/resources"
	"terraform-provider-hashicups-pf/services/baseservice/models"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"jiraserverfatih_projectrole_resource": resources.ProjectRoleResource(),
		},
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The domain for your jira server",
				Required:    true,
			},
			"authorization_method": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The authorization method in the request header, valid values: Bearer or Basic",
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The token for the authorization header",
			},
		},
		ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			var diags diag.Diagnostics

			domain := data.Get("domain").(string)
			authMethod := data.Get("authorization_method").(string)
			token := data.Get("token").(string)
			if domain == "" || authMethod == "" || token == "" {
				return nil, diag.FromErr(errors.New("domain or auth_method or token is empty"))
			}
			return models.JiraServerBase{
				Domain:              domain,
				Token:               token,
				AuthorizationMethod: authMethod,
			}, diags
		},
	}
}
