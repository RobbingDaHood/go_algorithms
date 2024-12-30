package btree

import (
	"testing"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name               string
		insert             []interface{}
		expectedLength     int
		comparator         func(first, second interface{}) ComparatorStatus
		search             interface{}
		allowedSearchError string
		allowedInsertError string
	}{
		{"ExistingValue", []interface{}{1}, 1, nil, 1, "", ""},
		{"NonExistingValue", []interface{}{1}, 1, nil, 2, "did not find the value", ""},
		{"DuplicatedInserts", []interface{}{1, 1}, 1, nil, 1, "", "value already exists in tree"},
		{"MixedTypes", []interface{}{1, "2", true}, 3, nil, 1, "", ""},
		{"DoNotInsertIfComparatorDoesNotMatch", []interface{}{1, "1", 2, true}, 2, ComparatorExpectingInts, 2, "", "value not comparable with given comparator"},
		{"DoNotSearchIfComparatorDoesNotMatch", []interface{}{1, 2}, 2, ComparatorExpectingInts, "2", "value not comparable with given comparator", ""},
		{"ElementsOutOfOrder", []interface{}{3, 2, 1}, 3, nil, 1, "", ""},
		{"InsertMoreThanDefaultMaxSizeAndThenSearch", get1To1000(), 1000, ComparatorExpectingInts, 333, "", ""},
	}

	// TODO closest match
	// TODO visitor pattern
	// TODO sum

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := CreateTreeDefaultValues()
			tree.nodeMaxSize = 100
			if len(tree.values) != 0 {
				t.Errorf("CreateTreeDefaultValues() returned tree with values")
			}
			if tt.comparator != nil {
				tree.comparator = tt.comparator
			}
			for _, v := range tt.insert {
				err := tree.Insert(v)
				if err != nil && tt.allowedInsertError != err.Error() {
					t.Errorf("Insert() returned error '%v' but expected: '%v'", err, tt.allowedInsertError)
				}
			}
			if tree.nodeCount != tt.expectedLength {
				t.Errorf("Expected tree length %d, got %d", tt.expectedLength, tree.nodeCount)
			}
			result, searchError := tree.Search(tt.search)
			if searchError != nil {
				if searchError.Error() != tt.allowedSearchError {
					t.Fatalf("Search() = %v, %v; want %v, %v", result, searchError, tt.search, tt.allowedSearchError)
				}
				if result != nil {
					t.Errorf("When there are search errors then result should be nil but it were: %v", tt.search)
				}
			}
			if searchError == nil && result != tt.search {
				t.Errorf("When no search errors then result should be %v but it were: %v", tt.search, result)
			}
		})
	}
}

func BenchmarkInsert10_000(b *testing.B) {
	b.N = 10_000
	tree := CreateTreeDefaultValues()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}
}

func BenchmarkInsert100_000(b *testing.B) {
	b.N = 100_000
	tree := CreateTreeDefaultValues()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}
}

func BenchmarkSearch10_000(b *testing.B) {
	b.N = 10_000
	tree := CreateTreeDefaultValues()
	for i := 0; i < b.N; i++ {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := tree.Search(i)
		if err != nil {
			b.Fatalf("Search() returned error: %v", err)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	tree := CreateTreeDefaultValues()
	for i := 0; i < b.N; i++ {
		err := tree.Insert(i)
		if err != nil {
			b.Fatalf("Insert() returned error: %v", err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := tree.Search(i)
		if err != nil {
			b.Fatalf("Search() returned error: %v", err)
		}
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
