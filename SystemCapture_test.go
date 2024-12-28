package main

import (
	"testing"
	"time"
)

func TestCaptureCommand(t *testing.T) {

	cmd := "df -h"
	tm := time.Now()
	tf := tm.Format("2006/01/02 15:04:05")
	res := captureCommand(tf, cmd)
	expect := true

	if res != expect {
		t.Errorf("got %t wanted %t", res, expect)
	}
}
