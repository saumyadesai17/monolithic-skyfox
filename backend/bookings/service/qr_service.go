package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"skyfox/bookings/model"
	"skyfox/common/logger"
	ae "skyfox/error"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

// QRBookingRepository is the repository slice that QRService needs.
type QRBookingRepository interface {
	FindBookingByID(ctx context.Context, id string) (*model.BookingRecord, error)
	FindShowByID(ctx context.Context, id string) (*model.ShowRecord, error)
	FindSeatsByBookingID(ctx context.Context, bookingID string) ([]string, error)
	UpdateQRCodeURL(ctx context.Context, id string, qrURL string) error
}

// QRService generates and caches QR codes for bookings.
type QRService interface {
	GenerateQR(ctx context.Context, bookingID string) (string, error)
}

// qrPayload is the data encoded inside the QR code.
// Fields are kept short to minimise QR complexity.
// The payload is self-contained: a scanner can verify authenticity and display
// booking details (movie, show time, seats) without any additional DB call.
type qrPayload struct {
	BookingID  string   `json:"bid"`   // booking UUID
	ShowID     string   `json:"sid"`   // show UUID
	CustomerID string   `json:"cid"`   // customer UUID
	MovieID    string   `json:"mid"`   // movie IMDB ID
	TheatreID  string   `json:"tid"`   // theatre UUID
	ShowStart  string   `json:"sst"`   // show start time (RFC3339)
	ShowEnd    string   `json:"set"`   // show end time (RFC3339)
	Seats      []string `json:"seats"` // seat UUIDs for this booking
	IssuedAt   int64    `json:"iat"`   // Unix timestamp of QR generation
}

type qrService struct {
	repo   QRBookingRepository
	secret string // HMAC-SHA256 signing key from config
}

func NewQRService(repo QRBookingRepository, secret string) QRService {
	return &qrService{repo: repo, secret: secret}
}

// GenerateQR returns a base64 data URL for the booking's QR code.
// If the booking already has a stored QR, it is returned directly (cache hit).
// Otherwise the QR is generated with a self-contained signed payload,
// stored in the DB, and returned.
func (s *qrService) GenerateQR(ctx context.Context, bookingID string) (string, error) {
	booking, err := s.repo.FindBookingByID(ctx, bookingID)
	if err != nil {
		logger.Error("qr: error fetching booking %s: %v", bookingID, err)
		return "", err
	}
	if booking == nil {
		return "", ae.NotFoundError("BookingNotFound", "booking not found", fmt.Errorf("booking %s not found", bookingID))
	}

	// Cache hit — QR already generated for this booking
	if booking.QRCodeURL != "" {
		return booking.QRCodeURL, nil
	}

	// Fetch the associated show for date/time/movie/theatre details
	show, err := s.repo.FindShowByID(ctx, booking.ShowId)
	if err != nil {
		logger.Error("qr: error fetching show %s: %v", booking.ShowId, err)
		return "", err
	}
	if show == nil {
		return "", ae.NotFoundError("ShowNotFound", "show not found for booking", fmt.Errorf("show %s not found", booking.ShowId))
	}

	// Fetch seat IDs reserved under this booking
	seatIDs, err := s.repo.FindSeatsByBookingID(ctx, bookingID)
	if err != nil {
		logger.Error("qr: error fetching seats for booking %s: %v", bookingID, err)
		return "", err
	}

	// Build the signed QR data string
	qrData, err := s.buildSignedPayload(booking, show, seatIDs)
	if err != nil {
		logger.Error("qr: failed to build signed payload for booking %s: %v", bookingID, err)
		return "", ae.InternalServerError("InternalServerError", "failed to generate QR code", err)
	}

	// Generate QR code PNG (256×256 px, medium error correction)
	pngBytes, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		logger.Error("qr: failed to encode QR PNG for booking %s: %v", bookingID, err)
		return "", ae.InternalServerError("InternalServerError", "failed to generate QR code", err)
	}

	// Encode as a browser-ready data URL
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)

	// Persist so future requests get a cache hit
	if err := s.repo.UpdateQRCodeURL(ctx, bookingID, dataURL); err != nil {
		logger.Error("qr: failed to store QR URL for booking %s: %v", bookingID, err)
		return "", err
	}

	return dataURL, nil
}

// buildSignedPayload creates the string that the QR code will encode.
//
// Format: base64(json_payload) + "." + hex(HMAC-SHA256(secret, base64_payload))
//
// The scanner can verify authenticity by recomputing the HMAC without a DB call.
// The payload is self-contained with movie, show time, theatre, and seat data.
func (s *qrService) buildSignedPayload(booking *model.BookingRecord, show *model.ShowRecord, seatIDs []string) (string, error) {
	payload := qrPayload{
		BookingID:  booking.Id,
		ShowID:     booking.ShowId,
		CustomerID: booking.CustomerId,
		MovieID:    show.MovieImdbId,
		TheatreID:  show.TheatreId,
		ShowStart:  show.StartTime.UTC().Format(time.RFC3339),
		ShowEnd:    show.EndTime.UTC().Format(time.RFC3339),
		Seats:      seatIDs,
		IssuedAt:   time.Now().Unix(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(payloadJSON)

	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(encodedPayload))
	sig := hex.EncodeToString(mac.Sum(nil))

	return encodedPayload + "." + sig, nil
}
