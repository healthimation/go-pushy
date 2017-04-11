package pushy

type pushRequest struct {
	To               string               `json:"to,omitempty"`
	Data             interface{}          `json:"data"`
	Tokens           []string             `json:"tokens,omitempty"`
	TimeToLive       *int64               `json:"time_to_live,omitempty"`
	Notification     *NotificationOptions `json:"notification,omitempty"`
	ContentAvailable *bool                `json:"content_avaiable,omitempty"`
	MutableContent   *bool                `json:"mutable_content,omitempty"`
}

type response struct {
	ID      *string `json:"id"`
	Success *bool   `json:"success"`
	Error   *string `json:"error"`
}
