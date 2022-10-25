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
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Could not make request: %v", err))
	}

	fmt.Println(FormResponse(resp))
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
		log.Fatalf(fmt.Sprintf("Could not make request: %v", err))
	}
	fmt.Println(FormResponse(resp))
}

func FormResponse(resp *http.Response) string {
	switch resp.StatusCode {
	case http.StatusOK:
		// Not a great approach, should rework in the future
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Sprintf("Could not read from response body: %v", err)
		}
		var res server.User
		err = json.Unmarshal(body, &res)
		if err != nil {
			return string(body)
		}

		return fmt.Sprintf("%s is %d years old\n", res.Name, res.Age)
	case http.StatusCreated:
		return "User created!"
	case http.StatusNotFound:
		return "User was not found!"
	case http.StatusForbidden:
		return "Could not create user. User already exists!"
	default:
		return fmt.Sprintf("Could not get results! Server responded with: \n%s", resp.Status)
	}
}

func Execute() {
	rootCmd.AddCommand(getAgeCmd)
	rootCmd.AddCommand(createUserCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("There was an error while executing your CLI %v", err)
	}
}
