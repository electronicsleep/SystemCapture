/*
// Author: https://github.com/electronicsleep
// Purpose: Golang application to capture system info when CPU thresholds is reached
// Released under the BSD license
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
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

func captureCommand(tf string, cmd string) {

	cmdOut, cmdErr := exec.Command(cmd).Output()

	if cmdErr != nil {
		log.Fatal("ERROR: cmd", cmdErr)
	}
	sCmd := string(cmdOut[:])
	cmdU := strings.ToUpper(cmd)
	logOutput(tf, cmdU+":", sCmd)
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
	data, err := ioutil.ReadFile("./SystemCapture.log")
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

func runCapture() {
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
				fmt.Println("INFO: Load over threshold: Running checks")
				log.Println("INFO: Load over threshold: Running checks")
				time.Sleep(3 * time.Second)

				// CMD: Top
				var topOut []byte
				var topErr error

				if runtime.GOOS == "linux" {
					// CMD: Linux specific top
					fmt.Println("INFO: OS: Linux")
					topOut, topErr = exec.Command("top", "-bn1").Output()
				} else {
					// CMD: MacOS specific top
					fmt.Println("INFO: OS: MacOS")
					topOut, topErr = exec.Command("top", "-l1").Output()
				}

				checkFatal("ERROR: top:", topErr)
				sTop := string(topOut[:])
				logOutput(tf, "TOP:", sTop)

				// CMD: netstat -ta
				netstatOut, netstatErr := exec.Command("netstat", "-ta").Output()
				checkFatal("ERROR: netstat:", netstatErr)
				sNetstat := string(netstatOut[:])
				logOutput(tf, "NETSTAT:", sNetstat)

				// CMD: ps aux
				psOut, psErr := exec.Command("ps", "aux").Output()
				checkFatal("ERROR: ps:", psErr)
				sPS := string(psOut[:])
				logOutput(tf, "PS:", sPS)

				// CMD: df -h
				dfOut, dfErr := exec.Command("df", "-h").Output()
				checkFatal("ERROR: df:", dfErr)
				sDF := string(dfOut[:])
				logOutput(tf, "DFH:", sDF)

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

			} else {
				fmt.Println("INFO: System load: Ok")
			}
			break
		}
		fmt.Println("INFO: Checking again in:", time.Minute*sleepInterval)
		time.Sleep(time.Minute * sleepInterval)
	}
}

func main() {
	fmt.Println("INFO: Starting SystemCapture")

	verboseFlag := flag.Bool("v", false, "Verbose checks")
	thresholdFlag := flag.Int("t", 0, "Set CPU threshold manually")
	webserverFlag := flag.Bool("w", false, "Run webserver")

	flag.Parse()

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
		go runCapture()
		fmt.Println("INFO: Running webserver mode: http://localhost:8080/logs")
		http.Handle("/", http.FileServer(http.Dir("./src")))
		http.HandleFunc("/logs", httpLogs)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			checkFatal("ERROR: webserver: ", err)
		}
	} else {
		fmt.Println("INFO: Running console mode")
		runCapture()
	}

}
