package models

type PermissionSchemeUpdateRequestModel struct {
	Id          int64  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
