package permissionschemeservice

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
	"terraform-provider-hashicups-pf/services/permissionschemeservice/models"
	"time"
)

type IPermissionSchemeService interface {
	Get(ctx context.Context, model models.PermissionSchemeGetRequestModel) (models.PermissionSchemeGetResponseModel, error)
	List(ctx context.Context, model models.PermissionSchemeListRequestModel) (models.PermissionSchemeListResponseModel, error)
	Create(ctx context.Context, model models.PermissionSchemeCreateRequestModel) (models.PermissionSchemeCreateResponseModel, error)
	Update(ctx context.Context, model models.PermissionSchemeUpdateRequestModel) (models.PermissionSchemeUpdateResponseModel, error)
	Delete(ctx context.Context, model models.PermissionSchemeDeleteRequestModel) (models.PermissionSchemeDeleteResponseModel, error)
}

type PermissionSchemeService struct {
	JiraServerBase models2.JiraServerBase `json:"jiraServerBase"`
}

func (p PermissionSchemeService) Get(ctx context.Context, model models.PermissionSchemeGetRequestModel) (models.PermissionSchemeGetResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start get permission scheme w. data: %s", model))
	permissionSchemes, err := p.List(ctx, models.PermissionSchemeListRequestModel{})
	if err != nil {
		log.Println("failed to list permission schemes")
		return *new(models.PermissionSchemeGetResponseModel), errors.New("failed to list roles")
	}

	foundPermissionScheme := models.PermissionSchemeGetResponseModel{}
	for _, ps := range permissionSchemes.PermissionSchemes {
		if ps.Id == model.Id {
			foundPermissionScheme = ps
		}
	}
	if foundPermissionScheme.Name == "" {
		tflog.Info(ctx, "permission scheme not found")
		return *new(models.PermissionSchemeGetResponseModel), errors.New("failed to list roles")
	}

	tflog.Info(ctx, "permission scheme found")
	return foundPermissionScheme, nil
}

func (p PermissionSchemeService) List(ctx context.Context, model models.PermissionSchemeListRequestModel) (models.PermissionSchemeListResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start list permission schemes w. data: %s", model))

	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/permissionscheme"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.PermissionSchemeListResponseModel), errors.New("error building http request")
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
		return *new(models.PermissionSchemeListResponseModel), errors.New("http request returned error")
	}

	tflog.Info(ctx, res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.PermissionSchemeListResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "response body: "+string(body))

	result := models.PermissionSchemeListResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.PermissionSchemeListResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success list permission schemes")
	return result, nil
}

func (p PermissionSchemeService) Create(ctx context.Context, model models.PermissionSchemeCreateRequestModel) (models.PermissionSchemeCreateResponseModel, error) {
	log.Printf("start create permission scheme w. data: %s", model)

	serialized, err := json.Marshal(model)
	if err != nil {
		tflog.Info(ctx, "error when json marshaling")
		return *new(models.PermissionSchemeCreateResponseModel), errors.New("error json marshal request data")
	}
	bodyReader := bytes.NewReader(serialized)
	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/permissionscheme"
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.PermissionSchemeCreateResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", p.JiraServerBase.AuthorizationMethod+" "+p.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, "error result from http request")
		return *new(models.PermissionSchemeCreateResponseModel), errors.New("error result from http request " + err.Error())
	}
	defer res.Body.Close()

	tflog.Info(ctx, res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.PermissionSchemeCreateResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "response body: "+string(body))

	result := models.PermissionSchemeCreateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.PermissionSchemeCreateResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success create permission scheme")
	return result, nil
}

func (p PermissionSchemeService) Update(ctx context.Context, model models.PermissionSchemeUpdateRequestModel) (models.PermissionSchemeUpdateResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start update permission scheme w. data: %s", model))

	permissionScheme, err := p.Get(ctx, models.PermissionSchemeGetRequestModel{
		Id: model.Id,
	})
	if err != nil {
		tflog.Info(ctx, "failed to find permission scheme named "+model.Name)
		return *new(models.PermissionSchemeUpdateResponseModel), err
	}

	serial, err := json.Marshal(models.PermissionSchemeUpdateRequestModel{
		Name:        permissionScheme.Name,
		Description: model.Description,
	})
	if err != nil {
		tflog.Info(ctx, "failed to marshal permission scheme")
		return *new(models.PermissionSchemeUpdateResponseModel), err
	}
	bodyReader := bytes.NewReader(serial)

	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/permissionscheme/" + strconv.FormatInt(permissionScheme.Id, 10)
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		return *new(models.PermissionSchemeUpdateResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", p.JiraServerBase.AuthorizationMethod+" "+p.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.PermissionSchemeUpdateResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read all body failed")
		return *new(models.PermissionSchemeUpdateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))

	result := models.PermissionSchemeUpdateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("failed to unmarshal response body")
		return *new(models.PermissionSchemeUpdateResponseModel), errors.New("error unmarshalling response body")
	}

	log.Println("success update permission scheme")
	return result, nil
}

func (p PermissionSchemeService) Delete(ctx context.Context, model models.PermissionSchemeDeleteRequestModel) (models.PermissionSchemeDeleteResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start delete permission scheme w. data: %s", model))

	permissionScheme, err := p.Get(ctx, models.PermissionSchemeGetRequestModel{
		Id: model.Id,
	})
	url := "https://" + p.JiraServerBase.Domain + "/rest/api/2/permissionscheme/" + strconv.FormatInt(permissionScheme.Id, 10)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Println("error building http request")
		return *new(models.PermissionSchemeDeleteResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", p.JiraServerBase.AuthorizationMethod+" "+p.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.PermissionSchemeDeleteResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	log.Println("delete permission scheme success")
	return models.PermissionSchemeDeleteResponseModel{}, nil
}
