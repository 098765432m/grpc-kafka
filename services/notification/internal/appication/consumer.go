package notification_appication

import (
	notification_domain "github.com/098765432m/grpc-kafka/notification/internal/domain"
	notification_infrastructure "github.com/098765432m/grpc-kafka/notification/internal/infrastructure"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	EmailSender *notification_infrastructure.EmailSender
}

func NewNotificationHandler(emailSender *notification_infrastructure.EmailSender) *NotificationHandler {
	return &NotificationHandler{
		EmailSender: emailSender,
	}
}

func (nh *NotificationHandler) HandleBookingCreated(event notification_domain.BookingCreatedEvent) {
	subject := ""
	body := ""

	if err := nh.EmailSender.SendEmail(event.UserEmail, subject, body); err != nil {
		zap.S().Infoln("Failed to send email: ", err)
	} else {
		zap.S().Infof("Email send to %s\n", event.UserEmail)
	}
}
