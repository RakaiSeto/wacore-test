package send

type VideoRequest struct {
	Phone     string `json:"phone" form:"phone"`
	Caption   string `json:"caption" form:"caption"`
	Video     string `json:"video" form:"video"`
	ViewOnce  bool   `json:"view_once" form:"view_once"`
	Compress  bool   `json:"compress"`
	TraceCode string `json:"trace_code" form:"trace_code"`
}
