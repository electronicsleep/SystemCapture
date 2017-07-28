# System Capture

Golang application to capture system information when thresholds are reached.

Ever want to capture system info during a spike? This is for you.

Can set automatically to the number of CPU cores on the system.

Checks: w, top, vmstat

Verbose Checks: w, top, netstat -ta, ps -ef, lsof, vmstat, iostat

Raise threshold to desired level or use auto NumCPU option to log details.

Tested/working with Debian Linux and MacOS Sierra.

Should work on all Linux and MacOS versions.

```
# Run
go run SystemCapture.go

# Run in background
nohup go run SystemCapture.go
```

https://golang.org
