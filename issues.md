## Add tests to handle lack of comma between expressions.

Currently, a file with two literals 

```
1 "Hello"
```
Would parse as a single literal because at the absence of comma, 
the parser finishes parsing the block of code and returns the result.

```js
{ 1 2 } 
```
Should be a syntax error because the parser expects a comma between the literals.

```js
{ 1, } 
```
Is probably fine

The problem is a file creates an implicit block so we cannot check for
closing brackets to determine the end of the block.

```
1 , a Int
```



## Create an example of the generated Go code

Complete [generated_go_structures_sample.go](internal/testdata/generated_go_structures_sample.go) to include asynchronous calls. 

Include a recursive method call.
## Pass build options to the compiler

We can keep the generated source files or discard them. 
At this moment we'll keep them for debugging purposes, but eventually we will need and option to discard them.

## Think about `=` as operator or as part of an identifier

The `=` token be the assigment operator and also can be part of an identifier.

To disambiguate we can:
1. Force to use empty space around the `=` token to make it the assignment operator (easy and ugly)
1. Perform and analysis of the surrounding tokens to determine if the `=` token is an assignment operator or part of an identifier. (hard and confusing)
1. Prohibit the use of `=` as part of an identifier. (easy and restrictive)
1. Force to call it with a dot `.` before the `=` token. (easy and restrictive but non-ambiguous)

e.g. 
```js
a = 1 // [ID, ASSIGN, NUMBER]
//vs
a=1 // [ID]
```

Probable use case, override `==`

```
Point : {
    x: 0
    y: 0
    == : {
        other Point
        x == other x && y == other y
    }
}
a: Point(1, 2)
b: Point(1, 2)
// option 1
a == b // true
// option 2
a==b // unknown symbol a==b
// option 3 (would need to have `equals` in the `Point` type)
a.equals(b) // true
// option 4
a.==b 
```
Option 1 could be the best option.

