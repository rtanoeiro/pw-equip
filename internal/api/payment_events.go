package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"pw-equip-change/internal/database"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func (config *Config) HandlePaymentEvents(writer http.ResponseWriter, req *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	webhookSigningSecret := os.Getenv("STRIPE_WEBHOOK_SIGNING_SECRET")

	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(writer, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}
	if errUnmarshal := json.Unmarshal(payload, &event); errUnmarshal != nil {
		log.Printf("⚠️  Webhook error while parsing basic request. %s/n", errUnmarshal.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	signatureHeader := req.Header.Get("Stripe-Signature")
	event, errConstructEvent := webhook.ConstructEvent(payload, signatureHeader, webhookSigningSecret)
	if errConstructEvent != nil {
		log.Printf("⚠️  Webhook signature verification failed. %v\n", errConstructEvent)
		writer.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	switch event.Type {
	case "charge.succeeded":
		processChargeSucceededEvent(&event, writer, config)
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}

	writer.WriteHeader(http.StatusOK)
}

// When a charge succeeds, we need to add to the database both payment and charges.
// It also needs to create an user if it doesn't exist.
func processChargeSucceededEvent(event *stripe.Event, writer http.ResponseWriter, config *Config) {
	var charge stripe.Charge
	errUnmarshal := json.Unmarshal(event.Data.Raw, &charge)
	if errUnmarshal != nil {
		log.Printf("Error parsing webhook JSON: %v\n", errUnmarshal)
		writer.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("Successful charge for %d.", charge.Amount)
	errCharge := config.DB.AddCharge(context.Background(), database.AddChargeParams{
		UserEmail:       charge.BillingDetails.Email,
		PaymentIntentID: charge.PaymentIntent.ID,
		Amount:          int32(charge.Amount),
		Status:          *stripe.String(charge.Status),
	})
	if errCharge != nil {
		log.Printf("Error adding charge to database: %v", errCharge)
	}
	log.Printf("HandlePaymentEvents: Successfully added charge to database")

	if errUser := config.UserExists(charge.BillingDetails.Email); !errUser {
		user := database.CreateUserParams{
			Email: charge.BillingDetails.Email,
			Hwid:  sql.NullString{Valid: false},
		}
		errCreate := config.DB.CreateUser(context.Background(), user)
		if errCreate != nil {
			log.Printf("HandlePaymentEvents: Error creating user %s: %v", charge.BillingDetails.Email, errCreate)
			writer.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("HandlePaymentEvents: Successfully created user %s", charge.BillingDetails.Email)
	}
}
