package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skyfox/bookings/dto/response"
	servicemocks "skyfox/bookings/service/mocks"
	ae "skyfox/error"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const (
	testBookingID    = "550e8400-e29b-41d4-a716-446655440000"
	testQRDataURL    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg=="
	qrRoutePattern   = "/api/v1/bookings/:bookingId/qr"
	qrRouteFormatted = "/api/v1/bookings/" + testBookingID + "/qr"
)

func setupQRRouter(svc QRService) (*gin.Engine, *QRController) {
	router := gin.New()
	ctrl := NewQRController(svc)
	router.GET(qrRoutePattern, ctrl.GetQRCode)
	return router, ctrl
}

func TestQRController_GetQRCode(t *testing.T) {
	notFoundErr := ae.NotFoundError("BookingNotFound", "booking not found", nil)
	internalErr := ae.InternalServerError("InternalServerError", "something went wrong", nil)

	tests := []struct {
		name         string
		bookingID    string
		setupMock    func(svc *servicemocks.MockQRService)
		wantHTTPCode int
		wantBody     interface{}
	}{
		{
			name:      "Should return 400 when bookingId is not a valid UUID",
			bookingID: "not-a-uuid",
			setupMock: func(svc *servicemocks.MockQRService) {
				// service must NOT be called for an invalid path param
			},
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:      "Should return 200 with qrCodeUrl when booking exists",
			bookingID: testBookingID,
			setupMock: func(svc *servicemocks.MockQRService) {
				svc.On("GenerateQR", mock.Anything, testBookingID).
					Return(testQRDataURL, nil).Once()
			},
			wantHTTPCode: http.StatusOK,
			wantBody:     &response.QRCodeResponse{QRCodeURL: testQRDataURL},
		},
		{
			name:      "Should return 404 when booking is not found",
			bookingID: testBookingID,
			setupMock: func(svc *servicemocks.MockQRService) {
				svc.On("GenerateQR", mock.Anything, testBookingID).
					Return("", notFoundErr).Once()
			},
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:      "Should return 500 when QR service returns an internal error",
			bookingID: testBookingID,
			setupMock: func(svc *servicemocks.MockQRService) {
				svc.On("GenerateQR", mock.Anything, testBookingID).
					Return("", internalErr).Once()
			},
			wantHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svcMock := servicemocks.NewMockQRService(t)
			tc.setupMock(svcMock)

			router, _ := setupQRRouter(svcMock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/bookings/"+tc.bookingID+"/qr", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantHTTPCode, w.Code)

			if tc.wantBody != nil {
				var got response.QRCodeResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tc.wantBody, &got)
			}

			if tc.wantHTTPCode >= 400 {
				var errBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errBody)
				assert.NoError(t, err)
				assert.NotEmpty(t, errBody["Code"], "expected non-empty error Code in response")
			}
		})
	}
}
