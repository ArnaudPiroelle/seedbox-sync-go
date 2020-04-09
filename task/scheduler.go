package task

import (
	"seedbox-sync/downloader"
	"seedbox-sync/model"
	"seedbox-sync/notifier"
	"seedbox-sync/provider"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	sync          *Sync
	configuration model.SchedulerConfiguration
	c             *cron.Cron
}

func (s *Scheduler) Execute() {
	s.c.AddFunc(s.configuration.Cron, s.sync.Execute)
	s.c.Run()
}

func NewSchedule(configuration model.SchedulerConfiguration, folders []model.Folder, provider provider.Provider, downloader downloader.Downloader, notifier notifier.Notifier) *Scheduler {
	sync := &Sync{
		folders:    folders,
		provider:   provider,
		downloader: downloader,
		notifier:   notifier,
	}
	return &Scheduler{
		configuration: configuration,
		sync:          sync,
		c:             cron.New(),
	}

}
