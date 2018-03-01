# System Capture

Golang application to capture system information when thresholds are reached.

Ever want to capture system info during a spike? This is for you.

Threshold automatically set to the number of CPU cores on the system.

**Checks:** w, top, vmstat, netstat -ta, ps -ef, df

**Verbose Checks:** w, top, vmstat, netstat -ta, ps -ef, lsof, iostat, df

Raise threshold to desired level or use auto NumCPU option to log details.

Tested/working with Debian Linux and MacOS Sierra.

Should work on all Linux and MacOS versions.

```
# Build
go build SystemCapture.go

# Build for Linux
GOOS=linux go build SystemCapture.go

# Run threshold is number of CPU cores
go run SystemCapture.go

# Run Verbose
go run SystemCapture.go -v

# Manually set threshold
go run SystemCapture.go -t

# Run in background
nohup go run SystemCapture.go
```

# Resources

https://golang.org
