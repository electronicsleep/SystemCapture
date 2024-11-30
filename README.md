# System Capture

SystemCapture - Go app to capture sysinfo on MacOS/Linux

Capture system details when thresholds are reached good for troubleshooting

Default will set threshold to number of CPU cores

**Regular Checks:** w, top, netstat -ta, ps -ef, ps, df

**Verbose Checks:** Regular Checks + vmstat, lsof, iostat

```
# Run
go run SystemCapture.go

# Set CPU threshold to 4
go run SystemCapture.go -t 4

# Set threshold -1 to always capture
go run SystemCapture.go -t -1

# Run with webserver and verbose:
# http://localhost:8080/logs
go run SystemCapture.go -t 4 -w -v

# Build for Linux
GOOS=linux go build SystemCapture.go

# Docker env Alpine
bash docker-run.sh

# Run background
nohup go run SystemCapture.go

# Example
git clone https://github.com/electronicsleep/SystemCapture.git && cd SystemCapture && go run SystemCapture.go

# Notifications and custom commands
# Config: config.yaml
slack_url: https://hooks.slack.com/services/
slack_msg: "<@user>"
commands:
  - "/bin/ls -l"
  - "/bin/df -i"
```

Update with your own commands to capture info, search for CMD:

Ever want to capture system info during a spike? This is for you

Useful for keeping an eye on processes running using CPU (similar to SAR report)

Threshold automatically set to the CPU cores on the system (will only capture when load is high)

Using docker for testing/verifying Linux

Raise threshold to desired level or use auto NumCPU option to use defaults to log details

Should work on all Linux and MacOS versions, if you find any issues let me know

# Resources

https://golang.org
