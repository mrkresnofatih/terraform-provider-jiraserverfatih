package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/issuetypeservice"
	models2 "terraform-provider-hashicups-pf/services/issuetypeservice/models"
)

func IssueTypeResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)
			description := data.Get("description").(string)
			avatar_id := data.Get("avatar_id").(int)

			issueTypeService := issuetypeservice.IssueTypeService{
				JiraServerBase: client,
			}

			createdIssueType, err := issueTypeService.Create(ctx, models2.IssueTypeCreateRequestModel{
				Name:        name,
				Description: description,
				AvatarId:    int64(avatar_id),
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", createdIssueType.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", createdIssueType.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("avatar_id", avatar_id); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(createdIssueType.Name)
			log.Println("success create issue type")
			return diags
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)
			issueTypeService := issuetypeservice.IssueTypeService{
				JiraServerBase: client,
			}

			foundIssueType, err := issueTypeService.Get(ctx, models2.IssueTypeGetRequestModel{
				Name: name,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", foundIssueType.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", foundIssueType.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("avatar_id", int(foundIssueType.AvatarId)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(foundIssueType.Name)
			log.Println("success get issue type")

			return diags
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)
			description := data.Get("description").(string)
			avatar_id := data.Get("avatar_id").(int)

			issueTypeService := issuetypeservice.IssueTypeService{
				JiraServerBase: client,
			}

			updatedIssueType, err := issueTypeService.Update(ctx, models2.IssueTypeUpdateRequestModel{
				Name:        name,
				Description: description,
				AvatarId:    int64(avatar_id),
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("name", updatedIssueType.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("description", updatedIssueType.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = data.Set("avatar_id", int(updatedIssueType.AvatarId)); err != nil {
				return diag.FromErr(err)
			}

			data.SetId(updatedIssueType.Name)
			log.Println("success update issue type")

			return diags
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			var diags diag.Diagnostics
			client := i.(models.JiraServerBase)

			name := data.Get("name").(string)
			issueTypeService := issuetypeservice.IssueTypeService{
				JiraServerBase: client,
			}

			_, err := issueTypeService.Delete(ctx, models2.IssueTypeDeleteRequestModel{
				Name: name,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			data.SetId("")
			log.Println("success delete issue type")

			return diags
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of issue type",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "description of issue type",
			},
			"avatar_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "avatar id of issue type",
			},
		},
	}
}
