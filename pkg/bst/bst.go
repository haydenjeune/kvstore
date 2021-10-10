package bst

// For now, we'll use a BST to store the sorted KV pairs. This is not ideal, because
// BSTs are not necessarily very well balanced. This leads to worst case time complexity
// of O(n) for Searches and Inserts. Ideally we'd use Red-Black trees or AVL trees.
// This would give O(log(n)) complexity.

// BinarySearchTree is a thin wrapper around BinaryNode to help with initialisation
type BinarySearchTree struct {
	root *BinaryNode
	size uint
}

func (t *BinarySearchTree) Insert(key string, value string) {
	t.size += 1
	if t.root == nil {
		t.root = &BinaryNode{key: key, value: value}
	} else {
		t.root.Insert(key, value)
	}
}

func (t *BinarySearchTree) Search(key string) (string, bool) {
	if t.root == nil {
		return "", false
	} else {
		return t.root.Search(key)
	}
}

func (t *BinarySearchTree) Size() uint {
	return t.size
}

type BinaryNode struct {
	key   string
	value string
	left  *BinaryNode
	right *BinaryNode
}

func (n *BinaryNode) Insert(key string, value string) {
	if key < n.key {
		if n.left != nil {
			n.left.Insert(key, value)
		} else {
			n.left = &BinaryNode{key: key, value: value}
		}
	} else if key > n.key {
		if n.right != nil {
			n.right.Insert(key, value)
		} else {
			n.right = &BinaryNode{key: key, value: value}
		}
	} else {
		n.value = value
	}
}

func (n *BinaryNode) Search(key string) (string, bool) {
	if key < n.key && n.left != nil {
		return n.left.Search(key)
	} else if key > n.key && n.right != nil {
		return n.right.Search(key)
	} else if key == n.key {
		return n.value, true
	} else {
		return "", false
	}
}

type InOrderTraversalIterator struct {
	q []*BinaryNode
}

func NewInOrderTraversalIterator(tree *BinarySearchTree) *InOrderTraversalIterator {
	i := &InOrderTraversalIterator{
		q: make([]*BinaryNode, 0, tree.Size()),
	}
	i.addNode(tree.root)
	return i
}

func (i *InOrderTraversalIterator) addNode(root *BinaryNode) {
	if root.left != nil {
		i.addNode(root.left)
	}
	if root != nil {
		i.q = append(i.q, root)
	}
	if root.right != nil {
		i.addNode(root.right)
	}
}

func (i *InOrderTraversalIterator) Next() *BinaryNode {
	if len(i.q) == 0 {
		return nil
	}

	next := i.q[0]
	i.q = i.q[1:]

	return next
}
