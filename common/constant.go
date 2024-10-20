package common

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
