package internal

type ruleToken = tokenType

const (
	BODY = ILLEGAL + iota
	EXPR
	STMT
	PARENTHESIS_EXPR
	INTEGER_EXPR
	DECIMAL_EXPR
	BLOCK_EXPR
	EMPTY_BLOCK_EXPR
)

type Trie struct {
	tt        tokenType
	children  []*Trie
	isComplex bool // if is body, expression or statement
	isLeaf    bool // if is the last element of the trie
}

func find(tries []*Trie, tt tokenType) (*Trie, bool) {
	for _, trie := range tries {
		if trie.tt == tt {
			return trie, true
		}
	}
	return nil, false
}
func filterComplex(tries []*Trie) ([]*Trie, bool) {
	var r []*Trie
	for _, trie := range tries {
		if trie.isComplex {
			r = append(r, trie)
		}
	}
	return r, len(r) > 0
}
func buildTrie(rules [][]ruleToken) *Trie {

	root := new(Trie)
	root.children = []*Trie{}
	for _, rule := range rules {
		node := root
		for _, tt := range rule {
			var nt *Trie
			var ok bool
			if nt, ok = find(node.children, tt); !ok {
				nt = new(Trie)
				nt.tt = tt
				nt.children = []*Trie{}
				nt.isComplex = tt == BODY || tt == EXPR || tt == STMT
				node.children = append(node.children, nt)
			}
			node = nt
		}
		node.isLeaf = true
	}
	return root
}

func expressionTrie() *Trie {
	return buildTrie([][]ruleToken{
		{
			LPAREN, EXPR, RPAREN, PARENTHESIS_EXPR,
		}, {
			INTEGER, INTEGER_EXPR,
		}, {
			DECIMAL, DECIMAL_EXPR,
		}, {
			LBRACE, BODY, RBRACE, BLOCK_EXPR,
		}, {
			LBRACE, RBRACE, EMPTY_BLOCK_EXPR,
		},
	})
}
