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
			n.insertNodeReferenceBeforeIndex(nodeReference[T]{maxValue: value}, index)
		}
	case NoElementMatchedOrWereBigger:
		if n.isLeaf {
			n.values = append(n.values, nodeReference[T]{maxValue: value})
		}
	case FoundMatch:
		return errors.New("value already exists in tree")
	case ValueNotComparable:
		return errors.New("value not comparable with given comparator")
	}

	if !n.isLeaf {
		relevantChildNode := &n.values[index]
		switch status {
		case NoElementMatchedOrWereBigger:
			relevantChildNode.maxValue = value
		case ElementAtIndexIsBigger:
			// Do nothing
		case FoundMatch:
			// Do nothing
		case ValueNotComparable:
			// Do nothing
		case NoElementsInList:
			// Do nothing
		}
		err := relevantChildNode.node.insert(value, comparator, nodeMaxSize, index)
		if err != nil {
			return err
		}
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
	valuesFromNode := nodeToSplit.node.values
	halfWayIndex := len(valuesFromNode) / 2

	tmpSmallerValues := valuesFromNode[:halfWayIndex]
	smallerValues := make([]nodeReference[T], len(tmpSmallerValues))
	copy(smallerValues, tmpSmallerValues)
	smallerValuesMaxValue := getMaxValue(smallerValues)
	tmpBiggerValues := valuesFromNode[halfWayIndex:]
	biggerValues := make([]nodeReference[T], len(tmpBiggerValues))
	copy(biggerValues, tmpBiggerValues)
	biggerValueMaxValue := getMaxValue(biggerValues)

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

	n.insertNodeReferenceBeforeIndex(smallerNodeReference, index)
}

func getMaxValue[T any](biggerValues []nodeReference[T]) T {
	return biggerValues[len(biggerValues)-1].maxValue
}

func (n *innerNode[T]) splitWithoutParent() {
	halfWayIndex := len(n.values) / 2

	tmpSmallerValues := n.values[:halfWayIndex]
	smallerValues := make([]nodeReference[T], len(tmpSmallerValues))
	copy(smallerValues, tmpSmallerValues)
	smallerValuesMaxValue := getMaxValue(smallerValues)
	tmpBiggerValues := n.values[halfWayIndex:]
	biggerValues := make([]nodeReference[T], len(tmpBiggerValues))
	copy(biggerValues, tmpBiggerValues)
	biggerValueMaxValue := getMaxValue(biggerValues)

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

func (n *innerNode[T]) insertNodeReferenceBeforeIndex(value nodeReference[T], index int) {
	// TODO enable this after checking performamnce, need to check for cap or else we get runtime errors
	//n.values = append(n.values, 0)
	//copy(n.values[index+1:], n.values[index:])
	//n.values[index] = value
	n.values = slices.Insert(n.values, index, value)
	//smaller := n.values[:index]
	//bigger := n.values[index:]
	//valueSlice := []nodeReference[T]{value}
	//n.values = append(n.values, nodeReference[T]{})
	//result := append(smaller, append(valueSlice, bigger...)...)
	//n.values = result
}
