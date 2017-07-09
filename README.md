# System Capture

Golang application to capture system information when thresholds are reached.

Ever want to capture system info during a spike? This is for you.

Set automatically to the number of CPU cores on the system.

Checks: w, top, netstat -ta, ps -ef, lsof, vmstat, iostat

Tested/working with Debian Linux and MacOS.

```
# Run
go run SystemCapture.go

# Run in background
nohup go run SystemCapture.go
```

https://golang.org
