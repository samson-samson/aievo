package feedback

type FeedbackInfo struct {
	Type  FeedbackType `json:"type"`
	Msg   string       `json:"msg"`
	Token int          `json:"token"`
}
