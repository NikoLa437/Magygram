package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator"
	"post-service/domain/model"
	"post-service/domain/repository"
	"post-service/domain/service-contracts"
	"post-service/service/intercomm"
)



type postService struct {
	repository.PostRepository
	intercomm.MediaClient
	intercomm.UserClient
}

func NewPostService(r repository.PostRepository, ic intercomm.MediaClient, uc intercomm.UserClient) service_contracts.PostService {
	return &postService{r , ic, uc}
}

func (p postService) CreatePost(ctx context.Context, bearer string, postRequest *model.PostRequest) (string, error) {

	userInfo, err := p.UserClient.GetLoggedUserInfo(bearer)
	if err != nil { return "", err}

	media, err := p.MediaClient.SaveMedia(postRequest.Media)
	if err != nil { return "", err}

	post, err := model.NewPost(postRequest, *userInfo, "REGULAR", media)

	if err != nil { return "", err}

	if err := validator.New().Struct(post); err!= nil {
		return "", err
	}

	result, err := p.PostRepository.Create(ctx, post)

	if err != nil { return "", err}
	if postId, ok := result.InsertedID.(string); ok {
		if err != nil { return "", err}
		return postId, nil
	}

	return "", err
}

func (p postService) GetPostsFirstImage(ctx context.Context, postId string) (*model.Media, error) {

	post, err := p.PostRepository.GetByID(ctx, postId)

	if err != nil {
		return nil, errors.New("invalid post id")
	}
	if len(post.Media) > 0 {
		return &post.Media[0], nil
	}
	return nil, nil
}

