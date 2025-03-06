package handlers

import (
	"encoding/json"
	"fmt"
	"lgm8-notification-service/internal/email"
	"log"
)

// UserVerificationEmailMessage represents the structure of the message
type UserVerificationEmailMessage struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// UserVerificationEmailHandler processes messages for user email verification
type UserVerificationEmailHandler struct {
	EmailSender *email.EmailSender
}

// NewUserVerificationEmailHandler creates a new instance of the handler
func NewUserVerificationEmailHandler(emailSender *email.EmailSender) *UserVerificationEmailHandler {
	return &UserVerificationEmailHandler{EmailSender: emailSender}
}

// HandleMessage processes incoming messages
func (h *UserVerificationEmailHandler) HandleMessage(body []byte) error {
	var msg UserVerificationEmailMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}

	// Construct verification email
	subject := "Verify Your Account"
	verificationLink := fmt.Sprintf("http://localhost/api/auth/v1/users/verify?token=%s", msg.Token)
	bodyText := fmt.Sprintf("Hello %s,\n\nPlease verify your account by clicking on the link below:\n%s\n\nThank you!", msg.Username, verificationLink)

	// Send the email
	err := h.EmailSender.SendEmail(msg.Email, subject, bodyText)
	if err != nil {
		log.Printf("Failed to send verification email to [%s]: [%s]", msg.Email, err)
		return err
	}

	log.Printf("Verification email sent to: [%s]", msg.Email)
	return nil
}
