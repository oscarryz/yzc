## How to handle duplicate definitions in different directories

If the source path is: `[src, lib]` 

The following file structure: 
```
   src/
        foo/
            bar.yz
    lib/
        foo/
            bar.yz
```
Would create the valid bocs `foo.bar` in the `src` directory and `foo.bar` in the `lib` directory.
Should this be allowed? Should the compiler throw an error? Should the compiler allow it and use the first one found?. 
Possible uses cases:

1. The user wants to augment the functionality of a library

In this case we would merge the definitions to create a single `foo.bar` boc. We would probably need to define this 
in the project configuration file e.g. 
```
src_path: [src, lib]
merge_src_path: true
```


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

