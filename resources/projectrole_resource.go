package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/projectroleservice"
	models2 "terraform-provider-hashicups-pf/services/projectroleservice/models"
)

func ProjectRoleResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)
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

			if err = data.Set("project_role_id", int(createdRole.Id)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(createdRole.Id, 10))
			log.Println("success create project role")
			return diags
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			id := data.Get("project_role_id").(int)

			projectRoleService := projectroleservice.ProjectRoleService{
				JiraServerBase: client,
			}

			projectRole, err := projectRoleService.GetRole(ctx, models2.ProjectRoleGetRequestModel{
				Id: int64(id),
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

			if err = data.Set("project_role_id", int(projectRole.Id)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(projectRole.Id, 10))
			log.Println("success get project role")
			return diags
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			id := data.Get("project_role_id").(int)
			name := data.Get("name").(string)
			description := data.Get("description").(string)

			projectRoleService := projectroleservice.ProjectRoleService{
				JiraServerBase: client,
			}

			updatedRole, err := projectRoleService.UpdateRole(ctx, models2.ProjectRoleUpdateRequestModel{
				Id:          int64(id),
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

			if err = data.Set("project_role_id", int(updatedRole.Id)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(updatedRole.Id, 10))
			log.Println("success update project role")
			return diags
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			id := data.Get("project_role_id").(int)

			projectRoleService := projectroleservice.ProjectRoleService{
				JiraServerBase: client,
			}

			_, err := projectRoleService.DeleteRole(ctx, models2.ProjectRoleDeleteRequestModel{
				Id: int64(id),
			})
			if err != nil {
				return diag.FromErr(err)
			}
			data.SetId("")

			log.Println("success delete project role")
			return diags
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of project role",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "description of project role",
			},
			"project_role_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the project role id",
			},
		},
	}
}
