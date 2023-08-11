package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"golang.org/x/crypto/ssh"
)

type Session struct {
	*ssh.Session
	client    *ssh.Client
	stdout    io.Reader
	stderr    io.Reader
	ctx       context.Context
	closeChan chan struct{}
}

func NewSession(ctx context.Context, client *ssh.Client) (*Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, err
	}
	return &Session{
		Session:   session,
		client:    client,
		stdout:    stdout,
		stderr:    stderr,
		ctx:       ctx,
		closeChan: make(chan struct{}, 1),
	}, nil
}

func (s *Session) ListenOut(writer io.Writer) {
	go func() {
		title := fmt.Sprintf("====== %s ======\r\n", s.client.RemoteAddr().String())
		for {
			select {
			case <-s.closeChan:
				log.Println("chan closed")
				return
			case <-s.ctx.Done():
				log.Println("context done")
				return
			default:
			}
			buffer := make([]byte, 1024)
			_, err := s.stdout.Read(buffer)
			if err != nil {
				if err != io.EOF {
					writer.Write([]byte(fmt.Sprintf("%sError reading:%v", title, err)))
				}
				return
			}
			writer.Write([]byte(fmt.Sprintf("%s%s", title, buffer)))
		}
	}()
}

func (s *Session) ListenErr(writer io.Writer) {
	go func() {
		title := fmt.Sprintf("====== %s ======\r\n", s.client.RemoteAddr().String())
		for {
			select {
			case <-s.closeChan:
				log.Println("chan closed")
				return
			case <-s.ctx.Done():
				log.Println("context done")
				return
			default:
			}
			buffer := make([]byte, 1024)
			_, err := s.stderr.Read(buffer)
			if err != nil {
				if err != io.EOF {
					writer.Write([]byte(fmt.Sprintf("%sError reading:%v", title, err)))
				}
				return
			}
			writer.Write([]byte(fmt.Sprintf("%serror message:%s", title, buffer)))
		}
	}()
}

func (s *Session) Close() {
	s.Session.Close()
	s.closeChan <- struct{}{}
}

type exporter struct {
	mutex *sync.RWMutex
}

func NewExporter() *exporter {
	return &exporter{
		mutex: &sync.RWMutex{},
	}
}

func (e *exporter) Write(p []byte) (n int, err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.Println(string(p))
	return len(p), nil
}
