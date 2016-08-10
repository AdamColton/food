package food

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"os"
	"strings"
)

// takes a filename and returns a channel that will stream the file one line at
// a time
func readFileByLine(f *os.File) <-chan string {
	ch := make(chan string)
	go lineReader(f, ch)
	return ch
}

// only for use by readFileByLine, do not call directly
func lineReader(f *os.File, ch chan<- string) {
	r := bufio.NewReader(f)
	for {
		lineBytes, _, e := r.ReadLine()
		if e != nil {
			break
		}
		ch <- string(lineBytes)
	}
	close(ch)
}

// split on ^ and remove leading and trailing ~
// I have no idea why the data is formatted this way
func splitLine(line string) []string {
	data := strings.Split(line, "^")
	for i := 0; i < len(data); i++ {
		data[i] = strings.Trim(data[i], "~")
	}
	return data
}

// gob encodes an object
func enc(obj interface{}) []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(obj)
	return buf.Bytes()
}

// decodes a gob object
func dec(b []byte, obj interface{}) {
	gob.NewDecoder(bytes.NewBuffer(b)).Decode(obj)
}
