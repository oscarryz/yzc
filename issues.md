## Types. Simplify Array type or types for that matter

- ~~It uses an expression which currently is the same expression (or almost) as the first element
probably there could be some other way~~
- ~~Probably well get there when we start doing declaration statements.~~
- ~~Arrays~~
- Dictionaries
- Bocs

## Change Boc type to expression
But first simplify types. Do I need a type struct? (kind)

## Add Dict type to the struct

Just like ArrayList has a type, Dict should have a type, it currently uses `[]` as value.
First handle types. 

## Create an example of the generated Go code

Complete [generated_go_structures_sample.go](internal/testdata/generated_go_structures_sample.go) to include
asynchronous calls.

Include a recursive method call.

## Pass build options to the compiler

We can keep the generated source files or discard them.
At this moment we'll keep them for debugging purposes, but eventually we will need and option to discard them.