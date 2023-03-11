package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/projectroleservice"
	models2 "terraform-provider-hashicups-pf/services/projectroleservice/models"
)

func ProjectRoleResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Id()
			description := data.Get("description").(string)

			projectRoleService := projectroleservice.ProjectRoleService{
				JiraServerBase: client,
			}

			createdRole, err := projectRoleService.CreateRole(ctx, models2.ProjectRoleCreateRequestModel{
				Name:        name,
				Description: description,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", createdRole.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", createdRole.Description); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(createdRole.Name)
			log.Println("success create project role")
			return diags
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Id()

			projectRoleService := projectroleservice.ProjectRoleService{
				JiraServerBase: client,
			}

			projectRole, err := projectRoleService.GetRole(ctx, models2.ProjectRoleGetRequestModel{
				Name: name,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", projectRole.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", projectRole.Description); err != nil {
				return diag.FromErr(err)
			}

			log.Println("success get project role")
			return diags
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Id()
			description := data.Get("description").(string)

			projectRoleService := projectroleservice.ProjectRoleService{
				JiraServerBase: client,
			}

			updatedRole, err := projectRoleService.UpdateRole(ctx, models2.ProjectRoleUpdateRequestModel{
				Name:        name,
				Description: description,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", updatedRole.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", updatedRole.Description); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(updatedRole.Name)

			log.Println("success update project role")
			return diags
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Id()

			projectRoleService := projectroleservice.ProjectRoleService{
				JiraServerBase: client,
			}

			_, err := projectRoleService.DeleteRole(ctx, models2.ProjectRoleDeleteRequestModel{
				Name: name,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			data.SetId("")

			log.Println("success delete project role")
			return diags
		},
		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "description of project role",
			},
		},
	}
}
