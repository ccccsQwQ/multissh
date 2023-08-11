package utils

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

// RunCmd
func RunCmd(ctx context.Context, client *ssh.Client, cmd string, writer io.Writer) {
	session, err := NewSession(ctx, client)
	if err != nil {
		log.Printf("fail to new session: %s\r\n", err)
		return
	}
	defer session.Close()

	err = session.Start(cmd)
	if err != nil {
		log.Printf("Failed to start command: %s\r\n", err)
		return
	}
	session.ListenOut(writer)
	session.ListenErr(writer)
	err = session.Wait()
	if err != nil {
		log.Printf("Command execution error: %s\r\n", err)
		return
	}
}

func NewClient(network string, config *Config) (*ssh.Client, error) {
	return ssh.Dial(network, config.Addr, config.ClientConfig)
}

func CheckHelp() bool {
	help := flag.Bool("help", false, "显示帮助信息")
	flag.Parse()
	if *help {
		fmt.Fprintf(os.Stderr, "用法: %s [执行命令] [配置文件路径]\r\n", os.Args[0])
		flag.PrintDefaults()
		return true
	}
	return false
}
