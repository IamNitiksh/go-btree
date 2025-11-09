package btree

// import (
// 	"fmt"
// 	"reflect"
// 	"testing"
// )

// // lessInt is a helper function for comparing integers.
// func lessInt(a, b int) bool {
// 	return a < b
// }

// // TestBTree_PutAndGet tests the basic functionality of inserting and retrieving items.
// func TestBTree_PutAndGet(t *testing.T) {
// 	// We'll use a 2-3-4 Tree (degree = 2)
// 	// Max items per node = 2*2 - 1 = 3
// 	tree := New[int](2, lessInt)

// 	// Insert some items
// 	itemsToInsert := []int{10, 20, 5, 6, 12, 30, 7, 17}
// 	for _, item := range itemsToInsert {
// 		tree.Put(item)
// 	}

// 	// Print the tree structure (useful for debugging)
// 	fmt.Println(tree.String())

// 	// Check if all inserted items are found
// 	for _, item := range itemsToInsert {
// 		val, ok := tree.Get(item)
// 		if !ok {
// 			t.Errorf("Get(%d): item not found, but should have been", item)
// 		}
// 		if val != item {
// 			t.Errorf("Get(%d): expected %d, got %d", item, item, val)
// 		}
// 	}

// 	// Check total length
// 	expectedLen := len(itemsToInsert)
// 	if tree.Len() != expectedLen {
// 		t.Errorf("Len(): expected %d, got %d", expectedLen, tree.Len())
// 	}

// 	// Check for a non-existent item
// 	val, ok := tree.Get(999)
// 	if ok {
// 		t.Errorf("Get(999): item found (%d), but should not have been", val)
// 	}

// 	// Test replacement
// 	replaced := tree.Put(10) // 10 already exists
// 	if !replaced {
// 		t.Errorf("Put(10): expected replaced=true, got false")
// 	}
// 	if tree.Len() != expectedLen {
// 		t.Errorf("Len() after replace: expected %d, got %d", expectedLen, tree.Len())
// 	}
// }

// // TestBTree_Split tests that the root split logic works.
// func TestBTree_Split(t *testing.T) {
// 	// degree=2 means max 3 items. 4th item should cause a split.
// 	tree := New[int](2, lessInt)

// 	tree.Put(10)
// 	tree.Put(20)
// 	tree.Put(30)

// 	// At this point, root node should be full: [10, 20, 30]
// 	if tree.root.isLeaf() == false {
// 		t.Fatal("Root should be a leaf before split")
// 	}
// 	if tree.root.isFull() == false {
// 		t.Fatal("Root should be full before split")
// 	}

// 	// This insert (the 4th item) should trigger a root split
// 	tree.Put(15)

// 	// After split, the new root should NOT be a leaf
// 	if tree.root.isLeaf() {
// 		t.Fatal("Root should not be a leaf after split")
// 	}
// 	// New root should have 1 item (the median key)
// 	if len(tree.root.items) != 1 {
// 		t.Fatalf("New root should have 1 item, got %d", len(tree.root.items))
// 	}
// 	// New root should have 2 children
// 	if len(tree.root.children) != 2 {
// 		t.Fatalf("New root should have 2 children, got %d", len(tree.root.children))
// 	}

// 	// The median key '20' should have been promoted.
// 	// (Note: insert order was 10,20,30, then 15. The split happens
// 	// on [10, 15, 20, 30], median is 15 or 20. Our logic promotes 20).
// 	// Let's re-run with a simpler order: 10, 20, 30.
// 	// Root is [10, 20, 30]. Insert 40.
// 	// Split [10, 20, 30, 40]. Median 20 or 30.
// 	// Let's trace our logic for 10, 20, 30. Root=[10, 20, 30].
// 	// Insert 15. Root full. Split.
// 	// OldRoot = [10, 20, 30]. NewRoot = []. children=[OldRoot]
// 	// splitChild(0) called. y=[10, 20, 30]. t=2. median=y[1]=20.
// 	// z (new right node) gets y[t:] = [30]
// 	// y (left node) becomes y[:t-1] = [10]
// 	// newRoot items = [20]
// 	// newRoot children = [y, z]
// 	// So newRoot is [20], children are [10] and [30].
// 	// Now, we insert 15.
// 	// n=root, item=15.
// 	// find(15) in [20] -> i=0, found=false.
// 	// child[0] ([10]) is not full.
// 	// call put(15) on child [10].
// 	// n=[10], item=15.
// 	// find(15) in [10] -> i=1, found=false.
// 	// it's a leaf. insert at i=1.
// 	// node [10] becomes [10, 15].
// 	//
// 	// Final tree:
// 	//   Node (items: 1): [20]
// 	//     Node (items: 2): [10, 15]
// 	//     Node (items: 1): [30]
// 	//
// 	// This all looks correct. Let's add the check.
// 	if tree.root.items[0] != 20 {
// 		t.Fatalf("Root item should be 20, got %d", tree.root.items[0])
// 	}
// 	fmt.Println("\n--- After Split ---")
// 	fmt.Println(tree.String())
// }

// // TestBTree_MinMaxAscend tests the Min, Max, and Ascend functions.
// func TestBTree_MinMaxAscend(t *testing.T) {
// 	tree := New[int](2, lessInt)

// 	// Test on an empty tree
// 	if _, ok := tree.Min(); ok {
// 		t.Error("Min() on empty tree should return ok=false")
// 	}
// 	if _, ok := tree.Max(); ok {
// 		t.Error("Max() on empty tree should return ok=false")
// 	}

// 	// Insert items
// 	// Use a permutation to create a non-trivial tree
// 	items := []int{5, 9, 1, 0, 8, 3, 7, 2, 4, 6}
// 	for _, item := range items {
// 		tree.Put(item)
// 	}

// 	fmt.Println("\n--- TestMinMaxAscend Tree ---")
// 	fmt.Println(tree.String())

// 	// Test Min
// 	min, ok := tree.Min()
// 	if !ok || min != 0 {
// 		t.Errorf("Min(): expected 0, got %d (ok=%v)", min, ok)
// 	}

// 	// Test Max
// 	max, ok := tree.Max()
// 	if !ok || max != 9 {
// 		t.Errorf("Max(): expected 9, got %d (ok=%v)", max, ok)
// 	}

// 	// Test Ascend (full iteration)
// 	var collected []int
// 	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
// 	tree.Ascend(func(item int) bool {
// 		collected = append(collected, item)
// 		return true // Continue iterating
// 	})

// 	if !reflect.DeepEqual(collected, expected) {
// 		t.Errorf("Ascend():\nExpected: %v\nGot:      %v", expected, collected)
// 	}

// 	// Test Ascend (early stop)
// 	collected = nil
// 	expectedStop := []int{0, 1, 2, 3, 4, 5}
// 	tree.Ascend(func(item int) bool {
// 		collected = append(collected, item)
// 		return item != 5 // Stop when we hit 5
// 	})

// 	if !reflect.DeepEqual(collected, expectedStop) {
// 		t.Errorf("Ascend (early stop):\nExpected: %v\nGot:      %v", expectedStop, collected)
// 	}
// }

// // TestBTree_Descend tests the Descend function.
// func TestBTree_Descend(t *testing.T) {
// 	tree := New[int](2, lessInt)
// 	items := []int{5, 9, 1, 0, 8, 3, 7, 2, 4, 6}
// 	for _, item := range items {
// 		tree.Put(item)
// 	}

// 	// Test Descend (full iteration)
// 	var collected []int
// 	expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
// 	tree.Descend(func(item int) bool {
// 		collected = append(collected, item)
// 		return true // Continue iterating
// 	})

// 	if !reflect.DeepEqual(collected, expected) {
// 		t.Errorf("Descend():\nExpected: %v\nGot:      %v", expected, collected)
// 	}

// 	// Test Descend (early stop)
// 	collected = nil
// 	expectedStop := []int{9, 8, 7, 6, 5}
// 	tree.Descend(func(item int) bool {
// 		collected = append(collected, item)
// 		return item != 5 // Stop when we hit 5
// 	})

// 	if !reflect.DeepEqual(collected, expectedStop) {
// 		t.Errorf("Descend (early stop):\nExpected: %v\nGot:      %v", expectedStop, collected)
// 	}
// }

// // TestBTree_Delete basic cases.
// func TestBTree_Delete(t *testing.T) {
// 	// We'll create a tree with known structure
// 	// degree=2 (2-3-4 Tree). min items = 1.
// 	tree := New[int](2, lessInt)
// 	items := []int{10, 20, 30, 40, 50, 60, 70, 80}
// 	for _, item := range items {
// 		tree.Put(item)
// 	}
// 	// Tree should look something like:
// 	//         [30, 60]
// 	//    [10, 20] [40, 50] [70, 80]

// 	fmt.Println("\n--- TestDelete Initial Tree ---")
// 	fmt.Println(tree.String())

// 	// Case 1: Delete from a leaf (easy)
// 	// Delete 80. Leaf [70, 80] becomes [70]. No rebalance needed.
// 	deleted, ok := tree.Delete(80)
// 	if !ok || deleted != 80 {
// 		t.Fatalf("Delete(80): expected %d, got %d (ok=%v)", 80, deleted, ok)
// 	}
// 	if _, ok := tree.Get(80); ok {
// 		t.Fatal("Get(80): item still exists after delete")
// 	}
// 	if tree.Len() != len(items)-1 {
// 		t.Fatalf("Len(): expected %d, got %d", len(items)-1, tree.Len())
// 	}
// 	fmt.Println("\n--- After deleting 80 ---")
// 	fmt.Println(tree.String())

// 	// Case 2: Delete from an internal node (successor replacement)
// 	// Delete 60. Replaced by successor 70.
// 	// Leaf [70] becomes []. This is underfull!
// 	// This will trigger rebalancing.
// 	//
// 	// Let's trace:
// 	// 1. delete(60) from root [30, 60]
// 	// 2. Found. Right child [70] has enough items (1 >= 1)? No, minItems=1.
// 	//    Let's check the logic.
// 	//    degree=2, minItems=1.
// 	//    len(n.children[i+1].items) >= n.tree.degree
// 	//    len([70]) >= 2 ? No.
// 	//    len(n.children[i].items) >= n.tree.degree
// 	//    len([40, 50]) >= 2 ? Yes.
// 	// 3. Predecessor replacement. Find max in [40, 50] -> 50.
// 	// 4. Root becomes [30, 50].
// 	// 5. Recursively delete 50 from child [40, 50].
// 	// 6. delete(50) from leaf [40, 50].
// 	// 7. Found. Splice out. Leaf becomes [40]. This is legal (minItems=1).
// 	deleted, ok = tree.Delete(60)
// 	if !ok || deleted != 60 {
// 		t.Fatalf("Delete(60): expected %d, got %d (ok=%v)", 60, deleted, ok)
// 	}
// 	if _, ok := tree.Get(60); ok {
// 		t.Fatal("Get(60): item still exists after delete")
// 	}
// 	if val, ok := tree.Get(50); !ok || val != 50 {
// 		t.Fatal("Get(50): successor 50 was not promoted or found")
// 	}
// 	fmt.Println("\n--- After deleting 60 (replaced by 50) ---")
// 	fmt.Println(tree.String())
// 	// Expected:
// 	//     [30, 50]
// 	// [10, 20] [40] [70]

// 	// Case 3: Trigger a merge
// 	// Delete 70.
// 	// 1. delete(70) from root [30, 50].
// 	// 2. Not found. Descend to child[2] ([70]).
// 	// 3. Before descending, check if child[2] ([70]) is underfull.
// 	//    len=1. minItems=1. It is *at* minimum.
// 	//    Our rebalance logic: `len(n.children[i].items) == n.tree.minItems()`
// 	//    This is true!
// 	// 4. rebalance(2) called on root.
// 	// 5. Try steal from left (i=2): child[1] ([40]). len=1. Fails (len > minItems).
// 	// 6. Try steal from right (i=2): No right child.
// 	// 7. Must merge. `i == len(n.children)-1`. True.
// 	// 8. mergeChildren(i-1) -> mergeChildren(1)
// 	// 9. Merge child[2] ([70]) into child[1] ([40])
// 	//    left=[40], right=[70]
// 	//    parent key n.items[1] is 50.
// 	//    left.items = [40, 50, 70]
// 	//    root.items = [30]
// 	//    root.children = [[10, 20], [40, 50, 70]]
// 	// 10. Restart delete(70) from root [30].
// 	// 11. Not found. Descend to child[1] ([40, 50, 70]).
// 	// 12. child[1] is not underfull.
// 	// 13. call delete(70) on [40, 50, 70].
// 	// 14. Found at i=2. Leaf. Splice out.
// 	// 15. Node becomes [40, 50].
// 	deleted, ok = tree.Delete(70)
// 	if !ok || deleted != 70 {
// 		t.Fatalf("Delete(70): expected %d, got %d (ok=%v)", 70, deleted, ok)
// 	}
// 	fmt.Println("\n--- After deleting 70 (triggered merge) ---")
// 	fmt.Println(tree.String())

// 	// Check final state
// 	expectedItems := []int{10, 20, 30, 40, 50}
// 	if tree.Len() != len(expectedItems) {
// 		t.Fatalf("Final Len(): expected %d, got %d", len(expectedItems), tree.Len())
// 	}
// 	collected := []int{}
// 	tree.Ascend(func(item int) bool {
// 		collected = append(collected, item)
// 		return true
// 	})
// 	if !reflect.DeepEqual(collected, expectedItems) {
// 		t.Errorf("Final Ascend():\nExpected: %v\nGot:      %v", expectedItems, collected)
// 	}

// 	// Case 4: Delete all remaining items
// 	for _, item := range expectedItems {
// 		tree.Delete(item)
// 	}
// 	if tree.Len() != 0 {
// 		t.Fatalf("Len() after deleting all: expected 0, got %d", tree.Len())
// 	}
// 	if tree.root != nil {
// 		t.Fatal("tree.root should be nil after deleting all items")
// 	}
// }
