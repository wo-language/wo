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

There is a priority to features. First: ones that allow experimenting, then by the most important / needed ones.
I'll indicate it in the specification.

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

if you run the default go compiler on a wo main file, you'd probably get something like "runtime.main_main·f: function 
main is 
undeclared in the main package"

compile errors about set

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

src/cmd/compile/internal/ir/node.go

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


#### Other soft errors

[] Implemented - need to test
- Only ignoring unused

Any of these could be ignored easily in theory:

unused var in switch
no new variables on left side of :=
generic function is missing function body
label %s declared and not used
main - func %s must have no type parameters, func %s must have no arguments and no return values
init - missing function body
%q imported and not used, %q imported as %s and not used
can only use ... with final parameter in list
cannot range over, range over %s permits no iteration variables, range over %s permits only one iteration variable, range clause permits at most two iteration variables
cannot use iteration variable of type %s
cannot use type %s outside a type constraint: interface is (or embeds) comparable
in interface - overlapping terms %s and %s


src/cmd/compile/internal/typecheck/stmt.go
src/internal/types/errors/codes.go
src/cmd/compile/internal/types2/errors.go





