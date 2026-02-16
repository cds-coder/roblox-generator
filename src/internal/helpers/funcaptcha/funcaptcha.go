package funcaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{}

func GetToken(api_key, http_version, browser_version, blob, proxy, cookies string, solvePOW bool) (*string, error) {

	body := map[string]interface{}{
		"api_key":      api_key,
		"site_key":     "A2A14B1D-1AF3-C791-9BBC-EE33CC7A0A6F",
		"proxy":        proxy,
		"locale":       "en-US",
		"blob":         blob,
		"cookies":      cookies,
		"http_version": http_version,
		"solve_pow":    solvePOW,
	}

	if browser_version != "" {
		body["browser_version"] = browser_version
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://cds-solver.com/createTask", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed createTask")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed createTask")
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var taskResp map[string]interface{}
	if err := json.Unmarshal(data, &taskResp); err != nil {
		return nil, fmt.Errorf("failed jsonParse createTask")
	}

	status, ok := taskResp["status"].(string)
	if !ok {
		fmt.Println(resp.StatusCode)
		return nil, fmt.Errorf("failed get status")
	}

	if status == "started" {

		taskID, ok := taskResp["task_id"].(string)

		if !ok {
			fmt.Println(resp.StatusCode)
			return nil, fmt.Errorf("failed get taskID")
		}

		payload := map[string]interface{}{
			"api_key": api_key,
			"task_id": taskID,
		}

		for count := 0; count <= 120; count++ {

			jsonPayload, _ := json.Marshal(payload)
			req2, _ := http.NewRequest("POST", "https://cds-solver.com/getTask", bytes.NewBuffer(jsonPayload))
			req2.Header.Set("Content-Type", "application/json")

			resp2, err := client.Do(req2)
			if err != nil {
				time.Sleep(600 * time.Millisecond)
				continue
			}
			defer resp2.Body.Close()

			respData, _ := ioutil.ReadAll(resp2.Body)
			var tokenResp map[string]interface{}
			if err := json.Unmarshal(respData, &tokenResp); err != nil {
				time.Sleep(600 * time.Millisecond)
				continue
			}

			status, _ := tokenResp["status"].(string)

			switch status {
			case "processing":
				time.Sleep(600 * time.Millisecond)

			case "success":
				tok, ok := tokenResp["token"].(string)
				if ok {
					return &tok, nil
				}
				return nil, fmt.Errorf("failed get token")

			case "failed":
				tok, ok := tokenResp["error"].(string)
				if ok {
					return nil, fmt.Errorf("failed: %s", tok)
				}
				return nil, fmt.Errorf("failed get token")

			default:
				fmt.Println(tokenResp)
				return nil, fmt.Errorf("unexpected error")
			}
		}
	} else {
		if status == "failed" {
			err, ok := taskResp["error"].(string)
			if !ok {
				fmt.Println(resp.StatusCode)
				return nil, fmt.Errorf("failed get err")
			}
			return nil, fmt.Errorf("failed solve captcha - %s", err)
		}
		return nil, fmt.Errorf("failed solve captcha - status: %s", status)
	}

	return nil, fmt.Errorf("failed get token")
}
