/*
// Author: https://github.com/electronicsleep
// Date: 07/03/2017
// Purpose: Golang application to capture system information when thresholds are reached
// Released under the BSD license
*/

package main

import (
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
var threshold int = runtime.NumCPU()

// Manually set
//const threshold int = 4

// Minutes to sleep between runs
const sleep_interval time.Duration = 1

// Number of top lines to capture
const top_lines int = 25

// Check CMDs
const check_netstat bool = true
const check_ps bool = true

func captureCommand(cmd string) {

	cmd_out, cmd_err := exec.Command(cmd).Output()

	if cmd_err != nil {
		fmt.Println("ERROR:")
		log.Fatal(cmd_err)
	}

	s_cmd := string(cmd_out[:])
	fmt.Printf("%s: %s", strings.ToUpper(cmd), s_cmd+"\n")
	log.Println(s_cmd)

}

func main() {

	for {
		log.Println("Starting SystemCapture")
		fmt.Println("OS:", runtime.GOOS)
		fmt.Println("CPU:", runtime.NumCPU())

		f, err := os.OpenFile("SystemCapture.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)

		fmt.Println("Checking System: Load")
		out, err := exec.Command("w").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("w: %s\n", out)
		s := string(out[:])
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			fmt.Println("W:", line)
			log.Println(line)
			s := strings.Split(line, " ")
			items_len := len(s)
			load15 := items_len - 1
			load5 := items_len - 2
			load1 := items_len - 3
			fmt.Println("Threshold:", threshold)
			//fmt.Println("Load:", s[load5])
			s_load15 := strings.Split(s[load15], ".")
			s_load5 := strings.Split(s[load5], ".")
			s_load1 := strings.Split(s[load1], ".")
			int_load15, err := strconv.Atoi(s_load15[0])
			int_load5, err := strconv.Atoi(s_load5[0])
			int_load1, err := strconv.Atoi(s_load1[0])
			if err != nil {
				fmt.Println("Conversion issue")
			}
			fmt.Println("load: ", int_load1, " ", int_load5, " ", int_load15)
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
					fmt.Println("ERROR:")
					log.Fatal(err)
				}

				s_top := string(top_out[:])
				lines_top := strings.Split(s_top, "\n")
				line_num := 0
				for _, line_top := range lines_top {
					line_num += 1
					fmt.Printf("TOP: %s", line_top+"\n")
					log.Println(line_top)
					if line_num == top_lines {
						break
					}
				}

				// netstat -ta
				if check_netstat {

					netstat_out, netstat_err := exec.Command("netstat", "-ta").Output()

					if netstat_err != nil {
						fmt.Println("ERROR:")
						log.Fatal(err)
					}

					s_netstat := string(netstat_out[:])
					fmt.Printf("NETSTAT: %s", s_netstat+"\n")
					log.Println(s_netstat)
				}

				// ps -ef
				if check_ps {

					cmd_out, cmd_err := exec.Command("ps", "-ef").Output()

					if cmd_err != nil {
						fmt.Println("ERROR:")
						log.Fatal(err)
					}

					s_cmd := string(cmd_out[:])
					fmt.Printf("PSEF: %s", s_cmd+"\n")
					log.Println(s_cmd)
				}

				// ps
				captureCommand("ps")
				// lsof
				captureCommand("lsof")
				// vmstat
				if runtime.GOOS == "linux" {
					captureCommand("vmstat")
				} else {
					captureCommand("vm_stat")
				}
				// iostat
				captureCommand("iostat")

			} else {
				fmt.Println("System load ok")
			}
			break
		}
		fmt.Println("Sleep for:", time.Minute*sleep_interval)
		time.Sleep(time.Minute * sleep_interval)
	}
}
