package models

type IssueTypeUpdateRequestModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarId    int64  `json:"avatarId"`
}
