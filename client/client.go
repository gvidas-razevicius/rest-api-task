package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	server "github.com/gvidas-razevicius/rest-api-task/server"
	cobra "github.com/spf13/cobra"
)

const endpoint string = "http://localhost:8080"

var rootCmd = &cobra.Command{}

var getAgeCmd = &cobra.Command{
	Use:   "getage",
	Short: "Gets the age of a person by name",
	Args:  cobra.ExactArgs(1),
	Run:   getAge,
}

// TODO: implement proper query builder
func getAge(cmd *cobra.Command, args []string) {
	r, err := http.Get(endpoint + "/users/" + args[0])
	if err != nil {
		fmt.Println("Error getting the age!", err)
		return
	}
	var res server.User
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		fmt.Println("Error decoding results!", err)
		return
	}

	fmt.Printf("%s age is %d\n", res.Name, res.Age)
}

func main() {
	rootCmd.AddCommand(getAgeCmd)
	if err := getAgeCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
