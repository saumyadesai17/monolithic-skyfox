package response

// QRCodeResponse is the response body for GET /api/v1/bookings/:bookingId/qr.
type QRCodeResponse struct {
	QRCodeURL string `json:"qrCodeUrl" example:"data:image/png;base64,iVBORw0KGgo="`
}

func NewQRCodeResponse(qrCodeURL string) *QRCodeResponse {
	return &QRCodeResponse{QRCodeURL: qrCodeURL}
}
