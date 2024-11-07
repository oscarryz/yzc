
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
1. Force to call it with a dot `.` before the `=` token. (easy and restrictive but non-ambiguous but creates an special case)

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
 == b // compiles to a.==(b) no error
// option 2
a==b // unknown symbol a==b
// option 3 (would need to have `equals` in the `Point` type)
a.equals(b) // true
// option 4
a.==b 


With option 1 we would allow: 
some=thing // [ID]
some=thing = 1 // [ID, ASSIGN, NUMBER]

It would't be a good idea to create such identifier but the compiler would allow it.
```
Option 1 could be the best option.

