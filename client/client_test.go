package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gvidas-razevicius/rest-api-task/server"
	"github.com/stretchr/testify/require"
)

func TestFormResponse(t *testing.T) {
	type test struct {
		response *http.Response
		output   string
	}

	uJson := server.User{
		Name: "Test",
		Age:  server.StringInt(100),
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
			output: fmt.Sprintf("%s is %d years old\n", uJson.Name, uJson.Age),
		},
		{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString("User deleted!")),
			},
			output: "User deleted!",
		},
		{
			response: &http.Response{
				StatusCode: http.StatusCreated,
			},
			output: "User created!",
		},
		{
			response: &http.Response{
				StatusCode: http.StatusNotFound,
			},
			output: "User was not found!",
		},
		{
			response: &http.Response{
				StatusCode: http.StatusForbidden,
			},
			output: "Could not create user. User already exists!",
		},
		{
			response: &http.Response{
				Status:     http.StatusText(http.StatusTeapot),
				StatusCode: http.StatusTeapot,
			},
			output: fmt.Sprintf("Could not get results! Server responded with: \n%s", http.StatusText(http.StatusTeapot)),
		},
	}

	for _, testX := range tests {
		require.Equal(t, testX.output, FormResponse(testX.response))
	}
}
