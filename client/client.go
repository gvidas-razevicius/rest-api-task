package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	server "github.com/gvidas-razevicius/rest-api-task/server"
	cobra "github.com/spf13/cobra"
)

// TODO: generate mocks for MakeRequest to test other functions
//
//go:generate go run github.com/golang/mock/mockgen -destination mocks/requester.go github.com/gvidas-razevicius/rest-api-task Requester
type Requester interface {
	MakeRequest(method string, endpoint string, values url.Values, json []byte) (*http.Response, error)
}

const (
	APIroot  = "http://localhost:8080"
	APIusers = APIroot + "/users"
)

var rootCmd = &cobra.Command{}

var getAgeCmd = &cobra.Command{
	Use:   "get-age <name>",
	Short: "Gets the age of a person by name",
	Args:  cobra.ExactArgs(1),
	Run:   GetAge,
}

var createUserCmd = &cobra.Command{
	Use:   "cr-user <name> <age>",
	Short: "Creates user in server",
	Args:  cobra.ExactArgs(2),
	Run:   CreateUser,
}

// Builds and performs a request given parameters.
func MakeRequest(method string, endpoint string, values url.Values, json []byte) (*http.Response, error) {
	// Form body for request
	var body io.Reader
	if json != nil {
		body = bytes.NewBuffer(json)
	}
	req, err := http.NewRequest(method, APIusers, body)
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
	req.URL = req.URL.JoinPath(endpoint)

	return http.DefaultClient.Do(req)
}

// Gets run for the get-age command and returns the age of the given persons name
func GetAge(cmd *cobra.Command, args []string) {
	resp, err := MakeRequest(http.MethodGet, args[0], nil, nil)
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var res server.User
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil && err != io.EOF {
			log.Fatalf("Could not decode results into json: %v", err)
		}

		fmt.Printf("%s is %d years old\n", res.Name, res.Age)
	case http.StatusNotFound:
		fmt.Printf("No user named %q was found\n", args[0])
	default:
		fmt.Println("Could not get results! Server responded with: ")
		fmt.Println(resp.Status, resp.Header)
	}
}

// Gets run for the cr-user command and returns the age of the given persons name
func CreateUser(cmd *cobra.Command, args []string) {
	// Convert the string argument into int
	age, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("Age argument is not a valid number: %v", err)
	}
	// Form the json struct
	user := server.User{
		Name: args[0],
		Age:  server.StringInt(age),
	}

	// Marshal data struct
	userBytes, err := json.Marshal(user)
	if err != nil && err != io.EOF {
		log.Fatalf("Could not encode data: %v", err)
	}

	resp, err := MakeRequest(http.MethodPost, "", nil, userBytes)
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("User created!")
	case http.StatusForbidden:
		fmt.Println("Could not create user. User already exists!")
	default:
		fmt.Println("Could not create user! Server responded with: ")
		fmt.Println(resp.Status, resp.Header)
	}
}

func Execute() {
	rootCmd.AddCommand(getAgeCmd)
	rootCmd.AddCommand(createUserCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("There was an error while executing your CLI %v", err)
	}
}
