package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nasermirzaei89/env"
	"net/http"
)

type State string

// const (
//	StateOK       = "ok"
//	StatePaused   = "paused"
//	StateAlerting = "alerting"
//	StatePending  = "pending"
//	StateNoData   = "no_data"
//)

type RequestBody struct {
	DashboardID int                    `json:"dashboardId"`
	EvalMatches []interface{}          `json:"evalMatches"`
	ImageURL    string                 `json:"imageUrl"`
	Message     string                 `json:"message"`
	OrgID       int                    `json:"orgId"`
	PanelID     int                    `json:"panelId"`
	RuleID      int                    `json:"ruleId"`
	RuleName    string                 `json:"ruleName"`
	RuleURL     string                 `json:"ruleUrl"`
	State       State                  `json:"state"`
	Tags        map[string]interface{} `json:"tags"`
	Title       string                 `json:"title"`
}

const u = "https://rest.nexmo.com/sms/json"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data RequestBody
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body", "error": err.Error()})
			return
		}

		text := fmt.Sprintf("%s\n%s", data.Title, data.Message)

		body := bytes.NewBufferString(
			fmt.Sprintf(
				"api_key=%s&api_secret=%s&from=%s&to=%s&text=%s",
				env.MustGetString("API_KEY"),
				env.MustGetString("API_SECRET"),
				env.MustGetString("FROM"),
				env.MustGetString("TO"),
				text,
			),
		)

		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, u, body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "error on new request", "error": err.Error()})
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "error on post", "error": err.Error()})
			return
		}

		defer func() { _ = res.Body.Close() }()

		if res.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "error on post status", "error": fmt.Sprintf("response status code is %d", res.StatusCode)})
			return
		}
	})

	err := http.ListenAndServe(env.GetString("API_ADDRESS", ":8000"), nil)
	if err != nil {
		panic(fmt.Errorf("error on listen and serve http: %w", err))
	}
}
