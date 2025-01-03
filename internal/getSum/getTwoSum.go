package getSum

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// Time analysis: O(n)
// Worst case: Last element in nums is part of the pairs: So have to itterate through the whole array.
// Every element in the array triggers a constant lookup in the hashtable and maybe a constant insert

// Memory usage: O(n)
// The list of nums and the hashmap. The map can maximum be n-1 elements.

func twoSum(nums []int, target int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	mapNeededValueToIndexOfKnownpair := make(map[int]int)
	for i, v := range nums {
		// Not atomic
		if knownIndex, exists := mapNeededValueToIndexOfKnownpair[v]; exists {
			return []int{knownIndex, i}
		} else {
			mapNeededValueToIndexOfKnownpair[target-v] = i
		}
	}

	return []int{}
}

// Setup convinient helpers
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

// Setup debug logs
var debugLogEnabled = false

func debugLog(input string, args ...interface{}) {
	if debugLogEnabled {
		fmt.Printf(input, args...)
	}
}

// Setup test code
func main() {
	runTests()
}

func runTests() {
	testing.Main(
		nil,
		[]testing.InternalTest{
			{"TwoSums", TestTwoSums},
		},
		nil, nil,
	)
}

func TestTwoSums(t *testing.T) {
	tests := []struct {
		name           string
		inputs         []int
		targetSum      int
		expectedResult []int
	}{
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, targetSum: 9, expectedResult: []int{0, 1}},
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, targetSum: 26, expectedResult: []int{2, 3}},
		{name: "SimpleTest", inputs: []int{3, 2, 4}, targetSum: 6, expectedResult: []int{1, 2}},
		{name: "Duplicate entries", inputs: []int{3, 3}, targetSum: 6, expectedResult: []int{0, 1}},
		{name: "OutOfOrder", inputs: []int{15, 11, 7, 2}, targetSum: 9, expectedResult: []int{2, 3}},
		{name: "OutOfOrder", inputs: []int{15, 7, 2, 11}, targetSum: 9, expectedResult: []int{1, 2}},
		{name: "NoResult", inputs: []int{2, 7, 11, 15}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: -1, expectedResult: []int{}},
		{name: "Negative values", inputs: []int{-2, -7, -11, -15}, targetSum: -9, expectedResult: []int{0, 1}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("Nums: %v TargetSum: %v", test.inputs, test.targetSum)
		t.Run(testname, func(t *testing.T) {
			result := twoSum(test.inputs, test.targetSum)
			if !reflect.DeepEqual(result, test.expectedResult) {
				t.Errorf("The result were %v but we expected %v", result, test.expectedResult)
			}
		})
	}

	_ = []struct { // Convience: This will be the total sum of tests so I can copy paste it up above, instead of commenting out each line.
		name           string
		inputs         []int
		targetSum      int
		expectedResult []int
	}{
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, targetSum: 9, expectedResult: []int{0, 1}},
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15}, targetSum: 26, expectedResult: []int{2, 3}},
		{name: "SimpleTest", inputs: []int{3, 2, 4}, targetSum: 6, expectedResult: []int{1, 2}},
		{name: "Duplicate entries", inputs: []int{3, 3}, targetSum: 6, expectedResult: []int{0, 1}},
		{name: "OutOfOrder", inputs: []int{15, 11, 7, 2}, targetSum: 9, expectedResult: []int{2, 3}},
		{name: "OutOfOrder", inputs: []int{15, 7, 2, 11}, targetSum: 9, expectedResult: []int{1, 2}},
		{name: "NoResult", inputs: []int{2, 7, 11, 15}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: -1, expectedResult: []int{}},
		{name: "Negative values", inputs: []int{-2, -7, -11, -15}, targetSum: -9, expectedResult: []int{0, 1}},
	}
}
