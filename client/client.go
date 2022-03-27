package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func main() {

	var addr = os.Getenv("ADDRESS")

	if addr == "" {
		flag.StringVar(&addr, "addr", "127.0.0.1:7001", "")
		flag.Parse()
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed listen tcp-server"))
	}

	wg := &sync.WaitGroup{}

	go func() {
		for {
			// Чтение входных данных от stdin
			text, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Println("(ERROR) " + err.Error())
				continue
			}
			// Отправляем в socket
			_, _ = fmt.Fprintf(conn, text+"\n")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			// Прослушиваем ответ
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("(ERROR) " + err.Error())
				continue
			}
			log.Printf("<<-- " + message)
		}
	}()

	wg.Wait()

}
