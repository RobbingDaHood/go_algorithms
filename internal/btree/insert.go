package btree

import (
	"errors"
	"fmt"
)

func (n *node) Insert(value interface{}, comparator func(first, second interface{}) ComparatorStatus, nodeMaxSize int) (bool, error) {
	index, status := n.getIndex(value, comparator)
	switch status {
	case NoElementsInList:
		n.values = append(n.values, value)
		return false, nil
	case ElementAfterIndexIsBigger:
		if n.isLeaf {
			n.insertElementBeforeIndex(value, index)
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
	case NoElementMatchedOrWereBigger:
		if n.isLeaf {
			n.values = append(n.values, value)
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

func (n *node) splitChildNode(nodeToInsertInCasted nodeReference, index int) {
	valuesFromNode := nodeToInsertInCasted.node.values
	halfWayIndex := len(valuesFromNode) / 2

	smallerValues := valuesFromNode[:halfWayIndex]
	smallerValuesMaxValue := smallerValues[len(smallerValues)-1]
	biggerValues := valuesFromNode[halfWayIndex:]
	biggerValueMaxValue := biggerValues[len(biggerValues)-1]

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
