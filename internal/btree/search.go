package btree

import "errors"

func (n *innerNode[T]) Search(value T, comparator func(first, second T) ComparatorStatus) (T, error) {
	index, status := n.getIndex(value, comparator)
	if status == FoundMatch {
		return n.values[index].maxValue, nil
	} else if status == ElementAfterIndexIsBigger {
		return n.values[index].node.Search(value, comparator)
	}

	var zeroValue T
	if status == ValueNotComparable {
		return zeroValue, errors.New("value not comparable with given comparator")
	}
	return zeroValue, errors.New("did not find the value")
}
