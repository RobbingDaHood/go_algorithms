package btree

import (
	"errors"
	"fmt"
)

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
	FoundIndexToInsert
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
			return i, FoundIndexToInsert
		} else if comparatorResult == ArgumentsNotComparable {
			return -1, ValueNotComparable
		}
		index = i
	}
	return index, FoundIndexToInsert
}

func (n *node) Insert(value interface{}, comparator func(first, second interface{}) ComparatorStatus, nodeMaxSize int) (bool, error) {
	index, status := n.getIndex(value, comparator)
	switch status {
	case NoElementsInList:
		n.values = append(n.values, value)
		return false, nil
	case FoundIndexToInsert:
		if n.isLeaf {
			n.insertElementAfterIndex(value, index)
		} else {
			nodeToInsertInCasted := n.values[index].(nodeReference)
			shouldSplitNode, err := nodeToInsertInCasted.node.Insert(value, comparator, nodeMaxSize)
			if err != nil {
				return false, err
			}
			if shouldSplitNode {
				n.splitChildNode(nodeToInsertInCasted, index)
			}
		}

		if len(n.values) > nodeMaxSize {
			if n.parent == nil {
				n.splitWithoutParent()
			} else {
				return true, nil
			}
		}

		return false, nil
	case FoundMatch:
		return false, errors.New("value already exists in tree")
	case ValueNotComparable:
		return false, errors.New("value not comparable with given comparator")
	default:
		return false, errors.New("unexpected status from getIndex: " + fmt.Sprint(status))
	}
}

func (n *node) insertElementAfterIndex(value interface{}, index int) {
	// TODO enable this after checking performamnce, need to check for cap or else we get runtime errors
	//n.values = append(n.values, 0)
	//copy(n.values[index+1:], n.values[index:])
	//n.values[index] = value
	smaller := n.values[:index+1]
	bigger := n.values[index+1:]
	valueSlice := []interface{}{value}
	result := append(smaller, append(valueSlice, bigger...)...)
	n.values = result
}

func (n *node) splitChildNode(nodeToInsertInCasted nodeReference, index int) {
	valuesFromNode := nodeToInsertInCasted.node.values
	halfWayIndex := len(valuesFromNode) / 2

	smallerValues := valuesFromNode[:halfWayIndex]
	smallerValuesMaxValue := smallerValues[len(smallerValues)-1]
	biggerValues := valuesFromNode[halfWayIndex:]
	biggerValueMaxValue := biggerValues[len(biggerValues)-1]

	nodeToInsertInCasted.node.values = smallerValues
	nodeToInsertInCasted.maxValue = smallerValuesMaxValue

	biggerNodeReference := nodeReference{
		maxValue: biggerValueMaxValue,
		node: &node{
			values: biggerValues,
			parent: n,
			isLeaf: nodeToInsertInCasted.node.isLeaf,
		},
	}

	n.insertElementAfterIndex(biggerNodeReference, index)
}

func (n *node) splitWithoutParent() {
	halfWayIndex := len(n.values) / 2

	smallerValues := n.values[:halfWayIndex]
	smallerValuesMaxValue := smallerValues[len(smallerValues)-1]
	biggerValues := n.values[halfWayIndex:]
	biggerValueMaxValue := biggerValues[len(biggerValues)-1]

	n.values = []interface{}{
		nodeReference{
			maxValue: smallerValuesMaxValue,
			node: &node{
				values: smallerValues,
				parent: n,
				isLeaf: true,
			},
		},
		nodeReference{
			maxValue: biggerValueMaxValue,
			node: &node{
				values: biggerValues,
				parent: n,
				isLeaf: true,
			},
		},
	}
	n.isLeaf = false
}

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
	} else if status == FoundIndexToInsert {
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
