package oncron

import (
	"github.com/OliveTin/OliveTin/internal/acl"
	"github.com/OliveTin/OliveTin/internal/config"
	"github.com/OliveTin/OliveTin/internal/executor"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func Schedule(cfg *config.Config, ex *executor.Executor) {
	scheduler := cron.New(cron.WithSeconds())

	for _, action := range cfg.Actions {
		for _, cronline := range action.ExecOnCron {
			scheduleAction(cfg, scheduler, cronline, ex, action)
		}
	}

	scheduler.Start()
}

func scheduleAction(cfg *config.Config, scheduler *cron.Cron, cronline string, ex *executor.Executor, action config.Action) {
	log.WithFields(log.Fields{
		"action":   action.Title,
		"cronline": cronline,
	}).Infof("Scheduling Action for cron")

	_, err := scheduler.AddFunc(cronline, func() {
		req := &executor.ExecutionRequest{
			ActionName: action.Title,
			Cfg:        cfg,
			Tags:       []string{"cron"},
			AuthenticatedUser: &acl.AuthenticatedUser{
				Username: "cron",
			},
		}

		ex.ExecRequest(req)
	})

	if err != nil {
		log.WithFields(log.Fields{
			"action":    action.Title,
			"cronError": err,
		}).Errorf("CRON schedule error")
		return
	}

}
