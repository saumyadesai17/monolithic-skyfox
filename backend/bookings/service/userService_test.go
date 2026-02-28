package service

import (
	"context"
	"skyfox/_mocks/repomocks"
	"skyfox/bookings/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService(t *testing.T) {

	expected := model.User{Id: 1, Username: "john", Password: "john"}

	t.Run("UserDetails", func(t *testing.T) {
		repo := repomocks.UserRepository{}
		repo.On("FindByUsername", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).
			Return(expected, nil).
			Once()

		service := NewUserService(&repo)

		got, err := service.UserDetails(context.Background(), "")

		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})
}
