package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/sethigeet/gql-go-auth-backend/graph/generated"
	"github.com/sethigeet/gql-go-auth-backend/graph/model"
	"github.com/sethigeet/gql-go-auth-backend/util"
	"github.com/sethigeet/gql-go-auth-backend/validator"
)

func (r *mutationResolver) ConfirmEmail(ctx context.Context, token string) (*model.UserResponse, error) {
	userID, err := util.GetUserIDFromEmailToken(r.RDB, token)
	if err != nil || userID == "" {
		return &model.UserResponse{
			Errors: []*model.FieldError{
				{
					Field:   "token",
					Message: validator.GetInvalidTokenMessage("token"),
				},
			},
			User: nil,
		}, nil
	}

	var user model.User
	result := r.DB.First(&user, "id = ?", userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &model.UserResponse{
				Errors: []*model.FieldError{
					{
						Field:   "token",
						Message: validator.GetInvalidTokenMessage("token"),
					},
				},
				User: nil,
			}, nil
		}

		return nil, result.Error
	}

	result = r.DB.Model(&user).Update("confirmed", true)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.UserResponse{
		Errors: nil,
		User:   &user,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, credentials model.LoginInput) (*model.UserResponse, error) {
	if validationErrors := validator.Validate(credentials); validationErrors != nil {
		return &model.UserResponse{
			User:   nil,
			Errors: validationErrors,
		}, nil
	}

	var user model.User
	var result *gorm.DB
	if strings.ContainsRune(credentials.UsernameOrEmail, '@') {
		result = r.DB.First(&user, "email = ?", credentials.UsernameOrEmail)
	} else {
		result = r.DB.First(&user, "username = ?", credentials.UsernameOrEmail)
	}

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &model.UserResponse{
				Errors: []*model.FieldError{
					{
						Field:   "usernameOrEmail",
						Message: validator.GetDoesNotExistMessage("username/email"),
					},
				},
				User: nil,
			}, nil
		}

		return nil, result.Error
	}

	if !user.Confirmed {
		return &model.UserResponse{
			Errors: []*model.FieldError{
				{
					Field:   "email",
					Message: "Please confirm your email first!",
				},
			},
			User: nil,
		}, nil
	}

	verified := util.ComparePasswords(user.Password, credentials.Password)
	if !verified {
		return &model.UserResponse{
			Errors: []*model.FieldError{
				{
					Field:   "password",
					Message: validator.GetIncorrectMessage("password"),
				},
			},
			User: nil,
		}, nil
	}

	err := r.SessionManager.Create(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		Errors: nil,
		User:   &user,
	}, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	var err error
	sessionID, err := r.SessionManager.Retrieve(true)
	if err != nil {
		return false, err
	}

	err = r.SessionManager.Delete(sessionID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) Register(ctx context.Context, credentials model.RegisterInput) (*model.UserResponse, error) {
	if validationErrors := validator.Validate(credentials); validationErrors != nil {
		return &model.UserResponse{
			User:   nil,
			Errors: validationErrors,
		}, nil
	}

	user := model.User{
		Email:    credentials.Email,
		Username: credentials.Username,
		Password: credentials.Password,
	}
	result := r.DB.Create(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			return &model.UserResponse{
				Errors: []*model.FieldError{
					{
						Field:   "email",
						Message: validator.GetAlreadyExistsMessage("email"),
					},
				},
				User: nil,
			}, nil
		}

		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint \"users_username_key\"") {
			return &model.UserResponse{
				Errors: []*model.FieldError{
					{
						Field:   "username",
						Message: validator.GetAlreadyExistsMessage("username"),
					},
				},
				User: nil,
			}, nil
		}
		return nil, result.Error
	}

	err := util.SendConfirmEmailEmail(r.RDB, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		Errors: nil,
		User:   &user,
	}, nil
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, credentials model.ForgotPasswordInput) (*model.ResetPasswordResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChangePassword(ctx context.Context, credentials model.ChangePasswordInput) (*model.ResetPasswordResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) CreatedAt(ctx context.Context, obj *model.User) (int, error) {
	return int(obj.CreatedAt.Unix()), nil
}

func (r *userResolver) UpdatedAt(ctx context.Context, obj *model.User) (int, error) {
	return int(obj.UpdatedAt.Unix()), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
