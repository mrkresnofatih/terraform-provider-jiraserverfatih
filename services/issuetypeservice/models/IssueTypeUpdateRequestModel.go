package models

type IssueTypeUpdateRequestModel struct {
	Id          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarId    int64  `json:"avatarId"`
}
