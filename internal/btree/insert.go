package btree

import (
	"errors"
	"slices"
)

func (n *innerNode[T]) insert(value T, comparator func(first, second T) ComparatorStatus, nodeMaxSize int, indexFromParent int) error {
	index, status := n.getIndex(value, comparator)
	switch status {
	case NoElementsInList:
		n.values = append(n.values, nodeReference[T]{maxValue: value})
		return nil
	case ElementAtIndexIsBigger:
		if n.isLeaf {
			n.values = slices.Insert(n.values, index, nodeReference[T]{maxValue: value})
		} else {
			relevantChildNode := &n.values[index]
			err := relevantChildNode.node.insert(value, comparator, nodeMaxSize, index)
			if err != nil {
				return err
			}
		}
	case NoElementMatchedOrWereBigger:
		if n.isLeaf {
			n.values = append(n.values, nodeReference[T]{maxValue: value})
		} else {
			relevantChildNode := &n.values[index]
			relevantChildNode.maxValue = value
			err := relevantChildNode.node.insert(value, comparator, nodeMaxSize, index)
			if err != nil {
				return err
			}
		}
	case FoundMatch:
		return errors.New("value already exists in tree")
	case ValueNotComparable:
		return errors.New("value not comparable with given comparator")
	}

	n.handlePossibleSplit(nodeMaxSize, indexFromParent)
	return nil
}

func (n *innerNode[T]) handlePossibleSplit(nodeMaxSize int, indexFromParent int) {
	if len(n.values) > nodeMaxSize {
		if n.parent == nil {
			n.splitWithoutParent()
		} else {
			n.parent.splitChildNode(&n.parent.values[indexFromParent], indexFromParent)
		}
	}
}

func (n *innerNode[T]) splitChildNode(nodeToSplit *nodeReference[T], index int) {
	smallerValues, smallerValuesMaxValue, biggerValues, biggerValueMaxValue := nodeToSplit.node.splitIntoTwo()

	nodeToSplit.node.values = biggerValues
	nodeToSplit.maxValue = biggerValueMaxValue

	smallerNodeReference := nodeReference[T]{
		maxValue: smallerValuesMaxValue,
		node: &innerNode[T]{
			values: smallerValues,
			parent: n,
			isLeaf: nodeToSplit.node.isLeaf,
		},
	}

	n.values = slices.Insert(n.values, index, smallerNodeReference)
}

func getMaxValue[T any](biggerValues []nodeReference[T]) T {
	return biggerValues[len(biggerValues)-1].maxValue
}

func (n *innerNode[T]) splitWithoutParent() {
	smallerValues, smallerValuesMaxValue, biggerValues, biggerValueMaxValue := n.splitIntoTwo()

	n.values = []nodeReference[T]{
		{
			maxValue: smallerValuesMaxValue,
			node: &innerNode[T]{
				values: smallerValues,
				parent: n,
				isLeaf: true,
			},
		},
		{
			maxValue: biggerValueMaxValue,
			node: &innerNode[T]{
				values: biggerValues,
				parent: n,
				isLeaf: true,
			},
		},
	}
	n.isLeaf = false
}

func (n *innerNode[T]) splitIntoTwo() ([]nodeReference[T], T, []nodeReference[T], T) {
	halfWayIndex := len(n.values) / 2

	smallerValues, smallerValuesMaxValue := n.copySlice(n.values[:halfWayIndex])
	biggerValues, biggerValueMaxValue := n.copySlice(n.values[halfWayIndex:])
	return smallerValues, smallerValuesMaxValue, biggerValues, biggerValueMaxValue
}

func (n *innerNode[T]) copySlice(subSetSlice []nodeReference[T]) ([]nodeReference[T], T) {
	subsetCopy := make([]nodeReference[T], len(subSetSlice))
	copy(subsetCopy, subSetSlice)
	smallerValuesMaxValue := getMaxValue(subsetCopy)
	return subsetCopy, smallerValuesMaxValue
}
