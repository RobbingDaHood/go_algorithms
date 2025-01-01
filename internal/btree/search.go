package btree

import "errors"

func (n *innerNode[T]) Search(value T, comparator func(first, second T) ComparatorStatus) (T, error) {
	index, status := n.getIndex(value, comparator)
	var zeroValue T
	if status == FoundMatch {
		return n.values[index].maxValue, nil
	} else if status == ElementAtIndexIsBigger {
		childNode := n.values[index].node
		if childNode != nil {
			return childNode.Search(value, comparator)
		} else {
			return zeroValue, errors.New("did not find the value")
		}
	}

	if status == ValueNotComparable {
		return zeroValue, errors.New("value not comparable with given comparator")
	}
	return zeroValue, errors.New("did not find the value")
}
