package service

import (
	"context"
	"skyfox/bookings/model"
	servicemocks "skyfox/bookings/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService(t *testing.T) {
	validUser := model.User{Id: 1, Username: "john", Password: "john"}

	tests := []struct {
		name      string
		setupMock func(repo *servicemocks.MockUserRepository)
		wantUser  model.User
		wantErr   bool
	}{
		{
			name: "Should return user when username exists",
			setupMock: func(repo *servicemocks.MockUserRepository) {
				repo.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(validUser, nil).Once()
			},
			wantUser: validUser,
		},
		{
			name: "Should return empty user when username is not found",
			setupMock: func(repo *servicemocks.MockUserRepository) {
				repo.On("FindByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(model.User{}, nil).Once()
			},
			wantUser: model.User{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := servicemocks.NewMockUserRepository(t)
			tc.setupMock(repo)

			svc := NewUserService(repo)
			got, err := svc.UserDetails(context.Background(), "")

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantUser, got)
			}
		})
	}
}
