package common

import "time"

func FormatDate(timeDate time.Time) string {
	formattedDate := timeDate.Format("01-02-2006")
	return formattedDate
}
