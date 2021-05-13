package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sethigeet/gql-go-auth-backend/graph/generated"
	"github.com/sethigeet/gql-go-auth-backend/graph/model"
	"github.com/sethigeet/gql-go-auth-backend/util"
	"github.com/sethigeet/gql-go-auth-backend/validator"
	"gorm.io/gorm"
)

func (r *mutationResolver) ConfirmEmail(ctx context.Context, token string) (*model.ConfirmEmailResponse, error) {
	success := false
	userID, deleteToken, err := util.GetUserIDFromToken(r.RDB, token, util.ConfirmEmailPrefix)
	if err != nil || userID == "" {
		return &model.ConfirmEmailResponse{
			Errors: []*model.FieldError{
				{
					Field:   "token",
					Message: validator.GetInvalidTokenMessage("token"),
				},
			},
			Successful: &success,
		}, nil
	}

	var user model.User
	result := r.DB.First(&user, "id = ?", userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &model.ConfirmEmailResponse{
				Errors: []*model.FieldError{
					{
						Field:   "token",
						Message: validator.GetInvalidTokenMessage("token"),
					},
				},
				Successful: &success,
			}, nil
		}

		return nil, result.Error
	}

	result = r.DB.Model(&user).Update("confirmed", true)
	if result.Error != nil {
		return nil, result.Error
	}

	deleteToken()

	success = true
	return &model.ConfirmEmailResponse{
		Errors:     nil,
		Successful: &success,
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

	err := util.SendEmail(r.RDB, user.ID, user.Email, util.ConfirmEmailPrefix, "/confirm-email")
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		Errors: nil,
		User:   &user,
	}, nil
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, credentials model.ForgotPasswordInput) (*model.ResetPasswordResponse, error) {
	success := false
	if validationErrors := validator.Validate(credentials); validationErrors != nil {
		return &model.ResetPasswordResponse{
			Successful: &success,
			Errors:     validationErrors,
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
			return &model.ResetPasswordResponse{
				Errors: []*model.FieldError{
					{
						Field:   "usernameOrEmail",
						Message: validator.GetDoesNotExistMessage("username/email"),
					},
				},
				Successful: &success,
			}, nil
		}

		return nil, result.Error
	}

	err := util.SendEmail(r.RDB, user.ID, user.Email, util.ForgotPasswordPrefix, "/forgot-password")
	if err != nil {
		return nil, err
	}

	success = true
	return &model.ResetPasswordResponse{
		Errors:     nil,
		Successful: &success,
	}, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, credentials model.ChangePasswordInput) (*model.ResetPasswordResponse, error) {
	success := false
	if validationErrors := validator.Validate(credentials); validationErrors != nil {
		return &model.ResetPasswordResponse{
			Successful: &success,
			Errors:     validationErrors,
		}, nil
	}

	userID, deleteToken, err := util.GetUserIDFromToken(r.RDB, credentials.Token, util.ForgotPasswordPrefix)
	if err != nil || userID == "" {
		return &model.ResetPasswordResponse{
			Successful: &success,
			Errors: []*model.FieldError{
				{
					Field:   "token",
					Message: validator.GetInvalidTokenMessage("token"),
				},
			},
		}, nil
	}

	var user model.User
	result := r.DB.First(&user, "id = ?", userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &model.ResetPasswordResponse{
				Successful: &success,
				Errors: []*model.FieldError{
					{
						Field:   "token",
						Message: validator.GetInvalidTokenMessage("token"),
					},
				},
			}, nil
		}

		return nil, result.Error
	}

	hashedPasswd, err := util.HashPassword(credentials.NewPassword)
	if err != nil {
		return nil, err
	}

	result = r.DB.Model(&user).Update("password", hashedPasswd)
	if result.Error != nil {
		return nil, result.Error
	}

	deleteToken()

	success = true
	return &model.ResetPasswordResponse{
		Errors:     nil,
		Successful: &success,
	}, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userID, err := r.SessionManager.Retrieve(false)
	if err != nil {
		return nil, err
	}

	var user model.User
	result := r.DB.First(&user, "id = ?", userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not authenticated")
		}

		return nil, result.Error
	}

	return &user, nil
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
