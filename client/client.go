package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	server "github.com/gvidas-razevicius/rest-api-task/server"
	cobra "github.com/spf13/cobra"
)

const (
	APIroot  = "http://localhost:8080"
	APIusers = APIroot + "/users"
	APIapp   = APIroot + "/app"
)

// Builds and performs a request given parameters.
// -method: the REST method to use in request
// -endpoint: the API endpoint url
// -file: the file name to add at the end of the endpoint url
// -values: url values to add to the query
// -json: the json payload
func MakeRequest(method string, endpoint string, file string, values url.Values, json []byte) (*http.Response, error) {
	// Form body for request
	var body io.Reader
	if json != nil {
		body = bytes.NewBuffer(json)
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	// Set content type if body exists
	if json != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add values if any
	if values != nil {
		req.URL.RawQuery = values.Encode()
	}

	// Add additional path to URL
	req.URL = req.URL.JoinPath(file)

	return http.DefaultClient.Do(req)
}

// Gets run for the get-age command and returns the age of the given persons name
func GetAge(cmd *cobra.Command, args []string) error {
	resp, err := MakeRequest(http.MethodGet, APIusers, args[0], nil, nil)
	defer resp.Body.Close()
	if err != nil {
		return ErrMakeRequest{RequestErr: err}
	}

	if err := CheckStatus(resp); err != nil {
		return err
	}

	var res server.User
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return ErrDecodeJson{JsonError: err}
	}

	fmt.Printf(green("%s is %d years old\n"), res.Name, res.Age)
	return nil
}

// Gets run for the cr-user command and returns the age of the given persons name
func CreateUser(cmd *cobra.Command, args []string) error {
	// Convert the string argument into int
	age, err := strconv.Atoi(args[1])
	if err != nil {
		return ErrTypeInvalid{ArgName: cmd.ValidArgs[1], TypeErr: err}
	}
	// Form the json struct
	user := server.User{
		Name: args[0],
		Age:  server.StringInt(age),
	}

	// Marshal data struct
	userBytes, err := json.Marshal(user)
	if err != nil && err != io.EOF {
		return ErrEncodeJson{JsonError: err}
	}

	resp, err := MakeRequest(http.MethodPost, APIusers, "", nil, userBytes)
	if err != nil {
		return ErrMakeRequest{RequestErr: err}
	}

	if err := CheckStatus(resp); err != nil {
		return err
	}

	fmt.Println(green("User created!"))
	return nil
}

func DeleteUser(cmd *cobra.Command, args []string) error {
	resp, err := MakeRequest(http.MethodDelete, APIusers, args[0], nil, nil)
	defer resp.Body.Close()
	if err != nil {
		return ErrMakeRequest{RequestErr: err}
	}

	if err := CheckStatus(resp); err != nil {
		return err
	}

	fmt.Println(green("User deleted!"))
	return nil
}

// Gets run for the get-age command and returns the age of the given persons name
func GetApp(cmd *cobra.Command, args []string) error {
	resp, err := MakeRequest(http.MethodGet, APIapp, args[0], nil, nil)
	defer resp.Body.Close()
	if err != nil {
		return ErrMakeRequest{RequestErr: err}
	}

	if err := CheckStatus(resp); err != nil {
		return err
	}

	var res server.App
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return ErrDecodeJson{JsonError: err}
	}

	fmt.Printf("%s costs %f and was created in %d\n", res.Name, res.Price, res.Created)
	return nil
}

// Gets run for the cr-user command and returns the age of the given persons name
func CreateApp(cmd *cobra.Command, args []string) error {
	// Convert the string argument into int
	price, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return ErrTypeInvalid{ArgName: cmd.ValidArgs[1], TypeErr: err}
	}
	// Form the json struct
	app := server.App{
		Name:  args[0],
		Price: server.StringFloat(price),
	}

	// Marshal data struct
	userBytes, err := json.Marshal(app)
	if err != nil && err != io.EOF {
		return ErrEncodeJson{JsonError: err}
	}

	resp, err := MakeRequest(http.MethodPost, APIapp, "", nil, userBytes)
	if err != nil {
		return ErrMakeRequest{RequestErr: err}
	}

	if err := CheckStatus(resp); err != nil {
		return err
	}

	fmt.Println(green("App created!"))
	return nil
}

func DeleteApp(cmd *cobra.Command, args []string) error {
	resp, err := MakeRequest(http.MethodDelete, APIapp, args[0], nil, nil)
	defer resp.Body.Close()
	if err != nil {
		return ErrMakeRequest{RequestErr: err}
	}

	if err := CheckStatus(resp); err != nil {
		return err
	}

	fmt.Println(green("App deleted!"))
	return nil
}

func CheckStatus(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent, http.StatusCreated:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusForbidden:
		return ErrForbidden
	default:
		return ErrBadStatus{StatusMsg: resp.Status}
	}
}
