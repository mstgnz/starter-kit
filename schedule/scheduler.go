package schedule

import (
	"log"
	"time"
	_ "time/tzdata"

	"github.com/cemilsahin/arabamtaksit/internal/config"
	"github.com/robfig/cron/v3"
)

// https://crontab.guru/
func CallSchedule(c *cron.Cron) {
	// set location
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		log.Println(err)
	}

	cron.WithLocation(loc)

	// Table and Column Set
	// At 02:00 on day-of-month 1.
	if _, err = c.AddFunc("0 2 1 * *", func() {
		config.ShuttingWrapper(func() {
			//SetTableColumn()
		})

	}); err != nil {
		log.Println("AddFunc SetTableColumn", err)
	}

	// Set Permission For Center Admin
	// At 03:00 on day-of-month 1.
	if _, err = c.AddFunc("0 3 1 * *", func() {
		config.ShuttingWrapper(func() {
			//SetPermissionForCenterAdmin()
		})

	}); err != nil {
		log.Println("AddFunc SetPermissionForCenterAdmin", err)
	}
}
