package main

type StackNode[T any] struct {
	value *T
	below *StackNode[T]
}

type Stack[T any] struct {
	top  *StackNode[T]
	size int
}

func (stack *Stack[T]) push(element *T) {
	stack.top = &StackNode[T]{element, stack.top}
	stack.size++
}

func (stack *Stack[T]) pop() *T {
	if stack.top == nil {
		return nil
	}
	element := stack.top.value
	stack.top = stack.top.below
	stack.size--
	return element
}

func (stack *Stack[T]) peek() *T {
	if stack.top == nil {
		return nil
	}
	return stack.top.value
}

func (stack *Stack[T]) length() int {
	return stack.size
}

func (stack *Stack[T]) empty() bool {
	return stack.size == 0
}
