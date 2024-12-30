package btree

import (
	"errors"
	"fmt"
)

func (n *node) insert(value interface{}, comparator func(first, second interface{}) ComparatorStatus, nodeMaxSize int, indexFromParent int) error {
	index, status := n.getIndex(value, comparator)
	switch status {
	case NoElementsInList:
		n.values = append(n.values, value)
		return nil
	case ElementAfterIndexIsBigger:
		if n.isLeaf {
			n.insertElementBeforeIndex(value, index)
		}
	case NoElementMatchedOrWereBigger:
		if n.isLeaf {
			n.values = append(n.values, value)
		}
	case FoundMatch:
		return errors.New("value already exists in tree")
	case ValueNotComparable:
		return errors.New("value not comparable with given comparator")
	default:
		return errors.New("unexpected status from getIndex: " + fmt.Sprint(status))
	}

	if !n.isLeaf {
		nodeToInsertInCasted := n.values[index].(nodeReference)
		err := nodeToInsertInCasted.node.insert(value, comparator, nodeMaxSize, index)
		if err != nil {
			return err
		}
	}

	n.handlePossibleSplit(nodeMaxSize, indexFromParent)
	return nil
}

func (n *node) handlePossibleSplit(nodeMaxSize int, indexFromParent int) {
	if len(n.values) > nodeMaxSize {
		if n.parent == nil {
			n.splitWithoutParent()
		} else {
			n.parent.splitChildNode(n.parent.values[indexFromParent].(nodeReference), indexFromParent)
		}
	}
}

func (n *node) splitChildNode(nodeToInsertInCasted nodeReference, index int) {
	valuesFromNode := nodeToInsertInCasted.node.values
	halfWayIndex := len(valuesFromNode) / 2

	smallerValues := valuesFromNode[:halfWayIndex]
	smallerValuesMaxValue := getMaxValue(smallerValues)
	biggerValues := valuesFromNode[halfWayIndex:]
	biggerValueMaxValue := getMaxValue(biggerValues)

	nodeToInsertInCasted.node.values = biggerValues
	nodeToInsertInCasted.maxValue = biggerValueMaxValue

	biggerNodeReference := nodeReference{
		maxValue: smallerValuesMaxValue,
		node: &node{
			values: smallerValues,
			parent: n,
			isLeaf: nodeToInsertInCasted.node.isLeaf,
		},
	}

	n.insertElementBeforeIndex(biggerNodeReference, index)
}

func getMaxValue(biggerValues []interface{}) interface{} {
	biggerValue := biggerValues[len(biggerValues)-1]
	switch biggerValue.(type) {
	case nodeReference:
		getMaxValue(biggerValue.(nodeReference).node.values)
	default:
		return biggerValue
	}
	panic("unexpected type in values")
}

func (n *node) splitWithoutParent() {
	halfWayIndex := len(n.values) / 2

	smallerValues := n.values[:halfWayIndex]
	smallerValuesMaxValue := getMaxValue(smallerValues)
	biggerValues := n.values[halfWayIndex:]
	biggerValueMaxValue := getMaxValue(biggerValues)

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

func (n *node) insertElementBeforeIndex(value interface{}, index int) {
	// TODO enable this after checking performamnce, need to check for cap or else we get runtime errors
	//n.values = append(n.values, 0)
	//copy(n.values[index+1:], n.values[index:])
	//n.values[index] = value
	smaller := n.values[:index]
	bigger := n.values[index:]
	valueSlice := []interface{}{value}
	result := append(smaller, append(valueSlice, bigger...)...)
	n.values = result
}
