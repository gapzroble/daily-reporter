package log

import L "log"

// Debug func
func Debug(thread, msg string, args ...interface{}) {
	L.Printf("["+thread+"] "+msg+"\n", args...)
}
