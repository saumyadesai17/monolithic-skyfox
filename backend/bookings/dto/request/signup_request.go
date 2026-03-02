package request

type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"required,phoneNumber"`
	Password string `json:"password" binding:"required,passwordStrength"`
	Email    string `json:"email" binding:"omitempty,email"`
}
