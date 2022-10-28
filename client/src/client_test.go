package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	server "github.com/gvidas-razevicius/rest-api-task/server/src"
	"github.com/stretchr/testify/require"
)

func TestCheckStatus(t *testing.T) {
	type test struct {
		response *http.Response
		output   error
	}

	uJson := server.User{
		Name: "Test",
		Age:  100,
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(uJson)
	uJsonb := ioutil.NopCloser(&buf)

	tests := []test{
		{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       uJsonb,
			},
			output: nil,
		},
		{
			response: &http.Response{
				StatusCode: http.StatusNoContent,
			},
			output: nil,
		},
		{
			response: &http.Response{
				StatusCode: http.StatusCreated,
			},
			output: nil,
		},
		{
			response: &http.Response{
				StatusCode: http.StatusNotFound,
			},
			output: ErrNotFound,
		},
		{
			response: &http.Response{
				StatusCode: http.StatusForbidden,
			},
			output: ErrForbidden,
		},
		{
			response: &http.Response{
				Status:     http.StatusText(http.StatusTeapot),
				StatusCode: http.StatusTeapot,
			},
			output: ErrBadStatus{StatusMsg: http.StatusText(http.StatusTeapot)},
		},
	}

	for _, testX := range tests {
		require.ErrorIs(t, testX.output, CheckStatus(testX.response))
	}
}
