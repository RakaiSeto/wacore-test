package send

type LinkRequest struct {
	Phone     string `json:"phone" form:"phone"`
	Caption   string `json:"caption"`
	Link      string `json:"link"`
	TraceCode string `json:"trace_code" form:"trace_code"`
}
