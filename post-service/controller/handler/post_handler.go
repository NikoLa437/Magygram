package handler

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"mime/multipart"
	"net/http"
	"post-service/domain/model"
	"post-service/domain/service-contracts"
)


type PostHandler interface {
	CreatePost(c echo.Context) error
	GetPostsFirstImage(c echo.Context) error

}

type postHandler struct {
	PostService service_contracts.PostService
}

func NewPostHandler(p service_contracts.PostService) PostHandler {
	return &postHandler{p}
}

func (p postHandler) CreatePost(c echo.Context) error {

	location := c.FormValue("location")
	description := c.FormValue("description")
	tagsString := c.FormValue("tags")

	mpf, _ := c.MultipartForm()
	var tags []string
	json.Unmarshal([]byte(tagsString), &tags)

	var headers []*multipart.FileHeader
	for _, v := range mpf.File {
		headers = append(headers, v[0])
	}

	postRequest := &model.PostRequest{
		Description: description,
		Location:    location,
		Media:       headers,
		Tags:        tags,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	bearer := c.Request().Header.Get("Authorization")

	postId, err := p.PostService.CreatePost(ctx, bearer, postRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, postId)
}

func (p postHandler) GetPostsFirstImage(c echo.Context) error {
	postId := c.Param("postId")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	postImage, err := p.PostService.GetPostsFirstImage(ctx, postId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, postImage)
}


