package dtos

// SendNotificationRequest is the body for manually sending a notification (staff)
// @Description Channel (email/sms), type, optional user_id, user_email (for email), user_phone (for sms), title, subject, message, optional token (e.g. reset link).
type SendNotificationRequest struct {
	Channel   string `json:"channel" binding:"required"`   // email, sms
	Type      string `json:"type" binding:"required"`      // notification type
	UserID    string `json:"user_id"`                      // optional; for in-app targeting
	UserEmail string `json:"user_email"`                   // required for email
	UserPhone string `json:"user_phone"`                   // required for sms
	Title     string `json:"title"`
	Subject   string `json:"subject" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Token     string `json:"token"`                        // optional; e.g. for reset links
}
