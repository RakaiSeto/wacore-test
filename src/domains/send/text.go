package send

type MessageRequest struct {
	Phone          string  `json:"phone" form:"phone"`
	Message        string  `json:"message" form:"message"`
	TraceCode      string  `json:"trace_code" form:"trace_code"`
	ReplyMessageID *string `json:"reply_message_id" form:"reply_message_id"`
}
