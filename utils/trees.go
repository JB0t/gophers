package trees

import (
	"bytes"
	"fmt"
)

type Node struct {
	parent *Node
	data []byte
	children []*Node
}

type Trie struct {
	rootNode *Node
}

func (t *Trie) init() {
	t.rootNode=&Node{parent:nil,data:make([]byte,0)}
}

func (t *Trie) printBranches() {
	branchesToPrint:=make([][]byte,0)
	for _,child:=range t.rootNode.children{
		printBranchFromTrie(child, make([]byte, 0), branchesToPrint)
	}
}

func (t *Trie) printNodes() {
	for _,child:=range t.rootNode.children{
		printNodeFromTrie(child)
	}
}

func printNodeFromTrie(node *Node) {
	fmt.Printf("Node: '%s'",string(node.data))
	if node.children != nil {
		for _,child := range node.children{
			printNodeFromTrie(child)
		}
	}
}

func printBranchFromTrie(node *Node, currentBranch []byte, branchesToPrint [][]byte){
	currentBranch = append(currentBranch, node.data...)
	if node.children == nil {
		fmt.Printf("Valid word: %s\n",currentBranch)
	} else {
		for _,child:=range node.children{
			 printBranchFromTrie(child, currentBranch, branchesToPrint)
		}
	}
}

func (t *Trie) addDataToTrie(data []byte) {
	addData(t.rootNode,data)
}

func (t *Trie) convertToPrefixTrie() {
	convertToPrefixTrie(t.rootNode)
}

func addData(node *Node, data []byte) {
	if node.children == nil {
		fmt.Printf("No more children, adding '%s' under '%s'\n",string(data),func(s []byte)string{if string(s)==""{return "root"};return string(s)}(node.data))
		insertDataUnderNode(node, data)
	} else {
		checkByte:=false
		for _,child:=range node.children {
			if bytes.Equal(child.data,data[:len(child.data)]){
				checkByte = true
				addData(child,data[len(child.data):])
			}
		}
		if checkByte == false {
			fmt.Printf("No adjacent matches, adding '%s' under '%s'\n",string(data),func(s []byte)string{if string(s)==""{return "root"};return string(s)}(node.data))
			insertDataUnderNode(node, data)
		}
	}
}

func insertDataUnderNode(nodeStart *Node, data []byte) {
	currData:=make([]byte,0)
	currNode:=&Node{parent:nodeStart,data:append(currData,data[0])}
	nodeStart.children = append(nodeStart.children,currNode)
	for i:=1;i<len(data);i++ {
		nextNode:=&Node{parent:currNode,data:append(currData,data[i])}
		currNode.children=append(currNode.children,nextNode)
		currNode=nextNode
	}
	currNode.children=nil
}

func convertToPrefixTrie(node *Node){
	for _,child := range node.children{
		if child.children == nil {
			if len(child.parent.children) == 1 {
				child.parent.data = append(child.parent.data, child.data...)
				child.parent.children = nil
				convertToPrefixTrie(child.parent)
			}
		} else {
			convertToPrefixTrie(child)
		}
		fmt.Printf("Node prefix of child: '%s'\n",string(node.data))
	}
}

