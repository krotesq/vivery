package mediamtx

type authDTO struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	IP        string `json:"ip"`
	Action    string `json:"action"`
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	ID        string `json:"id"`
	Query     string `json:"query"`
	UserAgent string `json:"userAgent"`
}
