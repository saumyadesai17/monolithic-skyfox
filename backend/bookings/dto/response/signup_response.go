package response

type SignupResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func NewSignupResponse(id string, name string, phone string) *SignupResponse {
	return &SignupResponse{
		Id:    id,
		Name:  name,
		Phone: phone,
	}
}
