package models

type GrantGetRequestModel struct {
	PermissionSchemeName string           `json:"permissionSchemeName"`
	Permission           string           `json:"permission"`
	Holder               GrantHolderModel `json:"holder"`
}
