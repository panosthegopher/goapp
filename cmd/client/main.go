package main

import (
	"flag"
	"fmt"
	"goapp/internal/pkg/config"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.Lshortfile)
}

/*
	Feature #C:
		A command line client is created through existing Makefile (and make rule) as a separate application
		which opens a requested number of sessions simultaneously.
*/

func start(clientAddr string) {
	resp, err := http.Get(clientAddr)
	if err != nil {
		log.Printf("Failed to open session: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Session opened to %s with status: %s", clientAddr, resp.Status)
}

// Get the return code of the curl command which returns the status of the server (goapp/http API)
func getCurlReturnCode(url string) int {
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

	/*
		Enhancement:
			Reading configuration from the config file and setting the pprof server configuration.
	*/
	clientHost, clientPort := config.GetClientConfig()
	clientAddr := clientHost + clientPort
	httpHost, httpPort := config.GetConfig()
	httpAddr := httpHost + httpPort

	// Check if the server is running
	if getCurlReturnCode(httpAddr) != 0 {
		log.Fatalf("Server is not running on %s", httpAddr)
	}

	// Reading the number of clients provided from the client
	reps := flag.Int("n", 1, "Number of clients to start")
	flag.Parse()

	fmt.Printf("Starting %d clients\n", *reps)

	// Start the provided number of client in goroutines
	for i := 0; i < *reps; i++ {
		go start(clientAddr)
	}

	// Wait for the clients to finish
	select {}
}
