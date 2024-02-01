package main

const (
	RED = iota
	BLACK
)

func NewNode(score int, value string) *Node {
	return &Node{
		Color: RED,
		Value: value,
		Score: score,
		Parent: nil,
		Left: nil,
		Right: nil,
	}
}

func (t *RBTree) leftRotate(node *Node) {
	if node.Right == nil {
		return
	}
	rnode := node.Right
	node.Right = rnode.Left
	if rnode.Left != nil {
		rnode.Left.Parent = node
	}
	rnode.Parent = node.Parent

	// Update parent node links
	if node.Parent == nil {
		t.Root = rnode
	} else if node == node.Parent.Left {
		node.Parent.Left = rnode
	} else {
		node.Parent.Right = rnode
	}

	rnode.Left = node
	node.Parent = rnode
}

func (t *RBTree) rightRotate(node *Node) {
	if node.Left == nil {
		return
	}	
	lnode := node.Left
	node.Left = lnode.Right
	if lnode.Right != nil {
		lnode.Right.Parent = node
	}
	lnode.Parent = node.Parent

	if node.Parent == nil {
		t.Root = lnode
	} else if node == node.Parent.Left {
		node.Parent.Left = lnode
	} else {
		node.Parent.Right = lnode
	}

	lnode.Right = node
	node.Parent = lnode
}

func (t *RBTree) insert(score int, value string) {
	node := NewNode(score, value)

	var temp *Node
	root := t.Root

	for root != nil {
		temp = root
		if node.Score < root.Score {
			root = root.Left
		} else {
			root = root.Right
		}
	}
	node.Parent = temp

	if temp == nil {
		t.Root = node
	} else if node.Score < temp.Score {
		temp.Left = node
	} else {
		temp.Right = node
	}
}

func (t *RBTree) fixInsert(node *Node) {
	for node.Parent.Color == RED {
		if node.Parent == node.Parent.Parent.Left {
			uncle := node.Parent.Parent.Right
			if uncle.Color == RED {
				// Both uncle and parent are red, grandparent must be black
				// Repaint parent and uncle to black and grandparent to red
				node.Parent.Color = BLACK
				uncle.Color = BLACK
				node.Parent.Parent.Color = RED
				node = node.Parent.Parent
			} else {
				if node == node.Parent.Right {
					node = node.Parent
					t.leftRotate(node)
				}
				node.Parent.Color = BLACK
				node.Parent.Parent.Color = RED
				t.rightRotate(node.Parent.Parent)
			}
		} else {
			uncle := node.Parent.Parent.Left
			if uncle.Color == RED {
				node.Parent.Color = BLACK
				uncle.Color = BLACK
				node.Parent.Parent.Color = RED
				node = node.Parent
			} else {
				if node == node.Parent.Left {
					node = node.Parent
					t.rightRotate(node)
				}
				node.Parent.Color = BLACK
				node.Parent.Parent.Color = RED
				t.leftRotate(node.Parent.Parent)
			}
		}
	}
	t.Root.Color = BLACK
}
