## Need to modify the concept of `program`

Currently, in the grammar a program has a blockbody and the blockbody has a list of expressions of statements. 

What I need is to define that a program is actually a block of code (like everything else)

So, an empty file `hi.yz` would define the `boc` : `hi #() {}` that is a block named `hi` with and empty body.

If a file has a main method then a full binary program would be the result, otherwise the compilation would only 
cache the library? 

I need to modify the grammar to reflect this.


## `=` as operator or as part of an identifier

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

