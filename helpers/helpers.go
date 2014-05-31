package helpers

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type Helperers interface {
	GetSites() (int64, error)
	GetInt(intChan chan int64)
	Close()
}

type Helper struct {
	scanner *bufio.Scanner
	file    *os.File
}

func New(file *os.File) *Helper {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	h := &Helper{scanner: scanner}
	return h
}

func (h *Helper) GetSites() (int64, error) {
	h.scanner.Scan()
	return strconv.ParseInt(h.scanner.Text(), 10, 32)
}

func (h *Helper) GetInt(intChan chan<- int64) {
	h.scanner.Split(bufio.ScanWords)
	for h.scanner.Scan() {
		val, err := strconv.ParseInt(h.scanner.Text(), 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		intChan <- val
	}
	close(intChan)
}

func (h *Helper) Close() {
	h.file.Close()
}
