package model

import (
	"errors"
	"github.com/beevik/guid"
	"html"
)

type User struct {
	Id string `bson:"_id,omitempty"`
	Username  string `bson:"username" validate:"required,min=1"`
	Name  string `bson:"name" validate:"required,min=2"`
	Email string `bson:"email" validate:"required,email"`
	Surname string `bson:"surname" validate:"required,min=2"`
	FavouritePosts map[string][]IdWithMedia `bson:"favouritePosts"`
}

type UserInfo struct {
	Id string `json:"id"`
	Username  string `json:"username"`
	ImageURL  string `json:"imageUrl"`
}

type UserRequest struct {
	Name  string `json:"name"`
	Surname string `json:"surname"`
	Username  string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	RepeatedPassword string `json:"repeatedPassword"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ActivateLinkRequest struct {
	Email string `json:"email"`
}

type ChangeNewPasswordRequest struct {
	ResetPasswordId string `json:"resetPasswordId"`
	Password string `json:"password"`
	PasswordRepeat string `json:"passwordRepeat"`
}

func validateMediaTypeEnums(md Media) error {

	if md.MediaType != "IMAGE" && md.MediaType !="VIDEO" {
		return errors.New("Invalid media type")
	}
	return nil
}

type IdWithMedia struct {
	Id string `json:"id"`
	Media Media `json:"media"`
}

type FavouritePostRequest struct {
	PostId string `json:"postId"`
	CollectionName string `json:"collectionName"`
}


type Media struct {
	Url string
	MediaType string
}

type MediaType string

const(
	IMAGE = iota
	VIDEO
)

var (
	DefaultCollection = "all posts"
)

func NewUser(userRequest *UserRequest) *User {
	return &User{Id: guid.New().String(),
		Name:     html.EscapeString(userRequest.Name),
		Surname:  html.EscapeString(userRequest.Surname),
		Username: html.EscapeString(userRequest.Username),
		Email:    html.EscapeString(userRequest.Email),
		FavouritePosts: map[string][]IdWithMedia{ DefaultCollection: {}}}
}