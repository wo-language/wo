
### `set`
- no syntax for it
- also impl's `setiter`, set_faststr, set_fast64, set_fast32, sets, sets/iter, make()
- https://dave.cheney.net/2018/05/29/how-the-go-runtime-implements-maps-efficiently-without-generics
```go
m := map[kType]vType // init
v := m[k]     // mapaccess1(m, k, &v)
v, ok := m[k] // mapaccess2(m, k, &v, &ok)
m[k] = 9001   // mapinsert(m, k, 9001)
delete(m, k)  // mapdelete(m, k)

s := set[eType] // init
---            // setaccess1 - disabled
ok := s[e]     // setaccess2(s, e, &ok) - similar signature of mapaccess1
s.insert(9001) // setinsert(s, e, 9001)
delete(s, e)   // setdelete(s, e)

```
meaning:
remove mapaccess1,
modify signature of mapaccess2,
change the parser's syntax

- it should really barely be much faster than map[type]struct{}
    - since I only removed the element field and any calculations for it (which were constant time ones)
    - the overall time complexity should be the same

relevant files structure:
- src
  - cmd/compile/internal
    - ir/node.go
    - reflectdata/reflect.go
    - typecheck.go
  - go/ast
    - ast.go
  - internal
    - abi
  - runtime
    - set.go
    - set_fast32.go
    - set_fast64.go
    - set_faststr.go
    - sets.go

relevant files relationships:

- ir/node.go
- reflect.go
- typecheck.go
- abi
- set.go
- set_fast32.go
- set_fast64.go
- set_faststr.go
- sets.go
