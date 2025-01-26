// Author: https://github.com/electronicsleep
// Purpose: Golang application to capture system info when CPU thresholds are reached
// Released under the BSD license

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// CPU auto detect CPU threshold based on number of CPU cores
var cpuCores = runtime.NumCPU()

// CPU threshold manually set [use -t] [set to -1 to always capture]
var threshold = -1

// Minutes to sleep between runs
const sleepInterval time.Duration = 1

// Verbose: Check vmstat, lsof, iostat
var verbose = false

// Webserver: Run webserver to show output (experimental)
var webserver = false

type configStruct struct {
	SlackURL string   `yaml:"slack_url"`
	SlackMsg string   `yaml:"slack_msg"`
	Commands []string `yaml:"commands"`
}

type stateStruct struct {
	Hostname string
}

func (config *configStruct) getConfig() *configStruct {

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("ERROR: YAML file not found #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return config
}

func (state *stateStruct) setState() *stateStruct {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	} else {
		state.Hostname = hostname
	}
	return state
}

func captureCommand(tf string, cmd string) bool {
	cmdOut, cmdErr := exec.Command("bash", "-c", cmd).Output()
	if cmdErr != nil {
		fmt.Println("ERROR: cmd", cmdErr)
		log.Fatal("ERROR: cmd", cmdErr)
	}
	sCmd := string(cmdOut[:])
	cmdU := strings.ToUpper(cmd)
	logOutput(tf, "CMD: "+cmdU+":", sCmd)
	return true
}

func logOutput(date string, cmd string, cmdOut string) {
	sCmd := string(cmdOut[:])
	linesCmd := strings.Split(sCmd, "\n")
	lineNum := 0
	for _, lineCmd := range linesCmd {
		lineNum++
		fmt.Println(date, cmd, lineCmd)
		log.Println(cmd, lineCmd)
	}
}

func httpLogs(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("./SystemCapture.log")
	if err != nil {
		fmt.Println("ERROR: reading file")
	}
	if data != nil {
		w.Write([]byte(data))
	}
	w.Write([]byte("END:"))
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Println("ERROR: ", msg, err)
		log.Println("ERROR: ", msg, err)
	}
}

func checkFatal(msg string, err error) {
	if err != nil {
		fmt.Println("FATAL: "+msg, err)
		log.Println("FATAL: "+msg, err)
		log.Fatal()
	}

}

func runCapture(state stateStruct, config configStruct) {
	loop := 0
	for {
		loop++
		fmt.Println("INFO: Runtime:", loop, "minutes")
		t := time.Now()
		tf := t.Format("2006/01/02 15:04:05")

		fmt.Println("INFO: Checking System: Load")
		out, err := exec.Command("w").Output()
		checkFatal("invalid command w:", err)
		fmt.Printf("INFO: W: %s\n", out)
		s := string(out[:])
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			log.Println("INFO: W: " + line)
			s := strings.Split(line, " ")
			itemsLen := len(s)
			load15 := itemsLen - 1
			load5 := itemsLen - 2
			load1 := itemsLen - 3
			fmt.Println("INFO: Threshold:", threshold)
			sLoad15 := strings.Split(s[load15], ".")
			sLoad5 := strings.Split(s[load5], ".")
			sLoad1 := strings.Split(s[load1], ".")
			intLoad15, err := strconv.Atoi(sLoad15[0])
			checkError("ERROR: conversion issue load 15", err)
			intLoad5, err := strconv.Atoi(sLoad5[0])
			checkError("ERROR: conversion issue load 5", err)
			intLoad1, err := strconv.Atoi(sLoad1[0])
			checkError("ERROR: conversion issue load 1", err)
			fmt.Println("INFO: Load: ", intLoad1, " ", intLoad5, " ", intLoad15)
			if intLoad1 > threshold || intLoad5 > threshold || intLoad15 > threshold {
				sendMessage("INFO: Load over threshold: Hostname: "+state.Hostname, config)
				log.Println("INFO: Load over threshold: Running checks")
				time.Sleep(3 * time.Second)

				if runtime.GOOS == "linux" {
					// CMD: Linux specific top
					fmt.Println("INFO: OS: Linux")
					captureCommand(tf, "top -bn1")
				} else {
					// CMD: MacOS specific top
					fmt.Println("INFO: OS: MacOS")
					captureCommand(tf, "top -l1")
				}

				// CMD: netstat -ta
				captureCommand(tf, "netstat -ta")

				// CMD: ps aux
				captureCommand(tf, "ps aux")

				// CMD: df -h
				captureCommand(tf, "df -h")

				if verbose {
					// CMD: lsof
					captureCommand(tf, "lsof")

					// CMD: vmstat
					if runtime.GOOS == "linux" {
						captureCommand(tf, "vmstat")
					} else {
						captureCommand(tf, "vm_stat")
					}

					// CMD: iostat
					captureCommand(tf, "iostat")
				}

				if len(config.Commands) != 0 {
					fmt.Println("INFO: Config: User defined Commands:", config.Commands)
					for idx, cmd := range config.Commands {
						line := fmt.Sprintf("User Command: %d: "+cmd, idx)
						fmt.Println(line)
						captureCommand(tf, cmd)
					}
				}

			} else {
				fmt.Println("INFO: System load: Ok")
			}
			break
		}
		fmt.Println("INFO: Checking again in:", time.Minute*sleepInterval)
		time.Sleep(time.Minute * sleepInterval)
	}
}

func sendMessage(send_text string, config configStruct) {
	if config.SlackURL != "" {
		fmt.Println("SlackURL is set: sending message")
		postSlack(send_text, config)
	} else {
		fmt.Println("INFO: SlackURL is not set: no messages will be sent")
	}
}

func postSlack(message string, config configStruct) {
	t := time.Now()
	tf := t.Format("2006/01/02 15:04:05")
	fmt.Println(tf + " INFO: postSlack:" + message)
	send_text := tf + " " + message + ": " + config.SlackMsg

	var jsonData = []byte(`{
		"text": "` + send_text + `",
        }`)

	if is_connected() {
		request, error := http.NewRequest("POST", config.SlackURL, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()

		fmt.Println("INFO: response Status:", response.Status)
		fmt.Println("INFO: response Headers:", response.Header)
		body, _ := io.ReadAll(response.Body)
		fmt.Println("INFO: response Body:", string(body))
	} else {
		fmt.Println("ERROR: Not connected to the net")
	}
}

func is_connected() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func main() {
	fmt.Println("INFO: Starting SystemCapture")

	verboseFlag := flag.Bool("v", false, "Verbose checks")
	thresholdFlag := flag.Int("t", 0, "Set CPU threshold manually")
	webserverFlag := flag.Bool("w", false, "Run webserver")

	flag.Parse()

	var config configStruct
	config.getConfig()

	var state stateStruct
	state.setState()
	fmt.Println("INFO: Hostname: " + state.Hostname)
	sendMessage("INFO: Starting SystemCapture on hostname: "+state.Hostname, config)

	verbose = *verboseFlag
	webserver = *webserverFlag
	thresholdFlagSet := *thresholdFlag

	if thresholdFlagSet == 0 {
		fmt.Println("INFO: Setting threshold to numCPU")
		threshold = cpuCores
	} else {
		fmt.Println("INFO: Using flag set threshold")
		threshold = thresholdFlagSet
	}

	if verbose {
		fmt.Println("INFO: Verbose:", verbose)
		fmt.Println("INFO: Webserver:", webserver)
		fmt.Println("INFO: Threshold:", threshold)
	}

	// Start logging
	f, err := os.OpenFile("SystemCapture.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkFatal("ERROR: opening file", err)

	defer f.Close()
	log.SetOutput(f)

	fmt.Println("INFO: Detect OS:", runtime.GOOS)
	fmt.Println("INFO: CPU Cores:", runtime.NumCPU())

	if webserver {
		go runCapture(state, config)
		fmt.Println("INFO: Running webserver mode: http://localhost:8080/logs")
		http.Handle("/", http.FileServer(http.Dir("./src")))
		http.HandleFunc("/logs", httpLogs)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			checkFatal("ERROR: webserver: ", err)
		}
	} else {
		fmt.Println("INFO: Running console mode")
		runCapture(state, config)
	}
}
