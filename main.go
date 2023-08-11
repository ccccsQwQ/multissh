package main

import (
	"context"
	"multissh/utils"
	"os"
	"sync"
)

func main() {
	isHelp := utils.CheckHelp()
	if isHelp {
		return
	}
	args := os.Args
	cmd := args[1]
	filepath := args[2]
	file, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	writer := utils.NewExporter()
	ctx := context.Background()
	clients, err := utils.GetClientsFromJson(file)
	if err != nil {
		panic(err.Error())
	}
	wg := &sync.WaitGroup{}
	for _, client := range clients {
		client := client
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer client.Close()
			utils.RunCmd(ctx, client, cmd, writer)
		}()
	}
	wg.Wait()
}
