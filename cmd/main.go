// clai/cmd/main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		os.Exit(1)
	}

	command := os.Args[1]

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you want to run: %s ? [Y/n] ", command)

	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "" || response == "y" {
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to execute command:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Suggested command rejected.")
	}
}