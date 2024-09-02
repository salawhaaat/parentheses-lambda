package parentheses

import (
	"fmt"
	"testing"
)

// TestGenerate tests the Generate function to ensure it returns a string of the correct length.
func TestGenerate(t *testing.T) {
	tests := []struct {
		length int
	}{
		{length: 2},
		{length: 4},
		{length: 6},
		{length: 10},
	}

	for _, test := range tests {
		result := Generate(test.length)
		if len(result) != test.length {
			t.Errorf("Generate(%d) returned a string of length %d, expected %d", test.length, len(result), test.length)
		}
	}
}

// TestIsBalanced tests the isBalanced function with various cases to check if it correctly identifies balanced and unbalanced strings.
func TestIsBalanced(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "()", expected: true},
		{input: "([])", expected: true},
		{input: "([{}])", expected: true},
		{input: "([)", expected: false},
		{input: "[(])", expected: false},
		{input: "(((", expected: false},
		{input: "(([]))", expected: true},
	}

	for _, test := range tests {
		result := IsBalanced(test.input)
		if result != test.expected {
			t.Errorf("isBalanced(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// ExampleGenerate provides a testable example of the Generate function.
func ExampleGenerate() {
	result := Generate(6)
	fmt.Println(len(result))
	// Output:
	// 6
}

// ExampleIsBalanced provides a testable example of the isBalanced function.
func ExampleIsBalanced() {
	// balanced parentheses
	fmt.Println(IsBalanced("([]){}"))
	// unbalanced parentheses
	fmt.Println(IsBalanced("([)]"))
	// Output:
	// true
	// false
}
