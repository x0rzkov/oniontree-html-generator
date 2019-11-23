package main

import "time"

func tfTarget() string {
	return target
}

func tfLastUpdate() string {
	return time.Now().Format("2006-01-02 15:04:05 MST")
}
