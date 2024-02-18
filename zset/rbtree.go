package zset

import "sync"

const (
	RED = iota
	BLACK
)

type Node struct {
	Color  int
	Value  string
	Parent *Node
	Left   *Node
	Right  *Node
}

type RBTree struct {
	Mutex *sync.RWMutex
	Root  *Node
	Count uint
}

func NewNode(value string) *Node {
	return &Node{
		Color:  RED,
		Value:  value,
		Parent: nil,
		Left:   nil,
		Right:  nil,
	}
}

func NewRBTree() *RBTree {
	return &RBTree{Root: nil, Count: 0}
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

func (t *RBTree) search(value string) (*Node, bool) {
	temp := t.Root
	for temp != nil {
		if temp.Value == value {
			return temp, true
		}
		if temp.Value > value {
			temp = temp.Left
		} else {
			temp = temp.Right
		}
	}
	return nil, false
}

func (t *RBTree) fixInsert(node *Node) {
	for node.Parent != nil && node.Parent.Color == RED {
		if node.Parent == node.Parent.Parent.Left {
			uncle := node.Parent.Parent.Right
			if uncle != nil && uncle.Color == RED {
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
			if uncle != nil && uncle.Color == RED {
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

func (t *RBTree) insert(value string) {
	if _, ok := t.search(value); ok {
		return
	}
	node := NewNode(value)

	var temp *Node
	root := t.Root

	for root != nil {
		temp = root
		if node.Value < root.Value {
			root = root.Left
		} else {
			root = root.Right
		}
	}
	node.Parent = temp

	if temp == nil {
		t.Root = node
	} else if node.Value < temp.Value {
		temp.Left = node
	} else {
		temp.Right = node
	}
	t.Count++
	t.fixInsert(node)
}

func (t *RBTree) min(node *Node) *Node {
	if node == nil {
		return nil
	}

	for node.Left != nil {
		node = node.Left
	}
	return node
}

func (t *RBTree) max(node *Node) *Node {
	if node == nil {
		return nil
	}

	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (t *RBTree) successor(node *Node) *Node {
	if node == nil {
		return nil
	}

	if node.Right != nil {
		return t.min(node.Right)
	}

	successor := node.Parent
	for successor != nil && node == successor.Right {
		node = successor
		successor = successor.Parent
	}
	return successor
}

func (t *RBTree) fixDelete(node *Node) {
	for node != t.Root && node.Color == BLACK {
		if node == node.Parent.Left {
			cousin := node.Parent.Right
			if cousin.Color == RED {
				cousin.Color = BLACK
				node.Parent.Color = RED
				t.leftRotate(node.Parent)
				cousin = node.Parent.Right
			}
			if cousin.Left.Color == BLACK && cousin.Right.Color == BLACK {
				cousin.Color = RED
				node = node.Parent
			} else {
				if cousin.Right.Color == BLACK {
					cousin.Left.Color = BLACK
					cousin.Color = RED
					t.rightRotate(cousin)
					cousin = node.Parent.Right
				}
				cousin.Color = node.Parent.Color
				node.Parent.Color = BLACK
				cousin.Right.Color = BLACK
				t.leftRotate(node.Parent)
				node = t.Root
			}
		} else {
			cousin := node.Parent.Left
			if cousin.Color == RED {
				cousin.Color = BLACK
				node.Parent.Color = RED
				t.rightRotate(node.Parent)
				cousin = node.Parent.Left
			}
			if cousin.Left.Color == BLACK && cousin.Right.Color == BLACK {
				cousin.Color = RED
				node = node.Parent
			} else {
				if cousin.Left.Color == BLACK {
					cousin.Right.Color = BLACK
					cousin.Color = RED
					t.leftRotate(cousin)
					cousin = node.Parent.Left
				}
				cousin.Color = node.Parent.Color
				node.Parent.Color = BLACK
				node.Left.Color = BLACK
				t.rightRotate(node.Parent)
				node = t.Root
			}
		}
	}
	node.Color = BLACK
}

func (t *RBTree) delete(value string) {
	node, ok := t.search(value)
	if !ok {
		return
	}
	
	var temp *Node
	if node.Left == nil || node.Right == nil {
		temp = node
	} else {
		temp = t.successor(node)
	}
	var tchild *Node
	if temp.Left != nil {
		tchild = temp.Left
	} else {
		tchild = temp.Right
	}

	tchild.Parent = temp.Parent
	if temp.Parent == nil {
		t.Root = tchild
	} else if temp == temp.Parent.Left {
		temp.Parent.Left = tchild
	} else {
		temp.Parent.Right = tchild
	}

	if temp != node {
		node.Value = temp.Value
	}
	if temp.Color == BLACK {
		t.fixDelete(tchild)
	}
	t.Count--
}

func inorder(node *Node, elements chan string) {
	if node == nil {
		return
	}
	inorder(node.Left, elements)
	elements <- node.Value
	inorder(node.Right, elements)
}

func (t *RBTree) members() []string {
	var result []string
	elements := make(chan string, t.Count)

	go func() {
		defer close(elements)
		inorder(t.Root, elements)
	}()

	for element := range elements {
		result = append(result, element)
	}
	return result
}
