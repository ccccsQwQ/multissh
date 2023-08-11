package utils

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	file, err := os.Open("../config.json")
	assert.NoError(t, err)
	defer file.Close()
	ctx := context.Background()
	clients, err := GetClientsFromJson(file)
	assert.NoError(t, err)
	wg := sync.WaitGroup{}
	writer := NewExporter()
	for _, client := range clients {
		c := client
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer c.Close()
			RunCmd(ctx, c, "tail -f /tmp/QAQ.txt", writer)
		}()
	}
	wg.Wait()
}
