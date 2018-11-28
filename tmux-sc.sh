#!/bin/bash
set -e
go build SystemCapture.go
tmux new-session 'htop' \; split-window -v './SystemCapture -w' \; split-window -h \;
