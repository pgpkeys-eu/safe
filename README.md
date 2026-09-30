# Safe

`safe` is a Go package that provides tagged-union types for handling nils and errors more safely.
It is inspired by Rust's `Option` and `Result` types, and informed by previous efforts to implement them in Go that exposed sharp edges in unexpected places.

## Overview

Due to the architecture of Go, there is no way to seamlessly implement Rust-style tagged unions.
`safe` attempts to expose the inevitable seams in places that minimise the risk of accidental unsafe practice.

The design brief of `safe` is:

* to eliminate nil dereferencing errors
* to reduce the occurrence of unchecked error return values

The chosen design forces consumers to use either type matching or helper functions, to prevent unsafe access to nils:

* tagged unions are implemented as generic struct types, not interfaces
* the zero literals `Option[T]{}` and `Result[T]{}` represent safe and meaningful default values
* all members of tagged unions are private
* tagged unions must be explicitly converted to a match type before type matching
* the match types `MatchNone`, `MatchSome`, `MatchOK` and `MatchErr` are used ephemerally for type matching, and must not be directly constructed
* `MatchNone` is inert, while the other match types expose the public members `Some`, `OK` or `Err` as appropriate

This design choice has several advantages:

* uninitialised variables are safe by default
* violations of the usage guidelines are highly visible, and can be checked statically
* temptation to stray from the path is encountered when values are generated (rare), but not when they are consumed (common)
* struct literals have an unintuitive format and scary inline documentation, discouraging their use

Return values should therefore be difficult or impossible to use unsafely, so long as `Option` and `Result` are used consistently in function signatures.
For convenience, helper functions are provided to enable use of the embedded values without type matching (e.g. `UnwrapOr()`).
There is no `Unwrap()` function by design, therefore there are no panics at consumption time.

This design differs from earlier efforts, which implement tagged unions as interfaces.
While this is conceptually intuitive, interfaces are nillable types and so cannot have safe default values.

## Usage

`safe` can be imported and used like a normal package:

```
package foo
import "github.com/pgpkeys-eu/safe"

func Foo() safe.Result[safe.Option[string]] {
    return safe.OK(safe.Some("foo"))
}
```

It can also be dot-imported.
While dot-importing runs the risk of namespace collisions, it reduces boilerplate considerably:

```
package foo
import . "github.com/pgpkeys-eu/safe"

func Foo() Result[Option[string]] {
    return OK(Some("foo"))
}
```

# Types

## Option[T]

Creating an `Option[T]`:

```
o1 := Option[T]{}
o2 := None[T]()
o3 := Some[T](value)
```

The zero literal `Option[T]{}` represents None.
If `value` is nil, `Some(value)` will (perhaps counterintuitively) return a None `Option`, otherwise it returns a Some.
This ensures that a nil value can never be obtained from a Some `Option`.

Consuming an `Option` is done by calling `Match()` and type-matching the return value:

```
switch x := o1.Match().(type){
case MatchSome[T]:
    fmt.printf("Some: %v", x.Some)
default:
    fmt.printf("None")
}
```

Note that in practice there is no need to explicitly handle `case MatchNone`.

`MatchSome[T]` contains a single public member `Some` which contains the original non-nil value.
`MatchNone` contains no members and represents a nil value.

Other helper functions defined for `Option[T]`:

* IsSome() bool
* IsNone() bool
* UnwrapOr(*T) *T
* UnwrapOrElse(func() *T) *T

These mirror their Rust equivalents.

### JSON

`Option` is compitible with `encoding/json`.
A Some value will be serialised identically to a `T` value, while a None value will be represented by `{}`.

Note that a Some `Option` containing a struct with no serialisable members will be serialised identically to None,
and so will deserialise to None rather than the original Some.
This is a limitation of JSON, and a similar caveat applies to nil pointers.

`Option` does not (currently) support `omitzero` or `omitempty`.

## Result[T]

Creating a `Result[T]`:

```
r1 := Result[T]{}
r2 := OK[T](t_value)
r3 := Err[T](err_value)
```

The zero literal `Result[T]{}` represents an OK containing the zero value of `T`.
If `err_value` is nil, `Err()` returns an Err `Result` containing an empty error message.
This ensures that a nil error can never be obtained from an Err `Result`.

Consuming a `Result` is done by calling `Match()` and type-matching the return value:

```
switch x := r1.Match().(type){
case MatchOK[T]:
    fmt.printf("OK: %v", x.OK)
case MatchErr[T]:
    fmt.printf("Err: %v", x.Err)
default:
    fmt.printf("Should not get here")
}
```

`MatchOK[T]` contains a single public member `OK` which contains the original value.
`MatchErr[T]` contains a single public member `Err` which contains the original error.

Other helper functions defined for Result:

* IsOK() bool
* IsErr() bool
* UnwrapOr(*T) *T
* UnwrapOrElse(func() *T) *T

These mirror their Rust equivalents.

### JSON

`Result` is not serialisable.
This is intentional.

### Result[Option[T]]

The common Rust return value type `Result[Option[T]]` can be easily generated:

```
return OK(None[T]())
return Result[Option[T]]{}      // ==OK(None[T]())
return OK(Some[T](nil))         // ==OK(None[T]())
return OK(Some(foo))
return Err(errors.New("help"))
```

## Option\<n\>[T1, T2, ...]

`Option<n>[T1, T2, ...]` is the generalisation of `Option[T]` to `n` Some types.
It is the equivalent of a generic (i.e. Rust-style) tagged union, with the exception that it always permits None (the zero value).

It comes with some additional caveats, due to [limitations of Go's type system](LIMITATIONS.md).
Unlike `Option[T]` and `Result[T]` it MAY panic on construction if the two type rules (given below) are not followed.
A basic unit test suite should trigger such a panic immediately, and the error message will describe the solution.

`Option<n>` is created similarly to `Option`, but with additional type parameters.
Using `Option2` as a concrete example:

```
o1 := Option2[T1, T2]{}
o2 := None2[T1, T2]()
o3 := Some2[T0, T1, T2](value) // where value is of type T0
```

Beware that unlike `Option<n>` and `None<n>`, `Some<n>` takes `n+1` type parameters.
This is because Go's type system cannot restrict at compile time that only a value of type `T1`, `T2` etc. can be supplied as input,
so we must pass the input type as an additional type parameter `T0`.
`Some<n>` will check at runtime whether `T0` is one of `T1`, `T2` etc., and if not it will panic with a friendly error message.

The above is the first type rule of `Option<n>`.

Go also cannot tell at compile time what size `T1` and `T2` are - this can only be done at runtime using reflection.
This means that it cannot tell which of `T1` or `T2` is bigger, and therefore how much memory to allocate to the Option struct.
Instead, it allocates enough memory to hold a `T1`, and assumes that `T2` is the same size or smaller.
`Some<n>` will check at runtime whether `T1` is the larger type, and if not it will panic with a friendly error message.

The above is the second type rule of `Option<n>`.

Note that these runtime panics depend only on the type parameters, not on any values.
They should therefore fire reliably in a test suite, so long as the code paths are covered.
Unit tests are your friend.

Note also that `reflect.DeepEquals` works fine on the happy path, but [can panic on the unhappy one](LIMITATIONS.md)

Consuming an `Option<n>` is done the same way as for `Option`:

```
switch x := o1.Match().(type){
case MatchSome[T1]:
    fmt.printf("SomeT1: %v", x.Some)
case MatchSome[T2]:
    fmt.printf("SomeT2: %v", x.Some)
default:
    fmt.printf("None")
}
```

None is represented by `MatchNone`, but as with `Option` there is no need to explicitly match it.

Other helper functions defined for `Option<n>`:

* IsSome() bool
* IsNone() bool

There is no `UnwrapOr` or `UnwrapOrElse` because there is more than one Some return type,
so type matching is unavoidable.
There is also no `IsSome[T]() bool`, due to [limitations in Go](LIMITATIONS.md).

Only `Option2` is currently implemented.

### JSON

`Option<n>` is compitible with `encoding/json`, similarly to `Option`.
When deserialising, each type is tried in the order None, `T1`, `T2` etc. and the first type to successfully deserialise is returned.

# References

Several packages exist for `Option` and/or `Result`, but do not provide nil dereference safety:

* [go-rust-std](https://avivatedgi.github.io/go-rust-std/) implements `Option`, `Result` and `Collection` interfaces
* [safetypes](https://github.com/eminarican/safetypes) implements `Option` and `Result` interfaces
* [go-option](https://github.com/sdwillbrand/go-option) implements an `Option` interface

There are also existing tools for implementing real union types, but which require a preprocessor stage:

* [MkUnion](https://widmogrod.github.io/mkunion/) implements tagged unions, but also reimplements type matching
* [unionize](https://github.com/zyedidia/unionize) implements C-style untagged unions
