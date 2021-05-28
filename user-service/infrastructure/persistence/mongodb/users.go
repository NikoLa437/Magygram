package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"user-service/domain/model"
	"user-service/domain/repository"
)

type userRepository struct {
	Col *mongo.Collection
}

func NewUserRepository(Col *mongo.Collection) repository.UserRepository {
	return &userRepository{Col}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*mongo.InsertOneResult, error) {
	return r.Col.InsertOne(ctx, user)
}

func (r *userRepository) Update(ctx context.Context, user *model.User) (*mongo.UpdateResult, error) {
	return r.Col.UpdateOne(ctx, bson.M{"_id":  user.Id},bson.D{{"$set", bson.D{{"email" , user.Email},
																{"active" , user.Active},
																{"name" , user.Name},
																{"password" , user.Password},
																{"surname" , user.Surname}}}})
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {

	var user = model.User{}
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {

	var user = model.User{}
	err := r.Col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllRolesByUserId(ctx context.Context, userId string) ([]model.Role, error) {
	user, err := r.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}


//func (r *userRepository) GetByEmailEagerly(ctx context.Context, email string) (*model.User, error) {
//	u := &model.User{Email: email}
//	err := r.Conn.Preload("Roles").Preload("Roles.Permissions").First(&u, "email = ?", email).Error
//	return u, err
//}
