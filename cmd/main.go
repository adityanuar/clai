package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"log"
    "path/filepath"
	"net/http"
	"github.com/atotto/clipboard"
)

const (
	Suggest string = "suggest"
)

type Config struct {
    ServerURL string `json:"server_url"`
}

func getServerURL() (string, error) {
    exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exePath)

	configPath := filepath.Join(exeDir, "config.json")
	file, err := os.Open(configPath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return "", err
    }

    var config Config
    if err := json.Unmarshal(bytes, &config); err != nil {
        return "", err
    }

    return config.ServerURL, nil
}

func main() {

	serverURL, err := getServerURL()
    if err != nil {
        fmt.Println(err)
        return
    }

	fmt.Println("Server URL: ", serverURL)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		os.Exit(1)
	}

	command := strings.Join(os.Args[1:], " ")

	// Create a map to hold the JSON body data
    data := map[string]string{
        "content": command,
    }
	// Marshal the map into JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err)
        return
    }
	// Make the POST request
    resp, err := http.Post(serverURL+"/"+Suggest, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

	// Read the response body
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

	// Unmarshal the JSON response
	var resultData map[string]string
	if err := json.Unmarshal(respBody, &resultData); err != nil {
		fmt.Println(err)
		return
	}

	// Get the "result" value
	result := resultData["result"]

	commandResult := string(result)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you want to run: %s ? [Y/n] ", commandResult)

	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "" || response == "y" {

		parts := strings.Fields(commandResult)
		
		// List of shell built-in commands
		builtInCommands := map[string]bool{
			"cd":      true,
			"exit":    true,
			"history": true,
			// Add other built-in commands if needed
		}
		if builtInCommands[parts[0]] {
			if err := clipboard.WriteAll(commandResult); err != nil {
				fmt.Println("Failed to copy to clipboard:", err)
				os.Exit(1)
			}
			fmt.Println("Command is a built in shell and can't be executed within this program, but rest assured it has been copied to your clipboard. Please paste it and execute it in your terminal.")
		} else {
			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Failed to execute command:", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("Suggested command rejected.")
	}
}