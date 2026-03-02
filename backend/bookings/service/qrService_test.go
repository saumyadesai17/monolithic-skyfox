package service

import (
	"context"
	"net/http"
	"skyfox/bookings/model"
	servicemocks "skyfox/bookings/service/mocks"
	ae "skyfox/error"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testQRSecret = "test-secret-key"

var validBooking = &model.BookingRecord{
	Id:         "booking-uuid-1234",
	ShowId:     "show-uuid-5678",
	CustomerId: "customer-uuid-9012",
}

var validShow = &model.ShowRecord{
	Id:          "show-uuid-5678",
	MovieImdbId: "tt1234567",
	TheatreId:   "theatre-uuid-aaaa",
	ScreenId:    "screen-uuid-bbbb",
	StartTime:   time.Date(2026, 3, 10, 18, 0, 0, 0, time.UTC),
	EndTime:     time.Date(2026, 3, 10, 20, 30, 0, 0, time.UTC),
}

var validSeatIDs = []string{"seat-uuid-0001", "seat-uuid-0002"}

func TestQRService_GenerateQR(t *testing.T) {
	dbErr := ae.InternalServerError("InternalServerError", "something went wrong", nil)

	tests := []struct {
		name         string
		setupMock    func(repo *servicemocks.MockQRBookingRepository)
		wantHTTPCode int
		wantErrCode  string
		wantDataURL  bool
		wantCacheHit bool
	}{
		{
			name: "Should generate and store QR when booking exists with show and seat data",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(validBooking, nil).Once()
				repo.On("FindShowByID", mock.Anything, validBooking.ShowId).
					Return(validShow, nil).Once()
				repo.On("FindSeatsByBookingID", mock.Anything, validBooking.Id).
					Return(validSeatIDs, nil).Once()
				repo.On("UpdateQRCodeURL", mock.Anything, validBooking.Id,
					mock.MatchedBy(func(url string) bool {
						return strings.HasPrefix(url, "data:image/png;base64,")
					})).Return(nil).Once()
			},
			wantDataURL: true,
		},
		{
			name: "Should return cached QR without any further DB calls when qr_code_url is already set",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				cachedBooking := &model.BookingRecord{
					Id:         validBooking.Id,
					ShowId:     validBooking.ShowId,
					CustomerId: validBooking.CustomerId,
					QRCodeURL:  "data:image/png;base64,alreadycached==",
				}
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(cachedBooking, nil).Once()
				// FindShowByID, FindSeatsByBookingID, UpdateQRCodeURL must NOT be called
			},
			wantDataURL:  true,
			wantCacheHit: true,
		},
		{
			name: "Should return 404 when booking is not found",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(nil, nil).Once()
			},
			wantHTTPCode: http.StatusNotFound,
			wantErrCode:  "BookingNotFound",
		},
		{
			name: "Should return 500 when FindBookingByID returns a database error",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(nil, dbErr).Once()
			},
			wantHTTPCode: http.StatusInternalServerError,
			wantErrCode:  "InternalServerError",
		},
		{
			name: "Should return 404 when show is not found for the booking",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(validBooking, nil).Once()
				repo.On("FindShowByID", mock.Anything, validBooking.ShowId).
					Return(nil, nil).Once()
			},
			wantHTTPCode: http.StatusNotFound,
			wantErrCode:  "ShowNotFound",
		},
		{
			name: "Should return 500 when FindShowByID returns a database error",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(validBooking, nil).Once()
				repo.On("FindShowByID", mock.Anything, validBooking.ShowId).
					Return(nil, dbErr).Once()
			},
			wantHTTPCode: http.StatusInternalServerError,
			wantErrCode:  "InternalServerError",
		},
		{
			name: "Should return 500 when FindSeatsByBookingID returns a database error",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(validBooking, nil).Once()
				repo.On("FindShowByID", mock.Anything, validBooking.ShowId).
					Return(validShow, nil).Once()
				repo.On("FindSeatsByBookingID", mock.Anything, validBooking.Id).
					Return(nil, dbErr).Once()
			},
			wantHTTPCode: http.StatusInternalServerError,
			wantErrCode:  "InternalServerError",
		},
		{
			name: "Should return 500 when UpdateQRCodeURL returns a database error",
			setupMock: func(repo *servicemocks.MockQRBookingRepository) {
				repo.On("FindBookingByID", mock.Anything, validBooking.Id).
					Return(validBooking, nil).Once()
				repo.On("FindShowByID", mock.Anything, validBooking.ShowId).
					Return(validShow, nil).Once()
				repo.On("FindSeatsByBookingID", mock.Anything, validBooking.Id).
					Return(validSeatIDs, nil).Once()
				repo.On("UpdateQRCodeURL", mock.Anything, validBooking.Id, mock.AnythingOfType("string")).
					Return(dbErr).Once()
			},
			wantHTTPCode: http.StatusInternalServerError,
			wantErrCode:  "InternalServerError",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := servicemocks.NewMockQRBookingRepository(t)
			tc.setupMock(repo)

			svc := NewQRService(repo, testQRSecret)
			dataURL, err := svc.GenerateQR(context.Background(), validBooking.Id)

			if tc.wantDataURL {
				assert.Nil(t, err)
				assert.True(t, strings.HasPrefix(dataURL, "data:image/png;base64,"),
					"expected a data URL, got: %s", dataURL)
			} else {
				assert.Empty(t, dataURL)
				assert.NotNil(t, err)
				appErr := err.(*ae.AppError)
				assert.Equal(t, tc.wantHTTPCode, appErr.HTTPCode())
				if tc.wantErrCode != "" {
					assert.Equal(t, tc.wantErrCode, appErr.Code)
				}
			}
		})
	}
}
