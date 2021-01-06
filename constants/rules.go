package constants

const (
	operatorsValue             = ">=|<=|>|<|!~|!|~"
	IsAObjectId                = "^[0-9a-fA-F]{24}$"
	IsBetweenDoubleQuotes      = "(?:\")(.*)(?:\")"
	Operators                  = "^" + operatorsValue
	ContainsOperatorsOrFilters = "(^" + operatorsValue + ")|" + ContainsOrOperation
	ValueWithoutOperators      = "^(?:" + operatorsValue + ")?(.+?)(?:\\[\\d{0,3}?])?$"
	ContainsOrOperation        = "(\\[\\d{0,3}?]$)"
	ValueWithoutOrOperator     = "(.*(?=\\[\\d{0,3}?])|.*$)"
	StartsWithDeny             = "^!(?!~)"
	IsAnInteger                = "^[-]?\\d+$"
	IsADate                    = "(?:\\d{4}-\\d{2}-\\d{2}(?:[ T](?:\\d{2}:\\d{2}:\\d{2}|\\d{2}:\\d{2}$)([Z]$)?(?:[.+]?\\d{2,3})?(?:-?\\d{2}:?\\d{2}|:\\d{2})?)?)$"
	IsAFloat                   = "^[+-]?([0-9]*[.,])?[0-9]+$"
)

func IsBetween(valueBefore string, valueAfter string) string {
	return "\\[(" + valueBefore + ".*?" + valueAfter + ")\\]"
}
