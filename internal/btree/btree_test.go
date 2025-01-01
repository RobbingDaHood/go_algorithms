package btree

import (
	"fmt"
	"github.com/google/btree"
	"math/rand"
	"testing"
)

type TestSearchData struct {
	name               string
	insert             []interface{}
	expectedLength     int
	comparator         func(first, second interface{}) ComparatorStatus
	allowedSearchError string
	allowedInsertError string
}

func TestSearch(t *testing.T) {
	randomData := generateRandomDataInterface(10_000)
	tests := []TestSearchData{
		{"ExistingValue", []interface{}{1}, 1, nil, "", ""},
		{"NonExistingValue", []interface{}{1}, 1, nil, "did not find the value", ""},
		{"DuplicatedInserts", []interface{}{1, 1}, 1, nil, "", "value already exists in tree"},
		{"MixedTypes", []interface{}{1, "2", true}, 3, nil, "", ""},
		{"DoNotInsertIfComparatorDoesNotMatch", []interface{}{1, "1", 2, true}, 2, ComparatorExpectingInts, "", "value not comparable with given comparator"},
		{"DoNotSearchIfComparatorDoesNotMatch", []interface{}{1, 2}, 2, ComparatorExpectingInts, "value not comparable with given comparator", ""},
		{"ElementsOutOfOrder", []interface{}{3, 2, 1}, 3, nil, "", ""},
		{"InsertMoreThanDefaultMaxSizeAndThenSearchInt", get1To1000(), 1000, ComparatorExpectingInts, "", ""},
		{"InsertMoreThanDefaultMaxSizeAndThenSearchString", get1To1000(), 1000, nil, "", ""},
		{"InsertRandomDataInt", randomData, 6288, ComparatorExpectingInts, "", "value already exists in tree"},
		{"InsertRandomDataString", randomData, 6288, nil, "", "value already exists in tree"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := setup(t, tt)
			testInsertingEveryElement(t, tt, &tree)
			checkThatWeCanSearchEveryValidInsertedValue(t, tt, tree)
		})
	}
}

func setup(t *testing.T, tt TestSearchData) Root[interface{}] {
	tree := CreateTreeDefaultValues[interface{}]()
	tree.nodeMaxSize = 100
	if len(tree.values) != 0 {
		t.Errorf("CreateTreeDefaultValues() returned tree with values")
	}
	if tt.comparator != nil {
		tree.comparator = tt.comparator
	}
	return tree
}

func testInsertingEveryElement(t *testing.T, tt TestSearchData, tree *Root[interface{}]) {
	for _, v := range tt.insert {
		err := tree.Insert(v)
		if err != nil && tt.allowedInsertError != err.Error() {
			t.Errorf("Insert() returned error '%v' but expected: '%v'", err, tt.allowedInsertError)
		}

	}
	if tree.nodeCount != tt.expectedLength {
		t.Errorf("Expected tree length %d, got %d", tt.expectedLength, tree.nodeCount)
	}
	tree.checkForDuplicates()
}

func (n *Root[T]) checkForDuplicates() {
	allValues := n.getAllMaxValue()

	for i, maxValue := range allValues {
		count := 0
		for k, existingMax := range allValues {
			if n.comparator(maxValue, existingMax) == Equal {
				count++
				if count > 1 {
					panic(fmt.Sprintf("duplicate values in node: %v, %v", i, k))
				}
			}
		}
	}
}

func (n *innerNode[T]) getAllMaxValue() []T {
	var allValues []T
	for _, node := range n.values {
		if n.isLeaf {
			allValues = append(allValues, node.maxValue)
		} else {
			allValues = append(allValues, node.node.getAllMaxValue()...)
		}
	}
	return allValues
}

func checkThatWeCanSearchEveryValidInsertedValue(t *testing.T, tt TestSearchData, tree Root[interface{}]) {
	for _, v := range tt.insert {
		if tree.comparator(v, v) == ArgumentsNotComparable {
			continue
		}
		result, searchError := tree.Search(v)
		if searchError != nil {
			if searchError.Error() != tt.allowedSearchError {
				t.Fatalf("Search() = %v, %v; want %v, %v", result, searchError, v, tt.allowedSearchError)
			}
			if result != nil {
				t.Errorf("When there are search errors then result should be nil but it were: %v", v)
			}
		}
		if searchError == nil && result != v {
			t.Errorf("When no search errors then result should be %v but it were: %v", v, result)
		}
	}
}

func BenchmarkInsert(b *testing.B) {
	data := generateRandomDataInt(b.N)

	b.ResetTimer()
	b.ReportAllocs()
	tree := CreateTreeDefaultValues[interface{}]()
	tree.comparator = ComparatorExpectingInts
	for i := range data {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}
}

func BenchmarkInsertInteger(b *testing.B) {
	data := generateRandomDataInt(b.N)

	b.ResetTimer()
	b.ReportAllocs()
	tree := CreateTreeDefaultValues[int]()
	tree.comparator = isLesserInteger
	for i := range data {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}
}

func BenchmarkInsertGoogleGeneric(b *testing.B) {
	data := generateRandomDataInt(b.N)

	b.ResetTimer()
	b.ReportAllocs()
	tree := btree.NewG[int](10_000, func(a, b int) bool {
		return a < b
	})
	for i := range data {
		tree.ReplaceOrInsert(i)
	}
}

func BenchmarkInsertGoogle(b *testing.B) {
	data := generateRandomDataInt(b.N)

	b.ResetTimer()
	b.ReportAllocs()
	tree := btree.New(10000)
	for i := range data {
		tree.ReplaceOrInsert(btree.Int(i))
	}
}

func BenchmarkSearch(b *testing.B) {
	data := generateRandomDataInt(b.N)

	tree := CreateTreeDefaultValues[interface{}]()
	for i := range data {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := range data {
		_, err := tree.Search(i)
		if err != nil {
			b.Fatalf("Search() returned error: %v", err)
		}
	}
}
func isLesserInteger(first, second int) ComparatorStatus {
	if first > second {
		return FirstArgumentBigger
	} else if first < second {
		return SecondArgumentBigger
	}
	return Equal
}

func BenchmarkSearchInteger(b *testing.B) {
	data := generateRandomDataInt(b.N)
	tree := CreateTreeDefaultValues[int]()
	tree.comparator = isLesserInteger
	for i := range data {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := range data {
		_, err := tree.Search(i)
		if err != nil {
			b.Fatalf("Search() returned error: %v", err)
		}
	}
}

func BenchmarkSearchGoogle(b *testing.B) {
	data := generateRandomDataInt(b.N)
	tree := btree.New(10000)
	for i := range data {
		tree.ReplaceOrInsert(btree.Int(i))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := range data {
		item := tree.Get(btree.Int(i))
		if item == nil {
			b.Fatalf("Search() returned nil")
		}
	}
}

func BenchmarkSearchGoogleGeneric(b *testing.B) {
	data := generateRandomDataInt(b.N)
	tree := btree.NewG[int](10_000, func(a, b int) bool {
		return a < b
	})
	for i := range data {
		tree.ReplaceOrInsert(i)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := range data {
		_, found := tree.Get(i)
		if !found {
			b.Fatalf("Search() returned nil")
		}
	}
}

func BenchmarkDifferentNodeSizes(b *testing.B) {
	sizes := []int{100, 1_000, 10_000, 100_000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			data := generateRandomDataInt(b.N)
			b.ResetTimer()
			b.ReportAllocs()
			tree := CreateTreeDefaultValues[int]()
			tree.comparator = isLesserInteger
			tree.nodeMaxSize = size
			for i := range data {
				err := tree.Insert(i)
				if err != nil {
					b.Fatalf("Insert() returned error: %v", err)
				}
			}
			for i := range data {
				_, err := tree.Search(i)
				if err != nil {
					b.Fatalf("Search() returned error: %v", err)
				}
			}
		})
	}
}

func BenchmarkDifferentInputSizes(b *testing.B) {
	sizes := []int{1_000, 10_000, 100_000, 1_000_000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			data := generateRandomDataInt(size)
			b.ResetTimer()
			b.ReportAllocs()
			tree := CreateTreeDefaultValues[int]()
			tree.comparator = isLesserInteger
			for i := range data {
				err := tree.Insert(i)
				if err != nil {
					b.Fatalf("Insert() returned error: %v", err)
				}
			}
			for i := range data {
				_, err := tree.Search(i)
				if err != nil {
					b.Fatalf("Search() returned error: %v", err)
				}
			}
		})
	}
}

func BenchmarkIntCompare(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ComparatorExpectingInts(1, 2)
	}
}

func ComparatorExpectingInts(first, second interface{}) ComparatorStatus {
	switch first.(type) {
	case int:
		// Do nothing
	default:
		return ArgumentsNotComparable
	}
	switch second.(type) {
	case int:
		// Do nothing
	default:
		return ArgumentsNotComparable
	}

	if first.(int) > second.(int) {
		return FirstArgumentBigger
	} else if first.(int) < second.(int) {
		return SecondArgumentBigger
	}
	return Equal
}

func get1To1000() []interface{} {
	var numbers []interface{}
	for i := 1; i <= 1000; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}
func generateRandomDataInt(n int) []int {
	rand.Seed(42) // Fixed seed for reproducibility
	size := n     // Size of the slice
	data := make([]int, size)
	for i := range data {
		data[i] = rand.Intn(n) // Random integers between 0 and 999
	}
	return data
}

func generateRandomDataInterface(n int) []interface{} {
	rand.Seed(42) // Fixed seed for reproducibility
	size := n     // Size of the slice
	data := make([]interface{}, size)
	for i := range data {
		data[i] = rand.Intn(n) // Random integers between 0 and 999
	}
	return data
}
