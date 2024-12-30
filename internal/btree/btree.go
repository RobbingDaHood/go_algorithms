package btree

import (
	"fmt"
)

type Root struct {
	node
	comparator  func(first, second interface{}) ComparatorStatus
	nodeMaxSize int
	nodeCount   int
}

func CreateTreeDefaultValues() Root {
	return Root{
		node: node{
			isLeaf: true,
		},
		comparator:  ComparatorEverythingAsString,
		nodeMaxSize: 10000,
	}
}

func (n *Root) Insert(value interface{}) error {
	err := n.node.insert(value, n.comparator, n.nodeMaxSize, -1)
	if err == nil {
		n.nodeCount++
	}
	return err
}

func (n *Root) Search(value interface{}) (interface{}, error) {
	return n.node.Search(value, n.comparator)
}

func ComparatorEverythingAsString(first, second interface{}) ComparatorStatus {
	firstAsString := fmt.Sprint(first)
	secondAsString := fmt.Sprint(second)
	if firstAsString > secondAsString {
		return FirstArgumentBigger
	} else if firstAsString < secondAsString {
		return SecondArgumentBigger
	}
	return Equal
}
