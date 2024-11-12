## Fix dictionary parsing
## Think if empty lines are considered empty expressions
The following would result two expressions: BasicLit: 1 and EmptyExpr 
```
1

```
## Fix parser movement methods

- nextToken()
- peek()
- consume()
- expect()

They all interact in a weird way, fix it so they are easier to use to validate the grammar. 

See for instance dictionary validation in the [parser.go](internal/parser.go) file.

## Create an example of the generated Go code

Complete [generated_go_structures_sample.go](internal/testdata/generated_go_structures_sample.go) to include asynchronous calls. 

Include a recursive method call.

## Pass build options to the compiler

We can keep the generated source files or discard them. 
At this moment we'll keep them for debugging purposes, but eventually we will need and option to discard them.