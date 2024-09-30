package ignorant_parser

type stack[T any] struct {
	tracker []T
}

func (s *stack[T]) isEmpty() bool {
	return len((*s).tracker) == 0
}

// Push a new value onto the stack
func (s *stack[T]) push(item T) {
	(*s).tracker = append((*s).tracker, item) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *stack[T]) pop() (*T, bool) {
	if s.isEmpty() {
		return nil, false
	} else {
		index := len((*s).tracker) - 1      // Get the index of the top most element.
		element := ((*s).tracker)[index]    // Index into the slice and obtain the element.
		(*s).tracker = (*s).tracker[:index] // Remove it from the stack by slicing it off.
		return &element, true
	}
}

func (s *stack[T]) clear() {
	if s.isEmpty() {
		return
	}
	(*s).tracker = nil
}

func (s *stack[T]) height() int {
	return len((*s).tracker)
}

func (s *stack[T]) peek() (*T, bool) {
	if s.isEmpty() {
		return nil, false
	} else {
		index := len((*s).tracker) - 1   // Get the index of the top most element.
		element := ((*s).tracker)[index] // Index into the slice and obtain the element.
		return &element, true
	}
}
