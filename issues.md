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