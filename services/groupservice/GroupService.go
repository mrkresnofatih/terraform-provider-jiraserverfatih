package groupservice

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
	models2 "terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/groupservice/models"
	"time"
)

type IGroupService interface {
	Create(ctx context.Context, model models.GroupCreateRequestModel) (models.GroupCreateResponseModel, error)
	Get(ctx context.Context, model models.GroupGetRequestModel) (models.GroupGetResponseModel, error)
	List(ctx context.Context, model models.GroupListApiRequestModel) (models.GroupListResponseModel, error)
	Delete(ctx context.Context, model models.GroupDeleteRequestModel) (models.GroupDeleteResponseModel, error)
}

type GroupService struct {
	JiraServerBase models2.JiraServerBase `json:"jiraServerBase"`
}

func (g GroupService) Create(ctx context.Context, model models.GroupCreateRequestModel) (models.GroupCreateResponseModel, error) {
	log.Printf("start create group w. data: %s", model)
	serialized, err := json.Marshal(model)
	if err != nil {
		log.Println("failed to marshal create request body")
		return *new(models.GroupCreateResponseModel), errors.New("error json marshal req body")
	}
	bodyReader := bytes.NewReader(serialized)
	url := "https://" + g.JiraServerBase.Domain + "/rest/api/2/group"
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		tflog.Info(ctx, "error building http request")
		return *new(models.GroupCreateResponseModel), err
	}

	req.Header.Set("Authorization", g.JiraServerBase.AuthorizationMethod+" "+g.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("fail: response not ok: " + err.Error())
		return *new(models.GroupCreateResponseModel), err
	}
	defer res.Body.Close()

	tflog.Info(ctx, res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.GroupCreateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))
	tflog.Info(ctx, "response body: "+string(body))

	result := models.GroupCreateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.GroupCreateResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success create group")
	return result, nil
}

func (g GroupService) Get(ctx context.Context, model models.GroupGetRequestModel) (models.GroupGetResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start get group w. data: %s", model))
	grouplist, err := g.List(ctx, models.GroupListApiRequestModel{
		GroupName: model.Name,
	})
	if err != nil {
		log.Println(err.Error())
		tflog.Info(ctx, "error listing groups")
		return *new(models.GroupGetResponseModel), err
	}

	if grouplist.Groups == nil || len(grouplist.Groups) == 0 {
		tflog.Info(ctx, "error group list values is null or empty")
		return *new(models.GroupGetResponseModel), err
	}

	foundGroup := grouplist.Groups[0]
	log.Println("success get group")

	return foundGroup, nil
}

func (g GroupService) List(ctx context.Context, model models.GroupListApiRequestModel) (models.GroupListResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start list groups w. data: %s", model))
	url := "https://" + g.JiraServerBase.Domain + "/rest/api/2/groups/picker?query=" + model.GroupName
	log.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.GroupListResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", g.JiraServerBase.AuthorizationMethod+" "+g.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, err.Error())
		tflog.Info(ctx, "http request failed")
		return *new(models.GroupListResponseModel), errors.New("http request returned error")
	}

	tflog.Info(ctx, res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.GroupListResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "response body: "+string(body))
	log.Println(string(body))

	result := models.GroupListResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.GroupListResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success group list")
	return result, nil
}

func (g GroupService) Delete(ctx context.Context, model models.GroupDeleteRequestModel) (models.GroupDeleteResponseModel, error) {
	tflog.Info(ctx, fmt.Sprintf("start delete group w. data: %s", model))
	url := "https://" + g.JiraServerBase.Domain + "/rest/api/2/group?groupname=" + model.Name
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.GroupDeleteResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", g.JiraServerBase.AuthorizationMethod+" "+g.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, err.Error())
		tflog.Info(ctx, "http request failed")
		return *new(models.GroupDeleteResponseModel), errors.New("http request returned error")
	}

	tflog.Info(ctx, res.Status)
	_, err = io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.GroupDeleteResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "success group list")
	return models.GroupDeleteResponseModel{}, nil
}
