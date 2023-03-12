package models

type GrantListResponseModel struct {
	Grants []GrantGetResponseModel `json:"permissions"`
}
