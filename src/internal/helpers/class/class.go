package class

import "fmt"

type UserValidate struct {
	Username string `json:"username"`
	Context  string `json:"context"`
	Birthday string `json:"birthday"`
}

type UserValidator struct {
	Username string `json:"username"`
	Birthday string `json:"birthday"`
}

type UsernameResponse struct {
	DidGenerateNewUsername bool     `json:"didGenerateNewUsername"`
	SuggestedUsernames     []string `json:"suggestedUsernames"`
}

type PasswordValidate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecureAuth struct {
	ClientPublicKey      string `json:"clientPublicKey"`
	ClientEpochTimestamp int64  `json:"clientEpochTimestamp"`
	ServerNonce          string `json:"serverNonce"`
	SaiSignature         string `json:"saiSignature"`
}

type SignupPayload struct {
	Username                   string      `json:"username"`
	Password                   string      `json:"password"`
	Birthday                   string      `json:"birthday"`
	Gender                     int         `json:"gender"`
	IsTosAgreementBoxChecked   bool        `json:"isTosAgreementBoxChecked"`
	Locale                     string      `json:"locale"`
	AgreementIds               []string    `json:"agreementIds"`
	SecureAuthenticationIntent *SecureAuth `json:"secureAuthenticationIntent"`
}

type ArkoseResponse struct {
	DataExchangeBlob string `json:"dataExchangeBlob"`
	UnifiedCaptchaId string `json:"unifiedCaptchaId"`
}

type ChallengeMetadata struct {
	UnifiedCaptchaId string `json:"unifiedCaptchaId"`
	CaptchaToken     string `json:"captchaToken"`
	ActionType       string `json:"actionType"`
}

type ChallengeResponse struct {
	ChallengeId       string `json:"challengeId"`
	ChallengeType     string `json:"challengeType"`
	ChallengeMetadata string `json:"challengeMetadata"`
}

type Config struct {
	Captcha  CaptchaConfig  `yaml:"settings_captcha"`
	Register RegisterConfig `yaml:"settings_register"`
}

type CaptchaConfig struct {
	Api_Key      string `yaml:"api_key"`
	Http_Version string `yaml:"http_version"`
	Solve_POW    bool   `yaml:"solve_pow"`
}

type RegisterConfig struct {
	Threads        int `yaml:"threads"`
	Limit_Accounts int `yaml:"limit_accounts"`
}

func (c *Config) Validate() error {
	if c.Register.Threads <= 0 || c.Register.Threads >= 1000 {
		return fmt.Errorf("Register.Thread must be > 0 <= 1000")
	}
	if c.Register.Limit_Accounts <= 0 {
		return fmt.Errorf("Register.Limit_Accounts must be > 0")
	}
	if c.Captcha.Api_Key == "" {
		return fmt.Errorf("Captcha.Api_Key must not be empty")
	}
	if c.Captcha.Http_Version == "" {
		return fmt.Errorf("Captcha.Http_Version must not be empty")
	}
	return nil
}
