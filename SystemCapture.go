/*
// Author: https://github.com/electronicsleep
// Date: 07/03/2017
// Purpose: Golang application to capture system information when thresholds are reached
// Released under the BSD license
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// CPU threshold based on number of CPU cores
var cpu_cores int = runtime.NumCPU()

// CPU threshold manually set
var threshold int = -1

// Minutes to sleep between runs
const sleep_interval time.Duration = 1

// Number of top lines to capture
const top_lines int = 25

// Verbose: Check netstat, ps -ef, df -h, lsof, iostat
var verbose = false

func captureCommand(tf string, cmd string) {

	cmd_out, cmd_err := exec.Command(cmd).Output()

	if cmd_err != nil {
		fmt.Println("ERROR:")
		log.Fatal(cmd_err)
	}
	s_cmd := string(cmd_out[:])
	cmd_u := strings.ToUpper(cmd)
	logOutput(tf, cmd_u+":", s_cmd)
}

func logOutput(date string, cmd string, cmd_out string) {

	s_cmd := string(cmd_out[:])
	lines_cmd := strings.Split(s_cmd, "\n")
	line_num := 0
	for _, line_cmd := range lines_cmd {
		line_num += 1
		fmt.Println(date, cmd, line_cmd)
		log.Println(cmd, line_cmd)
	}
}

func main() {

	verboseFlag := flag.Bool("v", false, "Verbose checks")
	cpuFlag := flag.Bool("c", false, "Detect CPU cores")

	flag.Parse()

	verbose = *verboseFlag
	cpuFlagSet := *cpuFlag

	if cpuFlagSet == true {
		fmt.Println("Setting threshold to numCPU")
		threshold = cpu_cores
	}

	fmt.Println("Verbose:", verbose)
	fmt.Println("Threshold:", threshold)

	// Start logging
	f, err := os.OpenFile("SystemCapture.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Starting SystemCapture
	log.Println("--> Starting SystemCapture")
	fmt.Println("--> Starting SystemCapture")
	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("CPU Cores:", runtime.NumCPU())
	for {

		t := time.Now()
		tf := t.Format("2006/01/02 15:04:05")

		fmt.Println("--> Checking System: Load")
		out, err := exec.Command("w").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("W: %s\n", out)
		s := string(out[:])
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			log.Println(line)
			s := strings.Split(line, " ")
			items_len := len(s)
			load15 := items_len - 1
			load5 := items_len - 2
			load1 := items_len - 3
			fmt.Println("Threshold:", threshold)
			s_load15 := strings.Split(s[load15], ".")
			s_load5 := strings.Split(s[load5], ".")
			s_load1 := strings.Split(s[load1], ".")
			int_load15, err := strconv.Atoi(s_load15[0])
			int_load5, err := strconv.Atoi(s_load5[0])
			int_load1, err := strconv.Atoi(s_load1[0])
			if err != nil {
				fmt.Println("Conversion issue")
			}
			fmt.Println("Load: ", int_load1, " ", int_load5, " ", int_load15)
			if int_load1 > threshold || int_load5 > threshold || int_load15 > threshold {
				fmt.Println("Over threshold load5")

				// Top
				var top_out []byte
				var top_err error = nil

				if runtime.GOOS == "linux" {
					// Linux specific top
					fmt.Println("Linux")
					top_out, top_err = exec.Command("top", "-bn1").Output()
				} else {
					// MacOS specific top
					fmt.Println("MacOS")
					top_out, top_err = exec.Command("top", "-l1").Output()
				}

				if top_err != nil {
					fmt.Println("ERROR:", err)
					log.Fatal(err)
				}

				s_top := string(top_out[:])
				logOutput(tf, "TOP:", s_top)

				// netstat -ta

				netstat_out, netstat_err := exec.Command("netstat", "-ta").Output()

				if netstat_err != nil {
					fmt.Println("ERROR:", err)
					log.Fatal(err)
				}

				s_netstat := string(netstat_out[:])
				logOutput(tf, "NETSTAT:", s_netstat)

				// ps -ef

				cmd_out, cmd_err := exec.Command("ps", "-ef").Output()

				if cmd_err != nil {
					fmt.Println("ERROR:", err)
					log.Fatal(err)
				}

				s_cmd := string(cmd_out[:])
				logOutput(tf, "PSEF:", s_cmd)

				// df -h

				cmd_out, cmd_err = exec.Command("df", "-h").Output()

				if cmd_err != nil {
					fmt.Println("ERROR:", err)
					log.Fatal(err)
				}

				s_cmd = string(cmd_out[:])
				logOutput(tf, "DFH:", s_cmd)

				// ps
				captureCommand(tf, "ps")
				// lsof
				if verbose {
					captureCommand(tf, "lsof")
				}
				// vmstat
				if runtime.GOOS == "linux" {
					captureCommand(tf, "vmstat")
				} else {
					captureCommand(tf, "vm_stat")
				}
				// iostat
				if verbose {
					captureCommand(tf, "iostat")
				}

			} else {
				fmt.Println("--> System load: Ok")
			}
			break
		}
		fmt.Println("Sleep for:", time.Minute*sleep_interval, "\n")
		time.Sleep(time.Minute * sleep_interval)
	}
}
