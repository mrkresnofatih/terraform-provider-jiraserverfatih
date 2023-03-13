package grantservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	models2 "terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/grantservice/models"
	"terraform-provider-hashicups-pf/services/permissionschemeservice"
	models3 "terraform-provider-hashicups-pf/services/permissionschemeservice/models"
	"terraform-provider-hashicups-pf/services/projectroleservice"
	models4 "terraform-provider-hashicups-pf/services/projectroleservice/models"
	"time"
)

type IGrantService interface {
	Get(ctx context.Context, model models.GrantGetRequestModel) (models.GrantGetResponseModel, error)
	List(ctx context.Context, model models.GrantListRequestModel) (models.GrantListResponseModel, error)
	Create(ctx context.Context, model models.GrantCreateRequestModel) (models.GrantCreateResponseModel, error)
	Delete(ctx context.Context, model models.GrantDeleteRequestModel) (models.GrantDeleteResponseModel, error)
}

type GrantService struct {
	JiraServerBase models2.JiraServerBase `json:"jiraServerBase"`
}

func (g GrantService) Get(ctx context.Context, model models.GrantGetRequestModel) (models.GrantGetResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start list permission scheme grants w. data: %s", model))

	grantsResult, err := g.List(ctx, models.GrantListRequestModel{
		PermissionSchemeName: model.PermissionSchemeName,
	})
	if err != nil {
		tflog.Info(ctx, "failed to find permission scheme grant")
		return *new(models.GrantGetResponseModel), errors.New("failed to find permission scheme")
	}

	projectRoleService := projectroleservice.ProjectRoleService{
		JiraServerBase: g.JiraServerBase,
	}

	projectRoleFound, err := projectRoleService.GetRole(ctx, models4.ProjectRoleGetRequestModel{
		Name: model.Holder.Parameter,
	})
	if err != nil {
		tflog.Info(ctx, "failed to find project role with role name "+model.Holder.Parameter)
		return *new(models.GrantGetResponseModel), errors.New("failed to find project role with role name " + model.Holder.Parameter)
	}

	foundGrant := models.GrantGetResponseModel{}
	for _, grant := range grantsResult.Grants {
		if grant.Permission == model.Permission && strings.ToLower(grant.Holder.Type) == strings.ToLower(model.Holder.Type) && grant.Holder.Parameter == strconv.FormatInt(projectRoleFound.Id, 10) {
			log.Println(grant)
			foundGrant = grant
		}
	}
	if foundGrant.Permission == "" {
		tflog.Info(ctx, "failed to get permission scheme grant")
		return *new(models.GrantGetResponseModel), errors.New("failed to find permission scheme grant in returned grants list")
	}

	foundGrant.Holder.Parameter = model.Holder.Parameter
	tflog.Info(ctx, "success find permission scheme grant")
	return foundGrant, nil
}

func (g GrantService) List(ctx context.Context, model models.GrantListRequestModel) (models.GrantListResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start list permission scheme grants w. data: %s", model))

	permissionSchemeService := permissionschemeservice.PermissionSchemeService{
		JiraServerBase: g.JiraServerBase,
	}

	permissionSchemeFound, err := permissionSchemeService.Get(ctx, models3.PermissionSchemeGetRequestModel{
		Name: model.PermissionSchemeName,
	})
	if err != nil {
		tflog.Info(ctx, "failed to find permission scheme w. name "+permissionSchemeFound.Name)
		return *new(models.GrantListResponseModel), errors.New("failed to find permission scheme")
	}

	url := "https://" + g.JiraServerBase.Domain + "/rest/api/2/permissionscheme/" + strconv.FormatInt(permissionSchemeFound.Id, 10) + "/permission"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("error building http request")
		return *new(models.GrantListResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", g.JiraServerBase.AuthorizationMethod+" "+g.JiraServerBase.Token)
	req.Header.Set("Accept", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.GrantListResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read all body failed")
		return *new(models.GrantListResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))

	result := models.GrantListResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("failed to unmarshal response body")
		return *new(models.GrantListResponseModel), errors.New("error unmarshalling response body")
	}

	log.Println("success list permission scheme grants")
	return result, nil
}

func (g GrantService) Create(ctx context.Context, model models.GrantCreateRequestModel) (models.GrantCreateResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start create permission scheme grant w. data: %s", model))

	permissionSchemeService := permissionschemeservice.PermissionSchemeService{
		JiraServerBase: g.JiraServerBase,
	}

	projectRoleService := projectroleservice.ProjectRoleService{
		JiraServerBase: g.JiraServerBase,
	}

	projectRoleFound, err := projectRoleService.GetRole(ctx, models4.ProjectRoleGetRequestModel{
		Name: model.Holder.Parameter,
	})
	if err != nil {
		tflog.Info(ctx, "failed to find project role with role name "+model.Holder.Parameter)
		return *new(models.GrantCreateResponseModel), errors.New("failed to find project role")
	}

	permissionSchemeFound, err := permissionSchemeService.Get(ctx, models3.PermissionSchemeGetRequestModel{
		Name: model.PermissionSchemeName,
	})
	if err != nil {
		tflog.Info(ctx, "failed to find permission scheme w. name "+permissionSchemeFound.Name)
		return *new(models.GrantCreateResponseModel), errors.New("failed to find permission scheme")
	}

	url := "https://" + g.JiraServerBase.Domain + "/rest/api/2/permissionscheme/" + strconv.FormatInt(permissionSchemeFound.Id, 10) + "/permission"
	serial, err := json.Marshal(models.GrantCreateApiRequestModel{
		Permission: model.Permission,
		Holder: models.GrantApiHolderModel(struct {
			Type      string
			Parameter int64
		}{Type: model.Holder.Type, Parameter: projectRoleFound.Id}),
	})
	if err != nil {
		tflog.Info(ctx, "failed to marshal body request")
		return *new(models.GrantCreateResponseModel), err
	}
	bodyReader := bytes.NewReader(serial)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		return *new(models.GrantCreateResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", g.JiraServerBase.AuthorizationMethod+" "+g.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.GrantCreateResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read all body failed")
		return *new(models.GrantCreateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))

	result := models.GrantCreateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("failed to unmarshal response body")
		return *new(models.GrantCreateResponseModel), errors.New("error unmarshalling response body")
	}

	result.PermissionSchemeName = model.PermissionSchemeName
	result.Holder.Parameter = model.Holder.Parameter
	log.Println("success create permission scheme grant")
	return result, nil
}

func (g GrantService) Delete(ctx context.Context, model models.GrantDeleteRequestModel) (models.GrantDeleteResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start list permission scheme grants w. data: %s", model))

	permissionSchemeService := permissionschemeservice.PermissionSchemeService{
		JiraServerBase: g.JiraServerBase,
	}

	permissionSchemeFound, err := permissionSchemeService.Get(ctx, models3.PermissionSchemeGetRequestModel{
		Name: model.PermissionSchemeName,
	})
	if err != nil {
		tflog.Info(ctx, "failed to find permission scheme w. name "+permissionSchemeFound.Name)
		return *new(models.GrantDeleteResponseModel), errors.New("failed to find permission scheme")
	}

	foundGrant, err := g.Get(ctx, models.GrantGetRequestModel{
		Permission:           model.Permission,
		Holder:               model.Holder,
		PermissionSchemeName: model.PermissionSchemeName,
	})
	if err != nil {
		log.Println("failed to find grant")
		return *new(models.GrantDeleteResponseModel), errors.New("failed to find perm. scheme grant")
	}

	url := "https://" + g.JiraServerBase.Domain + "/rest/api/2/permissionscheme/" + strconv.FormatInt(permissionSchemeFound.Id, 10) + "/permission/" + strconv.FormatInt(foundGrant.Id, 10)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Println("error building http request")
		return *new(models.GrantDeleteResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", g.JiraServerBase.AuthorizationMethod+" "+g.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.GrantDeleteResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	log.Println("delete permission scheme grant success")
	return models.GrantDeleteResponseModel{}, nil
}
