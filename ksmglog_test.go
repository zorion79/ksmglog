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
	opts := Opts{
		URL:       []string{"/"},
		User:      "user",
		Password:  "pass",
		SleepTime: 1 * time.Millisecond,
		Timeout:   150 * time.Millisecond,
	}

	svc := NewService(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	svc.Run(ctx)
}

func TestService_userLogin(t *testing.T) {
	ht := httptest.NewServer(router(t))
	opts := Opts{
		SleepTime: 1 * time.Millisecond,
		Timeout:   150 * time.Millisecond,
	}

	svc := NewService(opts)

	userType, c2htoken, _, err := svc.userLogin(ht.URL + "/login")
	assert.Nil(t, err)
	assert.Equal(t, 1, userType)
	assert.Equal(t, "token", c2htoken)
}

func TestService_getCurrentTime(t *testing.T) {
	ht := httptest.NewServer(router(t))
	opts := Opts{
		SleepTime: 1 * time.Millisecond,
		Timeout:   150 * time.Millisecond,
	}

	svc := NewService(opts)

	action, id, _, err := svc.getCurrentTime(ht.URL+"/time", "token", nil)
	assert.Nil(t, err)
	assert.Equal(t, "getCurrentTime", action)
	assert.Equal(t, 2, id)
}

func TestService_getCurrentTimeWithActionID(t *testing.T) {
	ht := httptest.NewServer(router(t))
	opts := Opts{
		SleepTime: 1 * time.Millisecond,
		Timeout:   150 * time.Millisecond,
	}

	svc := NewService(opts)

	_, err := svc.getCurrentTimeWithActionID(ht.URL+"/timeID", "token", 2, nil)
	assert.Nil(t, err)
}

func TestService_eventLoggerJournalQuery(t *testing.T) {
	ht := httptest.NewServer(router(t))
	opts := Opts{
		SleepTime: 1 * time.Millisecond,
		Timeout:   150 * time.Millisecond,
	}

	svc := NewService(opts)
	id, err := svc.eventLoggerJournalQuery(ht.URL+"/eventLoggerJournalQuery", "token", nil)
	assert.Nil(t, err)
	assert.Equal(t, 3, id)
}

func TestService_eventLoggerJournalQueryWithActionID(t *testing.T) {
	ht := httptest.NewServer(router(t))
	opts := Opts{
		SleepTime: 1 * time.Millisecond,
		Timeout:   150 * time.Millisecond,
	}

	svc := NewService(opts)
	records, err := svc.eventLoggerJournalQueryWithActionID(ht.URL+"/eventLoggerJournalQueryWithActionID",
		"token", 4, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(records))
	record := records[0]
	err = record.Hash()
	assert.Nil(t, err)
	assert.Equal(t, "11ca1b99161504cb44a1a6fdd7ccbb80", record.HashString)
	//t.Logf("%+v", record)
}

func TestService_doRequest(t *testing.T) {
	opts := Opts{
		Timeout: 100 * time.Millisecond,
	}

	svc := NewService(opts)

	httpSrv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(150 * time.Millisecond)
	}))
	request, err := http.NewRequest("GET", httpSrv.URL, nil)
	assert.Nil(t, err)

	_, err = svc.doRequest(request)
	assert.NotNil(t, err)
}

func router(t *testing.T) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, _ *http.Request) {
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
	})

	mux.HandleFunc("/time", func(w http.ResponseWriter, _ *http.Request) {
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
	})

	mux.HandleFunc("/timeID", func(w http.ResponseWriter, _ *http.Request) {
		result := struct {
			Action string `json:"action"`
			Data   struct {
				Tz   string `json:"tz"`
				Time int    `json:"time"`
			} `json:"data"`
		}{
			Action: "getCurrentTime",
			Data: struct {
				Tz   string `json:"tz"`
				Time int    `json:"time"`
			}{
				Tz:   "tz",
				Time: 11111,
			},
		}

		resByte, err := json.Marshal(&result)
		assert.Nil(t, err)
		_, err = w.Write(resByte)
		assert.Nil(t, err)
	})

	mux.HandleFunc("/eventLoggerJournalQuery", func(w http.ResponseWriter, _ *http.Request) {
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
	})

	mux.HandleFunc("/eventLoggerJournalQueryWithActionID", func(w http.ResponseWriter, _ *http.Request) {
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
						Description: "test description",
					},
				},
			},
		}

		resByte, err := json.Marshal(&resultFromResp)
		assert.Nil(t, err)
		_, err = w.Write(resByte)
		assert.Nil(t, err)
	})

	return mux
}
