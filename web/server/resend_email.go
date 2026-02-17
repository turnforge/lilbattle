package server

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
)

// ResendEmailSender sends emails via the Resend API
type ResendEmailSender struct {
	client   *resend.Client
	fromAddr string
}

func NewResendEmailSender(apiKey, fromAddr string) *ResendEmailSender {
	return &ResendEmailSender{
		client:   resend.NewClient(apiKey),
		fromAddr: fromAddr,
	}
}

func (r *ResendEmailSender) SendVerificationEmail(to string, verificationLink string) error {
	html := fmt.Sprintf(`<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Verify your email address</h2>
  <p>Please verify your email address by clicking the link below:</p>
  <p><a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #4f46e5; color: #ffffff; text-decoration: none; border-radius: 6px;">Verify Email</a></p>
  <p style="color: #666; font-size: 14px;">Or copy and paste this URL into your browser:<br>%s</p>
</div>`, verificationLink, verificationLink)

	params := &resend.SendEmailRequest{
		From:    r.fromAddr,
		To:      []string{to},
		Subject: "Verify your email address",
		Html:    html,
	}

	sent, err := r.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("resend: failed to send verification email to %s: %w", to, err)
	}
	log.Printf("Verification email sent to %s (id: %s)", to, sent.Id)
	return nil
}

func (r *ResendEmailSender) SendPasswordResetEmail(to string, resetLink string) error {
	html := fmt.Sprintf(`<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Reset your password</h2>
  <p>We received a request to reset your password. Click the link below to choose a new one:</p>
  <p><a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #4f46e5; color: #ffffff; text-decoration: none; border-radius: 6px;">Reset Password</a></p>
  <p style="color: #666; font-size: 14px;">Or copy and paste this URL into your browser:<br>%s</p>
  <p style="color: #666; font-size: 14px;">If you didn't request this, you can safely ignore this email.</p>
</div>`, resetLink, resetLink)

	params := &resend.SendEmailRequest{
		From:    r.fromAddr,
		To:      []string{to},
		Subject: "Reset your password",
		Html:    html,
	}

	sent, err := r.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("resend: failed to send password reset email to %s: %w", to, err)
	}
	log.Printf("Password reset email sent to %s (id: %s)", to, sent.Id)
	return nil
}
