This project is a fork of the "go" repo, which is the source for the Go compiler, library, and testing.
This means that this repo can get compiled into a modified version of the Go compiler.
If you were to follow the [full instructions for a source install](https://go.dev/doc/install/source#bootstrapFromCrosscompiledSource) with this repo, your computer would have this
repo's compiler as your Go compiler in many targets of your choosing.
That is how you can install Wo. Since it is a modification, it's not safe to rely on it as a Go builder.
This is because, although it is meant to run .go files, it modifies those exact same files with extra logic
(by checking if the current file is .wo), but that could fail, and then `go` could fail.

This file is meant to be changed per commit, associated with details that apply to that specific commit e.g. current outputs.

For syntax linking in a code editor, it will depend on your IDE, but it wasn't hard to add ".wo" as a file type association in goland

Development steps:

- running and testing this default repo [✅]
- run a modified compiler [✅]
    - run a .wo file (in a separate project) [✅]
- modifying the compiler code to support each kind of syntax [Doing...]
  - doing it for just one and test running it [Doing with set...]
  - make code formatter detect it
  - try to make some modular feature
- modifying the runner to support that transformation [Next] [Doing with set...]
- add appropriate tests
- possibly making the `wo` command separate, or at least instructing on how to make it an alias.
- checking if it still reaches all the targets
- setting up the website in Go then Wo
- dealing with versions and downloadable executable installations for other users to test
- perhaps offering an online playground
- a transpiler that converts them between each other

###

versioning:

should have independent wo versions correlating to go ones like

1.23.3 - wo 1.23.3A, 1.23.3B
1.19 - wo 1.19A, 1.19B

offering them by major section

letters separated for compatability

### other todo

1. refactor test/wo/*.wo -> test/wo_*.wo if they don't get ran, also need to change some bat file to include it in the tests maybe
2. add automated tests in my own run_wo.bat


99. 
100. better icon if anyone offers one


### run without tests

```
cd src
./make.bat
```

### current state

running go on a .wo file seems to give "function main is undeclared in the main package"

current expected output:
\wo\src\runtime\proc.go:6630:13: internal compiler error: type pMask has no receiver base type
it happens on a function called `set`, and the formatted was getting tripped up on anything called set earlier.
So I renamed it to hashset, but it still does that same error.
go clean -cache didn't do anything. Nor deleting the generated compilers.
I renamed the erroring function to something else, and it gave the same error.
Was the hint from the go fmt just a red herring, and there really is a problem with this part of the code?
I already checked, this, the upstream branch, and go's master branch all matched on this part of the code.
I must have really messed up the compiler, since I don't even see anything wrong with the code that created this error
now I tried go clean -cache -modcache. same error
with extra debugging: type: pMask, name: set, kind: SLICE, isPtr(): false
still nothing... however, the "SLICE" is interesting. I traced the error backwards and found out that I accidentally replaced "TSLICE" with "TSET" in the base receiver type checking.

---
I think the compiler runs in this order:
build { serialize scanner parser resolver walk } -> make exe pkg/compiler -> compile /runtime -> go.exe
creation of the compiler runs:

`go test -run=Generate -write=all`
to create custom types
which I probably must run too

current commit syntax attempt to add:

`!ident()` - fails bc order of operations

`interface` is `tie` - fails because token defined in multiple places

recognizes `->` - fails bc doesn't belong anywhere

`set`

after adding/removing any fundamental types, you have to run

1. switch to default compiler or both root and path to this one
  - don't install to /wo/scr/cmd/vendor, put it in default Go source
  - otherwise, this means you're deploying the compiler with some extra tool that should be optional
2. go get -u golang.org/x/tools/cmd/stringer
3. go mod vendor
4. go install stringer
4. switch back path
5. run commands:
   - cmd\compile\internal\types $ stringer -type Kind -trimprefix T type.go
   - src/cmd/compile/internal/ir/ $ stringer -type=Op -trimprefix=O node.go
     - creates op_string.go
6. switch back compiler

### how I added another reserved word steps

the locations of, say, "int8" and "interface" for compiling reasons don't really have a universal location.

Those types, per the many steps and parts of the compiler, show up in many, many places that need to be accounted for.

Some of them deal with errors, some deal with representing the structure of the syntax in the code, some deal with
comparisons, and some deal with the actual functionality behind it.

All of these should be updated, and it's not as simple as adding it to each list, as you'd have to implement it as how
it is within that file.

According to my current specifications,

I want to add these tokens:
`set`, `some()`, `none`, `->`,
`enum`, `export`.

and modify the meaning of:
`:`, `!`, `<`, `>`, `:=`, `var`, `?`.

and possibly remove (ignore):
`iota`, `range`, `any`.

however, I would actually keep these as tokens, as it should compile in both Go and Wo, and they share the same type specifications.
This would also better allow errors like "Wo doesn't use the range syntax, try : " for example.

other:

https://github.com/golang/go/blob/e6626dafa8de8a0efae351e85cf96f0c683e0a4f/doc/go_lang.txt








