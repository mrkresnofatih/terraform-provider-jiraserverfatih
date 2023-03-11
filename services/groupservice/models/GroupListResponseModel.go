package models

type GroupListResponseModel struct {
	Groups []GroupGetResponseModel `json:"groups"`
}
