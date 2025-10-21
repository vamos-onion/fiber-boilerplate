package util

type boolBlock struct{}

// variables
var (
	intMap    = map[bool]int{false: 0, true: 1}
	stringMap = map[bool]string{false: "false", true: "true"}
)

// ToInt :
func (b boolBlock) ToInt(value bool) int {
	return intMap[value]
}

// ToString :
func (b boolBlock) ToString(value bool) string {
	return stringMap[value]
}
