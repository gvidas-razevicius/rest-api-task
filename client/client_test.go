package client

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gvidas-razevicius/rest-api-task/server"
	"github.com/stretchr/testify/require"
)

func TestMakeRequestGET(t *testing.T) {
	rJohn, _ := MakeRequest(http.MethodGet, "John", nil, nil)
	rNo, _ := MakeRequest(http.MethodGet, "No", nil, nil)

	require.Equal(t, rJohn.StatusCode, http.StatusOK)
	require.Equal(t, rNo.StatusCode, http.StatusNotFound)
}

func TestMakeRequestPOST(t *testing.T) {
	user := server.User{
		Name: "Pete",
		Age:  server.StringInt(10),
	}

	// Marshal data struct
	userBytes, _ := json.Marshal(user)
	r, _ := MakeRequest(http.MethodPost, "", nil, userBytes)

	require.Equal(t, r.StatusCode, http.StatusOK)
}
