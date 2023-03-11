package models

type PermissionSchemeGetResponseModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Id          int64  `json:"id"`
}
