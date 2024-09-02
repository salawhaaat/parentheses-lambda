package parentheses

import (
	"math/rand"
	"time"
)

// isBalanced checks if the provided string of parentheses is balanced.
// A string is considered balanced if every opening bracket has a corresponding
// closing bracket in the correct order.
func IsBalanced(s string) bool {
	// Map of closing brackets to their corresponding opening brackets
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}
	var stack []rune
	for _, c := range s {
		switch c {
		case '(', '{', '[':
			// Push opening brackets onto the stack
			stack = append(stack, c)
		case ')', '}', ']':
			// If stack is empty or top of stack doesn't match closing bracket, it's unbalanced
			if len(stack) == 0 || stack[len(stack)-1] != pairs[c] {
				return false
			}
			// Pop the matching opening bracket off the stack
			stack = stack[:len(stack)-1]
		}
	}
	// The string is balanced if the stack is empty at the end
	return len(stack) == 0
}

// Generate creates a random sequence of parentheses, brackets, and braces
// with the specified length. The generated sequence may or may not be balanced.
func Generate(length int) (parentheses string) {
	options := "[](){}"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	sequence := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(options))
		sequence[i] = options[randomIndex]
	}
	parentheses = string(sequence)
	return
}
