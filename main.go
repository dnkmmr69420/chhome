// +build !windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Define command-line flags
	createFlag := flag.Bool("n", false, "Create new folder if it doesn't exist")
	flag.Parse()

	// Check if the new home directory is provided as a command-line argument
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: chhome [-n] <new_home_directory> [<command> [arguments...]]")
		return
	}

	// Extract the new home directory from the command-line arguments
	newHome := flag.Arg(0)

	// Check if the -n flag is used and create a new folder if needed
	if *createFlag {
		_, err := os.Stat(newHome)
		if os.IsNotExist(err) {
			err = os.Mkdir(newHome, 0755)
			if err != nil {
				log.Fatalf("Failed to create new home directory: %s\n", err)
			}
			fmt.Printf("Created new home directory: %s\n", newHome)
		}
	}

	// Check if the new home directory exists
	_, err := os.Stat(newHome)
	if os.IsNotExist(err) {
		log.Fatalf("Home directory does not exist: %s\n", newHome)
	}

	// Set the HOME environment variable
	err = os.Setenv("HOME", newHome)
	if err != nil {
		log.Fatalf("Failed to set HOME environment variable: %s\n", err)
	}

	// Export the HOME environment variable for child processes
	err = os.Setenv("EXPORT_HOME", "1")
	if err != nil {
		log.Fatalf("Failed to set EXPORT_HOME environment variable: %s\n", err)
	}

	// Check if additional commands are provided after the new home directory
	if len(flag.Args()) > 1 {
		// Execute the provided commands with the updated environment
		cmd := exec.Command(flag.Arg(1), flag.Args()[2:]...)

		// Set the same environment variables for the executed command
		cmd.Env = append(os.Environ(), "HOME="+newHome, "EXPORT_HOME=1")

		// Redirect the standard input/output/error of the executed command
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to execute command '%s': %s\n", flag.Arg(1), err)
		}
	} else {
		// No additional commands provided, run the default shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			log.Fatal("Default shell not found")
		}

		// Execute the default shell with the updated environment
		cmd := exec.Command(shell)

		// Set the same environment variables for the executed command
		cmd.Env = append(os.Environ(), "HOME="+newHome, "EXPORT_HOME=1")

		// Redirect the standard input/output/error of the executed command
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to execute default shell '%s': %s\n", shell, err)
		}
	}

	fmt.Println("Commands executed successfully.")
}
