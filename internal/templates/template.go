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

	// IMPLEMENT SOLUTION HERE

	return nums, target
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
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, searchTarget: 2, expectedSortedSlice: []int{2, 7, 11, 15}, expectedSearchIndex: 2},
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
var maxIntCalculated = (1 << (32 << (^uint(0) >> 63))) - 1

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

func InsertIntoSortedStringSlice(slice []string, value string) []string {
	index := sort.SearchStrings(slice, value)
	slice = append(slice, "")            // Make space for the new element
	copy(slice[index+1:], slice[index:]) // Shift elements to the right
	slice[index] = value
	return slice
}

func BinarySearchInStringSlice(slice []string, target string) int {
	index := sort.SearchStrings(slice, target)
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
