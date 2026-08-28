package cron

import (
	"testing"
)

func TestCronValidationAndSync(t *testing.T) {
	if err := ValidateSchedule("* * * * *"); err != nil {
		t.Errorf("Expected valid schedule, got error: %v", err)
	}

	if err := ValidateSchedule("invalid cron"); err == nil {
		t.Errorf("Expected error for invalid cron schedule, got nil")
	}

	mgr := NewManager(true)
	jobs := []JobItem{
		{Schedule: "* * * * *", Command: "php /var/www/app/artisan schedule:run", IsActive: true},
		{Schedule: "0 0 * * *", Command: "php /var/www/app/artisan backup:clean", IsActive: false},
	}

	if err := mgr.SyncUserCrontab("kp_test", jobs); err != nil {
		t.Fatalf("SyncUserCrontab failed: %v", err)
	}
}
