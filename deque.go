package main

import (
	"fmt"
	"strings"
)

type DequeNode[T any] struct {
	value *T
	prev  *DequeNode[T]
	next  *DequeNode[T]
}

type Deque[T any] struct {
	head *DequeNode[T]
	tail *DequeNode[T]
	size int
}

// Insert an element at the end of the deque
func (deque *Deque[T]) push(element *T) {
	if deque.tail != nil {
		deque.tail.next = &DequeNode[T]{element, deque.tail, nil}
		deque.tail = deque.tail.next
	} else {
		deque.head = &DequeNode[T]{element, nil, nil}
		deque.tail = deque.head
	}
	deque.size++
}

// Insert an element at the start of the deque
func (deque *Deque[T]) unshift(element *T) {
	h := &DequeNode[T]{element, nil, deque.head}
	deque.head.prev = h
	deque.head = h
	if deque.tail == nil {
		deque.tail = deque.head
	}
	deque.size++
}

// Remove the front element of the deque
func (deque *Deque[T]) shift() *T {
	if deque.head == nil {
		return nil
	}
	element := deque.head.value
	deque.head = deque.head.next
	if deque.head == nil {
		deque.tail = nil
	} else {
		deque.head.prev = nil
	}
	deque.size--
	return element
}

// Remove the back element of the deque
func (deque *Deque[T]) pop() *T {
	if deque.tail == nil {
		return nil
	}
	element := deque.tail.value
	deque.tail = deque.tail.prev
	if deque.tail == nil {
		deque.head = nil
	} else {
		deque.tail.next = nil
	}
	deque.size--
	return element
}

func (deque *Deque[T]) front() *T {
	if deque.head != nil {
		return deque.head.value
	}
	return nil
}

func (deque *Deque[T]) back() *T {
	if deque.tail != nil {
		return deque.tail.value
	}
	return nil
}

func (deque *Deque[T]) length() int {
	return deque.size
}

func (deque *Deque[T]) empty() bool {
	return deque.size == 0
}

func (deque *Deque[T]) string() string {
	str := strings.Builder{}
	str.WriteString("[")
	for current := deque.head; current != nil; current = current.next {
		if current != deque.head {
			str.WriteString(", ")
		}
		str.WriteString(fmt.Sprintf("%v", *(current.value)))
	}
	str.WriteString("]")
	return str.String()
}
