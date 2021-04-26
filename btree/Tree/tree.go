package Tree

import (
	"fmt"
	"strconv"
)

type bnode struct {
	data                 int
	prev, next           *bnode
	leftchid, rightchild *page
}

func newBnode(data int) *bnode {
	return &bnode{data, nil, nil, nil, nil}
}

type page struct {
	size       int
	start, end *bnode
}

func newpage() *page {
	return &page{0, nil, nil}
}
func startPage(data int) *page {
	node := newBnode(data)
	return &page{1, node, node}
}

type btree struct {
	root *page
}

func newBTree() *btree {
	return &btree{nil}
}
func startBTree(data int) *btree {
	return &btree{startPage(data)}
}
func (tree *btree) Insert(data int) {
	newBnode := newBnode(data)
	if tree.root == nil {
		tree.root = newpage()
		tree.root.start = newBnode
		tree.root.end = newBnode
		tree.root.size = 1
	} else {
		if tree.root.SearchToInsert(newBnode) == nil {
			tree.root = tree.root.Insert(newBnode)
		}
	}
}
func (root *page) Insert(newBNode *bnode) *page {
	toTop := newpage()
	if root != nil {
		aux := root.start
		if aux != nil {
			if aux.data > newBNode.data {
				if aux.leftchid != nil {
					toTop = aux.leftchid.Insert(newBNode)
					if toTop.start.leftchid != nil && toTop.size == 1 {
						newBNodeInTop := toTop.start
						newBNodeInTop.next = aux
						aux.prev = newBNodeInTop
						aux.leftchid = newBNodeInTop.rightchild
						root.start = newBNodeInTop
						root.size++
					} else if toTop.start.rightchild != nil && toTop.size == 1 {
						newBNodeInTop := toTop.start
						newBNodeInTop.next = aux
						aux.prev = newBNodeInTop
						aux.leftchid = newBNodeInTop.rightchild
						root.start = newBNodeInTop
						root.size++
					}
				} else {
					newBNode.next = aux
					aux.prev = newBNode
					root.start = newBNode
					root.size++
				}
			} else if aux.data < newBNode.data {
				for aux != nil && aux.next != nil {
					if aux.data < newBNode.data && aux.next.data > newBNode.data {
						if aux.rightchild != nil {
							toTop = aux.rightchild.Insert(newBNode)
							if toTop.start.leftchid != nil && toTop.size == 1 {
								newBNodeInTop := toTop.start
								newBNodeInTop.next = aux.next
								newBNodeInTop.prev = aux
								aux.next.prev = newBNodeInTop
								aux.next.leftchid = newBNodeInTop.rightchild
								aux.next = newBNodeInTop
								aux.rightchild = newBNodeInTop.leftchid
								root.size++
							} else if toTop.start.rightchild != nil && toTop.size == 1 {
								newBNodeInTop := toTop.start
								newBNodeInTop.next = aux.next
								newBNodeInTop.prev = aux
								aux.next.prev = newBNodeInTop
								aux.next.leftchid = newBNodeInTop.rightchild
								aux.next = newBNodeInTop
								aux.rightchild = newBNodeInTop.leftchid
								root.size++
							}
						} else {
							newBNode.next = aux.next
							newBNode.prev = aux
							aux.next.prev = newBNode
							aux.next = newBNode
							root.size++
						}
						break
					}
					aux = aux.next
				}
				if aux.next == nil {
					if aux.data < newBNode.data {
						if aux.rightchild != nil {
							toTop = aux.rightchild.Insert(newBNode)
							if toTop.size == 1 && toTop.start == toTop.end && toTop.start.rightchild != nil {
								newBNodeInTop := toTop.start
								aux.next = newBNodeInTop
								aux.rightchild = newBNodeInTop.leftchid
								newBNodeInTop.prev = aux
								root.size++
								root.end = newBNodeInTop
							}
						} else {
							newBNode.prev = aux
							aux.next = newBNode
							root.end = newBNode
							root.size++
						}
					}
				}
			}
		}
		if root.size == 5 {
			return root.Divide()
		} else {
			return root
		}
	}
	return nil
}
func (aux *page) SearchToInsert(newBNode *bnode) *bnode {
	if aux != nil {
		temp := aux.start
		if temp.data > newBNode.data {
			if temp.leftchid != nil {
				return temp.leftchid.SearchToInsert(newBNode)
			}
			return nil
		} else if temp.data == newBNode.data {
			return temp
		}
		for temp != nil && temp.next != nil {
			if temp.data == newBNode.data {
				return temp
			} else if temp.data > newBNode.data {
				if temp.leftchid != nil {
					return temp.leftchid.SearchToInsert(newBNode)
				}
				return nil
			} else if temp.data < newBNode.data {
				if temp.rightchild != nil {
					return temp.rightchild.SearchToInsert(newBNode)
				}
				return nil
			}
			temp = temp.next
		}
		if temp != nil && temp.next == nil {
			if temp.data < newBNode.data && temp.rightchild != nil {
				return temp.rightchild.SearchToInsert(newBNode)
			} else if temp.data == newBNode.data {
				return temp
			}
			return nil
		}
	}
	return nil
}

func (p *page) Divide() *page {
	media := p.start
	leftChld := newpage()
	rightChld := newpage()
	parent := newpage()
	count := 1
	for media != nil {
		media2 := newBnode(media.data)
		media2.rightchild = media.rightchild
		media2.leftchid = media.leftchid
		if count < 3 {
			if leftChld.start == nil {
				leftChld.start = media2
				leftChld.end = media2
			} else {
				leftChld.end.next = media2
				media2.prev = leftChld.end
				leftChld.end = media2
			}
			leftChld.size++
		} else if count > 3 && count <= 5 {
			if rightChld.start == nil {
				rightChld.start = media2
				rightChld.end = media2
			} else {
				rightChld.end.next = media2
				media2.prev = rightChld.end
				rightChld.end = media2
			}
			rightChld.size++
		} else {
			parent.start = media2
			parent.end = media2
			parent.size++
		}
		media = media.next
		count++
	}
	parent.start.rightchild = rightChld
	parent.start.leftchid = leftChld
	return parent
}

func (tree btree) WriteTree(root *page) {
	naux := root.start
	root.escribirPag()
	for naux != nil {
		if naux.leftchid != nil {
			tree.WriteTree(naux.leftchid)
		}
		naux = naux.next
	}
	naux = root.start
	for naux.next != nil {
		naux = naux.next
	}
	if naux.rightchild != nil {
		tree.WriteTree(naux.rightchild)
	}
}
func (pag *page) escribirPag() {
	aux := pag.start
	fmt.Print("|")
	for aux != nil {
		fmt.Print(aux.data)
		fmt.Print("|")
		aux = aux.next
	}
	fmt.Println("")
}
func (nodo *bnode) writeLabel() string {
	respuesta := ""
	respuesta = "Id: " + strconv.Itoa(nodo.data)
	return respuesta
}

func (pag *page) writePage() string {
	tempArrows := ""
	temp := "nodo" + strconv.Itoa(pag.start.data) + " [ label =\""
	tempBNode := pag.start
	i := 0
	detalles := ""
	for tempBNode != nil {
		temp += "<C" + strconv.Itoa(i) + ">DPI: " + strconv.Itoa(tempBNode.data) + " " + tempBNode.writeLabel() + "|"
		if tempBNode.leftchid != nil {
			tempArrows += "nodo" + strconv.Itoa(pag.start.data) + ":C" + strconv.Itoa(i) + "->nodo" + strconv.Itoa(tempBNode.leftchid.start.data) + "\n"
		}
		i++
		tempBNode = tempBNode.next
	}
	temp += "<C" + strconv.Itoa(i) + ">\" fillcolor=\"#FFFFFF\"];\n"
	tempBNode = pag.start
	for tempBNode.next != nil {
		tempBNode = tempBNode.next
	}
	if tempBNode.rightchild != nil {
		tempArrows += "nodo" + strconv.Itoa(pag.start.data) + ":C" + strconv.Itoa(i) + "->nodo" + strconv.Itoa(tempBNode.rightchild.start.data) + "\n"
	}
	temp += tempArrows
	temp += detalles
	return temp
}
func (bTree *btree) graph(page *page) string {
	bnodesString := ""
	if bTree.root == nil {
		return ""
	}
	if page == nil {
		return ""
	}
	bnodesString += page.writePage()
	aux := page.start
	for aux != nil {
		// if aux.leftchid != nil {
		bnodesString += bTree.graph(aux.leftchid)
		// }
		aux = aux.next
	}
	aux = page.start
	for aux.next != nil {
		aux = aux.next
	}
	// if aux.rightchild != nil {
	bnodesString += bTree.graph(aux.rightchild)
	// }
	return bnodesString
}

func (bTree *btree) getGraphvizCode() string {
	return "digraph grafica{\n" +
		"rankdir=TB;\n" +
		"node [shape = record, style=filled, fillcolor=seashell2];\n" +
		bTree.graph(bTree.root) +
		"}\n"
}

func TestBTree() {
	ab := newBTree()
	for i := 25; i > 0; i-- {
		ab.Insert(i)
	}
	for i := 75; i < 100; i++ {
		ab.Insert(i)
	}
	for i := 50; i > 25; i-- {
		ab.Insert(i)
	}
	for i := 51; i < 75; i++ {
		ab.Insert(i)
	}
	fmt.Print(ab.getGraphvizCode())
}
