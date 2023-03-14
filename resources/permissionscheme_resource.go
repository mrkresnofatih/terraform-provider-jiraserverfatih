package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/permissionschemeservice"
	models2 "terraform-provider-hashicups-pf/services/permissionschemeservice/models"
)

func PermissionSchemeResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)
			description := data.Get("description").(string)

			permissionSchemeService := permissionschemeservice.PermissionSchemeService{
				JiraServerBase: client,
			}

			createdPermSch, err := permissionSchemeService.Create(ctx, models2.PermissionSchemeCreateRequestModel{
				Name:        name,
				Description: description,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", createdPermSch.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", createdPermSch.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_scheme_id", int(createdPermSch.Id)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(createdPermSch.Id, 10))
			log.Println("success create permission scheme")
			return diags
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)
			description := data.Get("description").(string)

			permissionSchemeService := permissionschemeservice.PermissionSchemeService{
				JiraServerBase: client,
			}

			updatedPermSch, err := permissionSchemeService.Update(ctx, models2.PermissionSchemeUpdateRequestModel{
				Name:        name,
				Description: description,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", updatedPermSch.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", updatedPermSch.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_scheme_id", int(updatedPermSch.Id)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(updatedPermSch.Id, 10))
			log.Println("success update permission scheme")
			return diags
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			id := data.Get("permission_scheme_id").(int)

			permissionSchemeService := permissionschemeservice.PermissionSchemeService{
				JiraServerBase: client,
			}

			foundPermSch, err := permissionSchemeService.Get(ctx, models2.PermissionSchemeGetRequestModel{
				Id: int64(id),
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", foundPermSch.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", foundPermSch.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_scheme_id", int(foundPermSch.Id)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(foundPermSch.Id, 10))
			log.Println("success update permission scheme")
			return diags
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			id := data.Get("permission_scheme_id").(int)

			permissionSchemeService := permissionschemeservice.PermissionSchemeService{
				JiraServerBase: client,
			}
			_, err := permissionSchemeService.Delete(ctx, models2.PermissionSchemeDeleteRequestModel{
				Id: int64(id),
			})
			if err != nil {
				return diag.FromErr(err)
			}
			data.SetId("")
			log.Println("success delete permission scheme")
			return diags
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of permission scheme",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "description of permission scheme",
			},
			"permission_scheme_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of permission scheme",
			},
		},
	}
}
