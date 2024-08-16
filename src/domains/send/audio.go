package send

import "mime/multipart"

type AudioRequest struct {
	Phone     string                `json:"phone" form:"phone"`
	Audio     *multipart.FileHeader `json:"Audio" form:"Audio"`
	TraceCode string                `json:"trace_code" form:"trace_code"`
}
