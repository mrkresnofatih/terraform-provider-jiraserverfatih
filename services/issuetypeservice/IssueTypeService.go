package issuetypeservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"log"
	"net/http"
	models2 "terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/issuetypeservice/models"
	"time"
)

type IIssueTypeService interface {
	List(ctx context.Context, model models.IssueTypeListRequestModel) (models.IssueTypeListResponseModel, error)
	Get(ctx context.Context, model models.IssueTypeGetRequestModel) (models.IssueTypeGetResponseModel, error)
	Create(ctx context.Context, model models.IssueTypeCreateRequestModel) (models.IssueTypeCreateResponseModel, error)
	Update(ctx context.Context, model models.IssueTypeUpdateRequestModel) (models.IssueTypeUpdateResponseModel, error)
	Delete(ctx context.Context, model models.IssueTypeDeleteRequestModel) (models.IssueTypeDeleteResponseModel, error)
}

type IssueTypeService struct {
	JiraServerBase models2.JiraServerBase `json:"jiraServerBase"`
}

func (i IssueTypeService) List(ctx context.Context, model models.IssueTypeListRequestModel) (models.IssueTypeListResponseModel, error) {
	log.Printf("start list issue types w. data: %s", model)

	url := "https://" + i.JiraServerBase.Domain + "/rest/api/2/issuetype"
	log.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.IssueTypeListResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", i.JiraServerBase.AuthorizationMethod+" "+i.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, err.Error())
		tflog.Info(ctx, "http request failed")
		return *new(models.IssueTypeListResponseModel), errors.New("http request returned error")
	}

	tflog.Info(ctx, res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.IssueTypeListResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "response body: "+string(body))
	log.Println(string(body))

	result := models.IssueTypeListResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.IssueTypeListResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success group list")
	return result, nil
}

func (i IssueTypeService) Get(ctx context.Context, model models.IssueTypeGetRequestModel) (models.IssueTypeGetResponseModel, error) {
	log.Printf("start get issue type w. data: %s", model)

	issueTypes, err := i.List(ctx, models.IssueTypeListRequestModel{})
	if err != nil {
		log.Println("failed to list issue types")
		return *new(models.IssueTypeGetResponseModel), err
	}

	foundIssueType := models.IssueTypeGetResponseModel{}
	for _, it := range issueTypes {
		if it.Name == model.Name {
			foundIssueType = it
			break
		}
	}

	if foundIssueType.Id == "" {
		log.Println("issue type not found in issue type list")
		return *new(models.IssueTypeGetResponseModel), err
	}

	log.Println("success get issue type")
	return foundIssueType, nil
}

func (i IssueTypeService) Create(ctx context.Context, model models.IssueTypeCreateRequestModel) (models.IssueTypeCreateResponseModel, error) {
	log.Printf("start create issue type w. data: %s", model)

	model.Type = "standard"
	serialized, err := json.Marshal(model)
	if err != nil {
		log.Println("failed to marshal create request body")
		return *new(models.IssueTypeCreateResponseModel), errors.New("error json marshal req body")
	}
	bodyReader := bytes.NewReader(serialized)
	url := "https://" + i.JiraServerBase.Domain + "/rest/api/2/issuetype"
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		tflog.Info(ctx, "error building http request")
		return *new(models.IssueTypeCreateResponseModel), err
	}

	req.Header.Set("Authorization", i.JiraServerBase.AuthorizationMethod+" "+i.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("fail: response not ok: " + err.Error())
		return *new(models.IssueTypeCreateResponseModel), err
	}
	defer res.Body.Close()

	tflog.Info(ctx, res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.IssueTypeCreateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))
	tflog.Info(ctx, "response body: "+string(body))

	result := models.IssueTypeCreateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.IssueTypeCreateResponseModel), errors.New("error unmarshalling response body")
	}

	updatedIssueType, err := i.Update(ctx, models.IssueTypeUpdateRequestModel{
		Name:        model.Name,
		Description: model.Description,
		AvatarId:    model.AvatarId,
	})
	if err != nil {
		log.Println("failed to set issue type avatar")
		return *new(models.IssueTypeCreateResponseModel), errors.New("error unmarshalling response body")
	}

	log.Println(updatedIssueType)

	tflog.Info(ctx, "success create issue type")
	return result, nil
}

func (i IssueTypeService) Update(ctx context.Context, model models.IssueTypeUpdateRequestModel) (models.IssueTypeUpdateResponseModel, error) {
	log.Printf("start update issue type w. data: %s", model)

	foundIssueType, err := i.Get(ctx, models.IssueTypeGetRequestModel{
		Name: model.Name,
	})
	if err != nil {
		log.Println("fail to get issue type for update")
		return *new(models.IssueTypeUpdateResponseModel), err
	}

	serialized, err := json.Marshal(model)
	if err != nil {
		log.Println("failed to marshal create request body")
		return *new(models.IssueTypeUpdateResponseModel), errors.New("error json marshal req body")
	}
	bodyReader := bytes.NewReader(serialized)
	url := "https://" + i.JiraServerBase.Domain + "/rest/api/2/issuetype/" + foundIssueType.Id
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		tflog.Info(ctx, "error building http request")
		return *new(models.IssueTypeUpdateResponseModel), err
	}

	req.Header.Set("Authorization", i.JiraServerBase.AuthorizationMethod+" "+i.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("fail: response not ok: " + err.Error())
		return *new(models.IssueTypeUpdateResponseModel), err
	}
	defer res.Body.Close()

	tflog.Info(ctx, res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.IssueTypeUpdateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))
	tflog.Info(ctx, "response body: "+string(body))

	result := models.IssueTypeUpdateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.IssueTypeUpdateResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success update issue type")
	return result, nil
}

func (i IssueTypeService) Delete(ctx context.Context, model models.IssueTypeDeleteRequestModel) (models.IssueTypeDeleteResponseModel, error) {
	log.Printf("start delete issue type w. data: %s", model)

	foundIssueType, err := i.Get(ctx, models.IssueTypeGetRequestModel{
		Name: model.Name,
	})
	if err != nil {
		log.Println("failed to find issue type to be deleted")
		return *new(models.IssueTypeDeleteResponseModel), err
	}

	url := "https://" + i.JiraServerBase.Domain + "/rest/api/2/issuetype/" + foundIssueType.Id
	log.Println(url)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.IssueTypeDeleteResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", i.JiraServerBase.AuthorizationMethod+" "+i.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, err.Error())
		tflog.Info(ctx, "http request failed")
		return *new(models.IssueTypeDeleteResponseModel), errors.New("http request returned error")
	}

	tflog.Info(ctx, res.Status)

	tflog.Info(ctx, "success delete issue type list")
	return models.IssueTypeDeleteResponseModel{}, nil
}
