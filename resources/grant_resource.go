package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/grantservice"
	models2 "terraform-provider-hashicups-pf/services/grantservice/models"
)

func GrantResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			permissionSchemeId := data.Get("permission_scheme_id").(int64)
			permissionName := data.Get("permission_name").(string)
			holderType := data.Get("security_type").(string)
			holderParam := data.Get("security_param").(int)

			grantService := grantservice.GrantService{
				JiraServerBase: client,
			}

			createdGrant, err := grantService.Create(ctx, models2.GrantCreateRequestModel{
				PermissionSchemeId: permissionSchemeId,
				Holder: models2.GrantApiHolderModel{
					Type:      holderType,
					Parameter: int64(holderParam),
				},
				Permission: permissionName,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_scheme_id", createdGrant.PermissionSchemeId); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_name", createdGrant.Permission); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("security_type", createdGrant.Holder.Type); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("security_param", createdGrant.Holder.Parameter); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(createdGrant.Id, 10))
			log.Println("success create grant")
			return diags
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			permissionSchemeId := data.Get("permission_scheme_id").(int)
			permissionName := data.Get("permission_name").(string)
			holderType := data.Get("security_type").(string)
			holderParam := data.Get("security_param").(int)

			grantService := grantservice.GrantService{
				JiraServerBase: client,
			}

			foundGrant, err := grantService.Get(ctx, models2.GrantGetRequestModel{
				Permission:         permissionName,
				PermissionSchemeId: int64(permissionSchemeId),
				Holder: models2.GrantHolderModel{
					Type:      holderType,
					Parameter: strconv.Itoa(holderParam),
				},
			})

			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_scheme_id", foundGrant.PermissionSchemeId); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_name", foundGrant.Permission); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("security_type", foundGrant.Holder.Type); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("security_param", foundGrant.Holder.Parameter); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(foundGrant.Id, 10))
			log.Println("success get grant")
			return diags
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			permissionSchemeId := data.Get("permission_scheme_id").(int)
			permissionName := data.Get("permission_name").(string)
			holderType := data.Get("security_type").(string)
			holderParam := data.Get("security_param").(string)

			grantService := grantservice.GrantService{
				JiraServerBase: client,
			}

			foundGrant, err := grantService.Get(ctx, models2.GrantGetRequestModel{
				Permission:         permissionName,
				PermissionSchemeId: int64(permissionSchemeId),
				Holder: models2.GrantHolderModel{
					Type:      holderType,
					Parameter: holderParam,
				},
			})

			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_scheme_name", int(foundGrant.PermissionSchemeId)); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("permission_name", foundGrant.Permission); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("security_type", foundGrant.Holder.Type); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("security_param", foundGrant.Holder.Parameter); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(strconv.FormatInt(foundGrant.Id, 10))
			log.Println("success get grant")
			return diags
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			permissionSchemeId := data.Get("permission_scheme_id").(int)
			permissionName := data.Get("permission_name").(string)
			holderType := data.Get("security_type").(string)
			holderParam := data.Get("security_param").(string)

			grantService := grantservice.GrantService{
				JiraServerBase: client,
			}

			_, err := grantService.Delete(ctx, models2.GrantDeleteRequestModel{
				Permission:         permissionName,
				PermissionSchemeId: int64(permissionSchemeId),
				Holder: models2.GrantHolderModel{
					Type:      holderType,
					Parameter: holderParam,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
			data.SetId("")
			log.Println("success delete grant")
			return diags
		},
		Schema: map[string]*schema.Schema{
			"permission_scheme_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of target permission scheme",
			},
			"permission_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of permission name to be granted",
			},
			"security_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of security type, e.g. projectrole",
			},
			"security_param": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "value of security type input",
			},
		},
	}
}
