package api

import (
	"context"
	"log"
)

type UserInterface interface {
	UserPaid(email string) bool
	UserExists(email string) bool
}

func (config *Config) UserExists(email string) bool {
	_, err := config.DB.GetUserByEmail(context.Background(), email)
	return err == nil
}

func (config *Config) UserPaid(email string) bool {
	charges, err := config.DB.GetChargesByUserEmail(context.Background(), email)

	if err != nil {
		log.Printf("UserPaid: Error getting charges for user %s: %v", email, err)
		return false
	}

	if len(charges) == 0 {
		log.Printf("UserPaid: No charges found for user %s", email)
		return false
	}

	if charges[0].Amount < int32(paymentValueTreshold) {
		log.Printf("UserPaid: User %s has not paid or amount is less than 2000", email)
		return false
	}

	return true
}

func (*MockConfig) UserPaid(_ string) bool {
	return true
}

func (*MockConfig) UserExists(_ string) bool {
	return true
}
