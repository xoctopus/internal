package type_alias

type MySlice[T any] = []T // illegal before 1.24

type MyStruct[T any] struct {
	slice MySlice[T]
}
