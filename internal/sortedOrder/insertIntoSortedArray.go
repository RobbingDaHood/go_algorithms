package main

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// Time analysis:

// Memory usage:

func algorithmMain(nums []int, target int) ([]int, int) {
	if len(nums) == 0 {
		return []int{}, -1
	}

	var sortedSlice []int
	for _, v := range nums {
		sortedSlice = InsertIntoSortedIntSlice(sortedSlice, v)
	}

	searchResult := BinarySearchInSortedSlice(sortedSlice, target)

	return sortedSlice, searchResult
}

// Setup tests
func TestAlgorithmMain(t *testing.T) {
	tests := []struct {
		name                string
		inputs              []int
		searchTarget        int
		expectedSortedSlice []int
		expectedSearchIndex int
	}{
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, searchTarget: 2, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: 0},
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, searchTarget: 16, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: -1},
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, searchTarget: -2, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: -1},
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, searchTarget: 11, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: 2},
		{name: "SimpleTest", inputs: []int{15, 7, 2, 11}, searchTarget: 2, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: 0},
		{name: "SimpleTest", inputs: []int{15, 7, 2, 11}, searchTarget: 16, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: -1},
		{name: "SimpleTest", inputs: []int{15, 7, 2, 11}, searchTarget: -2, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: -1},
		{name: "SimpleTest", inputs: []int{15, 7, 2, 11}, searchTarget: 11, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: 2},
		{name: "SimpleTest", inputs: []int{-15, -7, 2, -11}, searchTarget: 2, expectedSortedSlice: []int{-15, -11, -7, 2}, expectedSearchIndex: 3},
		{name: "SimpleTest", inputs: []int{-15, -7, 2, -11}, searchTarget: -2, expectedSortedSlice: []int{-15, -11, -7, 2}, expectedSearchIndex: -1},
		{name: "SimpleTest", inputs: []int{-15, -7, 2, -11}, searchTarget: -7, expectedSortedSlice: []int{-15, -11, -7, 2}, expectedSearchIndex: 2},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("Nums: %v, SearchTarget: %v", test.inputs, test.searchTarget)
		t.Run(testname, func(t *testing.T) {
			sortedArrayResult, searchIndexResult := algorithmMain(test.inputs, test.searchTarget)
			if !reflect.DeepEqual(sortedArrayResult, test.expectedSortedSlice) {
				t.Errorf("The sortedArrayResult were %v but we expected %v", sortedArrayResult, test.expectedSortedSlice)
			}
			if searchIndexResult != test.expectedSearchIndex {
				t.Errorf("The searchIndexResult were %v but we expected %v", searchIndexResult, test.expectedSearchIndex)
			}
		})
	}
}

// Setup convenient helpers
var minIntCalculated = -1 << (32<<(^uint(0)>>63) - 1) // Respects if it is a 32 or 64 bit system
// ^uint(0) is pure 1's
// ^uint(0) >> 63 is 1 on a 64 bit system and 0 in a 32 bit system
// 32 << 1 - 1 is 64
// 32 << 0 - 1 is 32
// -1 << 62 is 64 bit minimum
// -1 << 32 is 32 bit minimum

func checkIfBetterMatch(candidate int, currentBest int, isNegative bool) bool {
	return (!isNegative && candidate > currentBest) || (isNegative && candidate < currentBest)
}

func InsertIntoSortedIntSlice(slice []int, value int) []int {
	index := sort.SearchInts(slice, value)
	slice = append(slice, 0)             // Make space for the new element
	copy(slice[index+1:], slice[index:]) // Shift elements to the right
	slice[index] = value
	return slice
}

func BinarySearchInSortedSlice(slice []int, target int) int {
	index := sort.SearchInts(slice, target)
	if 0 <= index && index < len(slice) && slice[index] == target {
		return index
	} else {
		return -1
	}
}

// Setup debug logs
var debugLogEnabled = false

func debugLog(input string, args ...interface{}) {
	if debugLogEnabled {
		fmt.Printf(input, args...)
	}
}

// Setup test infrastructure
func main() {
	runTests()
}

func runTests() {
	testing.Main(
		nil,
		[]testing.InternalTest{
			{"Test", TestAlgorithmMain},
		},
		nil, nil,
	)
}
