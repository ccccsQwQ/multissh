package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Config struct {
	Addr string
	*ssh.ClientConfig
}

func NewPasswordConfig(addr string, user string, secret string) *Config {
	return &Config{
		Addr: addr,
		ClientConfig: &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(secret),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
}

type jsonConfig struct {
	DnsArr []dns
}

type dns struct {
	Addr       string
	User       string
	AuthMethod string
	Password   string
}

func GetClientsFromJson(file *os.File) ([]*ssh.Client, error) {
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config := &jsonConfig{}
	err = json.Unmarshal(byteValue, config)
	if err != nil {
		return nil, err
	}
	var result []*ssh.Client
	for _, dns := range config.DnsArr {
		if dns.AuthMethod != "password" {
			return nil, errors.New("config error, the authMethod only support to 'password'")
		}
		if dns.Password == "" {
			password, err := loadPassword(dns)
			if err != nil {
				return nil, err
			}
			dns.Password = password
		}
		client, err := NewClient("tcp", NewPasswordConfig(dns.Addr, dns.User, dns.Password))
		if err != nil {
			log.Printf("connect fail: %s\r\n", dns.Addr)
			return nil, err
		}
		result = append(result, client)
	}
	return result, nil
}

func loadPassword(dns dns) (string, error) {
	log.Println("please enter password for ", dns.User, "@", dns.Addr)
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}
