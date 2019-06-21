package ksmglog

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/pkg/errors"
)

// Service create engine which collects logs from ksmg
type Service struct {
	Opts
}

// Opts collects parameters to initialize Service
type Opts struct {
	URL      string // url for ksmg web server https://ksmg01/ksmg/en-US/cgi-bin/klwi,https://ksmg02/ksmg/en-US/cgi-bin/klwi
	User     string
	Password string
}

// NewService initializes everything
func NewService(opts Opts) *Service {
	res := &Service{
		Opts: opts,
	}

	return res
}

// GetLogs return last audit logs
func (s *Service) GetLogs() (records []Record, err error) {
	records = make([]Record, 0)
	for _, ksmgURL := range strings.Split(s.URL, ",") {
		_, c2htoken, cookies, err := s.userLogin(ksmgURL)
		if err != nil {
			return nil, errors.Wrap(err, "could not login")
		}

		time.Sleep(100 * time.Millisecond)

		_, actionID, cookies, err := s.getCurrentTime(ksmgURL, c2htoken, cookies)
		if err != nil {
			return nil, errors.Wrap(err, "could not get current time")
		}

		time.Sleep(300 * time.Millisecond)

		cookies, err = s.getCurrentTimeWithActionID(ksmgURL, c2htoken, actionID, cookies)
		if err != nil {
			return nil, errors.Wrap(err, "could not get current time for action id")
		}

		time.Sleep(300 * time.Millisecond)

		actionID, err = s.eventLoggerJournalQuery(ksmgURL, c2htoken, cookies)
		if err != nil {
			return nil, errors.Wrap(err, "could not get event logger action id")
		}

		time.Sleep(2500 * time.Millisecond)

		recs, err := s.eventLoggerJournalQueryWithActionID(ksmgURL, c2htoken, actionID, cookies)
		if err != nil {
			return nil, errors.Wrap(err, "could not get records")
		}

		records = append(records, recs...)
	}

	return records, nil
}

func (s *Service) userLogin(ksmgURL string) (userType int, c2htoken string, cookie []*http.Cookie, err error) {
	requestBody := url.Values{}
	requestBody.Set("username", s.User)
	requestBody.Set("password", s.Password)
	body := strings.NewReader(requestBody.Encode())
	req, _ := http.NewRequest("POST", ksmgURL, body)
	query := req.URL.Query()
	query.Add("action", "userLogin")
	query.Add("cb", "332211")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := doRequest(req)
	if err != nil {
		return -1, "", []*http.Cookie{}, errors.Wrap(err, "could not request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] could not close body: %v", err)
		}
	}()

	result := struct {
		Action   string `json:"action"`
		UserType int    `json:"userType"`
		C2htoken string `json:"C2HToken"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return -1, "", []*http.Cookie{}, errors.Wrap(err, "could not unmarshal body")
	}

	log.Printf("[DEBUG] result from login: %v", result)

	return result.UserType, result.C2htoken, resp.Cookies(), nil
}

func (s *Service) getCurrentTime(ksmgURL string, c2htoken string, cookies []*http.Cookie) (action string, actionID int, cookie []*http.Cookie, err error) {
	req, _ := http.NewRequest("POST", ksmgURL, nil)
	query := req.URL.Query()
	query.Add("action", "getCurrentTime")
	query.Add("C2HToken", c2htoken)
	query.Add("cb", "332211")
	req.URL.RawQuery = query.Encode()
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := doRequest(req)
	if err != nil {
		return "", -1, []*http.Cookie{}, errors.Wrap(err, "could not request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] could not close body: %v", err)
		}
	}()

	result := struct {
		Action   string `json:"action"`
		ActionId int    `json:"action_id"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return "", -1, []*http.Cookie{}, errors.Wrap(err, "could not unmarshal body")
	}

	log.Printf("[DEBUG] result from getCurrentTime: %v", result)

	return result.Action, result.ActionId, resp.Cookies(), nil
}

func (s *Service) getCurrentTimeWithActionID(ksmgURL string, c2htoken string, actionID int, cookies []*http.Cookie) ([]*http.Cookie, error) {
	req, _ := http.NewRequest("POST", ksmgURL, nil)
	query := req.URL.Query()
	query.Add("action", "getCurrentTime")
	query.Add("C2HToken", c2htoken)
	query.Add("action_id", strconv.Itoa(actionID))
	query.Add("cb", "332211")
	req.URL.RawQuery = query.Encode()

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := doRequest(req)
	if err != nil {
		return []*http.Cookie{}, errors.Wrap(err, "could not request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] could not close body: %v", err)
		}
	}()

	result := struct {
		Action string `json:"action"`
		Data   struct {
			Tz   string `json:"tz"`
			Time int    `json:"time"`
		} `json:"data"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return []*http.Cookie{}, errors.Wrap(err, "could not unmarshal body")
	}

	log.Printf("[DEBUG] result from getCurrentTimeWithActionID: %v", result)

	return resp.Cookies(), nil
}

func (s *Service) eventLoggerJournalQuery(ksmgURL string, c2htoken string, cookies []*http.Cookie) (actionID int, err error) {
	req, _ := http.NewRequest("POST", ksmgURL, nil)
	query := req.URL.Query()
	query.Add("action", "eventLoggerJournalQuery")
	query.Add("C2HToken", c2htoken)
	query.Set("data", `{"filters":{"dateType":8}}`)
	req.URL.RawQuery = query.Encode()

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := doRequest(req)
	if err != nil {
		return -1, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] could not close body: %v", err)
		}
	}()

	result := struct {
		Action   string `json:"action"`
		ActionId int    `json:"action_id"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return -1, errors.Wrap(err, "could not unmarshal body")
	}

	log.Printf("[DEBUG] result from eventLoggerJournalQuery: %v", result)

	return result.ActionId, nil
}

func (s *Service) eventLoggerJournalQueryWithActionID(ksmgURL string, c2htoken string, actionID int, cookies []*http.Cookie) (res []Record, err error) {
	req, _ := http.NewRequest("POST", ksmgURL, nil)
	query := req.URL.Query()
	query.Add("action", "eventLoggerJournalQuery")
	query.Add("C2HToken", c2htoken)
	query.Add("data", `{"filters":{"dateType":8}}`)
	query.Add("action_id", strconv.Itoa(actionID))
	req.URL.RawQuery = query.Encode()

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := doRequest(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] could not close body: %v", err)
		}
	}()

	resultFromResp := struct {
		Action string `json:"action"`
		Data   struct {
			Count                int      `json:"count"`
			UnlimitedResultsSize int      `json:"unlimitedResultsSize"`
			Time                 int      `json:"time"`
			Items                []Record `json:"items"`
		} `json:"data"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&resultFromResp)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal body")
	}

	res = resultFromResp.Data.Items

	return res, nil
}

func doRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return resp, nil
}
