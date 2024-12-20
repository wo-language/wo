This project is a fork of the "go" repo, which is the source for the Go compiler, library, and testing. This means that this repo can get compiled into a modified version of the Go compiler. If you were to follow the [full instructions for a source install](https://go.dev/doc/install/source#bootstrapFromCrosscompiledSource) with this repo, your computer would have this repo's compiler as your Go compiler in many targets of your choosing. That is how you can install Wo. Since it is a modification, it's not safe to rely on it as a Go builder. This is because, although it is meant to run .go files, it modifies those exact same files with extra logic (by checking if the current file is .wo), but that could fail, and then `go` could fail.

For linking, it will depend on your IDE, but it wasn't hard to add ".wo" as a file type association in goland

Development steps:

- running and testing this default repo [✅]

- run a modified compiler [✅]
    - run a .wo file (in a separate project) [✅]

- modifying the compiler code to support each kind of syntax [Doing...]
  - doing it for just one and test running it
  - make code formatter detect it

modifying the runner to support that transformation

possibly making the `wo` command separate, or at least instructing on how to make it an alias.

checking if it still reaches all the targets

setting up the website in Go then Wo

dealing with versions and downloadable executable installations for other users to test, perhaps offering an online playground

a transpiler that converts them between each other

---

running go on a .wo file seems to give "function main is undeclared in the main package"

todo:

1. refactor test/wo/*.wo -> test/wo_*.wo if they don't get ran, also need to change some bat file to include it in the tests maybe
2. add automated tests in my own run_wo.bat



### run without tests:
```
cd src
./make.bat
```
creation of the compiler runs:
go test -run=Generate -write=all
to create custom types

I think the compiler runs in this order:
build serialize scanner parser resolver walk

current commit syntax attempt to add:

`!ident()` - fails bc order of operations

`interface` is `tie` - fails because token defined in multiple places

recognizes `->` - fails bc doesn't belong anywhere

### `set`
- no syntax for it
- also impl's `setiter`, set_faststr, set_fast64, set_fast32, sets, sets/iter
https://dave.cheney.net/2018/05/29/how-the-go-runtime-implements-maps-efficiently-without-generics
```go
m := map[kType]vType // init
v := m[k]     // mapaccess1(m, k, &v)
v, ok := m[k] // mapaccess2(m, k, &v, &ok)
m[k] = 9001   // mapinsert(m, k, 9001)
delete(m, k)  // mapdelete(m, k)

s := set[eType] // init
---            // setaccess1 - disabled
ok := s[e]     // setaccess2(s, e, &ok)
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


after adding/removing any fundamental types, you have to run

go get -u golang.org/x/tools/cmd/stringer
go install stringer
cmd\compile\internal\types $ stringer -type Kind -trimprefix T type.go



















