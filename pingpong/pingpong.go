package main

import (
	"flag"
	port "github.com/goerlang/port"
	"io"
	"log"
	"os"
)

var packetMode bool
var lineMode bool
var packetSize int
var logDest string

func init() {
	flag.BoolVar(&packetMode, "p", false, "Packet mode")
	flag.BoolVar(&lineMode, "l", false, "Line mode")
	flag.IntVar(&packetSize, "psize", 4, "Packet size")
	flag.StringVar(&logDest, "log", "/tmp/pingpong.log", "Log destination")
}

func main() {
	flag.Parse()

	var f *os.File
	var err error
	if f, err = os.Create(logDest); err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	var p port.Port

	if packetMode && !lineMode {
		log.Printf("working in packet mode %d", packetSize)
		p, err = port.Packet(os.Stdin, os.Stdout, packetSize)
	} else if lineMode {
		log.Printf("working in line mode")
		p, err = port.Line(os.Stdin, os.Stdout)
	} else {
		log.Fatal("unknown mode")
	}

	if err != nil {
		log.Fatalf("create: %s", err)
	}

	for {
		log.Printf("waiting for ping...")

		if data, err := p.ReadOne(); err == io.EOF {
			log.Printf("finished")
			break
		} else if err != nil {
			log.Fatalf("read: %s", err)
		} else {
			log.Printf("got %v (%s)", data, data)

			log.Printf("answering...")
			if size, err := p.Write(data); err != nil || size != len(data) {
				log.Fatalf("write: %s", err)
			}
			log.Printf("pong")
		}
	}
}
