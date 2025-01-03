package getSum

import (
	"fmt"
	"reflect"
	"testing"
)

// Correctnes: There are 2^N combinations of nums
// If a subset of nums is invalid then we do not need to consider that again

// Time analysis
// Worst case: Last element in nums is the result and all numbers before are smaller, so cannot be ignored
//  Before that we checked every combination !n
//  Then it is O(n^2)

// Memory usage: len(nums) + len(result) <= O(n)

func twoSum(nums []int, target int) []int {
	if target == 0 || len(nums) == 0 {
		return []int{}
	}

	result := twoSumReqClosestMatch(nums, target, 0, []int{})
	if getSum(nums, result) == target {
		return result
	} else {
		return []int{}
	}
}

func getSum(nums []int, tmpResult []int) int {
	result := 0

	if len(tmpResult) == 0 {
		return minIntCalculated
	}

	for _, v := range tmpResult {
		result += nums[v]
	}
	// debugLog("Returning sum %v for nums %v and indexes %v \n", result, nums, tmpResult)
	return result
}

func twoSumReqClosestMatch(nums []int, target int, startIndex int, currentSum []int) []int {
	debugLog("Checking this slice now %v \n", nums[startIndex:])
	debugLog("Checking this startIndex now %v \n", startIndex)
	debugLog("Checking this currentSum now %v \n", currentSum)
	currentSumSum := getSum(nums, currentSum)
	currentBestTmpResult := currentSum
	currentBestTmpSumResult := currentSumSum
	targetIsNegative := target < 0
	for i := startIndex; i < len(nums); i++ {
		v := nums[i]
		debugLog("Checking this v now %v \n", v)
		if checkIfBetterMatch(v+currentSumSum, target, targetIsNegative) {
			debugLog("currentSum %v + v %v were over target %v \n", currentSumSum, v, target)
			continue
		}

		if len(nums) > 1 {
			tmpResult := twoSumReqClosestMatch(nums, target, i+1, append(currentSum, i))
			tmpResultSum := getSum(nums, tmpResult)
			if checkIfBetterMatch(tmpResultSum, target, targetIsNegative) {
				debugLog("tmpResult %v were over target %v \n", getSum(nums, tmpResult), target)
				continue
			} else if tmpResultSum == target {
				debugLog("tmpResultSum equals target %v \n", tmpResult)
				return tmpResult
			} else if checkIfBetterMatch(tmpResultSum, currentBestTmpSumResult, targetIsNegative) {
				debugLog("Found new tmpResult %v \n", tmpResult)
				currentBestTmpResult = tmpResult
				currentBestTmpSumResult = tmpResultSum
			}
		}
	}
	debugLog("There were nothing here I could use, returning currentSum again %v \n", currentSum)
	return currentBestTmpResult
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
			{"GetSum", TestGetSum},
			{"TwoSums", TestTwoSums},
		},
		nil, nil,
	)
}

func TestGetSum(t *testing.T) {
	tests := []struct {
		name      string
		nums      []int
		indexs    []int
		targetSum int
	}{
		{name: "SimpleTest", nums: []int{2, 7, 11, 15}, indexs: []int{0, 1}, targetSum: 9},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getSum(test.nums, test.indexs)
			if result != test.targetSum {
				t.Errorf("The result were %v but we expected %v", result, test.targetSum)
			}
		})
	}
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
		{name: "SimpleTestAllResult", inputs: []int{2, 7, 11, 15}, targetSum: 35, expectedResult: []int{0, 1, 2, 3}},
		{name: "SimpleTestOneResult", inputs: []int{2, 7, 11, 15}, targetSum: 11, expectedResult: []int{2}},
		{name: "Duplicate entries", inputs: []int{3, 3}, targetSum: 6, expectedResult: []int{0, 1}},
		{name: "OutOfOrder", inputs: []int{15, 11, 7, 2}, targetSum: 9, expectedResult: []int{2, 3}},
		{name: "OutOfOrder", inputs: []int{15, 7, 2, 11}, targetSum: 9, expectedResult: []int{1, 2}},
		{name: "NoResult", inputs: []int{2, 7, 11, 15}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: -1, expectedResult: []int{}},
		{name: "Negative values", inputs: []int{-2, -7, -11, -15}, targetSum: -9, expectedResult: []int{0, 1}},
		{name: "SimpleTest", inputs: []int{2, 7, 11, 15, -1}, targetSum: 8, expectedResult: []int{0, 1, 4}},
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
		{name: "SimpleTestAllResult", inputs: []int{2, 7, 11, 15}, targetSum: 35, expectedResult: []int{0, 1, 2, 3}},
		{name: "SimpleTestOneResult", inputs: []int{2, 7, 11, 15}, targetSum: 11, expectedResult: []int{2}},
		{name: "Duplicate entries", inputs: []int{3, 3}, targetSum: 6, expectedResult: []int{0, 1}},
		{name: "OutOfOrder", inputs: []int{15, 11, 7, 2}, targetSum: 9, expectedResult: []int{2, 3}},
		{name: "OutOfOrder", inputs: []int{15, 7, 2, 11}, targetSum: 9, expectedResult: []int{1, 2}},
		{name: "NoResult", inputs: []int{2, 7, 11, 15}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: 23, expectedResult: []int{}},
		{name: "NoResult", inputs: []int{}, targetSum: -1, expectedResult: []int{}},
		{name: "Negative values", inputs: []int{-2, -7, -11, -15}, targetSum: -9, expectedResult: []int{0, 1}},
	}
}
