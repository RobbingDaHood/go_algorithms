package btree

type node struct {
	values []interface{}
}

type ComparatorStatus int

const (
	FirstArgumentBigger ComparatorStatus = iota
	SecondArgumentBigger
	Equal
	ArgumentsNotComparable
)

type Root struct {
	node
	comparator func(first, second interface{}) ComparatorStatus
}
