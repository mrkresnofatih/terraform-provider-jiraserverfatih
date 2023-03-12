package screenservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
	models2 "terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/screenservice/models"
	"time"
)

type IScreenService interface {
	Get(ctx context.Context, model models.ScreenGetRequestModel) (models.ScreenGetResponseModel, error)
	List(ctx context.Context, model models.ScreenListRequestModel) (models.ScreenListResponseModel, error)
	Create(ctx context.Context, model models.ScreenCreateRequestModel) (models.ScreenCreateResponseModel, error)
	Update(ctx context.Context, model models.ScreenUpdateRequestModel) (models.ScreenUpdateResponseModel, error)
	Delete(ctx context.Context, model models.ScreenDeleteRequestModel) (models.ScreenDeleteResponseModel, error)
}

type ScreenService struct {
	JiraServerBase models2.JiraServerBase `json:"jiraServerBase"`
}

func (s ScreenService) Get(ctx context.Context, model models.ScreenGetRequestModel) (models.ScreenGetResponseModel, error) {
	log.Printf("start get screen w. data %s", model)

	screens, err := s.List(ctx, models.ScreenListRequestModel{
		Name: model.Name,
	})
	if err != nil {
		log.Println("failed to list roles")
		return *new(models.ScreenGetResponseModel), errors.New("failed to list roles")
	}

	if len(screens.Values) == 0 {
		tflog.Info(ctx, "screens is empty, cannot find screen with provided name")
		return *new(models.ScreenGetResponseModel), errors.New("failed to list roles")
	}

	foundScreen := screens.Values[0]

	log.Println("success get screen")
	return foundScreen, nil
}

func (s ScreenService) List(ctx context.Context, model models.ScreenListRequestModel) (models.ScreenListResponseModel, error) {
	log.Printf("start list screen w. data %s", model)

	url := "https://" + s.JiraServerBase.Domain + "/rest/api/2/screens?querystring=" + url2.QueryEscape(model.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tflog.Info(ctx, "error building http request")
		return *new(models.ScreenListResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", s.JiraServerBase.AuthorizationMethod+" "+s.JiraServerBase.Token)
	req.Header.Set("Accept", "application/json")

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		tflog.Info(ctx, err.Error())
		tflog.Info(ctx, "http request failed")
		return *new(models.ScreenListResponseModel), errors.New("http request returned error")
	}

	tflog.Info(ctx, res.Status)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		tflog.Info(ctx, "read all body failed")
		return *new(models.ScreenListResponseModel), errors.New("error reading response body")
	}

	tflog.Info(ctx, "response body: "+string(body))

	result := models.ScreenListResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tflog.Info(ctx, "failed to unmarshal response body")
		return *new(models.ScreenListResponseModel), errors.New("error unmarshalling response body")
	}

	tflog.Info(ctx, "success project list roles")
	return result, nil
}

func (s ScreenService) Create(ctx context.Context, model models.ScreenCreateRequestModel) (models.ScreenCreateResponseModel, error) {
	log.Printf("start create screen w. data %s", model)

	url := "https://" + s.JiraServerBase.Domain + "/rest/api/2/screens"
	log.Println(url)
	serialized, err := json.Marshal(model)
	if err != nil {
		log.Println("error when json marshaling")
		return *new(models.ScreenCreateResponseModel), errors.New("error json marshal request data")
	}
	bodyReader := bytes.NewReader(serialized)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		return *new(models.ScreenCreateResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", s.JiraServerBase.AuthorizationMethod+" "+s.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.ScreenCreateResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read all body failed")
		return *new(models.ScreenCreateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))

	result := models.ScreenCreateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("failed to unmarshal response body")
		return *new(models.ScreenCreateResponseModel), errors.New("error unmarshalling response body")
	}

	log.Println("success create screen")
	return result, nil
}

func (s ScreenService) Update(ctx context.Context, model models.ScreenUpdateRequestModel) (models.ScreenUpdateResponseModel, error) {
	log.Printf("start update screen w. data %s", model)

	foundScreen, err := s.Get(ctx, models.ScreenGetRequestModel{
		Name: model.Name,
	})
	if err != nil {
		tflog.Info(ctx, "get screen failed")
		return *new(models.ScreenUpdateResponseModel), err
	}

	serial, err := json.Marshal(model)
	if err != nil {
		log.Println("error when json marshaling")
		return *new(models.ScreenUpdateResponseModel), errors.New("error json marshal request data")
	}
	bodyReader := bytes.NewReader(serial)

	url := "https://" + s.JiraServerBase.Domain + "/rest/api/2/screens/" + strconv.FormatInt(foundScreen.Id, 10)
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		log.Println("error building http request")
		return *new(models.ScreenUpdateResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", s.JiraServerBase.AuthorizationMethod+" "+s.JiraServerBase.Token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.ScreenUpdateResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read all body failed")
		return *new(models.ScreenUpdateResponseModel), errors.New("error reading response body")
	}

	log.Println("response body: " + string(body))

	result := models.ScreenUpdateResponseModel{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("failed to unmarshal response body")
		return *new(models.ScreenUpdateResponseModel), errors.New("error unmarshalling response body")
	}

	log.Println("success update screen")
	return result, nil
}

func (s ScreenService) Delete(ctx context.Context, model models.ScreenDeleteRequestModel) (models.ScreenDeleteResponseModel, error) {
	log.Printf("start delete screen w. data %s", model)

	foundScreen, err := s.Get(ctx, models.ScreenGetRequestModel{
		Name: model.Name,
	})
	if err != nil {
		tflog.Info(ctx, "error finding screen")
		return *new(models.ScreenDeleteResponseModel), err
	}

	url := "https://" + s.JiraServerBase.Domain + "/rest/api/2/screens/" + strconv.FormatInt(foundScreen.Id, 10)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Println("error building http request")
		return *new(models.ScreenDeleteResponseModel), errors.New("error building http request")
	}
	req.Header.Set("Authorization", s.JiraServerBase.AuthorizationMethod+" "+s.JiraServerBase.Token)

	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("error result from http request")
		return *new(models.ScreenDeleteResponseModel), errors.New("error result from http request")
	}
	defer res.Body.Close()

	log.Println(res.Status)

	log.Println("delete screen success")
	return models.ScreenDeleteResponseModel{}, nil
}
