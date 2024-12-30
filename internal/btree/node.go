package btree

type node struct {
	values []interface{}
	parent *node
	isLeaf bool // TODO consider working away from this field
}

type nodeReference struct {
	maxValue interface{}
	node     *node
}

type ComparatorStatus int

const (
	FirstArgumentBigger ComparatorStatus = iota
	SecondArgumentBigger
	Equal
	ArgumentsNotComparable
)

type GetIndexStatus int

const (
	FoundMatch GetIndexStatus = iota
	NoElementMatchedOrWereBigger
	ElementAfterIndexIsBigger
	ValueNotComparable
	NoElementsInList
)

func (n *node) getIndex(value interface{}, comparator func(first, second interface{}) ComparatorStatus) (int, GetIndexStatus) {
	if len(n.values) == 0 {
		return -1, NoElementsInList
	}
	var index int
	for i, v := range n.values {
		var comparatorResult ComparatorStatus
		switch v.(type) {
		case nodeReference:
			comparatorResult = comparator(value, v.(nodeReference).maxValue)
		default:
			comparatorResult = comparator(value, v)
		}
		if comparatorResult == Equal {
			return i, FoundMatch
		} else if comparatorResult == SecondArgumentBigger {
			return i, ElementAfterIndexIsBigger
		} else if comparatorResult == ArgumentsNotComparable {
			return -1, ValueNotComparable
		}
		index = i
	}
	return index, NoElementMatchedOrWereBigger
}
