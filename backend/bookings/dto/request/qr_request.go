package request

// QRCodeRequest holds the path parameters for GET /api/v1/bookings/:bookingId/qr.
type QRCodeRequest struct {
	BookingID string `uri:"bookingId" binding:"required,uuid"`
}
