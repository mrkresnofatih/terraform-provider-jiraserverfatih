package models

type IssueTypeCreateResponseModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarId    int64  `json:"avatarId"`
}
