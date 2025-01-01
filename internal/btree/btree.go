package btree

import (
	"fmt"
)

type Root[T any] struct {
	innerNode[T]
	comparator  func(first, second T) ComparatorStatus
	nodeMaxSize int
	nodeCount   int
}

func CreateTreeDefaultValues[T any]() Root[T] {
	defaultMaxSize := 10000
	return Root[T]{
		innerNode: innerNode[T]{
			//values: make([]nodeReference[T], 0, defaultMaxSize+1),
			isLeaf: true,
		},
		comparator:  ComparatorEverythingAsString[T],
		nodeMaxSize: defaultMaxSize,
	}
}

func (n *Root[T]) Insert(value T) error {
	err := n.innerNode.insert(value, n.comparator, n.nodeMaxSize, -1)
	if err == nil {
		n.nodeCount++
	}
	return err
}

func (n *Root[T]) Search(value T) (T, error) {
	return n.innerNode.Search(value, n.comparator)
}

func ComparatorEverythingAsString[T any](first, second T) ComparatorStatus {
	firstAsString := fmt.Sprint(first)
	secondAsString := fmt.Sprint(second)
	if firstAsString > secondAsString {
		return FirstArgumentBigger
	} else if firstAsString < secondAsString {
		return SecondArgumentBigger
	}
	return Equal
}
