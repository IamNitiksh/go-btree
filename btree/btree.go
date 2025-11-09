// Package btree implements an in-memory B-Tree using generics.
package btree

import (
	"sort"
)

// BTree is the main tree structure.
// 'T' is the type of item you want to store.
type BTree[T any] struct {
	root   *Node[T]
	degree int // The 't' value: min items = t-1, max items = 2t-1
	length int // Number of items in the tree

	// The user-provided function to compare two items.
	less func(a, b T) bool
}

// Node is an internal node in the tree.
type Node[T any] struct {
	items    []T        // The items stored in this node
	children []*Node[T] // The child nodes
	tree     *BTree[T]  // A pointer back to the tree for config (degree/less)
}

// New creates a new B-Tree.
// 'degree' is 't', the minimum degree. Must be 2 or greater.
// 'less' is the function used to order your items.
func New[T any](degree int, less func(a, b T) bool) *BTree[T] {
	if degree < 2 {
		// A B-Tree must have a degree of at least 2
		// to be able to split (2t-1 items).
		panic("btree: degree must be 2 or greater")
	}
	if less == nil {
		panic("btree: less function must not be nil")
	}
	return &BTree[T]{
		degree: degree,
		less:   less,
	}
}

// newNode creates a new node within the tree's context.
func (t *BTree[T]) newNode(leaf bool) *Node[T] {
	n := &Node[T]{tree: t}
	// Pre-allocate slices with max capacity for efficiency
	n.items = make([]T, 0, 2*t.degree-1)
	if !leaf {
		n.children = make([]*Node[T], 0, 2*t.degree)
	}
	return n
}

// isLeaf checks if the node is a leaf (has no children).
func (n *Node[T]) isLeaf() bool {
	return len(n.children) == 0
}

// isFull checks if the node has the maximum number of items.
func (n *Node[T]) isFull() bool {
	// 2*t - 1
	return len(n.items) == (2*n.tree.degree - 1)
}

// minItems returns the minimum number of items for a non-root node.
func (t *BTree[T]) minItems() int {
	return t.degree - 1
}

// Len returns the total number of items in the tree.
func (t *BTree[T]) Len() int {
	return t.length
}

// ### 1. PUT (INSERT) AND SPLIT LOGIC ###

// Put inserts or updates an item in the tree.
// If an item with the same key already exists, it's overwritten.
// Returns 'true' if an existing item was replaced.
func (t *BTree[T]) Put(item T) (replaced bool) {
	// Case 1: The tree is empty.
	if t.root == nil {
		t.root = t.newNode(true) // Create a new leaf root
		t.root.items = append(t.root.items, item)
		t.length++
		return false
	}

	// Case 2: The root node is full.
	// We must split the root *before* inserting.
	if t.root.isFull() {
		oldRoot := t.root
		newRoot := t.newNode(false) // New root is an internal node
		t.root = newRoot
		newRoot.children = append(newRoot.children, oldRoot)

		// Call the split logic
		newRoot.splitChild(0)
	}

	// Case 3: Root is not full. Insert into the (possibly new) root.
	replaced = t.root.put(item)
	if !replaced {
		t.length++
	}
	return replaced
}

// put is the internal recursive function to add an item.
// It assumes the current node 'n' is *not* full.
func (n *Node[T]) put(item T) (replaced bool) {
	// 1. Find the position for the item
	i, found := n.find(item)

	// 2. If found, update it
	if found {
		n.items[i] = item // Update
		return true
	}

	// 3. If this is a leaf node, insert it
	if n.isLeaf() {
		// Insert 'item' at index 'i'
		n.items = append(n.items, *new(T)) // Add a zero-value placeholder
		copy(n.items[i+1:], n.items[i:])   // Shift items to the right
		n.items[i] = item                  // Insert the new item
		return false
	}

	// 4. Not a leaf: must insert into a child.
	// Check if the child we're descending into is full.
	if n.children[i].isFull() {
		// Split the child *before* descending
		n.splitChild(i)

		// The split might have promoted a key.
		// We need to re-check which side to insert into.
		// if item > median key (which is now n.items[i])
		if n.tree.less(n.items[i], item) {
			i++ // Move to the new right child
		}
	}

	// 5. Recursively call put on the correct (and now not-full) child.
	return n.children[i].put(item)
}

// find finds the index where 'item' should go in this node.
// It returns the index and whether an *exact* match was found.
func (n *Node[T]) find(item T) (index int, found bool) {
	// sort.Search finds the first index 'i' where the condition is true.
	// We are looking for the first item that is NOT less than 'item'.
	// This is the insertion point.
	i := sort.Search(len(n.items), func(i int) bool {
		// return n.items[i] >= item
		return !n.tree.less(n.items[i], item)
	})

	// Check if we found an exact match
	if i < len(n.items) && !n.tree.less(item, n.items[i]) {
		// Found an item that is "equal" (not less than, and not greater than)
		return i, true
	}

	// No match, 'i' is the insertion point
	return i, false
}

// splitChild is the core splitting logic.
// 'n' is the parent node, and 'i' is the index of its *full* child.
// This function splits the child n.children[i] into two nodes
// and promotes the median key up into 'n'.
func (n *Node[T]) splitChild(i int) {
	// 't' is the degree
	t := n.tree.degree

	// y is the full child node we need to split
	y := n.children[i]

	// z is the new node (y's new sibling, to its right)
	z := n.tree.newNode(y.isLeaf())

	// The median key is at index t-1 (since it's 0-indexed)
	// Example: degree=2 (2t-1 = 3 items). Indices [0, 1, 2]. Median is at 1.
	// Example: degree=3 (2t-1 = 5 items). Indices [0, 1, 2, 3, 4]. Median is at 2.
	medianKey := y.items[t-1]

	// 1. Copy the "right half" of y's items to z
	// (items from index t to end)
	z.items = append(z.items, y.items[t:]...)

	// 2. If y is not a leaf, copy the "right half" of its children to z
	// (children from index t to end)
	if !y.isLeaf() {
		z.children = append(z.children, y.children[t:]...)
	}

	// 3. Truncate y (it's now the "left half")
	// (items from 0 to t-2)
	y.items = y.items[:t-1]
	if !y.isLeaf() {
		// (children from 0 to t-1)
		y.children = y.children[:t]
	}

	// 4. Insert z as a new child of the parent 'n' at index i+1
	n.children = append(n.children, nil)     // Add space for new child pointer
	copy(n.children[i+2:], n.children[i+1:]) // Shift
	n.children[i+1] = z                      // Insert z

	// 5. Promote the median key into the parent 'n' at index i
	n.items = append(n.items, *new(T)) // Add space for new item
	copy(n.items[i+1:], n.items[i:])   // Shift
	n.items[i] = medianKey             // Insert median key
}

// ### 2. GET (SEARCH) LOGIC ###

// Get finds an item in the tree matching 'key'.
// Returns the item and 'true' if found.
func (t *BTree[T]) Get(key T) (item T, ok bool) {
	if t.root == nil {
		return *new(T), false
	}
	return t.root.get(key)
}

// get is the internal recursive search function.
func (n *Node[T]) get(key T) (item T, ok bool) {
	// 1. Find the first item >= key
	i, found := n.find(key)

	// 2. Found it
	if found {
		return n.items[i], true
	}

	// 3. Not found, and this is a leaf
	if n.isLeaf() {
		return *new(T), false
	}

	// 4. Not found, descend into child
	return n.children[i].get(key)
}

// ### 3. READ-ONLY METHODS (MIN, MAX, ASCEND) ###

// Min returns the smallest item in the tree.
// ok is false if the tree is empty.
func (t *BTree[T]) Min() (item T, ok bool) {
	if t.root == nil {
		return *new(T), false
	}
	return t.root.findMin(), true
}

// findMin finds the smallest item in the subtree rooted at 'n'.
func (n *Node[T]) findMin() T {
	current := n
	// Keep descending to the leftmost node
	for !current.isLeaf() {
		current = current.children[0]
	}
	// The min item is the first item in the leftmost leaf
	return current.items[0]
}

// Max returns the largest item in the tree.
// ok is false if the tree is empty.
func (t *BTree[T]) Max() (item T, ok bool) {
	if t.root == nil {
		return *new(T), false
	}
	return t.root.findMax(), true
}

// findMax finds the largest item in the subtree rooted at 'n'.
func (n *Node[T]) findMax() T {
	current := n
	// Keep descending to the rightmost node
	for !current.isLeaf() {
		current = current.children[len(current.children)-1]
	}
	// The max item is the last item in the rightmost leaf
	return current.items[len(current.items)-1]
}

// Ascend iterates over all items in the tree in ascending order (A -> Z).
// The iterator function is called for each item. If the iterator
// returns 'false', the iteration stops.
func (t *BTree[T]) Ascend(iterator func(item T) bool) {
	if t.root != nil {
		t.root.ascend(iterator)
	}
}

// ascend is the internal recursive function for in-order traversal.
// It returns 'false' if iteration was stopped, 'true' otherwise.
func (n *Node[T]) ascend(iterator func(item T) bool) bool {
	if n.isLeaf() {
		// Leaf: iterate over items
		for _, item := range n.items {
			if !iterator(item) {
				return false // Stop
			}
		}
		return true // Continue
	}

	// Internal node:
	// We must visit: child[i], then item[i], then child[i+1]
	for i := 0; i < len(n.items); i++ {
		// 1. Visit left child (child[i])
		if !n.children[i].ascend(iterator) {
			return false // Stop
		}
		// 2. Visit item (item[i])
		if !iterator(n.items[i]) {
			return false // Stop
		}
	}
	// 3. Visit last (rightmost) child (child[len(items)])
	return n.children[len(n.children)-1].ascend(iterator)
}

// Descend iterates over all items in the tree in descending order (Z -> A).
// The iterator function is called for each item. If the iterator
// returns 'false', the iteration stops.
func (t *BTree[T]) Descend(iterator func(item T) bool) {
	if t.root != nil {
		t.root.descend(iterator)
	}
}

// descend is the internal recursive function for reverse-order traversal.
// It returns 'false' if iteration was stopped, 'true' otherwise.
func (n *Node[T]) descend(iterator func(item T) bool) bool {
	if n.isLeaf() {
		// Leaf: iterate over items in reverse
		for i := len(n.items) - 1; i >= 0; i-- {
			if !iterator(n.items[i]) {
				return false // Stop
			}
		}
		return true // Continue
	}

	// Internal node:
	// 1. Visit last (rightmost) child
	if !n.children[len(n.children)-1].descend(iterator) {
		return false // Stop
	}
	// Iterate from right to left
	for i := len(n.items) - 1; i >= 0; i-- {
		// 2. Visit item
		if !iterator(n.items[i]) {
			return false // Stop
		}
		// 3. Visit left child
		if !n.children[i].descend(iterator) {
			return false // Stop
		}
	}
	return true
}

// ### 4. DELETE AND REBALANCE LOGIC ###

// Delete removes an item matching 'key' from the tree.
// Returns the deleted item and 'true' if it was found.
func (t *BTree[T]) Delete(key T) (item T, ok bool) {
	if t.root == nil {
		// Tree is empty
		return *new(T), false
	}

	// Call the internal delete function
	item, ok = t.root.delete(key)
	if !ok {
		// Item not found
		return *new(T), false
	}

	// If found and deleted, decrement length
	t.length--

	// If the root node is now empty and has a child,
	// make that child the new root. This shrinks the tree height.
	if len(t.root.items) == 0 && !t.root.isLeaf() {
		t.root = t.root.children[0]
	}

	// If the last item was deleted, set root to nil
	if t.length == 0 {
		t.root = nil
	}

	return item, true
}

// delete is the internal recursive delete function.
// It returns the deleted item and 'true' if found.
func (n *Node[T]) delete(key T) (item T, ok bool) {
	// 1. Find the key in the current node
	i, found := n.find(key)

	if n.isLeaf() {
		if found {
			// Case 1: Key is in a leaf node.
			// Remove it from the items slice.
			item = n.items[i]
			// Splice out the item at index i
			n.items = append(n.items[:i], n.items[i+1:]...)
			return item, true
		}
		// Case 1b: Key not in leaf, and it's a leaf.
		// Key doesn't exist.
		return *new(T), false
	}

	// --- Not a leaf node ---

	if found {
		// Case 2: Key is in an internal node.
		item = n.items[i]

		// We need to replace it with a successor or predecessor
		// to maintain the B-Tree property.
		if len(n.children[i].items) >= n.tree.degree {
			// 2a: Left child (n.children[i]) has enough items.
			// Find the predecessor (max item in left child's subtree),
			// replace our key with it, and recursively delete the predecessor.
			predecessor := n.children[i].findMax()
			n.items[i] = predecessor
			_, ok = n.children[i].delete(predecessor)
			return item, ok
		}

		if len(n.children[i+1].items) >= n.tree.degree {
			// 2b: Right child (n.children[i+1]) has enough items.
			// Find the successor (min item in right child's subtree),
			// replace our key with it, and recursively delete the successor.
			successor := n.children[i+1].findMin()
			n.items[i] = successor
			_, ok = n.children[i+1].delete(successor)
			return item, ok
		}

		// 2c: Neither child has enough items.
		// Merge the left and right children (and the key) into a
		// single node. Then, recursively delete the key from the
		// new merged node.
		n.mergeChildren(i)
		// We delete from the *merged* child (which is now at n.children[i])
		return n.children[i].delete(key)
	}

	// --- Key is not in this node, must be in a child ---

	// Case 3: Key is not in this internal node. Descend to the
	// correct child (n.children[i]).

	// *** This is the core rebalancing logic ***
	// Before we descend, we must ensure the child node we are
	// visiting is *not* underfull (has fewer than 'degree' items).
	if len(n.children[i].items) == n.tree.minItems() {
		// Child is underfull! We must rebalance.
		n.rebalance(i)
		// After rebalancing, the key might have moved.
		// We must re-find the path.
		// This is a simple (though slightly inefficient) way.
		// We just restart the delete from the current node.
		return n.delete(key)
	}

	// Child is "safe" (not underfull), so we can descend.
	return n.children[i].delete(key)
}

// rebalance ensures that child 'i' of node 'n' has at least 'degree' items
// (i.e., is not underfull) so that a deletion can proceed safely.
func (n *Node[T]) rebalance(i int) {
	// Try to steal from left sibling
	if i > 0 && len(n.children[i-1].items) > n.tree.minItems() {
		n.stealFromLeft(i)
	} else if i < len(n.children)-1 && len(n.children[i+1].items) > n.tree.minItems() {
		// Try to steal from right sibling
		n.stealFromRight(i)
	} else {
		// Both siblings are at the minimum. We must merge.
		if i == len(n.children)-1 {
			// No right sibling, or right sibling is underfull.
			// Merge with left sibling.
			n.mergeChildren(i - 1)
		} else {
			// Merge with right sibling.
			n.mergeChildren(i)
		}
	}
}

// stealFromLeft moves one key from left sibling (i-1) to child (i).
func (n *Node[T]) stealFromLeft(i int) {
	child := n.children[i]
	sibling := n.children[i-1]

	// 1. Move parent's key down to child
	// Insert at the beginning of child's items
	child.items = append(child.items, *new(T)) // Make room
	copy(child.items[1:], child.items)         // Shift
	child.items[0] = n.items[i-1]              // Copy from parent

	// 2. Move sibling's max key up to parent
	n.items[i-1] = sibling.items[len(sibling.items)-1]

	// 3. Remove max key from sibling
	sibling.items = sibling.items[:len(sibling.items)-1]

	// 4. Move child pointer if internal node
	if !sibling.isLeaf() {
		child.children = append(child.children, nil)
		copy(child.children[1:], child.children)
		child.children[0] = sibling.children[len(sibling.children)-1]
		sibling.children = sibling.children[:len(sibling.children)-1]
	}
}

// stealFromRight moves one key from right sibling (i+1) to child (i).
func (n *Node[T]) stealFromRight(i int) {
	child := n.children[i]
	sibling := n.children[i+1]

	// 1. Move parent's key down to child
	// Append to the end of child's items
	child.items = append(child.items, n.items[i])

	// 2. Move sibling's min key up to parent
	n.items[i] = sibling.items[0]

	// 3. Remove min key from sibling
	sibling.items = sibling.items[1:]

	// 4. Move child pointer if internal node
	if !sibling.isLeaf() {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}
}

// mergeChildren merges child 'i+1' into child 'i'.
// This is the reverse of splitChild.
func (n *Node[T]) mergeChildren(i int) {
	left := n.children[i]
	right := n.children[i+1]

	// 1. Pull parent's key down into left child
	left.items = append(left.items, n.items[i])

	// 2. Append all items from right child
	left.items = append(left.items, right.items...)

	// 3. Append all children from right child (if not leaf)
	left.children = append(left.children, right.children...)

	// 4. Remove the key from the parent
	n.items = append(n.items[:i], n.items[i+1:]...)

	// 5. Remove the right child pointer from parent
	n.children = append(n.children[:i+1], n.children[i+2:]...)

	// Note: We don't free 'right' node, we let GC handle it.
}

// ### 5. STRING (FOR DEBUGGING) ###
// (Renumbered from 4)
