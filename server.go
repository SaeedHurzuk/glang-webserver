package main

import (
	"errors"
	f "fmt"
	"log"
	"net"
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
	ip, err := externalIP()
	if err != nil {
		f.Println(err)
	}
	f.Println(ip)
	f.Println("######### GoLang Webserver Configuration #########")
	f.Print("Enter Host: ")
	f.Scan(&host)
	f.Print("Enter Port (80|443): ")
	f.Scan(&port)
	// set url var with host and port
	if port == "443" {
		url = "https://" + host + ":" + port
	} else if port == "80" {
		url = "http://" + host + ":" + port
	} else {
		url = "http://" + host + ":" + port
	}
	clearConsole()

	// If conncted display connection details
	f.Println("#################### GoLang Webserver ####################")
	f.Println("GoLang Server Started")
	f.Println("Host: " + host)
	f.Println("Port: " + port)
	f.Println("URL: " + url)
	f.Println("#################### GoLang Webserver ####################")
	test()
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
				if oSystem != "windows" {
					pprof.StopCPUProfile()
					f.Println("#################### GoLang Webserver ####################")
					os.Exit(1)
				} else {
					f.Println("#################### GoLang Webserver ####################")
					os.Exit(1)
				}
			} else {
				clearConsole()
				f.Println("#################### GoLang Webserver ####################")
				f.Println("GoLang Server Started")
				f.Println("Host: " + host)
				f.Println("Port: " + port)
				f.Println("URL: " + url)
				f.Println("#################### GoLang Webserver ####################")
			}
		}
	}()
	// http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	// http.ListenAndServe(host+":"+port, nil)
	if port == "443" {
		// listen and server HTTPs
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
	} else {
		// listen via http on any other port other than 443
		if err := http.ListenAndServe(host+":"+port, nil); err != nil {
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

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
