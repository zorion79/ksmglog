package ksmglog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_Run(t *testing.T) {
	opts := Opts{}

	svc := NewService(opts)

	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)

	svc.Run(ctx)
}

func TestService_doRequest(t *testing.T) {
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
