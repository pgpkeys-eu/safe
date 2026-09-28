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
There is no `Unwrap()` function by design, therefore there are no panics.

This design differs from earlier efforts, which implement tagged unions as interfaces.
While this is conceptually intuitive, interfaces are nillable types and so cannot have safe default values.

## Usage

`safe` can be imported and used like a normal package:

```
package foo
import "safe"

func Foo() safe.Result[safe.Option[string]] {
    return safe.Result(safe.Some("foo"))
}
```

It can also be dot-imported.
While dot-importing runs the risk of namespace collisions, it reduces boilerplate considerably:

```
package foo
import . "safe"

func Foo() Result[Option[string]] {
    return Result(Some("foo"))
}
```

# Types

## Option

Creating an `Option[T]`:

```
o1 := Option[T]{}
o2 := None[T]()
o3 := Some[T](value)
```

The zero literal `Option[T]{}` represents None.
If `value` is nil, `o3` represents None, otherwise it represents Some.
This ensures that a nil value can never be obtained from a Some `Option`.

Consuming an `Option` is done by calling `Match()` and type-matching the return value:

```
switch x := o1.Match().(type){
case MatchSome[T]:
    fmt.printf("OK: %v", x.Some)
case MatchNone[T]:
    fmt.printf("None")
default:
    fmt.printf("Should not get here")
}
```

`MatchSome[T]` contains a single public member `Some` which contains the original non-nil value.
`MatchNone[T]` contains no public members and represents a nil value.

Other helper functions defined for Option:

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

`Option` does not support `omitzero` or `omitempty`.

## Result

Creating a `Result`:

```
r1 := Result[T]{}
r2 := OK[T](t_value)
r3 := Err[T](err_value)
```

The zero literal `Result[T]{}` represents an OK containing the zero value of `T`.
If `err_value` is nil, `r3` represents an Err containing an empty error message.
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

### Result[Option[T]]

The common Rust return value type `Result[Option[T]]` can be easily generated:

```
return OK(None[T]())
return Result[Option[T]]{}      // ==OK(None[T]())
return OK(Some[T](nil))         // ==OK(None[T]())
return OK(Some(foo))
return Err(errors.New("help"))
```

## References

Several similar packages exist, but do not provide nil dereference safety:

* https://avivatedgi.github.io/go-rust-std/
* https://github.com/eminarican/safetypes
* https://github.com/sdwillbrand/go-option
