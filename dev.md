This project is a fork of the "go" repo, which is the source for the Go compiler, library, and testing. This means that this repo can get compiled into a modified version of the Go compiler. If you were to follow the [full instructions for a source install](https://go.dev/doc/install/source#bootstrapFromCrosscompiledSource) with this repo, your computer would have this repo's compiler as your Go compiler in many targets of your choosing. That is how you can install Wo. Since it is a modification, it's not safe to rely on it as a Go builder. This is because, although it is meant to run .go files, it modifies those exact same files with extra logic (by checking if the current file is .wo), but that could fail, and then `go` could fail.

For linking, it will depend on your IDE, but it wasn't too hard to add ".wo" as a file type association in goland

Development steps:

running and testing this default repo

possibly making the `wo` command separate, or at least instructing on how to make it an alias.

modifying the compiler code to support that syntax

modifying te runner to support that transformation

checking if it still reaches all the targets like 

setting up the website in Go and/or Wo

dealing with versions and downloadable executable installations for other users to test, perhaps offering an online playground

a transpiler that converts them between each other

---

my own notes:

```
cd src
all.bat
```

current version syntax:

`!ident()`

`interface` is `#`

recognizes `->`

`set`
