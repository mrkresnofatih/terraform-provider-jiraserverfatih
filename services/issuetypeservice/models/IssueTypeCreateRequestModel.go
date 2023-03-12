package models

type IssueTypeCreateRequestModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	AvatarId    int64  `json:"-"`
}
