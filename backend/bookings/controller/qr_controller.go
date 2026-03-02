package controller

import (
	"context"
	"net/http"
	"skyfox/bookings/dto/request"
	"skyfox/bookings/dto/response"
	"skyfox/common/logger"

	ae "skyfox/error"

	"github.com/gin-gonic/gin"
)

type QRService interface {
	GenerateQR(ctx context.Context, bookingID string) (string, error)
}

type QRController struct {
	qrService QRService
}

func NewQRController(qrService QRService) *QRController {
	return &QRController{qrService: qrService}
}

// GetQRCode godoc
//
//	@Summary		Get QR Code for a Booking
//	@Description	Returns a base64-encoded PNG QR code (data URL) for the given booking.
//	@Description	The QR payload is HMAC-SHA256 signed and contains booking, show, and customer identifiers.
//	@Description	The result is cached in the database; subsequent calls for the same booking return the stored code instantly.
//	@Tags			QR Code
//	@Produce		json
//	@Param			bookingId	path		string						true	"Booking UUID (must be a valid UUID v4)"	example("550e8400-e29b-41d4-a716-446655440000")
//	@Success		200			{object}	response.QRCodeResponse	"QR code data URL"//	@Failure		400			{object}	ae.AppError					"Invalid booking UUID"//	@Failure		404			{object}	ae.AppError				"Booking not found"
//	@Failure		500			{object}	ae.AppError				"Internal server error"
//	@Router			/api/v1/bookings/{bookingId}/qr [GET]
func (qc *QRController) GetQRCode(c *gin.Context) {
	var req request.QRCodeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, ae.BadRequestError("InvalidBookingId", "bookingId must be a valid UUID", err))
		return
	}

	dataURL, err := qc.qrService.GenerateQR(c.Request.Context(), req.BookingID)
	if err != nil {
		appErr, ok := err.(*ae.AppError)
		if !ok {
			logger.Error("QRController: unexpected error type: %v", err)
			c.JSON(http.StatusInternalServerError,
				ae.InternalServerError("InternalServerError", "an unexpected error occurred", err))
			return
		}
		c.JSON(appErr.HTTPCode(), appErr)
		return
	}

	c.JSON(http.StatusOK, response.NewQRCodeResponse(dataURL))
}
