# System Capture

Golang application to capture system information when thresholds are reached.

Ever want to capture system info during a spike? This is for you.

Set automatically to the number of CPU cores on the system.

Checks: w, top, vmstat

Checks Verbose: w, top, netstat -ta, ps -ef, lsof, vmstat, iostat

Raise threshold to desired level or use auto NumCPU option.

Tested/working with Debian Linux and MacOS.

```
# Run
go run SystemCapture.go

# Run in background
nohup go run SystemCapture.go
```

https://golang.org
