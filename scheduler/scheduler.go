package scheduler

import (
	"billing_engine/services/billing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func NewScheduler(billingSvc billing.Service) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	c := cron.New(
		cron.WithLocation(loc))
	defer c.Stop()

	c.AddFunc("@daily", func() {
		err := billingSvc.UpdateDeliquency()
		if err != nil {
			logrus.Errorf("Error running cron job for update deliquency; err = %v", err)
		}
	})
	logrus.Info("Scheduler Started....")
}
