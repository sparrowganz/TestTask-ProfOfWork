package main

import (
	"TestTask-ProfOfWork/server/pkg/excerption"
	"TestTask-ProfOfWork/server/pkg/limiter"
	"TestTask-ProfOfWork/server/pkg/locks"
	"TestTask-ProfOfWork/server/pkg/tasks"
	"bufio"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	var port = os.Getenv("SERVER_PORT")

	if port == "" {
		flag.StringVar(&port, "port", "7001", "")
		flag.Parse()
	}

	lister, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed listen tcp-server"))
	}
	defer lister.Close()

	s := &Server{
		limiter:    limiter.New(),
		excerption: excerption.New(),
		locks:      locks.New(),
		tasks:      tasks.New(),
	}
	s.Start(lister)
}

type Server struct {
	limiter    limiter.Limiter
	excerption excerption.Service
	locks      locks.Service
	tasks      tasks.Service
}

func (s *Server) Start(l net.Listener) {
	log.Printf("-- Start Listening")
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(errors.Wrap(err, "failed accept lister"))
			return
		}

		go s.handleConnect(c)
	}
}

func (s *Server) handleConnect(c net.Conn) {
	defer func() {
		log.Printf("-- Stop Serving %s\n", c.RemoteAddr().String())
		_, _ = fmt.Fprint(c, "Closing connection. Bye\n")
		_ = c.Close()
	}()
	log.Printf("-- Serving %s\n", c.RemoteAddr().String())
	_, _ = fmt.Fprint(c, "Hi\n")

	for {
		if !s.handleUserInput(c) {
			break
		}
	}
}

func (s *Server) handleUserInput(c net.Conn) bool {
	command, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return false
		}
		log.Println(errors.Wrap(err, "failed get netData"))
		return false
	}

	command = strings.ReplaceAll(command, "\n", "")

	currentLimit := s.limiter.Set(c.RemoteAddr().String())

	switch command {
	case "GET":

		if currentLimit.Blocked {
			hasContinue, hasStop := s.handleLock(c, command)
			if hasStop {
				return false
			}

			if !hasContinue {
				return true
			}
		}

		str, err := s.excerption.Get()
		if err != nil {
			log.Println(errors.Wrap(err, "failed get excertion"))
			_, _ = fmt.Fprint(c, "Something went wrong\n")
			return true
		}

		_, _ = fmt.Fprint(c, str+"\n")
	case "STOP":
		return false
	default:
		if currentLimit.Blocked {
			hasContinue, hasStop := s.handleLock(c, command)
			if hasStop {
				return false
			}

			if !hasContinue {
				return true
			}

			str, err := s.excerption.Get()
			if err != nil {
				log.Println(errors.Wrap(err, "failed get excertion"))
				_, _ = fmt.Fprint(c, "Something went wrong\n")
				return true
			}

			_, _ = fmt.Fprint(c, str+"\n")
			return true
		}

		_, _ = fmt.Fprint(c, "Incorrect input\n")
	}

	return true
}

func (s *Server) handleLock(c net.Conn, input string) (hasContinue bool, hasStop bool) {
	lock, ok := s.locks.Get(c.RemoteAddr().String())
	if !ok {
		task := s.tasks.Get()
		s.locks.Add(c.RemoteAddr().String(), &locks.Lock{
			Task:          task,
			AnswerAttempt: 0,
		})
		_, _ = fmt.Fprintf(c, "You are blocked!!! Please send answer: %s (%s)\n", task.Description, task.Variants)
		return false, false
	}

	if lock.Task.Answer == input {
		s.locks.Reset(c.RemoteAddr().String())
		s.limiter.Unblock(c.RemoteAddr().String())
		return true, false
	}

	if 4-lock.AnswerAttempt > 1 {
		s.locks.AnswerAttempt(c.RemoteAddr().String())
		_, _ = fmt.Fprintf(c, "INVALID ANSWER: you have %d attempts\n", 4-lock.AnswerAttempt)
		return false, false
	}

	return false, true
}
