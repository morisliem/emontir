package validator

import (
	"e-montir/pkg/date"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
	"unicode"
)

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 128 {
		return fmt.Errorf("name cannot exceed 128 character")
	}

	return nil
}

func ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if len(email) > 128 {
		return fmt.Errorf("email cannot exceed 128 characters")
	}
	if !strings.Contains(email, ".") {
		return errors.New("email is missing . ")
	}

	if !strings.Contains(email, "@") {
		return errors.New("email is missing @")
	}

	i := strings.Index(email, "@")
	host := email[i+1:]

	_, err := net.LookupMX(host)
	if err != nil {
		return errors.New("could not find email's domain server")
	}
	return nil
}

func ValidateID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if len(id) > 36 {
		return fmt.Errorf("id cannot exceed 128 characters")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) > 60 {
		return fmt.Errorf("password cannot exceed 60 characters")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be more than 8 characters")
	}

	hasUpperCase := false
	hasLowerCase := false
	hasNumber := false
	hasSpecialCharacter := false

	for _, v := range password {
		if unicode.IsNumber(v) {
			hasNumber = true
		}
		if unicode.IsLower(v) {
			hasLowerCase = true
		}
		if unicode.IsUpper(v) {
			hasUpperCase = true
		}
		if unicode.IsPunct(v) {
			hasSpecialCharacter = true
		}
	}

	if hasUpperCase && hasLowerCase && hasNumber && hasSpecialCharacter {
		return nil
	}
	if !hasLowerCase {
		return fmt.Errorf("must have a lowercase character")
	}
	if !hasUpperCase {
		return fmt.Errorf("must have a uppercase character")
	}
	if !hasNumber {
		return fmt.Errorf("must have a numerical character")
	}
	if !hasSpecialCharacter {
		return fmt.Errorf("must have a special character")
	}
	return nil
}

func ValidatePage(page string) error {
	if strings.TrimSpace(page) == "" {
		return fmt.Errorf("page cannot be empty")
	}

	return nil
}

func ValidateLimit(limit string) error {
	if strings.TrimSpace(limit) == "" {
		return fmt.Errorf("limit cannot be empty")
	}

	return nil
}

func ValidateKeyword(keyword string) error {
	if strings.TrimSpace(keyword) == "" {
		return fmt.Errorf("keyword cannot be empty")
	}
	if len(keyword) > 128 {
		return fmt.Errorf("keyword cannot exceed 128 characters")
	}
	return nil
}

func ValidateDate(dateIn string) (string, error) {
	if strings.TrimSpace(dateIn) == "" {
		return "", fmt.Errorf("date cannot be empty")
	}

	// yyyy-mm-dd
	d, err := time.Parse(date.Format, dateIn)
	if err != nil {
		return "", fmt.Errorf("wrong date format")
	}
	return d.Local().Format(date.Format), nil
}

func ValidateTime(timeIn string) error {
	timeslot := map[string]bool{
		"07:00-10:00": true,
		"10:00-14:00": true,
		"14:00-18:00": true,
	}

	if strings.TrimSpace(timeIn) == "" {
		return fmt.Errorf("time cannot be empty")
	}

	if !timeslot[timeIn] {
		return fmt.Errorf("time must be 07:00-10:00 or 10:00-14:00 or 14:00-18:00")
	}

	return nil
}

func ValidateServiceID(service string) error {
	if strings.TrimSpace(service) == "" {
		return fmt.Errorf("service_id cannot be empty")
	}

	return nil
}

func ValidateLabel(label string) error {
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("label cannot be empty")
	}
	return nil
}

func ValidateAddress(address string) error {
	if strings.TrimSpace(address) == "" {
		return fmt.Errorf("address cannot be empty")
	}
	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	if strings.TrimSpace(phoneNumber) == "" {
		return fmt.Errorf("phone_number cannot be empty")
	}
	return nil
}

func ValidateRecipientName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("recipient_name cannot be empty")
	}
	return nil
}

func ValidatePaymentToken(token string) error {
	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("token cannot be empty")
	}
	return nil
}

func ValidateOrderID(orderID string) error {
	if strings.TrimSpace(orderID) == "" {
		return fmt.Errorf("order_id cannot be empty")
	}
	return nil
}

func ValidateTotalPrice(totalPrice string) error {
	if strings.TrimSpace(totalPrice) == "" {
		return fmt.Errorf("total_price cannot be empty")
	}
	return nil
}

func ValidateRating(rating string) error {
	if strings.TrimSpace(rating) == "" {
		return fmt.Errorf("rating cannot be empty")
	}
	return nil
}

func ValidateFeedback(feedback string) error {
	if strings.TrimSpace(feedback) != "" {
		if len(feedback) < 10 || len(feedback) > 300 {
			return fmt.Errorf("feedback must be between 10 and 300 characters")
		}
	}
	return nil
}
