package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/burakguven/hzip"
)

func main() {
	log.SetPrefix("hzip: ")
	log.SetFlags(0)

	bufw := bufio.NewWriter(os.Stdout)
	w := hzip.NewWriter(bufw)
	if _, err := io.Copy(w, os.Stdin); err != nil {
		log.Fatal(err)
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	bufw.Flush()
}
