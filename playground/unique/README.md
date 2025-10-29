# Testing the `unique` Package Introduced in Go 1.25

## Test Method

1. Create **1000 different strings**.
2. Store them **1000 times** in a data structure:
    - **`Unique` version:** use `unique.Make` when inserting data.
    - **`WithoutUnique` version:** insert directly without using `unique.Make`.
3. Observe and compare the **memory usage** of their Example.

---

## Conclusion

`unique.Make` implements **interning** for constant (comparable) types.  
Combined with `hashtrieMap`, fine-grained locking, and `runtime.GC`, it allows *
*interned values of comparable types** to be **garbage-collected** when no
longer referenced.

### Use Case

If your application stores a large number of constant-like values and
experiences high memory pressure,  
consider using `unique.Make` to reduce duplication and improve memory
efficiency.

---

## Read More

- [Interning strings](https://commaok.xyz/post/intern-strings/)
- [Wikipedia: String interning](https://en.wikipedia.org/wiki/String_interning)
- [Go issue #69370 – `unique` package discussion](https://github.com/golang/go/issues/69370)
- [go4org/intern](https://github.com/go4org/intern)