package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/burakguven/hzip"
)

func main() {
	log.SetPrefix("hunzip: ")
	log.SetFlags(0)

	var err error
	r, err := hzip.NewReader(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	bufw := bufio.NewWriter(os.Stdout)
	_, err = io.Copy(bufw, r)
	if err != nil {
		log.Fatal(err)
	}
	bufw.Flush()
}
