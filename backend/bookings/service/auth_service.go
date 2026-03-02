package service

import (
	"context"
	"fmt"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/model"
	"skyfox/common/logger"
	ae "skyfox/error"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// AuthUserRepository is the slice of the UserAccountRepository that AuthService needs.
// Keeping this as a local interface makes the service very easy to unit-test with mocks.
type AuthUserRepository interface {
	FindByPhone(ctx context.Context, phone string) (*model.UserAccount, error)
	FindByEmail(ctx context.Context, phone string) (*model.UserAccount, error)
	CreateUser(ctx context.Context, user *model.UserAccount) error
}

type authService struct {
	userRepo AuthUserRepository
}

func NewAuthService(userRepo AuthUserRepository) *authService {
	return &authService{userRepo: userRepo}
}

// isPasswordStrong enforces: min 8 chars, at least one uppercase letter.
// one digit, and one special character. This mirrors the HTTP binding validator
// and acs as a service-layer defence for non-HTTP callers.
func isPasswordStrong(password string) bool {
			if len(password) < 8 {
			return false
		}
		var hasUpper, hasDigit, hasSpecial bool
		for _, ch := range password {
			switch {
			case unicode.IsUpper(ch):
				hasUpper = true
			case unicode.IsDigit(ch):
				hasDigit = true
			case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
				hasSpecial = true
			}
		}
		return hasUpper && hasDigit && hasSpecial
}

// Signup validates the request, prevents duplicate phone numbers, hashes the password]
// with bcrypt, and persists the new user. Returns the created user
// (PasswordHash is empty - never put it on the struct returned here).
func (s *authService) Signup(ctx context.Context, req request.SignupRequest) (*model.UserAccount, error) {
	// Password strength check (service-layer defence - HTTP binding is the first layer)
	if !isPasswordStrong(req.Password) {
		return nil, ae.BadRequestError("WeakPassword", "password must be at least 8 characters and contain an uppercase letter, a digit, and a special character", fmt.Errorf("password failed strength check"))
	}

	// Duplicate phone check
	existingPhone, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		logger.Error("error checking phone existence: %v", err)
		return nil, err
	}
	if existingPhone != nil {
		logger.Error("signup attempt with duplicate phone: %s", req.Phone)
		return nil, ae.ConflictError("DuplicatePhone", "an account with this phone number already exists", fmt.Errorf("phone %s already exists", req.Phone))
	}

	// Duplicate email check
	existingEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("error checking email existence: %v", err)
		return nil, err
	}
	if existingEmail != nil {
		logger.Error("signup attempt with duplicate email: %s", req.Email)
		return nil, ae.ConflictError("DuplicateEmail", "an account with this email already exists", fmt.Errorf("email %s already exists", req.Email))
	}

	// Hash Password
	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if bcryptErr != nil {
		logger.Error("bcrypt error: %v", bcryptErr)
		return nil, ae.InternalServerError("InternalServerError", "something went wrong", fmt.Errorf("bcrypt failed: %w", bcryptErr))
	}

	// Build and persist the user
	user := &model.UserAccount{
		Name: req.Name,
		Phone: req.Phone,
		Email: req.Email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		logger.Error("error creating user: %v", err)
		return nil, err
	}

	// Clear the hash before returning so it never leaks up the stack
	user.PasswordHash = ""
	return user, nil
}