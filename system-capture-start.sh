#!/bin/bash
echo "### Start rsync"
#service rsync start
echo "### Start SystemCapture"
./SystemCapture &
echo "Add side proc with SystemCapture"
echo "### Start app"
sleep 10
