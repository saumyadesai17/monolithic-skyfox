package service

import (
	"context"
	"net/http"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/model"
	serviceMocks "skyfox/bookings/service/mocks"
	ae "skyfox/error"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var validReq = request.SignupRequest{
	Name:     "Alice",
	Phone:    "9876543210",
	Password: "Str0ng@Pass",
	Email:    "alice@example.com",
}

func TestAuthService_Signup(t *testing.T) {
	existingUser := &model.UserAccount{
		Id:    "existing-uuid",
		Phone: "9876543210",
	}
	dbErr := ae.InternalServerError("InternalServerError", "something went wrong", nil)

	var tests = []struct {
		name          string
		setupMock     func(repo *serviceMocks.MockAuthUserRepository)
		wantHTTPCode  int
		wantErrCode   string
		wantUser      bool // true = expect a non-nil user back
		wantEmptyHash bool // true = password hash must not be on the returned user
		weakPassword  bool // true= send a request with a weak password
	}{
		{
			name: "Should create user and not return password hash when phone is not registered",
			setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
				repo.On("FindByPhone", mock.Anything, validReq.Phone).Return(nil, nil).Once()
				repo.On("FindByEmail", mock.Anything, validReq.Email).Return(nil, nil).Once()
				repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.UserAccount) bool {
					return u.Phone == validReq.Phone && u.Name == validReq.Name
				})).Run(func(args mock.Arguments) {
					u := args.Get(1).(*model.UserAccount)
					u.Id = "new-uuid-1234"
				}).Return(nil).Once()
			},
			wantHTTPCode:  http.StatusCreated,
			wantErrCode:   "",
			wantUser:      true,
			wantEmptyHash: true,
		},
		{
			name: "Should store password as bcrypt hash when user is created",
			setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
				repo.On("FindByPhone", mock.Anything, validReq.Phone).Return(nil, nil).Once()
				repo.On("FindByEmail", mock.Anything, validReq.Email).Return(nil, nil).Once()
				repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.UserAccount) bool {
					// The hash must be a valid bcrypt hash of the original password
					err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(validReq.Password))
					return err == nil && u.PasswordHash != validReq.Password
				})).Return(nil).Once()
			},
			wantHTTPCode:  http.StatusCreated,
			wantUser:      true,
			wantEmptyHash: true,
		},
		{
			name: "Should return 400 when password does not meet strength requirements",
			setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
				// no repo calls expected - validation fails before any DB interaction
			},
			wantHTTPCode: http.StatusBadRequest,
			wantErrCode:  "WeakPassword",
			wantUser:     false,
			weakPassword: true,
		},
		{
			name: "Should return 409 Conflict when phone is already registered",
			setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
				repo.On("FindByPhone", mock.Anything, validReq.Phone).Return(existingUser, nil).Once()
			},
			wantHTTPCode:  http.StatusConflict,
			wantErrCode:   "DuplicatePhone",
			wantUser:      false,
			wantEmptyHash: false,
		},
		{
			name: "Should return 409 Conflict when email already exists",
			setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
				repo.On("FindByPhone", mock.Anything, validReq.Phone).Return(nil, nil).Once()
				repo.On("FindByEmail", mock.Anything, validReq.Email).Return(existingUser, nil).Once()
			},
			wantHTTPCode:  http.StatusConflict,
			wantErrCode:   "DuplicateEmail",
			wantUser:      false,
			wantEmptyHash: false,
		},
		{
			name: "Should return 500 when FindByPhone returns a database error",
			setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
				repo.On("FindByPhone", mock.Anything, validReq.Phone).Return(nil, dbErr).Once()
			},
			wantHTTPCode: http.StatusInternalServerError,
			wantErrCode:  "InternalServerError",
			wantUser:     false,
		},
		{
			name: "Should return 500 when FindByEmail returns a database error",
			setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
				repo.On("FindByPhone", mock.Anything, validReq.Phone).Return(nil, nil).Once()
				repo.On("FindByEmail", mock.Anything, validReq.Email).Return(nil, dbErr).Once()
			},
			wantHTTPCode: http.StatusInternalServerError,
			wantErrCode:  "InternalServerError",
			wantUser:     false,
		},
		// {
		// 	name: "Should Return Status Code 500 when bcrypt Returns An Error",
		// 	setupMock: func(repo *serviceMocks.MockAuthUserRepository) {
		// 		repo.On("FindByPhone", mock.Anything, validReq.Phone).Return(nil, nil).Once()
		// 		repo.On("FindByEmail", mock.Anything, validReq.Email).Return(nil, nil).Once()
		// 		repo.On()
		// 	},
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := serviceMocks.NewMockAuthUserRepository(t)
			tc.setupMock(repo)

			req := validReq
			if tc.weakPassword {
				req.Password = "23456"
			}

			svc := NewAuthService(repo)
			user, err := svc.Signup(context.Background(), req)

			if tc.wantUser {
				assert.Nil(t, err)
				assert.NotNil(t, user)
				if tc.wantEmptyHash {
					assert.Empty(t, user.PasswordHash, "password hash must never be returned from the service")
				}
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, user)
				appErr := err.(*ae.AppError)
				assert.Equal(t, tc.wantHTTPCode, appErr.HTTPCode())
				if tc.wantErrCode != "" {
					assert.Equal(t, tc.wantErrCode, appErr.Code)
				}
			}
		})
	}
}
