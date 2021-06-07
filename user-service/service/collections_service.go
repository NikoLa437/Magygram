package service

import (
	"context"
	"errors"
	"fmt"
	"user-service/domain/model"
	"user-service/domain/repository"
	"user-service/domain/service-contracts"
	"user-service/service/intercomm"
)

type collectionsService struct {
	repository.UserRepository
	intercomm.AuthClient
	intercomm.PostClient
}

func NewCollectionsService(r repository.UserRepository, ic 	intercomm.AuthClient, pc intercomm.PostClient) service_contracts.CollectionsService {
	return &collectionsService{r, ic, pc}
}

func (c collectionsService) CreateCollection(ctx context.Context, bearer string, collectionName string) error {

	userId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return err
	}

	user, err := c.UserRepository.GetByID(ctx, userId)
	if err != nil {
		return errors.New("invalid user id")
	}

	if _, ok := user.FavouritePosts[collectionName]; ok {
		return errors.New(fmt.Sprintf("collection with name %s already exist", collectionName))
	}

	user.FavouritePosts[collectionName] = []model.IdWithMedia{}
	_, err = c.UserRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (c collectionsService) AddPostToCollection(ctx context.Context, bearer string, favouritePostRequest *model.FavouritePostRequest) error {

	userId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return err
	}

	user, err := c.UserRepository.GetByID(ctx, userId)
	if err != nil {
		return errors.New("invalid user id")
	}
	if favouritePostRequest.CollectionName != ""{
		if _, ok := user.FavouritePosts[favouritePostRequest.CollectionName]; !ok {
			return errors.New(fmt.Sprintf("invalid %s collection", favouritePostRequest.CollectionName))
		}
	}


	for colName, _ := range user.FavouritePosts {
		if colName != model.DefaultCollection {
			for _, favMedia := range user.FavouritePosts[colName] {
				if favMedia.Id == favouritePostRequest.PostId {
					return errors.New(fmt.Sprintf("post with %s id already in favourites", favouritePostRequest.PostId))
				}
			}
		}
	}

	postImage, err := c.PostClient.GetPostsFirstImage(favouritePostRequest.PostId)
	if err != nil {
		return err
	}

	if !isPostInDefaultCollection(user.FavouritePosts[model.DefaultCollection], favouritePostRequest.PostId) {
		user.FavouritePosts[model.DefaultCollection] = append(user.FavouritePosts[model.DefaultCollection], model.IdWithMedia{
			Id:    favouritePostRequest.PostId,
			Media: *postImage,
		})
	}

	if favouritePostRequest.CollectionName != "" {
		user.FavouritePosts[favouritePostRequest.CollectionName] = append(user.FavouritePosts[favouritePostRequest.CollectionName], model.IdWithMedia{
			Id:    favouritePostRequest.PostId,
			Media: *postImage,
		})
	}

	_, err = c.UserRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func isPostInDefaultCollection(media []model.IdWithMedia, id string) bool {

	for _, med := range media {
		if med.Id == id {
			return true
		}
	}
	return false
}

func (c collectionsService) GetUsersCollections(ctx context.Context, bearer string, except string) (map[string][]model.IdWithMedia, error) {
	userId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	user, err := c.UserRepository.GetByID(ctx, userId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var collectionsWith4Media = make(map[string][]model.IdWithMedia)

	for colName, _ := range user.FavouritePosts {
		if colName != except {
			var i int
			collectionsWith4Media[colName] = []model.IdWithMedia{}
			if i < 4 {
				for _, favMedia := range user.FavouritePosts[colName] {
					collectionsWith4Media[colName] = append(collectionsWith4Media[colName], favMedia)
					i = i + 1
				}
			}
			i = 0
		}
	}
	return collectionsWith4Media, nil
}

func (c collectionsService) CheckIfPostsInFavourites(ctx context.Context, bearer string, postIds *[]string) ([]*model.PostIdFavouritesFlag, error) {
	userId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	user, err := c.UserRepository.GetByID(ctx, userId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var postsFavFlags []*model.PostIdFavouritesFlag
	for _, postId := range *postIds {
		fav := false
		for _, favMedia := range user.FavouritePosts[model.DefaultCollection] {
			if favMedia.Id == postId {
				fav = true
			}
		}
		postsFavFlags = append(postsFavFlags, &model.PostIdFavouritesFlag{
			Id:         postId,
			Favourites: fav,
		})
	}
	return postsFavFlags, nil
}

func (c collectionsService) DeletePostFromCollections(ctx context.Context, bearer string, postId string) error {
	userId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return err
	}

	user, err := c.UserRepository.GetByID(ctx, userId)
	if err != nil {
		return errors.New("invalid user id")
	}

	user.FavouritePosts[model.DefaultCollection] = deletePostFromCollection(user.FavouritePosts[model.DefaultCollection], postId)

	for colName, _ := range user.FavouritePosts {
		wentIn := false
		if colName != model.DefaultCollection && wentIn == false {
			user.FavouritePosts[colName] = deletePostFromCollection(user.FavouritePosts[colName], postId)
			wentIn = true
		}
	}
	_, err = c.UserRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func deletePostFromCollection(collection []model.IdWithMedia, postId string) []model.IdWithMedia {
	var colCpy []model.IdWithMedia

	for _, favMedia := range collection {
		if favMedia.Id != postId {
			colCpy = append(colCpy, favMedia)
		}
	}
	return colCpy
}


func (c collectionsService) GetCollectionPosts(ctx context.Context, bearer string, collectionName string) ([]*model.PostProfileResponse, error) {
	userId, err := c.AuthClient.GetLoggedUserId(bearer)
	if err != nil {
		return nil, err
	}

	user, err := c.UserRepository.GetByID(ctx, userId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if _, ok := user.FavouritePosts[collectionName]; !ok {
		return nil, errors.New(fmt.Sprintf("collection with name %s not exist", collectionName))
	}

	var postsFavFlags []*model.PostProfileResponse
	for _, fav := range user.FavouritePosts[collectionName] {
		postsFavFlags = append(postsFavFlags, &model.PostProfileResponse{
			Id:    fav.Id,
			Media: fav.Media,
		})
	}
	fmt.Println(len(postsFavFlags))
	return postsFavFlags, nil
}