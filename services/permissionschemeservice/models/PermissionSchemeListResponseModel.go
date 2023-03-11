package models

type PermissionSchemeListResponseModel struct {
	PermissionSchemes []PermissionSchemeGetResponseModel `json:"permissionSchemes"`
}
