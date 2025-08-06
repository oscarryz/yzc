## Finalize the grammar for `when` and `match` expressions
Decide if want to keep the `when` and `match` as separate expressions
with their own keywordm  or unified them into a single `match` expression 
with overloaded syntax.

A) Separete

```
// Checking if an expression is true
when 
{ n == 0 => print("n is zero"},
{ else => print("not zero" }

// Checking if a variable has a type
match obj 
{ Option => print("The value of obj is `obj`") },
{ None  => print("There was no value") } 
```

B) Same 
```
// Checking if an expression is true
match 
{ n == 0 => print("n is zero"},
{ else => print("not zero" }

// Checking if a variable has a type
match obj 
{ Option => print("The value of obj is `obj`") },
{ None  => print("There was no value") } 
```
## Commas and semicolons

Currently, a boc can parse two expressions separated by a comma. This should be a special case when more than one
variable is assigned

```js
// Good
a, b = 1, 2
// Bad
a = 1, b = 2
```

The second should be written as

```js
a = 1;
b = 2
```

## Update grammar to support Variants

Variant example

```
Option {
  Some(value T),
  None()
}
```
## A variable is not an expression

It currently implements it, but it should be removed.

## Support `when` and `match` expression
Add support for pattern when expression with two forms:

### Form 1: Conditional Matching `when`
```
when
{ condition => action1; action2 },
{ condition => action1; action2 },
{ _ => default_action }
```
- No leading expression
- Cases must be boolean expressions
- `_` is the default case (catch-all)
- Multiple actions allowed per case, separated by semicolons
- Returns the value of the last expression in the matched case

### Form 2: Type Matching `match`
```
expr match
{ Type1 => action1; action2 },
{ Type2 => action1; action2 },
{ _ => default_action }
```
- Has leading expression to match against
- Cases must be type identifiers
- `_` is the default case (catch-all)
- Multiple actions allowed per case, separated by semicolons
- Returns the value of the last expression in the matched case

### Notes
- Match is an expression, not just a statement
- Returns the value of the last expression in the executed case block
- Comma separates cases
- Each case is enclosed in curly braces
- `=>` separates pattern from actions
- Last case may optionally have a trailing comma
- Actions can be expressions or statements

### Examples with Return Values
```
x : when { a > 0 =>
  print("positive");
  42
}  // returns 42

y : some_value match {
    Some => print("exists"); "found",  // returns "found"
    None => print("empty"); "not_found"  // returns "not_found"
}
```

## Work on assignments, multiple and single.

```
a = 1
```

Is an expression that whose value is the rhs.

What about `a, b,c = 1, 2, 3`?

Should it be the same as `a, b, c : {1, 2, 3}()`?

Is that a list in the AST? Or is a sugared form of  ?
```
a = 1
b = 2
c = 3
```

## Create an example of the generated Go code

Complete [generated_go_structures_sample.go](internal/testdata/generated_go_structures_sample.go) to include
asynchronous calls.

Include a recursive method call.

## Pass build options to the compiler

We can keep the generated source files or discard them.
At this moment we'll keep them for debugging purposes, but eventually we will need and option to discard them.
