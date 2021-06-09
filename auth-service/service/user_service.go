package service

import (
	"auth-service/domain/model"
	"auth-service/domain/repository"
	"auth-service/domain/service-contracts"
	"auth-service/logger"
	"context"
	"errors"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

type userService struct {
	repository.UserRepository
	repository.LoginEventRepository
}

func NewUserService(r repository.UserRepository,a repository.LoginEventRepository) service_contracts.UserService {
	return &userService{r,a}
}

func (u userService) RegisterUser(ctx context.Context, userRequest *model.UserRequest) (string, error) {
	user, err := model.NewUser(userRequest)
	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return "", err
	}

	if err := validator.New().Struct(user); err!= nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Warn("User registration validation failure")
		return "", err
	}

	result, err := u.UserRepository.Create(ctx, user)
	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"email" : userRequest.Email}).Error("User database create failure")
		return "", err
	}

	if userId, ok := result.InsertedID.(string); ok {
		logger.LoggingEntry.WithFields(logrus.Fields{"user_id" : userId}).Info("User registered")
		return userId, nil
	}
	return "", err
}

func (u userService) ActivateUser(ctx context.Context, userId string) (bool, error) {

	user, err := u.UserRepository.GetByID(ctx, userId)
	if err != nil {
		return false, err
	}
	user.Active = true
	_, err = u.UserRepository.Update(ctx, user)
	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"user_id": userId}).Error("User database update failure")
		return false, err
	}

	u.HandleLoginEventAndAccountActivation(ctx, user.Email, true, model.ActivatedAccount)
	logger.LoggingEntry.WithFields(logrus.Fields{"user_id" : userId}).Info("User activated")

	return true, err
}

func (u userService) HandleLoginEventAndAccountActivation(ctx context.Context, userEmail string, successful bool, eventType string) {
	if successful {
		_, err := u.LoginEventRepository.Create(ctx, model.NewLoginEvent(userEmail, eventType, 0))
		if err != nil {
			logger.LoggingEntry.WithFields(logrus.Fields{"user_email" : userEmail}).Warn("Create success login event, database failure")
		}
		return
	}
}

func (u userService) DeactivateUser(ctx context.Context, userEmail string) (bool, error) {
	user, err := u.UserRepository.GetByEmail(ctx, userEmail)
	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"email": userEmail}).Warn("Invalid email address")
		return false, err
	}
	user.Active = false
	_, err = u.UserRepository.Update(ctx, user)
	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"user_id": user.Id}).Error("User database update failure")
		return false, err
	}

	logger.LoggingEntry.WithFields(logrus.Fields{"user_id": user.Id}).Info("User deactivated")

	return true, err
}

func (u userService) ResetPassword(ctx context.Context, changePasswordRequest *model.PasswordChangeRequest) (bool, error) {
	hashAndSalt, err := model.HashAndSaltPasswordIfStrongAndMatching(changePasswordRequest.Password, changePasswordRequest.PasswordRepeat)
	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"user_id" : changePasswordRequest.UserId}).Warn("Passwords not valid")
		return false, err
	}

	user, err := u.UserRepository.GetByID(ctx, changePasswordRequest.UserId)
	if err != nil {
		return false, err
	}
	user.Password = hashAndSalt
	user.Active = true
	_, err = u.UserRepository.Update(ctx, user)
	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"user_id": user.Id}).Error("User database update failure")
		return false, err
	}

	logger.LoggingEntry.WithFields(logrus.Fields{"user_id" : user.Id}).Info("Users password changed")
	return true, err
}

func (u userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := u.UserRepository.GetByEmail(ctx, email)

	if err != nil {
		logger.LoggingEntry.WithFields(logrus.Fields{"email": email}).Warn("Invalid email address")
		return nil, errors.New("invalid user id")
	}

	return user, err
}

func (u userService) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	user, err := u.UserRepository.GetByID(ctx, userId)

	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return user, err
}

func (u userService) GetAllRolesByUserId(ctx context.Context, userId string) ([]model.Role, error) {
	return u.UserRepository.GetAllRolesByUserId(ctx, userId)
}




