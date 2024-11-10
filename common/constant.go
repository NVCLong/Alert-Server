package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	ADMIN_ROLE = "user_admin"
	USER_ROLE  = "user_member"
)

var (
	NUMERIC_OPERATOR  = []string{">", "<", ">=", "<="}
	EQUAL_OPERATOR    = []string{"=", "!="}
	CONTAIN_OPERATOR  = []string{" IN ", " NOT IN "}
	CONDITION_KEYWORD = []string{"AND", "OR", "NOT"}
)

const (
	AT_LEAST_PATTERN = "AT LEAST ONE "
	TOTAL_PATERN     = "TOTAL "
	HAVE_MORE_PATERN = "HAVE MORE "
	HAVE             = "HAVE"
)

type Handler string

const (
	AT_LEAST_HANDLER  Handler = "handleAtLeastCase"
	TOTAL_HANDLER     Handler = "handleTotalCase"
	HAVE_MORE_HANDLER Handler = "handleHaveMoreCase"
	HAVE_HANDLER      Handler = "handleHaveCase"
)

const (
	ALERT_TYPE_WEEKLY   = "weekly"
	ALERT_TYPE_SCHEDULE = "schedule"
	ALERT_TYPE_DAILY    = "daily"
)

var (
	DATES = []string{"MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}
)

var DayOffsets = map[string]int{
	"MON":  1,
	"TUE":  2,
	"WED":  3,
	"THUR": 4,
	"FRI":  5,
	"SAT":  6,
	"SUN":  0,
}

const TracingIDKey = "TracingID"

func GetTracingIDFromContext(ctx *gin.Context) string {
	tracingID, exists := ctx.Get(TracingIDKey)
	if !exists {
		return "UnknownTracingID"
	}
	return tracingID.(string)
}

func SetTraceIDHeader(ctx *gin.Context, traceID any) {
	strID, ok := traceID.(string)
	if !ok {
		strID = fmt.Sprintf("%v", traceID)
	}
	ctx.Header("Trace-Id", strID)
}
