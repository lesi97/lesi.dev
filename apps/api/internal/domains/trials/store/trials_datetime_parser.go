package trials_store

import "time"

func parseTrialsUpdatedAt(data *TrialsData) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, data.Platforms.Num0.RecentStats.UpdatedAt)
}

func isTrialsFresh(updatedAt time.Time) bool {
	return time.Since(updatedAt) < 90*time.Minute
}