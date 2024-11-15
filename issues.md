## Fix parser movement methods

- nextToken()
- peek()
- consume()
- expect()

They all interact in a weird way, fix it so they are easier to use to validate the grammar.

See for instance dictionary validation in the [parser.go](internal/parser.go) file.
## Add Dict type to the struct

Just like ArrayList has a type, Dict should have a type, it currently uses `[]` as value.

## Change Boc type to expression

## Handle Array of Dictionaries

## Consider simplifying  Boc > BlockBody

Currently, a Boc has a BlockBody which is a list of expressios and statments, but probably a boc could be a list of expressions and statements itself
Do we need both?

## Does a Boc need a name?

The name is always the variable to which it gets assigned, so probably is not needed

## Simplify Array type or types for that matter

It uses an expression which currently is the same expression (or almost) as the first element
probably there could be some other way

Probably well get there when we start doing declaration statements.

## Create an example of the generated Go code

Complete [generated_go_structures_sample.go](internal/testdata/generated_go_structures_sample.go) to include
asynchronous calls.

Include a recursive method call.

## Pass build options to the compiler

We can keep the generated source files or discard them.
At this moment we'll keep them for debugging purposes, but eventually we will need and option to discard them.