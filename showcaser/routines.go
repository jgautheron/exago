package showcaser

import (
	"os/signal"
	"syscall"
	"time"
)

// catchInterrupt traps termination signals.
func catchInterrupt() {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-signals:
		logger.Warn("Termination signal caught, saving the showcaser entries...")
		err := data.save()
		if err != nil {
			logger.Errorf("Got error while saving: %v", err)
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
				logger.Errorf("Got error while updating the popular list: %v", err)
			}
			logger.Debug("Rebuilt the popular list")
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
				logger.Errorf("Error while serializing index: %v", err)
				continue
			}
			logger.Debug("Index persisted in database")
		}
	}
}
