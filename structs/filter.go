package structs

type Filter struct {
	Group           KeyValuePair
	Field           string
	Value           interface{}
	Operation       string
	IsSpecialFilter bool
}

