package models

type GrantCreateRequestModel struct {
	PermissionSchemeId int64               `json:"-"`
	Permission         string              `json:"permission"`
	Holder             GrantApiHolderModel `json:"holder"`
}

type GrantCreateApiRequestModel struct {
	PermissionSchemeName string              `json:"-"`
	Permission           string              `json:"permission"`
	Holder               GrantApiHolderModel `json:"holder"`
}

type GrantHolderModel struct {
	Type      string `json:"type"`
	Parameter string `json:"parameter"`
}

type GrantApiHolderModel struct {
	Type      string `json:"type"`
	Parameter int64  `json:"parameter"`
}
