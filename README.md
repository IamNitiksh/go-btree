# go-btree

*A B-Tree implementation for Go that's fast, generic, and not at all tree-hugging.*

***

## 🌲 What is this?
**go-btree** is an in-memory, generic B-Tree library for Go programs. It lets you store sortable data, do fast lookups, perform range queries, and feel smarter than using plain maps. If you ever wished slices sorted themselves or wanted an efficient ordered map — congrats, you’re on the right branch.

***

## 🚀 Features
- **Generics!** Store ints, structs, or even your deepest secrets (okay, probably just structs).
- **Fast insert, search, and delete.** No drama, always balanced.
- **Order guaranteed.** Traverse ascending or descending, because sometimes you want your data backwards.
- **Minimal dependencies.** Just Go, no fertilizer required.

***

## 🦸‍♂️ Example Usage
```go
package main

import (
    "fmt"
    "github.com/IamNitiksh/go-btree/btree"
)

type User struct {
    ID   int
    Name string
}

func userLess(a, b User) bool {
    return a.ID < b.ID
}

func main() {
    tree := btree.New[User](3, userLess)
    tree.Put(User{ID: 42, Name: "Alice"})
    tree.Put(User{ID: 7, Name: "Bob"})
    tree.Put(User{ID: 1337, Name: "Charlie"})
    fmt.Println("All Users in B-Tree:")
    tree.Ascend(func(u User) bool {
        fmt.Printf("  [%d] %s\n", u.ID, u.Name)
        return true
    })
}
```

***

## 📖 API Highlights
- `Put(item)` — Insert or update an item. Tree will split/merge as needed while staying balanced (unlike my diet).
- `Get(key)` — Retrieve an item. Finds your stuff so you don’t have to.
- `Delete(key)` — Remove items, like clearing browser history.
- `Ascend(fn)` — Iterate items in order. Sorted. Satisfying.
- `Descend(fn)` — Reverse-iterated. For when you just wanna go downhill.

***

## 📝 How to install?
```bash
go get github.com/IamNitiksh/go-btree/btree
```
Then import:
```go
import "github.com/IamNitiksh/go-btree/btree"
```

***

## 🤔 Why use B-Trees?
Because balanced trees are fast, efficient, and *orderly*. They power databases, filesystems, and now... your Go programs, too.

***

## 🦑 Contributing
Want to add features, report bugs, or just plant some new jokes in the docs?
- Open an issue or pull request on GitHub.
- Star the repo! (Stars help the ecosystem and my self-esteem)

***

## 📄 License
MIT — use it, share it, fork it, eat a snack, enjoy fun code.

***

## 🥇 Fun facts
- Real trees grow slow, but B-Trees balance instantly.
- No trees were harmed building this Go module.
- This README is longer than most of my Go functions.

***

Happy coding!
