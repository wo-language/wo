In here, I give explain, argue for, and document Wo's design decisions. Maybe it could be interesting to some people, but it serves as an important documentation to me, so I know what has already been considered and tested.

### Index

- [Theory](#Theory)
  1. The goals of code
  2. To restrict or to follow
  3. Conventions
- Conventions
  1. Variable naming
- Features
  1. variable
  2. interface{}
  2. set

## Theory

When someone makes a new programming language, it should solve a problem, not just do something that vaguely feels attractive because it combines that paradigm from that language is based on C so it's fast.

A similar situation was Scala's improvements over Java. It clearly improved the syntax and design, especially with pattern matching, and, importantly, interoped with Java.

I have not seen anything in the programming language landscape like this - a direct child of Go that addresses its design.

And let it be known, the internet is full of needless speculation as to "why Go did this and not that?", so I'll also make this transparent and be skeptical of any justification that people give to Go. Instead, prioritizing the objective best way of doing something regardless of what the theories originally purported. In other words, this is about what happens in practice, not how logically sound or nice the theory is behind it. For example, Vim sounds crazy on paper to people the first time they hear of it, thinking "but why can't you type by default!?", but only once they start to try it out do they realize that it's incredible to use in practice. Or, hypothetically, they end up realizing it's terrible, and they simply enjoy hurting their hands with the arrow keys.

Ultimately, certain improvements can be more valuable in different circumstances, while it doesn't matter in others. For example, I had a Java program that simplifies math expressions, and making that [one file](https://github.com/Branzz/DiscreteMath/blob/scala_integration/src/bran/tree/compositions/expressions/operators/OperatorExpression.scala#L452) into Scala out of the whole project shortened [that code](https://github.com/Branzz/DiscreteMath/blob/scala_integration/src/bran/tree/compositions/expressions/operators/OperatorExpression0.java#L223) by about 2.5 times as much because of pattern matching, but all the other files were fine being Java.

So it is just nice to have the option of a design that is attuned to your circumstances, as opposed to some forced grifting replacement for all of Go that must be better because the author thinks so.

This is also why I am planning to make the features modular / able to be swapped in and out.

### The goals of code

Code communicates and guarantees that it achieves something when ran by a computer. These two fight with each other in ways I won't be able to describe fully here. Intention seems to be an important part of Go's design, and I believe it is important. Just keep this in mind for later.

I believe that adding comments is to compensate for code that doesn't communicate. They should be rarely used in practice, only for things like magic constants and documentation. Even with documentation, it should be obvious what a function is going to do from its name. The syntax and style of a language along with the programmer's design of the code, such as identifier names and logical design, all contribute to the given "intent" of the program. Therefore:

- The compiler should not force you to be vague.
  - and preferably should help you to be clear.
- You should avoid being vague when given the choice.
  - Like in variable and function names.

It's our job to pay attention to details, but let's still make it as easy on ourselves as possible.

The current state of Go's compiler forces you to be vague, and their style designs recommend using vague variable and function names. This isn't exactly a criticism, but just a description of how Go appears to me.

### To restrict or to allow

Should we allow bad language and hope that users don't use it? Or should we completely ban it...

When I choose a syntax, it has to not already exist, or if it exists, in an unmistakably different context. If the most obvious or common choice is used elsewhere, or if it would confuse the compiler, then I will consider a syntax that seems completely foreign. All programming syntax was foreign at first, and plenty of people try out new ways of writing things that they hate at first and end up getting used to and enjoying. If you can't adapt to something like that (granted the syntax isn't atrocious), modern programming probably isn't right for you.

For example, Learning Go 2nd edition says:

> Note: The Go compiler won’t stop you from creating unread package-level variables. This is one more reason you should avoid creating package-level variables.

This is absolutely backwards logic to me from a compiler's perspective. It shouldn't allow you to do something that you shouldn't do.

### Why modularity

I am considering making different language features **modular**. That is, they can be enabled or disabled either through a compiler flag, in the module file, or with some header.

This isn't the newest idea, as languages all have versions one can choose of their liking. Rust has that capability with something like `#![allow(unused)]` which allows unused variables in the entire file.

If someone just likes only the interface syntax, and that's all they want, then they can still use Wo in that way without dealing with the parts they don't like.

For example, enforcing the type before the variable name is universally disagreed on, so this could just be an additional option, not the Wo default. If a feature isn't restrictive, then it doesn't need a flag.

That means there are these types of features: Those enforced without an option, those that are on by default, those that are off by default. All of them except experimental or "indifferent" ones would be enabled by default.

## Conventions

### Variable naming

One of the very first things I learned when I started programming is that using 1 or 2 length variable names in most situations was incredibly bad practice that leads to misunderstandings. You probably already know why, but in case you don't, I will explain below.

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

You check the docs or source code for `roleHierarchy` and find out it returns a `Node`. And you check the `Node` struct, which contains a `val Role` a `left Node` and a `right Node`.

It turns out, the `t` in `t.leftBranch().cut()` was just a tree.

Why should I have to analyze any of this when 3 characters would have explained enough. If the code used `tree` as the variable name, none or almost none of this would have been necessary - even with better documentation. We would have read that single word, and moved on to the next thing, rather than being disrupted.

> Good code is not vague.

`t` is more objectively more vague than `tree`.

`t` does not declare intent.

```go
tree.leftBranch().cut()
```

`tree` is vague to a certain level, as I could be referring to a literal tree or a programmatic tree, but `t` has much less meaning.

Go code is exactly that situation of searching for the meaning of shortened names over and over again.

The two principles of using more lines in code along with shortening variable names contradict each other.

Code either extends vertically (less functional abstraction) or horizontally (more function calls, longer names). Shortening names and using loads of null checking both go in the direction of vertical. Please, take your hand off the scroll wheel (or the `hjkl`). In between these two directions is a more square shaped code. And the other extreme typically happens with nested function calling, like some overly clever and lengthy Java streams solution.

![Image of 3 code editors of code that's tall, square, and wide](https://raw.githubusercontent.com/wo-language/wo-info/refs/heads/main/wo%20resources/code_rectangles_whiteborder.png)

As you can see, the first code editor has 8 lines and reaches the first line, then compressed to 6 lines and reaches the second line, then to just 4 lines. I tried to make them each have the exact same volume of "code".

I find shortened variable names and repeated 3-4 line checking less readable, so I chose to make Wo more towards the square.

There is one situation where shortened variable names might be acceptable, which is lambda function calls like `starfruits.map(s -> s.weight() * 2.2)`, or generally very short function calls, where you can easily see what `s` means.

In the same realm is shadowing and keyword overloading, which I go into later.

Wo also does not allow variable names like `π` or `__`.

And it's hard to do something about the inconsistent capitalization in functions. There certainly shouldn't be both "Init" and "init" in the same file, though. Nonetheless, Wo uses camelCase function and variable names.

## Syntax Features

### interface{}

I chose `<T>` for `interface{T}`. I considered something like `~`, but you can't wrap around with that. There was also `#{}`, but the shortness of `<>` was more attractive. Tags are a symbol that are not used in Go and is already associated with types.

```go
func f(a <>, b <bool>) {
}
```

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
which is a bit weird unless you really like C++. Using `interface` in the type declaration doesn't feel exasperating anyway.

So I'll keep `type NAME interface` for now.

### Unused variables

```go
func main() {
    var x = 3 // error: Unused variable 'x'
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

It's basically an enhanced for loop. `range` isn't used anywhere else, and it's not like you can assign a variable as a range. I think they needed `range` because `i, v := nums` alone is misleading since that doesn't actually return the index and value, but we can just get around that by doing what Java did back when I was a baby:

```go
for i, v : nums {
    sum += v
}
```
as `:` is given a new meaning.

Note that this will use the **value** by default for a single variable

`for value : nums {}`

In Go, when there is a single variable declared, it is for the index

`for index := range nums {}`

I see this as "memorized information"; it's arbitrary. There's no way of knowing whether it's "for the indices in the range of a collection" or "for the values in the range of a collection" without seeing it before. Since people are used to that way, switching it could be confusing, but I don't really want to rely on something like that when my goal is to offer an alternative design. Additionally, `for value : values` is the common pattern seen in other languages anyway, so it shouldn't really be surprising that `for i : values` isn't actually the index when taken out of the context of Go.

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
 e = 2        // possible, just assigns
 e, y := 2, 3 // possible
 e, _ := 2, 3 // not possible
 if d == 1 {
     d := 2   // possible
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

Here's every possibility in Wo according to that grid:

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

I assume one of the reasons it allows overloading reserved words (`int`, `nil`) is because of backwards compatability, but I simply don't need that since this is a fresh start for syntax. Allowing the ability to override those is always confusing and unsafe. Words spelled the same with different meanings used in the same exact contexts, which can be done by accident, is totally confusing. Enough said.

### Overloading package names

This is just allowed because a package could be called anything, but it shouldn't be allowed without some kind of error. I'm taking for granted people don't necessarily rely on IDEs here. For example,

```go
import { "strings" }

func combineThem(strings /* Wo Error */ []string) string {
    return strings.Join /* Go error */ (strings, ", ")
}
```

would not compile in Go, but not because of the existence of the variable name `strings`, but because `Join` is being called on that variable, when the author intended for it to be `Join` from the `strings`. It is because `strings` is overloading the `"strings"` package from the `import`.

By the way, in Wo, I plan to make it so that one could just skip importing `strings` and just be able to call `Join` on a `[]string` like `stringsVariable.Join(", ")`. This could help contribute to avoiding these situations, but it could still happen of course.

One way around it is to rename `strings`, but this is a perfectly good variable name that might be used frequently across the file. This means a better alternative would be to use the `"strings" as "string_util"` syntax, or to differentiate the formatting of packages used in code like `@strings.append` as a rudimentary example.

By the way, I dream of a language where all the reserved words have some symbol, and you write all your own stuff like regular words and spaces like `bird $get color` for `bird.get(color)`, and you get to define the meaning of all your own sentences by token order like some declaration `String A "with" String B -> concat(A, B)` or `Number A (Number B) -> A * B`. Or maybe Haskell has invaded my subconsciousness?

## Overloading functions

Preventing function overloading sounds like a good idea in theory, but in practice it results in artificially lengthening function names, when their original form was already the most descriptive. A description of something can be done by its contents; the parameters describe the function already, there is no necessity to change the name when the parameters change too. It aligns with language and nature.

Without overloading, shortened function names:

```go
payWithC(cash)
payWithCI(creditCardInfo)
payWithCNZ(creditCardNumber, zipCode)
payWithCNCS(creditCardNumber, city, state)
```

With long description function names (worst of both worlds):

```go
payWithCash(cash)
payWithCreditInfo(creditCardInfo)
payWithCreditNumberZip(creditCardNumber, zipCode)
payWithCreditNumberCityState(creditCardNumber, city, state)
```

With the same function names:

```go
payWith(cash)
payWith(creditCardInfo)
payWith(creditCardNumber, zipCode)
payWith(creditCardNumber, city, state)
```

which has much less redundant information.

I have used this aspect of programming many, many times. It has been far from the top of the list of things that could make my code vague, and I'm not convinced that it's ever a primary culprit.

The real problem is that, when it comes to compiling, we can't just know which one you're referring to when the type parameters are vague. It requires some type analysis, but I believe it is totally possible here. It is very close to the line of inheritance making sense, as a type being vague with another one implies some shared classification. The most obvious thing to be would be some structures which are just unsafe pointers underneath. Still, I think this can be done at compile time as these types were designed to be strict, e.g. `[3]int != [2]int != []int != []*int`.

Additionally, `[:]` already does this. It's equivalent to `slice(start=0, end=0, max=0)`.

## i

`x := 5 - 3i`

I vote to keep this since this is cool and kind of funny. It doesn't intersect with any other syntax.

## Ternary

There is `?:` and `if else`, but let's look at more possibilities
.
In Go,

```go
var hvac
if indoorTemp < outdoorTemp {
    hvac = heating
} else {
    hvac = ac
}
```

could also be represented with

```go
var hvac = if indoorTemp < outdoorTemp { heating } else { ac } // no parens
```

or

```go
var hvac = if (indoorTemp < outdoorTemp) heating else ac
```

or

```go
var hvac = indoorTemp < outdoorTemp ? heating : ac
```

Despite the first one being the longest, I think I'll go for that one since it is consistent with the `if cond` no parentheses style that Wo already has (which was inherited from Go). Another way around the parentheses problem is to do something like this

```go
var hvac = if indoorTemp < outdoorTemp then heating else ac // imagine then is highlighted
var hvac = if indoorTemp < outdoorTemp ? heating else ac
var hvac = if indoorTemp < outdoorTemp ? heating : ac
var hvac = heating if indoorTemp < outdoorTemp else ac
```

Since it is an expression, an important question is what these look like when applied at more depth

```go
var hvac
if indoorTemp < outdoorTemp {
    if thermostat > indoorTemp {
        hvac = heating
    }
    hvac = none
} else {
    if thermostat < indoorTemp {
        hvac = ac
    }
    hvac = none
}

var hvac = if indoorTemp < outdoorTemp { if thermostat > indoorTemp { heating } else { none } } else { if thermostat < indoorTemp { ac } else { none } }
var hvac = if (indoorTemp < outdoorTemp) if (thermostat > indoorTemp) heating else none else if (thermostat < indoorTemp) ac else none
var hvac = indoorTemp < outdoorTemp ? thermostat > indoorTemp ? heating : none : thermostat < indoorTemp ? ac : none
var hvac = if indoorTemp < outdoorTemp THEN if thermostat > indoorTemp THEN heating else none else if thermostat < indoorTemp THEN ac else none
var hvac = if indoorTemp < outdoorTemp ? if thermostat > indoorTemp ? heating else none else if thermostat < indoorTemp ? ac else none
var hvac = if indoorTemp < outdoorTemp ? if thermostat > indoorTemp ? heating : none : if thermostat < indoorTemp ? ac : none
var hvac = heating if thermostat > indoorTemp else none if indoorTemp < outdoorTemp else ac if thermostat < indoorTemp else none
```

and also with else if

```go
var hvac
if indoorTemp == outdoorTemp {
    hvac = off
} else if indoorTemp < outdoorTemp {
    hvac = heating
} else {
    hvac = ac
}

var hvac = if indoorTemp == outdoorTemp { off } else if indoorTemp < outdoorTemp { heating } else { ac }
var hvac = if (indoorTemp == outdoorTemp) off else if (indoorTemp < outdoorTemp) heating else ac
var hvac = indoorTemp == outdoorTemp ? off : indoorTemp < outdoorTemp ? heating : ac
var hvac = if indoorTemp == outdoorTemp then off else if indoorTemp < outdoorTemp then heating else ac
var hvac = if indoorTemp == outdoorTemp ? off else if indoorTemp < outdoorTemp ? heating else ac
var hvac = if indoorTemp == outdoorTemp ? off : if indoorTemp < outdoorTemp ? heating : ac
var hvac = off if indoorTemp == outdoorTemp else heating if indoorTemp < outdoorTemp else ac
```

Hmmmm... Surprisingly, I don't think any of these are vague to the compiler given you go by right to left associativity.

And no, you really shouldn't be making ternary statements ridiculously complicated, but I need to make sure those are still parsable and still readable.

I know the last one is weird, but it is actually very interesting. It goes like this: "if (A) {B} else {C}" -> "B if A else C". It ends up in a binary tree shape. You can still read it left to right in plain English. For example, "if you know C, your life is great" can also have it phrased as "your life is great if you know C", but places slightly more emphasis on the value than the condition.

It also avoids some of the vagueness of

`var hvac = indoorTemp == outdoorTemp...`

seeming like `var (hvac = indoorTemp)` at first glance.

I'll combine the `else if` with the further depth: `var hvac = off if indoorTemp == outdoorTemp else heating if thermostat > indoorTemp else none if indoorTemp < outdoorTemp else ac if thermostat < indoorTemp else none`

The hvac is off if it's the same temperature indoors as it is outdoors, otherwise it's heating if the thermostat is higher than the indoor temperature, otherwise it's no hvac if the indoor temperature is lower than outdoor temperature, otherwise it's AC if the thermostat is lower than the indoor temperature, otherwise it's none.

That syntax could also be nice for short assignments like

`value = dereference(input) if input.isPtr() else input`

I'm also kinda interested in the one with "then" since it does make reading it more obvious without requiring `()` or `{}`.

However, I think that first basic option is the most readable since you can clearly see the depth level with the curly braces.

But it's *just* the curly braces that make it easy to read for me. How about I apply them to the other ones despite being redundant (to be thorough):

I'll combine the two: `var hvac = { off } if indoorTemp == outdoorTemp else { { heating } if thermostat > indoorTemp else { { none } if indoorTemp < outdoorTemp else { { ac } if thermostat < indoorTemp else { none } } } }`

I find this kinda weird. Next:

`var hvac = indoorTemp == outdoorTemp ? { off } : { indoorTemp < outdoorTemp ? { thermostat > indoorTemp ? { heating } : { none } } : { thermostat < indoorTemp ? { ac } : { none } }`

That's ok, but why not just use `if`/`else` in place of those:

`var hvac = indoorTemp == outdoorTemp if { off } else { indoorTemp < outdoorTemp if { thermostat > indoorTemp if { heating } else { none } } else { thermostat < indoorTemp if { ac } else { none } }`

I should only stray from the most obvious variation of Go when it's a clear improvement over the status quo. Idk if these really are, I personally like them, but I can already hear the angry voices insisting any of these are evil. Which points to a bit of a reality here: nothing will only be praised, and nothing will only be shamed...

I'm going with `v = 2 * (if cond { a } else { b }) + ` for now, despite the curly braces feeling EXASPERATING to add in.

I'll do this by either adding an expression identical to the `if else` statement, or modify the statement to become an expression.

## Functional Features

### set

Implementing this really wasn't all that interesting or challenging.

Because of the way `map` is designed (a hashmap), the keys and values are stored rather insignificantly, and removing the values from its structure was pretty simple to do. There isn't anything special about the difference between key and value besides that one part gets hashed and one part doesn't. So, yes, I copied map and refactored it; I really don't think there is any faster way to do *this map* for *its* intended purposes without also improving hashmap. That wasn't really my intention with Wo, but if someone sees a way of seriously improving the native hashmap when it doesn't have values, or if I spot an obvious one, then let's go ahead. I haven't changed the time complexity, but it should technically be insignificantly faster than `map[]struct{}`.

So I made `map`'s keys as `set`'s elements, wiping any functionality with `map`'s values.

This does mean the removal of the `val = m[key]` method, as that doesn't really mean anything for sets. Instead, I modified and kept the `_, ok = m[key]` method, using it like `ok = s[elem]`.

```go
primes set[int] = { 2, 3, 5 }  // declaration
ok = primes[4]                 // is ok if contains elem
primes.insert[7]               // insert / add
primes.delete[3]               // delete / remove
```

I prefer `add` and `remove`, but the naming (from `map`) uses `insert`, so I don't want it to get too inconsistent and therefore unpredictable. It's not an impossible consideration, however, but I'd prefer renaming the `map` methods too in that case.

There are also fast versions of the map for `strings`, `int32`, and `int64`, which Wo also has implemented.

And Wo also support a `sets` package in the same ways that the `maps` package does.

Sets in math use `{ }` to mean "unordered, unique collection", but in Go, which uses EBNF, it means "ordered, repeatable collection". I think it is ok to use the curly brackets for sets, since it is programmatically ordered and repeatable data at first, but then it will become converted from that explicit representation into something which is guaranteed to be an actual set. I can actually still say `{ a, b, a, c }` in math, but it represents a set of `a`, `b` and `c` without order. It is also predictable with the formatting already used with arrays and maps. If someone made their own set or any kind of math collection, it'd use the curly brackets.

## Array

I've concluded that it's not feasible to use `arr[]` because of how it interacts with map.

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
map[byte[]]map[byte[]]string[][]     // arr[] --- vague
map[[]byte, map[[]byte, []string]][] // map[A, B]
map[byte[], map[byte[], string[]]][] // map[A, B] and arr[]
```

The second one is ambiguous, since it could mean a double array of strings, which doesn't happen when we use `map[A, B]`

The last one prefers depth, so it ends up pushing more symbols to the end.

For arrays, I say either keep []arr with map[A, B], or just don't make any changes

### Map[K, V]

I understand `map[key]value` is supposed to reflect the `func(input) val` pattern, as well as the `value = map[key]`, but there is nothing about the fundamental concept of maps that imply they should reflect the "return type afterwards" pattern. If anything, `map[key]` should not necessarily mean "get", it could have meant `contains` or `indexOf` as arrays do with `[index]`. `get(key K) V {}` will already represent the function format, since it is just a function. There aren't many other options besides `map[key, value]`. However, I think Go's is still better in practice.

I think this is too disruptive and unnecessary of a change as shown in the previous section, so I'll keep `map[key]value`

### enum

```go
type Days enum {
    Sunday("sun", false),
    Monday("moon", true)

    root    string,
    workday bool
}

Sunday.position
Monday.working
```

### union

```go
type Point union {

}
```

### type

`struct S {}`
`interface I {}`

---

## Scope control

There are over 100 "halls of shame" in Go's source code, which is a kind of comment they have that links to repos that used a function that it "shouldn't have".

It's not really a laughing matter at that point, programs should be able to represent who gets access to what. Or, they should be given proper solutions to the work-arounds that they had to use.


