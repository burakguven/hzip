# hzip

File compressor that uses the [Huffman coding algorithm](https://en.wikipedia.org/wiki/Huffman_coding). It can be used both as a library and on the command line with the `hzip` and `hunzip` commands.

## Installation

Install by running the following command ([Go](https://golang.org/dl/) needs to be installed):

```
go get github.com/burakguven/hzip/...
```

## Command Line Usage

Any data piped into `hzip` will be compressed and written to `stdout`.

Any data piped into `hunzip` will be decompressed and written to `stdout`.

Example:

    $ echo Hello World | hzip | hexdump -C
    00000000  0c 00 00 00 09 00 00 00  0a 04 d0 20 03 00 48 04  |........... ..H.|
    00000010  b0 57 04 a0 64 03 80 65  04 c0 6c 02 40 6f 03 e0  |.W..d..e..l.@o..|
    00000020  72 03 20 bc 5e 2b 96 68                           |r. .^+.h|
    00000028

    $ echo Hello World | hzip | hunzip
    Hello World

## Library Usage

Use `hzip.NewWriter` to get an `io.Writer` that will compress any data written to it.

Use `hzip.NewReader` to get an `io.Reader` that will do the opposite.

For example, the following program will write a hex dump of some compressed
data to `stdout`.

```go
package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/burakguven/hzip"
)

func main() {
	w := hzip.NewWriter(hex.Dumper(os.Stdout))
	fmt.Fprint(w, "Hello World")
	w.Close()
	fmt.Println()
}
```

