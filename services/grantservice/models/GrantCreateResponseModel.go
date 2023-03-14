package models

type GrantCreateResponseModel struct {
	Id                 int64            `json:"id"`
	PermissionSchemeId int64            `json:"permissionSchemeName"`
	Permission         string           `json:"permission"`
	Holder             GrantHolderModel `json:"holder"`
}
