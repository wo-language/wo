
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
    - ir/no
    - noder/
    - reflect/type
    - reflectdata/reflect
    - rrtype/rrtype
    - typecheck
    - types
      - fmt
      - identity
      - kind_string
      - size
      - type
      - universe
    - walk
      - assign
      - range
  - go
    - ast
      - ast
      - walk
    - build/build
    - parser
      - interface
      - parser
      - resolver
    - scanner
      - errors
      - scanner
    - token/token
  - internal
    - abi/type
    - pkgbits/codes
  - reflect/value
  - runtime
    - alg 
    - set
    - set_fast32
    - set_fast64
    - set_faststr
    - sets
    - type
  - sets
    - example_test
    - iter
    - iter_test
    - sets
    - sets_test


data graph:

hset <- value
abi.SetType <- settype, reflect.setType <- value <- rrtype
abi.Type <- reflect.rtype
hset, reflect.type.Set <- reflectdata
TSET <- range
SetType <- alg

OMAKESET













