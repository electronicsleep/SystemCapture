/*
// Author: https://github.com/electronicsleep
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
var cpuCores = runtime.NumCPU()

// CPU threshold manually set
var threshold = 1

// Minutes to sleep between runs
const sleepInterval time.Duration = 1

// Number of top lines to capture
const topLines int = 25

// Verbose: Check netstat, ps -ef, df -h, lsof, iostat
var verbose = false

func captureCommand(tf string, cmd string) {

	cmdOut, cmdErr := exec.Command(cmd).Output()

	if cmdErr != nil {
		fmt.Println("ERROR:")
		log.Fatal(cmdErr)
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
		if lineNum == topLines {
			break
		}
	}
}

func main() {

	verboseFlag := flag.Bool("v", false, "Verbose checks")
	thresholdFlag := flag.Bool("t", false, "Set threshold instead of using CPU cores")

	flag.Parse()

	verbose = *verboseFlag
	thresholdFlagSet := *thresholdFlag

	if thresholdFlagSet == false {
		fmt.Println("Setting threshold to numCPU")
		threshold = cpuCores
	} else {
		fmt.Println("Using manually set threshold")
	}

	fmt.Println("Verbose:", verbose)
	fmt.Println("Threshold:", threshold)

	// Start logging
	f, err := os.OpenFile("SystemCapture.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:")
	}
	defer f.Close()
	log.SetOutput(f)

	// Starting SystemCapture
	log.Println("--> Starting SystemCapture")
	fmt.Println("--> Starting SystemCapture")
	fmt.Println("Detect OS:", runtime.GOOS)
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
			log.Println("W: " + line)
			s := strings.Split(line, " ")
			itemsLen := len(s)
			load15 := itemsLen - 1
			load5 := itemsLen - 2
			load1 := itemsLen - 3
			fmt.Println("Threshold:", threshold)
			sLoad15 := strings.Split(s[load15], ".")
			sLoad5 := strings.Split(s[load5], ".")
			sLoad1 := strings.Split(s[load1], ".")
			intLoad15, err := strconv.Atoi(sLoad15[0])
			if err != nil {
				fmt.Println("Conversion issue")
			}

			intLoad5, err := strconv.Atoi(sLoad5[0])
			if err != nil {
				fmt.Println("Conversion issue")
			}
			intLoad1, err := strconv.Atoi(sLoad1[0])
			if err != nil {
				fmt.Println("Conversion issue")
			}
			fmt.Println("Load: ", intLoad1, " ", intLoad5, " ", intLoad15)
			if intLoad1 > threshold || intLoad5 > threshold || intLoad15 > threshold {
				fmt.Println("Over threshold load5")

				// CMD: Top
				var topOut []byte
				var topErr error

				if runtime.GOOS == "linux" {
					// CMD: Linux specific top
					fmt.Println("Linux")
					topOut, topErr = exec.Command("top", "-bn1").Output()
				} else {
					// CMD: MacOS specific top
					fmt.Println("MacOS")
					topOut, topErr = exec.Command("top", "-l1").Output()
				}

				if topErr != nil {
					fmt.Println("ERROR:", err)
					log.Fatal(err)
				}

				sTop := string(topOut[:])
				logOutput(tf, "TOP:", sTop)

				if verbose {

					// CMD: netstat -ta
					netstatOut, netstatErr := exec.Command("netstat", "-ta").Output()

					if netstatErr != nil {
						fmt.Println("ERROR:", err)
						log.Fatal(err)
					}

					sNetstat := string(netstatOut[:])
					logOutput(tf, "NETSTAT:", sNetstat)

					// CMD: ps -ef
					cmdOut, cmdErr := exec.Command("ps", "-ef").Output()

					if cmdErr != nil {
						fmt.Println("ERROR:", err)
						log.Fatal(err)
					}

					sCmd := string(cmdOut[:])
					logOutput(tf, "PSEF:", sCmd)

					// CMD: df -h
					cmdOut, cmdErr = exec.Command("df", "-h").Output()

					if cmdErr != nil {
						fmt.Println("ERROR:", err)
						log.Fatal(err)
					}

					sCmd = string(cmdOut[:])
					logOutput(tf, "DFH:", sCmd)

					// CMD: ps
					captureCommand(tf, "ps")

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
				fmt.Println("--> System load: Ok")
			}
			break
		}
		fmt.Println("Sleep for:", time.Minute*sleepInterval)
		time.Sleep(time.Minute * sleepInterval)
	}
}
