package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	server "github.com/gvidas-razevicius/rest-api-task/server/src"
	cobra "github.com/spf13/cobra"
)

const (
	APIroot  = "http://localhost:8080"
	APIusers = APIroot + "/users"
	APIapp   = APIroot + "/apps"
)

// Builds and performs a request given parameters.
// -method: the REST method to use in request
// -endpoint: the API endpoint url
// -values: url values to add to the query
// -json: the json payload
func MakeRequest(method string, endpoint string, values url.Values, json []byte) (*http.Response, error) {
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

	// Do request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return resp, ErrMakeRequest{RequestErr: err}
	}

	if err := CheckStatus(resp); err != nil {
		return resp, ErrServerNegative{NegativeMsg: err}
	}

	return resp, err
}

// Gets run for the get-age command and returns the age of the given persons name
func GetAge(cmd *cobra.Command, args []string) error {
	var res []server.User

	if err := get(args, &res, APIusers); err != nil {
		return err
	}

	for _, user := range res {
		fmt.Printf(green("%s is %d years old\n"), user.Name, user.Age)
	}

	return nil
}

// Gets run for the cr-user command and creates user object
func CreateUser(cmd *cobra.Command, args []string) error {
	// Convert the string argument into int
	age, err := strconv.Atoi(args[1])
	if err != nil {
		return ErrTypeInvalid{ArgName: cmd.ValidArgs[1], TypeErr: err}
	}

	obj := server.User{
		Name: args[0],
		Age:  age,
	}
	if err := create(APIusers, obj); err != nil {
		return err
	}

	fmt.Println(green("User created!"))
	return nil
}

// Gets run by the del-user command and deletes user object
func DeleteUser(cmd *cobra.Command, args []string) error {
	if err := del(args, APIusers); err != nil {
		return err
	}
	fmt.Println(green("User deleted!"))
	return nil
}

// Gets run for the get-app command and returns the app
func GetApp(cmd *cobra.Command, args []string) error {
	var res []server.App

	if err := get(args, &res, APIapp); err != nil {
		return err
	}

	for _, app := range res {
		fmt.Printf("%s costs %f and was created in %d\n", app.Name, app.Price, app.Created)
	}
	return nil
}

// Gets run for the cr-app command and creates an app object
func CreateApp(cmd *cobra.Command, args []string) error {
	// Convert the string argument into float
	price, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return ErrTypeInvalid{ArgName: cmd.ValidArgs[1], TypeErr: err}
	}

	obj := server.App{
		Name:    args[0],
		Price:   price,
		Created: time.Now().Year(),
	}

	if err := create(APIapp, obj); err != nil {
		return err
	}

	fmt.Println(green("App created!"))
	return nil
}

// Gets run by the del-app command
func DeleteApp(cmd *cobra.Command, args []string) error {
	if err := del(args, APIapp); err != nil {
		return err
	}
	fmt.Println(green("App deleted!"))
	return nil
}

// Populates result with T objects requested from server
func get[T server.Object](args []string, result *[]T, endpoint string) error {
	val := url.Values{
		"names": args,
	}
	resp, err := MakeRequest(http.MethodGet, endpoint, val, nil)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ErrDecodeJson{JsonError: err}
	}

	return nil
}

// Creates the T object by sending the request to server
func create[T server.Object](endpoint string, obj T) error {
	// Marshal data struct
	userBytes, err := json.Marshal(obj)
	if err != nil && err != io.EOF {
		return ErrEncodeJson{JsonError: err}
	}

	resp, err := MakeRequest(http.MethodPost, endpoint, nil, userBytes)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

// Deletes obj specified in args by sending request to server
func del(args []string, endpoint string) error {
	val := url.Values{
		"names": args,
	}
	resp, err := MakeRequest(http.MethodDelete, endpoint, val, nil)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

// Checks the status code of the response and returns an error
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
