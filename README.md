# System Capture

Simple Golang application to capture system information when thresholds are reached.

Update with your own commands to capture info, search for CMD:

Ever want to capture system info during a spike? This is for you.

Useful for keeping an eye on processes running using CPU

Threshold automatically set to the number of CPU cores on the system.

**Checks:** w, top, vmstat, netstat -ta, ps -ef, df

**Verbose Checks:** w, top, vmstat, netstat -ta, ps -ef, lsof, iostat, df

Using docker for tests.

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
