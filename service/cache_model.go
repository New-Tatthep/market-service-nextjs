package service

type AuthConfig struct {
	Concurrent                 int64 `json:"concurrent"`
	SessionDuration            int64 `json:"session_duration"`
	MaxLoginFail               int64 `json:"max_login_fail"`
	DisplaySessionTimeoutAlert int64 `json:"display_session_timeout_alert"`
}

type ActiveDirectorySetting struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TenantId     string `json:"tenant_id"`
	DomainName   string `json:"domain_name"`
}

type PasswordPolicySetting struct {
	Number      bool  `json:"number"`
	CharLower   bool  `json:"char_lower"`
	CharUpper   bool  `json:"char_upper"`
	LengthMin   int64 `json:"length_min"`
	LengthMax   int64 `json:"length_max"`
	CharSpecial bool  `json:"char_special"`
}
