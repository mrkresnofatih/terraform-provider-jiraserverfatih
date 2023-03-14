package projectroleservice

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
	models2 "terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/projectroleservice/models"
	"time"
)

type IProjectRoleService interface {
	GetRole(ctx context.Context, model models.ProjectRoleGetRequestModel) (models.ProjectRoleGetResponseModel, error)
	ListRoles(ctx context.Context, model models.ProjectRoleListRequestModel) (models.ProjectRoleListResponseModel, error)
	UpdateRole(ctx context.Context, model models.ProjectRoleUpdateRequestModel) (models.ProjectRoleUpdateResponseModel, error)
	CreateRole(ctx context.Context, model models.ProjectRoleCreateRequestModel) (models.ProjectRoleCreateResponseModel, error)
	DeleteRole(ctx context.Context, model models.ProjectRoleDeleteRequestModel) (models.ProjectRoleDeleteResponseModel, error)
}

type ProjectRoleService struct {
	JiraServerBase models2.JiraServerBase `json:"jiraServerBase"`
}

func (p ProjectRoleService) UpdateRole(ctx context.Context, model models.ProjectRoleUpdateRequestModel) (models.ProjectRoleUpdateResponseModel, error) {
	log.Printf("start update role w. data: %s", model)
	role, err := p.GetRole(ctx, models.ProjectRoleGetRequestModel{
		Id: model.Id,
	})
	if err != nil {
		log.Println("error get role not found")
		return *new(models.ProjectRoleUpdateResponseModel), errors.New("role for update not found")
	}

	serialized, err := json.Marshal(models.ProjectRoleCreateApiRequestModel{
		Name:        model.Name,
		Description: model.Description,
	})
	if err != nil {
		log.Println("error when json marshaling")
		return *new(models.ProjectRoleUpdateResponseModel), errors.New("error json marshal request data")
	}
	bodyReader := bytes.NewReader(serialized)

	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/role/" + strconv.FormatInt(role.Id, 10)
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		return *new(models.ProjectRoleUpdateResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", p.JiraServerBase.AuthorizationMethod+" "+p.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.ProjectRoleUpdateResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read all body failed")
		return *new(models.ProjectRoleUpdateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))

	result := models.ProjectRoleUpdateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("failed to unmarshal response body")
		return *new(models.ProjectRoleUpdateResponseModel), errors.New("error unmarshalling response body")
	}

	log.Println("success update role")
	return result, nil
}

func (p ProjectRoleService) GetRole(ctx context.Context, model models.ProjectRoleGetRequestModel) (models.ProjectRoleGetResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start get role w. data: %s", model))
	roles, err := p.ListRoles(ctx, models.ProjectRoleListRequestModel{})
	if err != nil {
		log.Println("failed to list roles")
		return *new(models.ProjectRoleGetResponseModel), errors.New("failed to list roles")
	}

	foundRole := models.ProjectRoleGetResponseModel{}
	for _, role := range roles {
		if role.Id == model.Id {
			foundRole = role
		}
	}
	if foundRole.Name == "" && foundRole.Description == "" {
		log.Println("role not found")
		return *new(models.ProjectRoleGetResponseModel), errors.New("failed to get role")
	}

	log.Println("success get role")
	return foundRole, nil
}

func (p ProjectRoleService) ListRoles(ctx context.Context, model models.ProjectRoleListRequestModel) (models.ProjectRoleListResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start list roles w. data: %s", model))
	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/role"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.ProjectRoleListResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", p.JiraServerBase.AuthorizationMethod+" "+p.JiraServerBase.Token)
	req.Header.Set("Accept", "application/json")

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, err.Error())
		tflog.Info(ctx, "http request failed")
		return *new(models.ProjectRoleListResponseModel), errors.New("http request returned error")
	}

	tflog.Info(ctx, res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.ProjectRoleListResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "response body: "+string(body))

	result := models.ProjectRoleListResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.ProjectRoleListResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success project list roles")
	return result, nil
}

func (p ProjectRoleService) CreateRole(ctx context.Context, model models.ProjectRoleCreateRequestModel) (models.ProjectRoleCreateResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("Start CreateRole w. data: %s", model))

	serialized, err := json.Marshal(model)
	if err != nil {
		tflog.Info(ctx, "error when json marshaling")
		return *new(models.ProjectRoleCreateResponseModel), errors.New("error json marshal request data")
	}
	bodyReader := bytes.NewReader(serialized)
	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/role"
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.ProjectRoleCreateResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", p.JiraServerBase.AuthorizationMethod+" "+p.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, "error result from http request")
		return *new(models.ProjectRoleCreateResponseModel), errors.New("error result from http request " + err.Error())
	}
	defer res.Body.Close()

	tflog.Info(ctx, res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.ProjectRoleCreateResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "response body: "+string(body))

	result := models.ProjectRoleCreateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.ProjectRoleCreateResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success create project role")
	return result, nil
}

func (p ProjectRoleService) DeleteRole(ctx context.Context, model models.ProjectRoleDeleteRequestModel) (models.ProjectRoleDeleteResponseModel, error) {
	log.Printf("start delete role w. data: %s", model)
	role, err := p.GetRole(ctx, models.ProjectRoleGetRequestModel{
		Id: model.Id,
	})
	if err != nil {
		log.Println("error get role not found")
		return *new(models.ProjectRoleDeleteResponseModel), errors.New("role for delete not found")
	}

	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/role/" + strconv.FormatInt(role.Id, 10)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Println("error building http request")
		return *new(models.ProjectRoleDeleteResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", p.JiraServerBase.AuthorizationMethod+" "+p.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.ProjectRoleDeleteResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	log.Println("delete role success")
	return models.ProjectRoleDeleteResponseModel{}, nil
}
