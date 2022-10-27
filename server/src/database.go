package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

var cache = Objects{
	Users: make(map[string]User),
	Apps:  make(map[string]App),
}

const filePath = "server/data/data.json"

func initDb() error {
	return readToMemory()
}

func writeToDisk() error {
	b, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("Could not marshal data: %v", err)
	}

	f, err := os.Create(filePath)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("Could not create file: %v", err)
	}

	w := bufio.NewWriter(f)

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("Could not write to file: %v", err)
	}
	w.Flush()

	return nil
}

func readToMemory() error {
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("Could not open file: %v", err)
	}

	err = json.NewDecoder(f).Decode(&cache)
	if err != nil {
		return fmt.Errorf("Could not decode json file: %v", err)
	}

	return nil
}
