# Design Limitations

## Option[T]

Unlike Rust, the underlying Some types are not considered to be instances of `Option` by the Go type system, and None is a special value rather than a type.
It is not possible therefore to construct a None value without a type parameter, and the parameterised constructor `None[T]()` is just syntactic sugar for `Option[T]{}`.
It is possible to infer the return type of `Some()` from its argument, but this does not generalise to `Some2()` and above.

## Option\<n\>[T1, T2, ...]

Implementing `Option<n>` stretches the definition of `safe`ty somewhat, mainly due to limitations in Go's type system.

**Naming Things**

In Go, the generic types `Option[T1]`, `Option[T1, T2]`, `Option[T1, T2, T3]` are not distinct.
We disambiguate them by name: `Option2`, `Option3` etc., each of which must be implemented separately.

Only `Option2` is currently implemented.
`Option3`, `Option4` and higher can be easily implemented in future, because the design constraints are the same for all `n>1`.

Also, struct methods cannot have type parameters, so we cannot implement `func (Option2[T1, T2]) IsSome[T0]() bool`.
We could in theory implement a top-level function `func IsSome[T0, T1, T2 any](Option2[T1, T2]) bool` but that's just *ugly*.

**Deep Reflection** 

If T1 is (or contains) a pointer, `reflect.DeepEquals()` will attempt to dereference that location in memory as a pointer, regardless of the type currently stored.
This can cause test suites to panic with a pointer error even if no actual pointers are being compared.
If you need to use `DeepEquals`, there are several precautions that you can take:

* Don't mix pointer and non-pointer types in the same `Option<n>`
* If your largest non-pointer is the same size as a pointer (`uintptr`), use it as `T1`
* If your largest non-pointer is smaller than a pointer, use a dummy `uintptr` as `T1`

**Panic at the Disco**

It may seem like `Some2` *should* be implementable as follows:

```
func Some2[T1, T2 any](value interface{ ~T1 | ~T2 }) (o Option2[T1, T2]) {
    ...
}
```

But this fails to compile because `interface{ ... }` can only contain static types, not inferred ones.

Similarly, if you are familiar with other languages such as C it may seem possible to define `Option2` without assuming that `T1` is larger:

```
type Option2[T1, T2 any] struct {
	t any
	v [unsafe.Sizeof(T1)<unsafe.Sizeof(T2)?unsafe.Sizeof(T1):unsafe.Sizeof(T2))]byte
}
```

This fails to compile because in Go, `Sizeof` is evaluated at runtime, unlike in C where it, and the conditional expression containing it, can be evaluated at compile time.

The following declaration also fails for the same reason:

```
type Option2[T1, T2 any] struct {
	t any
	v [unsafe.Sizeof(T1)]byte
}
```

Even this would have been an improvement, because then DeepEquals would compare `v` as raw bytes and wouldn't raise a pointer error (See "Deep Reflection" above).
But we can't have nice things.
