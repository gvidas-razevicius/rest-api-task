package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeRequestGET(t *testing.T) {
	rJohn, _ := MakeRequest(http.MethodGet, "John", nil, nil)
	rNo, _ := MakeRequest(http.MethodGet, "No", nil, nil)

	// TODO: look into what this is useful for
	//require.HTTPStatusCode(t, nil, http.MethodGet, APIusers+"/John", nil, http.StatusOK, nil)

	require.Equal(t, rJohn.StatusCode, http.StatusOK)
	require.Equal(t, rNo.StatusCode, http.StatusNotFound)
}

func TestMakeRequestPOST(t *testing.T) {

}
