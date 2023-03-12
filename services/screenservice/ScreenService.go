package screenservice

import (
	"context"
	models2 "terraform-provider-hashicups-pf/services/baseservice/models"
	"terraform-provider-hashicups-pf/services/screenservice/models"
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
	//TODO implement me
	panic("implement me")
}

func (s ScreenService) List(ctx context.Context, model models.ScreenListRequestModel) (models.ScreenListResponseModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s ScreenService) Create(ctx context.Context, model models.ScreenCreateRequestModel) (models.ScreenCreateResponseModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s ScreenService) Update(ctx context.Context, model models.ScreenUpdateRequestModel) (models.ScreenUpdateResponseModel, error) {
	//TODO implement me
	panic("implement me")
}

func (s ScreenService) Delete(ctx context.Context, model models.ScreenDeleteRequestModel) (models.ScreenDeleteResponseModel, error) {
	//TODO implement me
	panic("implement me")
}
