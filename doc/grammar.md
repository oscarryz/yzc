
- [Grammar BNF](./grammar.bnf)
- [Rule Trie for parsing diagram](https://www.figma.com/file/duG5ggrqoL8Hssl5rkdN8K/grammar?type=whiteboard&node-id=0%3A1&t=ooTKBNrtrJW1ToVY-1)

### Experimenting parsing with a Rule Trie.

For backtracking I'm going to experiment putting the grammar rules in a trie.

To check if a token sequence belongs to a given rule I can check the trie, if it is, I can now build the specific AST.

Consider the following variable declarations:
```
// var of type string
a1 String

// string variable init to 'a2'
a2 String = 'a2'

// string var inferred to 'a3'
a3 : 'a3'

// array var
a4 []

// typed var array
a5 [] String

// typed dictionary
a6 [ Int : String ]
```

They start with an identifier, and then they have different options to check before decide what kind of variable it is, some of them in turn have other nested options.

So, by creating a trie with the rules we can check if a given path belongs to this or that rule.

So, the rules for the examples above:

```
variable_definition ::= variable ":" (expression)+
                    |   variable type "=" expression

variable_declaration ::= variable type
                       | variable array_declaration
                       | variable dictionary_declaration
 ```

Could be encoded in a trie like this (var section, under statements) [grammar](https://www.figma.com/file/duG5ggrqoL8Hssl5rkdN8K/grammar?type=whiteboard&node-id=0%3A1&t=eYDMkISt5qONQSyD-1)