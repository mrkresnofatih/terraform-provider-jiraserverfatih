package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/groupservice"
	models2 "terraform-provider-hashicups-pf/services/groupservice/models"
)

func GroupResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)

			groupService := groupservice.GroupService{
				JiraServerBase: client,
			}

			createdGroup, err := groupService.Create(ctx, models2.GroupCreateRequestModel{
				Name: name,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", createdGroup.Name); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(createdGroup.Name)
			log.Println("success create group")
			return diags
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			name := data.Get("name").(string)

			if err := data.Set("name", name); err != nil {
				return diag.FromErr(err)
			}
			data.SetId(name)
			return diags
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)

			groupService := groupservice.GroupService{
				JiraServerBase: client,
			}

			_, err := groupService.Delete(ctx, models2.GroupDeleteRequestModel{
				Name: name,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			data.SetId("")
			log.Println("success delete group")
			return diags
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)

			groupService := groupservice.GroupService{
				JiraServerBase: client,
			}

			_, err := groupService.Get(ctx, models2.GroupGetRequestModel{
				Name: name,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			return diags
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of user group",
			},
		},
	}
}
