package watcher

import (
	"log"

	"github.com/example/go-streaming/watcher/domain/producer"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	FolderPath string
}

func (w *Watcher) Start(producer *producer.Producer) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan struct{})

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("Watcher channel closed")
					return
				}

				log.Println("Event:", event)
				producer.SendEvent(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(w.FolderPath)
	if err != nil {
		log.Println("Error adding folder to watcher:", err)
		log.Fatal(err)
		return err
	}
	<-done

	return nil
}

