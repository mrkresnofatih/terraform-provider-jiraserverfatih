package models

type GrantGetRequestModel struct {
	PermissionSchemeId int64            `json:"permissionSchemeId"`
	Permission         string           `json:"permission"`
	Holder             GrantHolderModel `json:"holder"`
}
