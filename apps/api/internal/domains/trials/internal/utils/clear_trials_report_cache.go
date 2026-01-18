package utils

func ClearTrialsReportCache() {
	trialsReportCacheMu.Lock()
	trialsReportCache = nil
	trialsReportCacheMu.Unlock()
}
