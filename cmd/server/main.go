package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	goapp "goapp/internal/app/server"
	"goapp/internal/pkg/config"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.Lshortfile)
}

func main() {

	/*
		Enhancement:
			Reading configuration from the config file and setting the pprof server configuration.
	*/
	pprofEnable, pprofHost, pprofPort := config.GetPprofConfig()
	pprofAddr := pprofHost + pprofPort

	// If pprof is enabled, start the pprof server and write the heap profile to a file.
	// The heap profile is written to a file named 'heap.prof' and it can be analyzed using the cmd: `go tool pprof heap.prof`.
	if pprofEnable {
		go func() {
			log.Println("Starting pprof server on", pprofAddr)
			log.Println(http.ListenAndServe(pprofPort, nil))
		}()

		new_file, err := os.Create("heap.prof")
		if err != nil {
			log.Fatal("Failed to create heap profile. Error: ", err)
		}

		defer new_file.Close()
		if err := pprof.WriteHeapProfile(new_file); err != nil {
			log.Fatal("Failed to write the already created heap profile. Error: ", err)
		}

	}

	// Register signal handlers for exiting
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)

	// Start.
	if err := goapp.Start(exitChannel); err != nil {
		log.Fatalf("fatal: %+v\n", err)
	}
}
