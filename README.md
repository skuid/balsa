# balsa

This is mostly intended to be used as a library for parsing, manipulating, and serializing a condition logic string.

For example, given a condition logic string like so `1 OR 2 AND (3 AND 4)`, this will allow you to remove any leaves by value and still get a valid condition logic string.

# Usage

## Parse

To parse an existing condition logic string, call `Parse(string) error`

```go
import {
	"github.com/skuid/balsa/parse"
}

logic := "1 OR 2 AND (3 OR 4)"

tree, err := parse.Parse(logic)

if err != nil {
	fmt.Print("Error: %v", err)
}

fmt.Printf("Built a binary expression tree that looks like this: %v", tree)
```

## Node.Eval

This will take a tree and write it to any `io.Writer`

```go
var b strings.Builder

err := tree.Eval(&b)

if err != nil {
	fmt.Print("Error: %v", err)
}

fmt.Printf("Serialized the tree into condition logic: %s", b.String())
```

## Node.Remove

This can be called multiple times to remove leafs from the tree by value. It will hoist any remaining expressions where needed and ignore any leafs it does not contain.

```go
logic := "1 OR (5 AND (1 OR 1)) AND (1 AND 2 OR (56 AND 1)) OR 4"

tree, _ := parse.Parse(logic)

n := tree.Remove(1)
n = tree.Remove(8)

var b strings.Builder
n.Eval(&b)
fmt.Println(b.String())
tree.Remove(8)

var b strings.Builder
tree.Eval(&b)
fmt.Println(b.String())
// 5 AND (2 OR 56 OR 4)
```

## Sequence

This will take a node and re-sequence all of the leaves based on ordinal positioning, starting at 0. For example, if you have a tree that is `5 AND 3`, this will re-sequence it as `1 AND 0`.

```go
logic := "5 AND 3"

tree, _ := parse.Parse(logic)

tree = parse.Sequence(tree)

var b strings.Builder
tree.Eval(&b)
fmt.Println(b.String())

// 2 AND 1
```

## Node.Index

This will re-index all of the leaves starting from the new starting index. If you have a tree that is `5 AND 3` and you want to start the indexing at 9, it will yield `10 AND 9`. This is typically used after the tree is sequenced.

```go
logic := "5 AND 3"

tree, _ := parse.Parse(logic)

tree = parse.Sequence(tree)
tree = tree.Index(9)

var b strings.Builder
tree.Eval(&b)
fmt.Println(b.String())
// 10 AND 9 
```

# Errors

`Parse` will throw `ParseError` errors mostly. These errors contain:
- **Position**: The location in the logic string the error occurred
- **Logic**: The logic string that the parser was tryign to parse
- **Reason**: The reason it failed to parse the logic at that location

`Eval` will throw `SerializeError` errors. They contain:
- **Op**: The operation that failed to serialize. Leaf values are unlikely to fail
- **Reason**: The reason it could not serialize that operation, which is typically due to missing nodes.
