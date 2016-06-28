package showcaser

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// catchInterrupt traps termination signals to save a snapshot.
func catchInterrupt() {
	select {
	case <-signals:
		log.Warn("Termination signal caught, saving the showcaser entries...")
		err := data.save()
		if err != nil {
			log.Errorf("Got error while saving: %v", err)
		}
		close(signals)
	}
}

func periodicallyRebuildPopularList() {
	for {
		select {
		case <-signals:
			return
		case <-time.After(10 * time.Minute):
			err := data.updatePopular()
			if err != nil {
				log.Errorf("Got error while updating the popular list: %v", err)
			}
			log.Debug("Rebuilt the popular list")
		}
	}
}

func periodicallySave() {
	for {
		select {
		case <-signals:
			return
		case <-time.After(30 * time.Minute):
			if err := data.save(); err != nil {
				log.Errorf("Error while serializing index: %v", err)
				continue
			}
			log.Debug("Index persisted in database")
		}
	}
}
