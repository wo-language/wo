
## Theory

When someone makes a new programming language, it should solve a problem, not just do something that vaguely feels attractive because it combines that paradigm from that language and it's based on C so it's fast.

A similar situation was Scala's improvements over Java. It clearly improved the syntax and design, especially with pattern matching, and importantly, interoped with Java.

I have not seen anything in the programming language landscape like this - a direct child of Go.

Ultimately, in practice, certain improvements can be more valuable in different circumstances while it doesn't matter in others. For example, I had a program that simplifies math expressions, and making that [one file](https://github.com/Branzz/DiscreteMath/blob/scala_integration/src/bran/tree/compositions/expressions/operators/OperatorExpression.scala#L452) Scala out of the whole project shortened [that code](https://github.com/Branzz/DiscreteMath/blob/scala_integration/src/bran/tree/compositions/expressions/operators/OperatorExpression0.java#L223) by about 2.5 times as much because of pattern matching, but all the other files were fine being Java.

So it is just nice to have the option of an improved design, not a forced grifting replacement for all of Go.

### The goals of code

Code communicates and guarantees that it achieves something when ran by a computer. These two fight with each other in ways I won't be able to describe fully here. Intention seems to be an important part of Go's philosophy. Just keep this in mind for later.

I believe that comments are to compensate for code that doesn't communicate. They should be rarely used in practice, only for things like magic constants and documentation. Even with documentation, it should be obvious what a function is going to do from its name. The syntax, style, as well as the programmer's design of the code, such as identifier names and logical design, all contribute to the given "intent" of the program. Therefore:

- The compiler should not force you to be vague.
- You should avoid being vague when given the choice.
  - Like in variable and function names

Go's compiler forces you to be vague, and their style designs recommend using vague variable and function names. This isn't really a criticism, it's just me pointing out a fact.

### To restrict or to allow

Should we allow bad language and hope that users don't use it? Or should we completely ban it...

When I choose a syntax, it has to not already exist, or if it exists, in an unmistakably different context. If the most obvious or common choice is used elsewhere, or if it would confuse the compiler, then I will consider a syntax that seems completely foreign. All programming syntax was foreign at first, and plenty of people try out new ways of writing things that they hate at first and end up getting used to and enjoying. If you can't adapt to something like that (granted the syntax isn't atrocious), modern programming probably isn't right for you (it's not right for most of the population anyway).

For example, Learning Go 2nd edition says:

> Note: The Go compiler won’t stop you from creating unread package-level variables. This is one more reason you should avoid creating package-level variables.


## Conventions

### Variable naming

One of the very first things I realized when I first started programming is that using 1 or 2 length variable names in most situations was incredibly bad practice that leads to misunderstanding. You probably already know why, but in case you don't, I will explain below.

Let's say you came across this, 40 lines deep into a function:

```go
t.leftBranch().cut()
```

What does this mean?

You go check the definition:

```go
t := roleHierarchy()
```

What?

You check the source code for `roleHierarchy` and find out it returns a `Node`. And you check the `Node` struct, which contains a `val Role` a `left Node` and a `right Node`.

It turns out, the `t` in `t.leftBranch().cut()` was just a tree.

If the code used `tree`, none or almost none of this would have been necessary - even with better documentation. We would have read the single word, and moved on to the next thing, rather than being disrupted.

```go
tree.leftBranch().cut()
```

`t` is more objectively more vague than `tree`.

`t` does not declare intent.

`tree` is already technically vague, as I could be referring to a literal tree or a programmatic tree, but `t` has much less meaning.

Go code is exactly that situation except over and over again.

This part of the "Go philosophy" is not anything more than (what most philosophy actually is) circular reasoning and non sequiturs bundled up to seem logically conclusive.

It claims both "use longer code for more meaning" and "use shorter code for more meaning".

Code either extends vertically (less functional abstraction) or horizontally (more function calls, longer names). Shortening names and using loads of null checking both go in the direction of vertical. Please, take your hand off the scroll wheel. In between these two directions is a more square shaped code. And the other extreme is some overly clever and lengthy Java streams solution. I find that more readable than shortened variable names and repeated 3-4 line null checking, because at least you can usually read it without checking definitions, so I chose to make Wo more towards a square.

There is one situation where shortened variable names might be acceptable, which is lambda function calls like `starfruits.map(s -> s.weight() * 2.2)`, or generally for very short function calls, since you can see on that or the previous line what `s` means.

In the same realm is shadowing and keyword overloading, which I go about later.

Wo also does not allow variable names like `π` or `__`

And it's hard to do something about the inconsistent capitalization in functions. There certainly shouldn't be both "Init" and "init" in the same file, though. Nonetheless, Wo always uses camelCase function and variable names.

## Features

### interface{}

I chose `<T>` for `interface{T}` because it can still wrap around some T, and it's a symbol associated with types. I was also considering something like `#{}`, but the shortness of `<>` was more attractive.

As for

```go
type I interface {
    bool
}
```

If I follow the same suit, it ends up being
```go
type I <
    bool
>
```
which is just unnecessary. Using `interface` in the type declaration doesn't feel exasperating anyway.

### Unused variables

```go
func main() {
    x := 3 // error: Unused variable 'x'
}
```

Wo simply allows this. It will become a warning and compiled away. The reason Go doesn't do this is probably because of how it optimizes variables. It does allow unused `const` for the same reason, since those are easier to optimize. However, it is not otherwise impossible to optimize unused variables away.

### For range

The default `for range` syntax is
```go
for i, v := range nums {
    sum += v
}
```

It's basically an enhanced for loop. I think they needed `range` because `i, v := nums` alone is misleading since it doesn't return the index, but we can just get around that by doing what Java did back when I was a baby:

```go
for i, v : nums {
    sum += v
}
```
as `:` gives a new meaning.

Note that this will use the **value** by default for a single variable

`for value : nums {}`

In Go, when there is a single variable declared, it is for the index

`for index := range nums {}`

I see this as "memorized information"; it's arbitrary. There's no way of knowing it's "for the indices in the range of a collection" or "for the values in the range of a collection" without seeing it before. Since people are used to that way, switching it could be confusing, but I don't really want to rely on that when offering alternative design. Additionally, `for value : values` is the common pattern seen in other languages anyway, so it shouldn't really be surprising that `for i : values` isn't actually the index when taken out of the context of Go.

I chose to make it the **value** by default as it would be more common and intuitive as one seems to want to ignore the index by nature of using the enhanced "for an item *in* items", possibly opting for a traditional `for i = 0; i < len; i++` otherwise, or just using `for i, _ : values {}` for access to the index.

- That could be problematic when frequently using this when trying to modify arrays by their index, so range could be kept to mean "range of indices over"

### The possibilities of variable declaration and assignment

Go offers these styles of declarations:
```go
//z := 1 // not possible at package level
func declares() {
var a int = 1
b int := 1 // not possible
var c = 1
d := 1
var e int
f int // not possible (unlike C style int f;)
var X // not possible
var (
    // all the same things it could already do
)
var m, n int = 1, 2
o, p int := 1, 2 // not possible
var q, r = 1, 2
s, t := 1, 2
var u, v int
w, x int // not possible

// note: e was already declared (we wouldn't need this note if shadowing were more explicit)
e := 2       // not possible
e = 2        //   assigns
e, y := 2, 3 //   possible
e, _ := 2, 3 // not possible
if d == 1 {
    d := 2 // possible
    m, n := 3, 4 // possible
}
// at this point,
// d == 1
// m, n == 1, 2 

fmt.Println(a, b, ... x, y) // haha
}
```

Let's reduce this down as much as possible to rules for describing what's above. I'll use a logical grid strategy.

```
var a, _ int = 8, 9  var ( b = 8 )
var c string         var ( d, e string )
f := 10              a, z := 10
const g int = 8      const ( h = 8 )
f, _ = 10, 11        c = "aoeu"
```

R - Required, ! - None, ? - Optional, # - Other

Each column will represent these in order. So the last column of this grid represents an assignment

```
const  ! ! ! R !
var    R R ! ! !
(...)  ? ? ! ? ? // var ( a = 1, ... )
names  R R R R R
  :=   ! ! R ! !
types  ? R ! ? !
litral ? ? ? R ?
  =    R ! ! R R
values R ! R R R
at pkg ? ? ! ? ?
multi  ? ? R ? ? // a, b =
shadow ! ! # ! ?
count  2 1 2 2 2

var names (type  |  (= values))  |
var \( (names (type  |  (= values)))... \)   |
name := value {not in package level}

// where names is really (name | (name | _)...)
```

That adds up to 9 main possibilities by ignoring literal, package, multi declaration, and `( ... )`

For C, it would be like this, ignoring anything that became all disallowed

```
prefxs ? ! (const, volatile, static)
types  R !
names  R R
  =    R R
value  ? R
```

the first column represents `long id = 16`

and the second is `id = 32`

These two systems both result in basically the same exact thing, except adapting to different needs for optimization.

Can I make it stricter without sacrificing Go's functionality?

```
const  ? ! R ! ! #
names  R R R R R R
types  R ! ! ! ! R
  =    R R R R ! R
values R R R R R R

var    ! R ! ! ! R
(...)  ! ! ! ! ! R

  :=   ! ! ! ! R !
shadow ! ! ! ! R !
count  2 1 1 1 1 2
```

I made const like a prefix, required the type, only allow var for multi line and untyped decl, allow everything in or out of the package, and actually added more specification than before by requiring `:=` for shadowing.

This is 6 possibilities without including `(...)` as before, but 8 otherwise, which takes the original 9, removes 2 redundant ones, then adds 2 restrictive one. It adds conditions so there are not multiple ways to do the same task. Only one assignment, only one const declaration, only one shadow, etc., but does not merge untyped declaration.

Here's every possibility:

```go
x int = 8
a var = 8
y const int = 9 // maybe: const y int = 9
b const = y // maybe: b const var = y
x = b
{ x := 0 }
var ( z, _ (int, error) = count() )
    // ~~ g const int        = 84 // mixed ~~
const ( b string = "永", e error = nil )
```

`var` and `const` are not placed at the start of the line to keep the variable names inline with each other, but puts it before the type since that's what it is a part of, so it is then readable with the type as "y is a constant integer" and "a is a variable".

TBD I could put both `const var` together to show that the **type** is variable.

TBD I'm not sure whether to require the type to be stated. The one time that could be annoying is for tuples, since their types are a bit bigger. I tried without the parentheses, it becomes hard to tell between the name and type.

Requiring the value does mean that zero values are removed. It might sound like a dramatic change, but zero values are unpredictable; they do not declare intent. string's default value is "", even though it is similar to a char* internally and most languages default strings to null, making that more expected, especially since nullptr is 0. I could just make string nullable / require their value to be defined, but I'm removing it for all the other types anyway because they're still vague.

I recycled the same symbols by keeping just some of their usage with `var ()` and `const ()`. As well as with `:=`.

I made `:=` as only shadow since that's what it already does, but now it will be seen as a rare symbol. Like the table implies, `=` can't be used for the exclusive case of shadowing. This serves as an alert. If either you make a new variable that is getting shadowed by later code, or if you name something the same as earlier code, this will error "`can't shadow a variable with =`". And if you were to unshadow something like as just described, then the reverse would happen "`can't assign a variable with :=`".

`:=` can theoretically happen at the package level, while it is not allowed in Go.

### Multi variable declaration

In Go, you can assign multiple variables like this
```go
getPoint := func() (x int, y int) { return 1, 2 }
a, b := getPoint()
_, c := getPoint()
d, _ := getPoint()
```
but you would not be able to do `p := getPoint()` nor could you do `y := getPoint()` to select just the `y` part of the return values. One pattern is to use a struct, but let's see what it could look like in a longer flow of functions

```go

```

Although it's more of an exception, since `for i := range stuff {}` is already allowed to shorten the `i, _` in Go, it's easy to justify it being allowed elsewhere.

TBD I will probably remove it for multiple valued assignments like

```go
a, b = 1, 2
```

as this is tricky to read and unnecessarily horizontal.

### Overloading reserved words
I assume one of the reasons it allows overloading reserved words (`int`, `nil`) is because of backwards compatability, but I simply don't need that since this is a fresh start for syntax. Allowing the ability to override those is always confusing and unsafe.

### Overloading package names
This is just allowed because a package could be called anything, but it shouldn't be allowed without some kind of error. I'm taking for granted people don't necessarily rely on IDEs here. For example,
```go
import { "strings" }

func toNames(wombats []Wombat) bool {
    strings := []string{}
    
    for _, wombat in range wombats {
        
    }
    return strings
}
```

would not compile in Go, but not because of the existence of the variable `strings`, but because `append` is being called on that variable.

In Wo, the syntax would look like this:

```go
import { "strings" }

func areAnyHairyNosed(wombats Wombat[]) bool {
    strings string[] = {} // error
    TODO
    
    return false
}
```

However, it would not compile for a different reason, because `strings` is overloading the `"strings"` package.

In Wo, one could skip importing `strings` and just use TODO

One way around it is to rename `strings`, but this is a perfectly good variable name that might be used frequently. This means a better alternative would be to use the `"strings" as "string_util"` syntax, or to differentiate the formatting of packages used in code like `@strings.append` as a rudimentary example.

## Map

I understand `map[key]value` is supposed to reflect the `func(input) return` pattern, as well as the `value = map[key]`, but there is nothing about the fundamental concept of maps that imply they should reflect the "return type afterwards" pattern. If anything, `map[key]` should not necessarily mean "get", it could have meant `contains` or `indexOf` as arrays do with `[index]`. `get(key K) V {}` will already represent the function format, since it is just a function. There aren't many other options besides `map[key, value]`. However, I think Go's is still better in practice.



I think this is too disruptive and unnecessary of a change, so I kept `map[key]value`

## Array

Map is declared and called like this:

`m := map[k]v`

`v := m[k]`

but array is declared and called like this:

`a := [x]arr`

`x := a[i]`

The odd one out is `[x]arr`, which has the array marks as a prefix. What if it were the suffix?

For example, what about an array of a map from keys of arrays of bytes to values of (maps with keys of byte arrays to values of arrays of strings)
```
[]map[[]byte]map[[]byte][]string     // Go
map[byte[]]map[byte[]]string[][]     // arr[] - vague
map[[]byte, map[[]byte, []string]][] // map[A, B]
map[byte[], map[byte[], string[]]][] // map[A, B] and arr[]
```

The second one is ambiguous, since it could mean a double array of strings, which doesn't happen if we use `map[A, B]`


The last one prefers depth, so it ends up pushing a lot of symbols to the end.

I say either keep []arr with map[A, B], or just don't make any changes

## i

`x := 5 - 3i`

I vote to keep this since this is cool and kind of funny. It doesn't intersect with any other syntax.

## ternary

`? :`

`if cond {} else {}`
## set

`set[E]`

## enum

```go
type Days Enum {
    Sunday("sun", false),
    Monday("moon", true) // no comma

    root    string,
    working bool
}

Sunday.position
Monday.working
```
## union
```go
type point union {

}
```
## type

`struct S {}`
`interface I {}`
