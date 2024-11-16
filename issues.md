## The initial split of filename into directories should create an AST of ShortDecl

For the input `a/b/c.yz`

The AST should be something like :

```
 ShortDcl { 
        key: a , 
        val: Boc{ 
                key: b,
                val: ShortDecl { 
                    key: c ,
                    val: {}     
        }
    }
}
```

## Remove the Boc.name property

Is not needed

## Add Dict type to the struct

Just like ArrayList has a type, Dict should have a type, it currently uses `[]` as value.

## Change Boc type to expression

## Handle Array of Dictionaries

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