package ksmglog

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_Run(t *testing.T) {
	ht := httptest.NewServer(router(t))
	opts := Opts{
		URL:       []string{ht.URL},
		User:      "user",
		Password:  "pass",
		SleepTime: 1 * time.Millisecond,
		Timeout:   150 * time.Millisecond,
	}

	svc := NewService(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	go svc.Run(ctx)

	ch := svc.Channel()
	for r := range ch {
		t.Logf("%+v", r)
		assert.Equal(t, 222, r.ID)
	}
	cancel()
}

func router(t *testing.T) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "test",
			Value: "test",
		})
		switch r.URL.Query().Get("action") {
		case "userLogin":
			result := struct {
				Action   string `json:"action"`
				UserType int    `json:"userType"`
				C2htoken string `json:"C2HToken"`
			}{
				Action:   "userLogin",
				UserType: 1,
				C2htoken: "token",
			}
			resultByte, err := json.Marshal(&result)
			assert.Nil(t, err)
			_, err = w.Write(resultByte)

			assert.Nil(t, err)
		case "getCurrentTime":
			if r.URL.Query().Get("action_id") == "" {
				result := struct {
					Action   string `json:"action"`
					ActionId int    `json:"action_id"`
				}{
					Action:   "getCurrentTime",
					ActionId: 2,
				}

				resByte, err := json.Marshal(&result)
				assert.Nil(t, err)
				_, err = w.Write(resByte)
				assert.Nil(t, err)
				return
			}

			if r.URL.Query().Get("action_id") == "2" {
				result := struct {
					Action   string `json:"action"`
					ActionId int    `json:"action_id"`
				}{
					Action:   "eventLoggerJournalQuery",
					ActionId: 3,
				}

				resByte, err := json.Marshal(&result)
				assert.Nil(t, err)
				_, err = w.Write(resByte)
				assert.Nil(t, err)
				return
			}
		case "eventLoggerJournalQuery":
			if r.URL.Query().Get("action_id") == "" {
				result := struct {
					Action   string `json:"action"`
					ActionId int    `json:"action_id"`
				}{
					Action:   "eventLoggerJournalQuery",
					ActionId: 3,
				}

				resByte, err := json.Marshal(&result)
				assert.Nil(t, err)
				_, err = w.Write(resByte)
				assert.Nil(t, err)
				return
			}

			if r.URL.Query().Get("action_id") == "3" {
				resultFromResp := struct {
					Action string `json:"action"`
					Data   struct {
						Count                int      `json:"count"`
						UnlimitedResultsSize int      `json:"unlimitedResultsSize"`
						Time                 int      `json:"time"`
						Items                []Record `json:"items"`
					} `json:"data"`
				}{
					Action: "eventLoggerJournalQuery",
					Data: struct {
						Count                int      `json:"count"`
						UnlimitedResultsSize int      `json:"unlimitedResultsSize"`
						Time                 int      `json:"time"`
						Items                []Record `json:"items"`
					}{
						Count:                1,
						UnlimitedResultsSize: 1,
						Time:                 1,
						Items: []Record{
							{
								ID:          111,
								Time:        int(time.Now().Unix()),
								Description: "test description",
							},
							{
								ID:          222,
								Time:        int(time.Now().AddDate(0, 0, -2).Unix()),
								Description: "second record",
							},
						},
					},
				}

				resByte, err := json.Marshal(&resultFromResp)
				assert.Nil(t, err)
				_, err = w.Write(resByte)
				assert.Nil(t, err)
			}
		}
	})

	return mux
}
