package utils

func FormatPlayTime(playTime int) PlayTime {
	hours := playTime / 60
	minutes := playTime % 60
	return PlayTime{
		Hours:   hours,
		Minutes: minutes,
	}
}
