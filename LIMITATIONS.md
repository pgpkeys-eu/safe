# Design Limitations

You cannot prevent people from using struct literals to generate arbitrary objects with inconsistent internal states in Go.
You can only provide big scary warnings in the contract and urge them to use constructors instead.

On the other hand, once you have a well-formed object, Go does a reasonable job of preventing unauthorised modifications.

## Option[T]

Unlike Rust, the underlying Some types are not considered to be instances of `Option` by the Go type system, and None is a special value rather than a type.
It is not possible therefore to construct a None value without a type parameter, and the parameterised constructor `None[T]()` is just syntactic sugar for `Option[T]{}`.
It is possible to infer the return type of `Some()` from its argument, but this does not generalise to `Some2()` and above.

Type matching of an `Option` requires two type matches in practice, so is less efficient than an interface.
The first type match is internal to `safe` and checks for nils, while the second is performed by the caller on a value that is guaranteed to be non-nil.

## Option\<n\>[T1, T2, ...]

Implementing `Option<n>` stretches the definition of `safe`ty somewhat, mainly due to limitations in Go's type system.

**Naming Things**

In Go, the generic types `Option[T1]`, `Option[T1, T2]`, `Option[T1, T2, T3]` are not distinct.
We disambiguate them by name: `Option2`, `Option3` etc., each of which must be declared separately.
We use type aliasing to share a single implementation between all `Option<n>` types, however each one must be explicitly mapped onto the generic implementation using glue methods.

Also, struct methods cannot have type parameters, so we cannot test for a specific type of Some using a method like `func (Option2[T1, T2]) IsSome[T0]() bool`.
We could in theory implement a top-level function `func IsSome[T0, T1, T2 any](Option2[T1, T2]) bool` but that's just *ugly*.

Type parameters cannot be nested, so for example we can't define `func Some2[T0 any, Option[T1, T2 any]](value T0) Option[T1, T2]`.
This makes it difficult to remember the order of type parameters, particularly when there are three or more.

**Deep Reflection**

If `T1` is (or contains) a pointer, `reflect.DeepEquals()` will attempt to dereference that location in memory as a pointer, regardless of the type currently stored.
This can cause test suites to panic with a pointer error even if no actual pointers are being compared.
If you need to use `DeepEquals`, there are several precautions that you can take:

* Don't mix pointer and non-pointer types in the same `Option<n>`
* If your largest non-pointer is the same size as a pointer (`uintptr`), use it as `T1`
* If your largest non-pointer is smaller than a pointer, use a dummy `uintptr` as `T1`

**Panic at the Disco**

Size and type mismatch panics would appear inevitable with the current feature set of the base Go language (as of 1.26).
The current design tries to expose them as soon and as consistently as possible, so that they can be caught in unit testing.

It may seem like `Some2` *should* be implementable using type constraints as follows:

```
func Some2[T1, T2 any](value interface{ ~T1 | ~T2 }) (o Option2[T1, T2]) {
    ...
}
```

But this fails to compile because interfaces can only be constrained to static types, not inferred ones.

Similarly, if you are familiar with other languages such as C it may seem possible to define `Option2` without assuming that `T1` is larger:

```
type Option2[T1, T2 any] struct {
	t any
	v [(unsafe.Sizeof(T1)<unsafe.Sizeof(T2))?unsafe.Sizeof(T1):unsafe.Sizeof(T2))]byte
}
```

This fails to compile because in Go, `Sizeof` is evaluated at runtime, unlike in C where it, and the conditional expression containing it, can be evaluated at compile time.

The following simpler declaration also fails for the same reason:

```
type Option2[T1, T2 any] struct {
	t any
	v [unsafe.Sizeof(T1)]byte
}
```

Even this would have been an improvement, because then `DeepEquals` would compare `v` as raw bytes and wouldn't raise a pointer error (See "Deep Reflection" above).
But we can't have nice things.

**Nice Things**

OK, maybe we can have some nice things.

Using type `any` for the type tag means we can point the tag interface at its own value field.
This may seem excessive - we know where the value field is without having to dereference the pointer - but it has some nice properties.

Firstly, by using pointer mangling, the tag interface implicitly remembers the dynamic type of the stored value.
We can type match on the tag to discover the dynamic type, instead of hand-rolling our own mapping of tags to types.

Second, pointing the tag interface at the value means there is always a valid pointer to the stored value for the lifetime of the `Option`.
Because interfaces are type-annotated, the garbage collector knows the dynamic type of the target,
including whether it has pointer members, and so won't double-free any second-order pointer targets.

And finally, it allows us to call `MarshalJSON` directly on the tag interface, which falls through to the value below, with the correct dynamic typing.
This avoids a LOT of boilerplate type matching code.

# Pretty Please With a Cherry on Top

There are a few small (and independently defensible) additions to the Go core language that would eliminate the panics in `Option<n>`:

* compile-time evaluation of `Sizeof` and simple expressions (as in C) would let us statically declare a value field of type `[]byte`, therefore:
    * no size check panics in `Some<n>`
    * no pointer panics in `DeepEquals`
* interface constraints using inferred types would let us statically limit the type of the input parameter to `Some<n>`
    * no type mismatch panics in `Some<n>`
