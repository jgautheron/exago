package showcaser

import (
	"os/signal"
	"syscall"
	"time"

	. "github.com/hotolab/exago-svc/config"
)

// catchInterrupt traps termination signals.
func (s *Showcase) catchInterrupt() {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case <-signals:
		logger.Warn("Termination signal caught, saving the showcaser entries...")
		err := s.save()
		if err != nil {
			logger.Errorf("Got error while saving: %v", err)
		}
		close(signals)
	}
}

func (s *Showcase) periodicallyRebuildPopularList() {
	for {
		select {
		case <-signals:
			logger.Debug("periodicallyRebuildPopularList closing")
			return
		case <-time.After(Config.ShowcaserPopularRebuildInterval):
			err := s.updatePopular()
			if err != nil {
				logger.Errorf("Got error while updating the popular list: %v", err)
			}
			logger.Debug("Rebuilt the popular list")
		}
	}
}

// func (s *Showcase) periodicallySave() {
// 	for {
// 		select {
// 		case <-signals:
// 			return
// 		case <-time.After(30 * time.Minute):
// 			if err := s.save(); err != nil {
// 				logger.Errorf("Error while serializing index: %v", err)
// 				continue
// 			}
// 			logger.Debug("Index persisted in database")
// 		}
// 	}
// }
