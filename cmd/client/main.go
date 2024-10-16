package main

import (
	"flag"
	"fmt"
	client "goapp/internal/app/client"
	"goapp/internal/pkg/config"
	"log"
	"os"
	"os/exec"
	"sync"
)

/*
	Feature #C:
		A command line client is created through existing Makefile (and make rule) as a separate application
		which opens a requested number of sessions simultaneously.
*/

// Initialize the logger as doing for the server as well
func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.Lshortfile)
}

// Get the return code of the curl command which returns the status of the server (goapp/http API)
func getHttpServerHealthCheckState(url string) int {
	cmd := exec.Command("curl", "-v", url)
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		return -1
	}
	return 0
}

func main() {
	var wg sync.WaitGroup

	// Get the server address from the config
	httpHost, httpPort := config.GetConfig()
	httpAddr := httpHost + httpPort

	// Check if the server is running by sending a health check request
	if getHttpServerHealthCheckState(httpAddr) != 0 {
		log.Fatalf("Server is not running on %s", httpAddr)
	}

	// Reading the number of clients provided from the client
	reps := flag.Int("n", 1, "Number of clients to start")
	flag.Parse()

	fmt.Printf("Starting %d clients\n", *reps)

	// Start the client where go routines included there
	client.StartClient(*reps, httpAddr, &wg)

	// Wait for all client routines to complete
	wg.Wait()

}
