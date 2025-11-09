package main

import (
	"fmt"

	"github.com/IamNitiksh/go-btree/btree"
)

// User is a demo struct stored in the BTree
type User struct {
	ID   int
	Name string
}

// userLess compares users by ID
func userLess(a, b User) bool {
	return a.ID < b.ID
}

func main() {
	tree := btree.New[User](3, userLess)
	tree.Put(User{ID: 1, Name: "Alice"})
	tree.Put(User{ID: 2, Name: "Bob"})
	fmt.Println("Users in BTree:")
	tree.Ascend(func(u User) bool {
		fmt.Printf("[%d] %s\n", u.ID, u.Name)
		return true
	})
}
