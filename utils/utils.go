package utils

import (
	"bufio"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// LoadConfigFromFile from file
func LoadConfigFromFile(filepath string, cfg interface{}) error {
	if confContent, err := ioutil.ReadFile(filepath); err != nil {
		return err
	} else if err := yaml.Unmarshal([]byte(confContent), cfg); err != nil {
		return fmt.Errorf("parser %s error. %v", filepath, err)
	}
	return nil
}

func RunTimeout(fn func(), millisecond int64) bool {
	var job sync.WaitGroup
	chTimeout := make(chan struct{})

	job.Add(1)
	go func() {
		fn()
		job.Done()
	}()
	go func() {
		job.Wait()
		chTimeout <- struct{}{}
	}()

	select {
	case <-time.After(time.Millisecond * time.Duration(millisecond)):
		return true
	case <-chTimeout:
		return false
	}
}

func GetInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	return scanner.Text()
}
