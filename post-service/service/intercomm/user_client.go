package intercomm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"post-service/conf"
	"post-service/domain/model"
	"post-service/logger"
)

type UserClient interface {
	GetLoggedUserInfo(bearer string) (*model.UserInfo,error)
	MapPostsToFavourites(bearer string, postIds []string) ([]*model.PostIdFavouritesFlag,error)
	IsUserPrivate(userId string) (bool, error)
	UpdateLikedPosts(bearer string, postId string) error
	AddComment(bearer string, postId string) error
	UpdateDislikedPosts(bearer string, postId string) error
	GetLikedPosts(bearer string) ([]string, error)
	GetDislikedPosts(bearer string) ([]string, error)
}

type userClient struct {}

func NewUserClient() UserClient {
	baseUsersUrl = fmt.Sprintf("%s%s:%s/api/users", conf.Current.Userservice.Protocol, conf.Current.Userservice.Domain, conf.Current.Userservice.Port)
	return &userClient{}
}

var (
	baseUsersUrl = ""
)


func (u userClient) AddComment(bearer string, postId string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/post/comment/%s", baseUsersUrl, postId), nil)
	req.Header.Add("Authorization", bearer)
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err
	}

	return nil
}

func (u userClient) GetLoggedUserInfo(bearer string) (*model.UserInfo, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/logged", baseUsersUrl), nil)
	req.Header.Add("Authorization", bearer)
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		if resp == nil {
			logger.LoggingEntry.Error("User-service not available")
			return &model.UserInfo{}, err
		}

		logger.LoggingEntry.Error("User-service get logged user info")

		return &model.UserInfo{}, errors.New("unauthorized")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &model.UserInfo{}, err
	}
	var userInfo model.UserInfo
	_ = json.Unmarshal(bodyBytes, &userInfo)

	return &userInfo, nil
}

func (u userClient) UpdateLikedPosts(bearer string, postId string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/post/like/" + postId, baseUsersUrl), nil)
	req.Header.Add("Authorization", bearer)
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err
	}

	return nil
}

func (u userClient) UpdateDislikedPosts(bearer string, postId string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/post/dislike/" + postId, baseUsersUrl), nil)
	req.Header.Add("Authorization", bearer)
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err
	}

	return nil
}

func (u userClient) MapPostsToFavourites(bearer string, postIds []string) ([]*model.PostIdFavouritesFlag, error) {

	jsonStr, err:= json.Marshal(postIds)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/collections/check-favourites", baseUsersUrl), bytes.NewReader(jsonStr))
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		if resp == nil {
			logger.LoggingEntry.WithFields(logrus.Fields{"post_ids" : postIds}).Error("User-service not available")
			return nil, err
		}

		logger.LoggingEntry.WithFields(logrus.Fields{"post_ids" : postIds}).Error("User-service map posts to favourites")
		return nil, errors.New("unauthorized")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var postIdFav []*model.PostIdFavouritesFlag
	_ = json.Unmarshal(bodyBytes, &postIdFav)

	return postIdFav, nil
}


func (u userClient) IsUserPrivate(userId string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/is-private", baseUsersUrl, userId), nil)
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		if resp == nil {
			logger.LoggingEntry.WithFields(logrus.Fields{"user_id" : userId}).Error("User-service not available")
			return false, err
		}

		logger.LoggingEntry.WithFields(logrus.Fields{"user_id" : userId}).Error("User-service check user privacy")
		return false, errors.New("user not found")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var isPrivate bool
	json.Unmarshal(bodyBytes, &isPrivate)

	return isPrivate, nil
}

func (u userClient) GetLikedPosts(bearer string) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/post/liked", baseUsersUrl), nil)
	req.Header.Add("Authorization", bearer)
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return []string{} ,err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var users []string
	_ = json.Unmarshal(bodyBytes, &users)

	fmt.Println(users)

	return users, nil
}

func (u userClient) GetDislikedPosts(bearer string) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/post/disliked", baseUsersUrl), nil)
	req.Header.Add("Authorization", bearer)
	hash, _ := bcrypt.GenerateFromPassword([]byte(conf.Current.Server.Secret), bcrypt.MinCost)
	req.Header.Add(conf.Current.Server.Handshake, string(hash))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return []string{} ,err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var users []string
	_ = json.Unmarshal(bodyBytes, &users)

	fmt.Println(users)

	return users, nil}