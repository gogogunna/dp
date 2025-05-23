package nullable

type Nullable[T any] struct {
	value T
	ok    bool
}

func NewValue[T any](value T) Nullable[T] {
	return Nullable[T]{value: value, ok: true}
}

func (n *Nullable[T]) Value() T {
	return n.value
}

func (n *Nullable[T]) IsDefined() bool {
	return n.ok
}
func (n *Nullable[T]) IsNil() bool {
	return !n.ok
}

func (n *Nullable[T]) SetValue(value T) {
	n.ok = true
	n.value = value
}

func NilDefaultValuee[T comparable](value T) Nullable[T] {
	if value == *new(T) {
		return Nullable[T]{}
	}

	return NewValue(value)
}
