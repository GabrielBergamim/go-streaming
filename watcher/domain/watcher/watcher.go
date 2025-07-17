package watcher

import (
	"log"
	"time"

	"github.com/example/go-streaming/watcher/domain/producer"
	"github.com/radovskyb/watcher"
)

type Watcher struct {
	FolderPath string
}

func (w *Watcher) Start(producer *producer.Producer) error {
	watch := watcher.New()
	watch.SetMaxEvents(1)
	watch.FilterOps(watcher.Create)

	if err := watch.Add(w.FolderPath); err != nil {
		log.Println("Error adding folder to watcher:", err)
		return err
	}

	go func() {
		for {
			select {
			case event := <-watch.Event:
				log.Println("Detected:", event) // Create, Remove, etc.
				if err := producer.SendEvent(event); err != nil {
				log.Println("Error sending event:", err)
					continue
				}
			case err := <-watch.Error:
				log.Println("Error:", err)
			case <-watch.Closed:
				log.Println("Watcher closed")
				return
			}
		}
	}()

	if err := watch.Start(time.Second * 10); err != nil {
		log.Println("Error starting watcher:", err)
		return err;
	}

	return nil
}
