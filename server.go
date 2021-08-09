package main

import (
	f "fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

var clear map[string]func() //create a map for storing clear funcs
var oSystem string          // set os var

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		oSystem = "linux"
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		oSystem = "windows"
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearConsole() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main() {
	// set time variable
	t := time.Now()

	// set public directory
	http.Handle("/", http.FileServer(http.Dir("./public_html")))

	// ask user for input on host & port
	var host, port, url string

	// Clear any previous data in console
	clearConsole()
	f.Println("######### GoLang Webserver Configuration #########")
	f.Print("Enter Host: ")
	f.Scan(&host)
	f.Print("Enter Port (80|443): ")
	f.Scan(&port)
	// set url var with host and port
	url = host + ":" + port
	clearConsole()

	// If conncted display connection details
	f.Println("#################### GoLang Webserver ####################")
	f.Println("GoLang Server Started")
	f.Println("Host: " + host)
	f.Println("Port: " + port)
	f.Println("URL: http://" + url)
	f.Println("#################### GoLang Webserver ####################")

	// check if user hits ctrl + c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			var confirmExit string
			f.Println("Stop GoLang Server ?")
			f.Print("Exit (y/n): ")
			f.Scan(&confirmExit)
			if strstr(confirmExit, "y") == "y" {
				clearConsole()
				f.Println("#################### GoLang Webserver ####################")
				log.Println("Stopping GoLang Server and exiting..")
				f.Println("#################### GoLang Webserver ####################")
				if oSystem != "windows" {
					pprof.StopCPUProfile()
					os.Exit(1)
				} else {
					log.Fatal("Please Check Your PORTS")
				}
			} else {
				clearConsole()
				f.Println("#################### GoLang Webserver ####################")
				f.Println("GoLang Server Started")
				f.Println("Host: " + host)
				f.Println("Port: " + port)
				f.Println("URL: http://" + url)
				f.Println("#################### GoLang Webserver ####################")
			}
		}
	}()
	// http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	// http.ListenAndServe(host+":"+port, nil)
	if err := http.ListenAndServeTLS(host+":"+port, "./certificates/server.crt", "./certificates/server.key", nil); err != nil {
		clearConsole()
		f.Println("#################### GoLang Webserver ####################")
		f.Print(t.Format("2006-01-02 15:04:05"))
		f.Print(" Unable to connect to PORT " + port)
		f.Print("\r\n#################### GoLang Webserver ####################\r\n")
		if oSystem == "windows" {
			os.Exit(1)
		}
	}
}

func strstr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}
