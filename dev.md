This project is a fork of the "go" repo, which is the source for the Go compiler, library, and testing.
This means that this repo can get compiled into a modified version of the Go compiler.
If you were to follow the [full instructions for a source install](https://go.dev/doc/install/source#bootstrapFromCrosscompiledSource) with this repo, your computer would have this
repo's compiler as your Go compiler in many targets of your choosing.
That is how you can install Wo. Since it is a modification, it's not safe to rely on it as a Go builder.
This is because, although it is meant to run .go files, it modifies those exact same files with extra logic
(by checking if the current file is .wo), but that could fail, and then `go` could fail.

This file is meant to be changed per commit, associated with details that apply to that specific commit e.g. current outputs.

For syntax linking in a code editor, it will depend on your IDE, but it wasn't hard to add ".wo" as a file type association in goland

#### Development timeline

- Compile, run, and test the default repo ✅
- Run some modified compiler ✅
  - run a .wo file (in a separate project) ✅

- For each feature in order of priority...
  - Optional: make a separate branch for it
  - Implement
    - Add its syntax in the compiler
      - With modularity
      - See [Operators](#Operators)
    - Reflection
    - Runtime functionality [✅ Set]
  - Create and run tests (for each target)
  - Update code formatter / syntax highlighting
  - Update specification and justifications
  - Merge

- Possibly make it into a `wo` command
  - or at least instructing on how to make an alias
- Website
  - Should showcase this project and just have the documentation
  - Write in Go
  - Then convert to Wo
  - Maybe an online playground
- Offer compiled binaries
- Transpiler that converts between each other

There is a priority to features to determine which ones to implement next. First: ones that allow experimenting, then by the most important / needed ones.
It is indicated it in the specification.

#### Versioning

should have independent wo versions correlating to go ones like

1.23.3 - wo 1.23.3A, 1.23.3B
1.19 - wo 1.19A, 1.19B

offering them by major section

letters separated for compatability

### Other todo

1. refactor test/wo/*.wo -> test/wo_*.wo if they don't get ran, also need to change some bat file to include it in the tests maybe
2. add automated tests in my own run_wo.bat


99. 
100. better icon if anyone offers one


### Commands to run this proj without tests

```
cd src
./make.bat
```

### Current state

If you run the default go compiler on a wo main file, you'd probably get something like "runtime.main_main·f: function main is undeclared in the main package"

Compile errors about set

---

I think the compiler runs in this order:
build { serialize scanner parser resolver walk } -> make exe pkg/compiler -> compile /runtime -> go.exe
creation of the compiler runs:

`go test -run=Generate -write=all`
to create custom types
which I probably must run too

#### Operators

After adding/removing any fundamental types, you have to run

1. switch to default compiler or both root and path to this one
   - don't install to /wo/scr/cmd/vendor, put it in default Go source
   - otherwise, this means you're deploying the compiler with some extra tool that should be optional
2. open new console
3. go get -u golang.org/x/tools/cmd/stringer
4. go mod vendor
5. go install stringer
6. switch back path
7. run commands:
   - cmd\compile\internal\types $ stringer -type Kind -trimprefix T type.go
   - src/cmd/compile/internal/ir/ $ stringer -type=Op -trimprefix=O node.go
     - creates op_string.go
8. switch back compiler

### How to add reserved words

the locations of, say, "int8" and "interface" for compiling reasons don't really have a universal location.
Those types, per the many steps and parts of the compiler, show up in many, many places that need to be accounted for.

Some layers deal with errors, some deal with representing the structure of the syntax in the code, some deal with
comparisons, some with reflection, and some deal with the actual functionality behind it.

All of these should be updated, and it's not as simple as adding it to each list, as you'd have to implement it as how
it is within that file. It is usually intuitive to model it off of how all the other types are being implemented.

According to my current specifications,

src/cmd/compile/internal/ir/node.go

I want to add these tokens:
`set`, `->`, `enum`, `export`.

and modify the meaning of:
`:`, `!`, `<`, `>`, `:=`, `var`, `?`.

and possibly remove (ignore):
`iota`, `range`.

However, still keep any removed ones as defined tokens as they already were.
It should compile in both Go and Wo, and they share the same type files and specifications but react to it differently.
This would also better allow compatibility tip errors like "Wo doesn't use `x := range xs` syntax, try `x : xs`" for example.

other:

https://github.com/golang/go/blob/e6626dafa8de8a0efae351e85cf96f0c683e0a4f/doc/go_lang.txt


#### Other soft errors

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

#### More

also see: [set.md](/set.md)



