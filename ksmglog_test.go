package ksmglog

import (
	"net/http"
	"net/http/httptest"

	//"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ClientTimeout(t *testing.T) {
	opts := Opts{
		Timeout: 100 * time.Millisecond,
	}

	svc := NewService(opts)

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
	}))
	request, err := http.NewRequest("GET", httpSrv.URL, nil)
	assert.Nil(t, err)

	_, err = svc.doRequest(request)
	assert.NotNil(t, err)
}
