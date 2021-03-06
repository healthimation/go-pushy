package pushy

// PushOptions holds all the optional settings for a push notification
type PushOptions struct {
	TimeToLive       *int64               `json:"time_to_live"`
	Notification     *NotificationOptions `json:"notification"`
	ContentAvailable *bool                `json:"content_available"`
	MutableContent   *bool                `json:"mutable_content"`
}

// NotificationOptions contains the iOS optional arguments for sending push notifications
type NotificationOptions struct {
	Body                   *string  `json:"body,omitempty"`
	Badge                  *int64   `json:"badge,omitempty"`
	Sound                  *string  `json:"sound,omitempty"`
	Title                  *string  `json:"title,omitempty"`
	Category               *string  `json:"category,omitempty"`
	LocalizationKey        *string  `json:"loc_key,omitempty"`
	LocalizationArgs       []string `json:"loc_args,omitempty"`
	TitleLocalizationKey   *string  `json:"title_loc_key,omitempty"`
	TitleLocalizationArges []string `json:"title_loc_args,omitempty"`
}
