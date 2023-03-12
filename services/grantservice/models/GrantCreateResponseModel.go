package models

type GrantCreateResponseModel struct {
	Id                   int64            `json:"id"`
	PermissionSchemeName string           `json:"permissionSchemeName"`
	Permission           string           `json:"permission"`
	Holder               GrantHolderModel `json:"holder"`
}
