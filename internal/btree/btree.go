package btree

import (
	"errors"
	"fmt"
)

func CreateTreeDefaultValues() Root {
	return Root{
		node:       node{},
		comparator: ComparatorEverythingAsString,
	}
}

func (n *Root) Insert(value interface{}) (int, error) {
	index, status := getIndex(value, n)
	switch status {
	case FoundIndexToInsert:
		n.values = append(append(n.values[index:], value), n.values[:index]...)
		return index, nil
	case FoundMatch:
		return -1, errors.New("value already exists in tree")
	case ValueNotComparable:
		return -1, errors.New("value not comparable with given comparator")
	default:
		return -1, errors.New("unexpected status from getIndex: " + fmt.Sprint(status))
	}
}

type GetIndexStatus int

const (
	FoundMatch GetIndexStatus = iota
	FoundIndexToInsert
	ValueNotComparable
)

func getIndex(value interface{}, n *Root) (int, GetIndexStatus) {
	var index int
	for i, v := range n.values {
		comparatorResult := n.comparator(value, v)
		if comparatorResult == Equal {
			return i, FoundMatch
		} else if comparatorResult == ArgumentsNotComparable {
			return -1, ValueNotComparable
		}
		index = i
	}
	return index, FoundIndexToInsert
}

func (n *Root) Search(value interface{}) (interface{}, error) {
	index, status := getIndex(value, n)
	if status == FoundMatch {
		return n.values[index], nil
	} else if status == ValueNotComparable {
		return nil, errors.New("value not comparable with given comparator")
	}
	return nil, errors.New("did not find the value")
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
