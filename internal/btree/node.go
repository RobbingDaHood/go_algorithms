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
	ElementAtIndexIsBigger
	ValueNotComparable
	NoElementsInList
)

func (n *innerNode[T]) getIndex(value T, comparator func(first, second T) ComparatorStatus) (int, GetIndexStatus) {
	if n.values == nil || len(n.values) == 0 {
		return -1, NoElementsInList
	}
	return Search(len(n.values), func(i int) ComparatorStatus {
		return comparator(value, n.values[i].maxValue)
	})
}

func Search(n int, f func(int) ComparatorStatus) (int, GetIndexStatus) {
	i, j := 0, n
	var didFindBigger = false
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		status := f(h)
		if status == ArgumentsNotComparable {
			return -1, ValueNotComparable
		} else if status == Equal {
			return h, FoundMatch
		} else if status == FirstArgumentBigger {
			i = h + 1
		} else {
			didFindBigger = true
			j = h
		}
	}

	if didFindBigger {
		return i, ElementAtIndexIsBigger
	} else {
		return n - 1, NoElementMatchedOrWereBigger
	}
}
