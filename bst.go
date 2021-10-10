package main

// For now, we'll use a BST to store the sorted KV pairs. This is not ideal, because
// BSTs are not necessarily very well balanced. Ideally we'd use Red-Black trees or AVL
// trees.

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
