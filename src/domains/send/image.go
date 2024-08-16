package send

type ImageRequest struct {
	Phone     string `json:"phone" form:"phone"`
	Caption   string `json:"caption" form:"caption"`
	Image     string `json:"image" form:"image"`
	ViewOnce  bool   `json:"view_once" form:"view_once"`
	Compress  bool   `json:"compress"`
	TraceCode string `json:"trace_code" form:"trace_code"`
}
