package main

import (
	"github.com/bizy01/scanport/scan"
)

func main() {
	s :=scan.NewScan("127.0.0.1", "3000-10000")

	s.Run()

	s.Output()
}
