package main

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	dat, err := ioutil.ReadFile("nova/tiqqe_june")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("nova/tiqqe_june.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var row []string
	columns := 7

	lines := strings.Split(string(dat), "\n")
	for _, line := range lines {
		if len(row) == 4 {
			line = strings.ReplaceAll(line, ",", ".")
		}
		row = append(row, line)
		if len(row) == columns {
			writer.Write(row)
			row = make([]string, 0, columns)
		}
	}

	if len(row) > 0 {
		writer.Write(row)
	}
}
