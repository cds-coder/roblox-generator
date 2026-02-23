package register

import (
	"RobloxRegister/src/internal/helpers/class"
	"RobloxRegister/src/internal/helpers/funcaptcha"
	"RobloxRegister/src/internal/helpers/roblox_profile"
	"RobloxRegister/src/internal/helpers/utils"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Noooste/azuretls-client"
	fhttp "github.com/Noooste/fhttp"
	"github.com/google/uuid"
)

type Container struct {
	HttpClient  *azuretls.Session
	Proxy       string
	Cookies     [][]string
	CapConfig   class.CaptchaConfig
	User        string
	Password    string
	Birthday    string
	Gender      int
	XCsrfToken  string
	Transparent string
}

var (
	order_cookies = []string{"rbx-ip2", "RBXEventTrackerV2", "GuestData", "RBXPaymentsFlowContext", "RBXcb"}
	regexCsrf     = regexp.MustCompile(`<meta\s+name=["']csrf-token["']\s+data-token=["']([^"']+)["']`)
)

const (
	maxRetries      = 3
	userAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36"
	sec_ch_ua       = `"Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"`
	accept_language = "en-US,en;q=0.9"
)

func RegistrationProcess(CaptchaConfig class.CaptchaConfig, worker_id int) bool {

	RegistrationContainer := &Container{
		Proxy:       utils.GetProxy(),
		Transparent: utils.GetTransparent(),
		HttpClient:  azuretls.NewSession(),
		CapConfig:   CaptchaConfig,
		User:        roblox_profile.GetUsername(),
		Password:    roblox_profile.GetPassword(),
		Birthday:    roblox_profile.GetBirthDay(),
		Gender:      roblox_profile.GetGender(),
	}

	utils.Output("INFO", fmt.Sprintf("Start generate - %s", RegistrationContainer.User))

	if err := RegistrationContainer.SetHttpSession(); err != nil {
		utils.Output("FAILED", fmt.Sprintf("%s - %s", RegistrationContainer.User, err))
		return false
	}

	if err := RegistrationContainer.BeforeSignUp(); err != nil {
		utils.Output("FAILED", fmt.Sprintf("%s - %s", RegistrationContainer.User, err))
		return false
	}

	if err := RegistrationContainer.SignUp(); err != nil {
		utils.Output("FAILED", fmt.Sprintf("%s - %s", RegistrationContainer.User, err))
		return false
	}

	return true

}

func (g *Container) SetHttpSession() error {

	if err := g.HttpClient.ApplyHTTP2("1:65536;2:0;4:6291456;6:262144|15663105|0|m,a,s,p"); err != nil {
		return fmt.Errorf("failed set HTTP2")
	}

	if err := g.HttpClient.SetProxy(g.Proxy); err != nil {
		return fmt.Errorf("failed set Proxy")
	}

	return nil

}

func (g *Container) DoRequest(method, url string, body []byte) (*azuretls.Response, error) {

	var resp *azuretls.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {

		req := &azuretls.Request{
			Method:   method,
			Url:      url,
			Body:     body,
			TimeOut:  10 * time.Second,
			NoCookie: true,
		}

		resp, err = g.HttpClient.Do(req)
		if err == nil {
			return resp, nil
		}
	}

	return nil, fmt.Errorf("failed send request")
}

func parseCookies(headers fhttp.Header) map[string]string {
	cookies := make(map[string]string)

	setCookies := headers["Set-Cookie"]
	for _, c := range setCookies {
		parts := strings.Split(c, ";")
		if len(parts) > 0 {
			kv := strings.SplitN(parts[0], "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				cookies[key] = value
			}
		}
	}

	return cookies
}

func (g *Container) BeforeSignUp() error {

	g.HttpClient.OrderedHeaders = azuretls.OrderedHeaders{
		{"upgrade-insecure-requests", "1"},
		{"user-agent", userAgent},
		{"accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		{"sec-fetch-site", "none"},
		{"sec-fetch-mode", "navigate"},
		{"sec-fetch-user", "?1"},
		{"sec-fetch-dest", "document"},
		{"sec-ch-ua", sec_ch_ua},
		{"sec-ch-ua-mobile", "?0"},
		{"sec-ch-ua-platform", "\"Windows\""},
		{"accept-encoding", "gzip, deflate, br, zstd"},
		{"accept-language", accept_language},
		{"priority", "u=0, i"},
	}

	response, err := g.DoRequest("GET", "https://www.roblox.com/", nil)

	if err != nil {
		return fmt.Errorf("getPage error")
	}

	match := regexCsrf.FindStringSubmatch(string(response.Body))
	if len(match) > 1 {
		rawToken := match[1]
		g.XCsrfToken = html.UnescapeString(rawToken)
	} else {
		return fmt.Errorf("x-csrf-token not found")
	}

	cookies := parseCookies(response.Header)

	if len(cookies) != 0 {

		cookies["RBXPaymentsFlowContext"] = fmt.Sprintf("%s,", uuid.New())
		cookies["RBXcb"] = "RBXViralAcquisition%3Dfalse%26RBXSource%3Dfalse%26GoogleAnalytics%3Dfalse"

		for _, key := range order_cookies {
			value, ok := cookies[key]
			if !ok {
				continue
			}
			g.Cookies = append(g.Cookies, []string{"cookie", key + "=" + value})
		}

	} else {
		return fmt.Errorf("cookies empty")
	}

	body := &class.UserValidate{
		Username: g.User,
		Context:  "Signup",
		Birthday: g.Birthday,
	}

	data, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("failed json")
	}

	g.HttpClient.OrderedHeaders = azuretls.OrderedHeaders{
		{"content-length", strconv.Itoa(len(data))},
		{"sec-ch-ua-platform", `"Windows"`},
		{"x-csrf-token", g.XCsrfToken},
		{"sec-ch-ua", sec_ch_ua},
		{"sec-ch-ua-mobile", "?0"},
		{"traceparent", g.Transparent},
		{"user-agent", userAgent},
		{"accept", "application/json, text/plain, */*"},
		{"content-type", "application/json;charset=UTF-8"},
		{"origin", "https://www.roblox.com"},
		{"sec-fetch-site", "same-site"},
		{"sec-fetch-mode", "cors"},
		{"sec-fetch-dest", "empty"},
		{"referer", "https://www.roblox.com/"},
		{"accept-encoding", "gzip, deflate, br, zstd"},
		{"accept-language", accept_language},
		{"priority", "u=1, i"},
	}

	response, err = g.DoRequest("POST", "https://auth.roblox.com/v1/usernames/validate", data)

	if err != nil {
		return fmt.Errorf("userValidate request error")
	}

	if string(response.Body) == `{"code":0,"message":"Token Validation Failed"}` {

		g.XCsrfToken = string(response.Header.Get("X-Csrf-Token"))

		g.HttpClient.OrderedHeaders.Set("x-csrf-token", g.XCsrfToken)

		response, err = g.DoRequest("POST", "https://auth.roblox.com/v1/usernames/validate", data)

		if err != nil {
			return fmt.Errorf("userValidate request error")
		}

	}

	if string(response.Body) != `{"code":0,"message":"Username is valid"}` {

		body := &class.UserValidator{
			Username: g.User,
			Birthday: g.Birthday,
		}

		data, err := json.Marshal(body)

		if err != nil {
			return fmt.Errorf("failed json")
		}

		g.HttpClient.OrderedHeaders = azuretls.OrderedHeaders{
			{"content-length", strconv.Itoa(len(data))},
			{"sec-ch-ua-platform", `Windows"`},
			{"x-csrf-token", g.XCsrfToken},
			{"sec-ch-ua", sec_ch_ua},
			{"sec-ch-ua-mobile", "?0"},
			{"traceparent", g.Transparent},
			{"user-agent", userAgent},
			{"accept", "application/json, text/plain, */*"},
			{"content-type", "application/json;charset=UTF-8"},
			{"origin", "https://www.roblox.com"},
			{"sec-fetch-site", "same-site"},
			{"sec-fetch-mode", "cors"},
			{"sec-fetch-dest", "empty"},
			{"referer", "https://www.roblox.com/"},
			{"accept-encoding", "gzip, deflate, br, zstd"},
			{"accept-language", accept_language},
			{"priority", "u=1, i"},
		}

		body.Username = roblox_profile.GetUsername()

		data, err = json.Marshal(body)

		if err != nil {
			return fmt.Errorf("failed json")
		}

		response, err = g.DoRequest("POST", "https://auth.roblox.com/v1/validators/username", data)

		if err != nil {
			return fmt.Errorf("userValidator request error")
		}

		var dataUsernameSuggestion class.UsernameResponse

		err = json.Unmarshal(response.Body, &dataUsernameSuggestion)
		if err != nil {
			return fmt.Errorf("failed unmarshal")
		}

		if len(dataUsernameSuggestion.SuggestedUsernames) == 0 {
			return fmt.Errorf("not found suggestion usernames")
		}

		g.User = dataUsernameSuggestion.SuggestedUsernames[0]

	}

	return nil

}

func (g *Container) SignUp() error {

	var secureAuth *class.SecureAuth

	g.HttpClient.OrderedHeaders = azuretls.OrderedHeaders{
		{"traceparent", g.Transparent},
		{"sec-ch-ua-platform", "\"Windows\""},
		{"user-agent", userAgent},
		{"accept", "application/json, text/plain, */*"},
		{"sec-ch-ua", sec_ch_ua},
		{"sec-ch-ua-mobile", "?0"},
		{"origin", "https://www.roblox.com"},
		{"sec-fetch-site", "same-site"},
		{"sec-fetch-mode", "cors"},
		{"sec-fetch-dest", "empty"},
		{"referer", "https://www.roblox.com/"},
		{"accept-encoding", "gzip, deflate, br, zstd"},
		{"accept-language", accept_language},
	}

	g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, g.Cookies...)

	g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, []string{"priority", "u=1, i"})

	response, err := g.DoRequest("GET", "https://apis.roblox.com/hba-service/v1/getServerNonce", nil)

	if err != nil {
		return fmt.Errorf("getServerNonce error")
	}

	nonce := strings.Trim(string(response.Body), "\"")

	if nonce == "" {
		return fmt.Errorf("nonce empty")
	}

	secureAuth, err = utils.GenerateSecureAuth(nonce)

	if err != nil {
		return fmt.Errorf("failed GenerateSecureAuth")
	}

	body := &class.SignupPayload{
		Username:                 g.User,
		Password:                 g.Password,
		Birthday:                 g.Birthday,
		Gender:                   g.Gender,
		IsTosAgreementBoxChecked: true,
		AgreementIds:             []string{"306cc852-3717-4996-93e7-086daafd42f6", "2ba6b930-4ba8-4085-9e8c-24b919701f15"},
		AuditContent: class.AuditSystemContent{
			CapturedAuditContent: map[string]class.AuditItem{
				"Authentication.SignUp.Label.Birthday": {
					TranslationKey:         "Label.Birthday",
					TranslationNamespace:   "Authentication.SignUp",
					TranslatedSourceString: "Birthday",
				},
				"Authentication.SignUp.Description.SignUpAgreement.FullCopy": {
					TranslationKey:         "Description.SignUpAgreement.FullCopy",
					TranslationNamespace:   "Authentication.SignUp",
					TranslatedSourceString: "By clicking Sign Up, you are agreeing...",
					Parameters: map[string]string{
						"termsOfUseLink":    "<a target=\"_blank\" href=\"https://www.roblox.com/info/terms\">Terms of Use</a>",
						"privacyPolicyLink": "<a target=\"_blank\" href=\"https://www.roblox.com/info/privacy\">Privacy Policy</a>",
					},
				},
			},
			AdditionalAuditContent: map[string]any{},
		},
		SecureAuthenticationIntent: secureAuth,
	}

	dataSignup, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("failed json")
	}

	g.HttpClient.OrderedHeaders = azuretls.OrderedHeaders{
		{"content-length", strconv.Itoa(len(dataSignup))},
		{"sec-ch-ua-platform", "\"Windows\""},
		{"x-csrf-token", g.XCsrfToken},
		{"sec-ch-ua", sec_ch_ua},
		{"sec-ch-ua-mobile", "?0"},
		{"traceparent", g.Transparent},
		{"user-agent", userAgent},
		{"accept", "application/json, text/plain, */*"},
		{"content-type", "application/json;charset=UTF-8"},
		{"origin", "https://www.roblox.com"},
		{"sec-fetch-site", "same-site"},
		{"sec-fetch-mode", "cors"},
		{"sec-fetch-dest", "empty"},
		{"referer", "https://www.roblox.com/"},
		{"accept-encoding", "gzip, deflate, br, zstd"},
		{"accept-language", accept_language},
	}

	g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, g.Cookies...)

	g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, []string{"priority", "u=1, i"})

	response, err = g.DoRequest("POST", "https://auth.roblox.com/v2/signup", dataSignup)

	if err != nil {
		return fmt.Errorf("getBlob error")
	}

	if string(response.Body) == `{"errors":[{"code":0,"message":"Challenge is required to authorize the request"}]}` {

		header := response.Header.Get("Rblx-Challenge-Metadata")
		if header == "" {
			return fmt.Errorf("failed Rblx-Challenge-Metadata")
		}

		ark, err := utils.ParseArkoseHeader(header)
		if err != nil {
			return fmt.Errorf("failed ParseArkoseHeader")
		}

		ArkoseBlob := ark.DataExchangeBlob
		UnifiedCaptchaId := ark.UnifiedCaptchaId

		token, err := funcaptcha.GetToken(g.CapConfig.Api_Key, g.CapConfig.Http_Version, g.CapConfig.Browser_Version, ArkoseBlob, g.Proxy, "test=test", g.CapConfig.Solve_POW)

		if err != nil {
			return err
		}

		captchaToken := *token

		if captchaToken == "" {
			return fmt.Errorf("failed get token")
		}

		utils.Output("CAPTCHA", fmt.Sprintf("Solved %s", captchaToken[:28]))

		ChallengeMeta := &class.ChallengeMetadata{
			UnifiedCaptchaId: UnifiedCaptchaId,
			CaptchaToken:     captchaToken,
			ActionType:       "Signup",
		}

		metaBytes, err := json.Marshal(ChallengeMeta)

		if err != nil {
			return fmt.Errorf("challenge metadata marshal error")
		}

		metaBase64 := base64.StdEncoding.EncodeToString(metaBytes)

		Challenge := &class.ChallengeResponse{
			ChallengeId:       UnifiedCaptchaId,
			ChallengeType:     "captcha",
			ChallengeMetadata: string(metaBytes),
		}

		body, err := json.Marshal(Challenge)

		if err != nil {
			return fmt.Errorf("challenge marshal error")
		}

		g.HttpClient.OrderedHeaders = azuretls.OrderedHeaders{
			{"content-length", strconv.Itoa(len(body))},
			{"sec-ch-ua-platform", "\"Windows\""},
			{"x-csrf-token", g.XCsrfToken},
			{"sec-ch-ua", sec_ch_ua},
			{"sec-ch-ua-mobile", "?0"},
			{"traceparent", g.Transparent},
			{"user-agent", userAgent},
			{"accept", "application/json, text/plain, */*"},
			{"content-type", "application/json;charset=UTF-8"},
			{"origin", "https://www.roblox.com"},
			{"sec-fetch-site", "same-site"},
			{"sec-fetch-mode", "cors"},
			{"sec-fetch-dest", "empty"},
			{"referer", "https://www.roblox.com/"},
			{"accept-encoding", "gzip, deflate, br, zstd"},
			{"accept-language", accept_language},
		}

		g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, g.Cookies...)

		g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, []string{"priority", "u=1, i"})

		response, err = g.DoRequest("POST", "https://apis.roblox.com/challenge/v1/continue", body)

		if err != nil {
			return fmt.Errorf("captchaContiniue error")
		}

		if response.HttpResponse.StatusCode != 200 {
			return fmt.Errorf("reject continiue by API")
		} else {

			for attempt := 1; attempt <= maxRetries; attempt++ {
				g.HttpClient.OrderedHeaders = azuretls.OrderedHeaders{
					{"content-length", strconv.Itoa(len(dataSignup))},
					{"rblx-challenge-metadata", metaBase64},
					{"sec-ch-ua-platform", "\"Windows\""},
					{"x-csrf-token", g.XCsrfToken},
					{"sec-ch-ua", sec_ch_ua},
					{"rblx-challenge-id", UnifiedCaptchaId},
					{"rblx-challenge-type", "captcha"},
					{"sec-ch-ua-mobile", "?0"},
					{"traceparent", g.Transparent},
					{"user-agent", userAgent},
					{"accept", "application/json, text/plain, */*"},
					{"content-type", "application/json;charset=UTF-8"},
					{"x-retry-attempt", "1"},
					{"origin", "https://www.roblox.com"},
					{"sec-fetch-site", "same-site"},
					{"sec-fetch-mode", "cors"},
					{"sec-fetch-dest", "empty"},
					{"referer", "https://www.roblox.com/"},
					{"accept-encoding", "gzip, deflate, br, zstd"},
					{"accept-language", accept_language},
				}

				g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, g.Cookies...)

				g.HttpClient.OrderedHeaders = append(g.HttpClient.OrderedHeaders, []string{"priority", "u=1, i"})

				response, err = g.DoRequest("POST", "https://auth.roblox.com/v2/signup", dataSignup)

				if err != nil {
					return fmt.Errorf("signup error")
				}

				if string(response.Body) == `{"code":0,"message":"Token Validation Failed"}` {
					csrf := string(response.Header.Get("X-Csrf-Token"))
					if csrf != "" {
						g.XCsrfToken = csrf
					}
					continue
				} else {
					break
				}
			}

			cookies := response.Header.Get("set-cookie")

			if cookies == "" {
				return fmt.Errorf("failed get .ROBLOSECURITY")
			} else {
				parts := strings.Split(cookies, ";")

				cookieValue := ""

				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, ".ROBLOSECURITY=") {
						cookieValue = strings.TrimPrefix(part, ".ROBLOSECURITY=")
					}
				}

				if cookieValue == "" {
					return fmt.Errorf("failed get .ROBLOSECURITY")
				}

				utils.SaveAccount(g.User, g.Password, cookieValue)

				utils.Output("SUCCESS", fmt.Sprintf("Successfully created - %s", g.User))

				return nil
			}

		}

	} else {
		return fmt.Errorf("getBlob failed")
	}

}
