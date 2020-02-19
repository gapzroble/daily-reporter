package nova

import (
	"os"
	"time"
)

var (
	novaURL                      = "https://nova.fmi.filemaker-cloud.com/fmi/webd/nova%205"
	novaHours                    = "8,5"
	username                     = "Randolph Roble"
	password                     = ""
	project                      = "Byggmax"
	details                      = "(autolog)"
	width          int64         = 1390
	height         int64         = 895
	waitScreenshot time.Duration = 15 // seconds
)

func init() {
	if val := os.Getenv("NOVA_USERNAME"); val != "" {
		username = val
	}
	password = os.Getenv("NOVA_PASSWORD")
	if val := os.Getenv("PROJECT"); val != "" {
		project = val
	}
	if val := os.Getenv("DETAILS"); val != "" {
		details = val
	}
}
