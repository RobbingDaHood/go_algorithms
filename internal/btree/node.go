package btree

type innerNode[T any] struct {
	values []nodeReference[T]
	parent *innerNode[T]
	isLeaf bool //TODO see if I can work out of this one
}

type nodeReference[T any] struct {
	maxValue T
	node     *innerNode[T]
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

func (n *innerNode[T]) getIndex(value T, comparator func(first, second T) ComparatorStatus) (int, GetIndexStatus) {
	if len(n.values) == 0 {
		return -1, NoElementsInList
	}
	var index int
	for i, v := range n.values {
		var comparatorResult = comparator(value, v.maxValue)
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
