package btree

import "errors"

func (n *node) Search(value interface{}, comparator func(first, second interface{}) ComparatorStatus) (interface{}, error) {
	index, status := n.getIndex(value, comparator)
	if status == FoundMatch {
		matchedValue := n.values[index]
		switch matchedValue.(type) {
		case nodeReference:
			return matchedValue.(nodeReference).maxValue, nil
		default:
			return matchedValue, nil
		}
	} else if status == ElementAfterIndexIsBigger {
		matchedValue := n.values[index]
		switch matchedValue.(type) {
		case nodeReference:
			return matchedValue.(nodeReference).node.Search(value, comparator)
		default:
			return nil, errors.New("did not find the value")
		}
	} else if status == ValueNotComparable {
		return nil, errors.New("value not comparable with given comparator")
	}
	return nil, errors.New("did not find the value")
}
