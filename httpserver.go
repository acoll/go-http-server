package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// COLORS
var green = ansi.ColorFunc("green")
var red = ansi.ColorFunc("red")
var cyan = ansi.ColorFunc("cyan")

// Logging
var Info = log.New(os.Stdout, green("INFO:  "), log.Ldate|log.Ltime)
var Error = log.New(os.Stderr, red("ERROR: "), log.Ldate|log.Ltime)

var ascii = `
  _________    __ _________________    ___________ _   _________ 
 / ___/ __ \  / // /_  __/_  __/ _ \  / __/ __/ _ \ | / / __/ _ \
/ (_ / /_/ / / _  / / /   / / / ___/ _\ \/ _// , _/ |/ / _// , _/
\___/\____/ /_//_/ /_/   /_/ /_/    /___/___/_/|_||___/___/_/|_| 
`

func main() {

	// Listen for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		Info.Println("Good bye")
		os.Exit(1)
	}()

	fmt.Println(cyan(ascii))

	// Flags
	var Port int
	var Directory = "."

	var RootCommand = &cobra.Command{
		Use:  "httpserver",
		Long: "Serves static files in a simple http server.",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) > 0 {
				Directory = args[0]
			}

			Info.Printf("Serving %s on port %v", Directory, Port)

			fs := http.FileServer(http.Dir(Directory))
			http.Handle("/", fs)
			err := http.ListenAndServe(fmt.Sprintf(":%v", Port), nil)

			if err != nil {
				Error.Println(err)
			}
		},
	}

	RootCommand.Flags().IntVarP(&Port, "port", "p", 3000, "Port to use")

	RootCommand.Execute()
}
