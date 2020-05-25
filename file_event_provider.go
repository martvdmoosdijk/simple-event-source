package simple_event_source

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

var _ IEventProvider = &FileEventProvider{}

type FileEventProvider struct {
	File string
}

func (self *FileEventProvider) ReadEvents() <-chan EventEntry {
	if self.File == "" {
		self.File = "events.log"
	}

	/* Create file if not exists */
	_, err := os.Stat(self.File)
	if os.IsNotExist(err) {
		err := ioutil.WriteFile(self.File, []byte(""), 0644)
		if err != nil {
			panic(err)
		}
	}

	c := make(chan EventEntry)
	go func() {
		defer close(c)

		file, err := os.Open(self.File)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		/* Read line by line and send into channel */
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			entry := EventEntry{}
			if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
				panic(err)
			}
			c <- entry
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	return c
}

func (self *FileEventProvider) SaveEvent(entry EventEntry) error {
	f, err := os.OpenFile(self.File, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(string(data) + "\n"); err != nil {
		return err
	}

	return nil
}
