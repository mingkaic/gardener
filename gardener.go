//// file: gardener.go

// Package gardener ...
// Is a package for generating random tree and graph structures
package gardener

import (
	"math/rand"
	"time"
)

// =============================================
//                    Declarations
// =============================================

// TreeNode ...
// Is the abstract output from the core generation routines
type TreeNode interface {
	NewChild() *TreeNode           // create and add node as child
	AddChild(child *TreeNode)      // create edge between existing nodes
	HasChild(child *TreeNode) bool // check for existing edge
}

// =============================================
//                    Globals
// =============================================

var gen = rand.New(rand.NewSource(time.Now().Unix()))

var tokens = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// =============================================
//                    Public
// =============================================

// Seed ...
// Seeds the local generator
func Seed(seed int64) {
	gen.Seed(seed)
}

// RandTree ...
// Builds an n node subtree below input root
func RandTree(root *TreeNode, n uint) {
	randMinSpanTree(root, n)
}

// RandGraph ...
// Builds an n node graph attached to root
func RandGraph(root *TreeNode, n uint) {
	mst := randMinSpanTree(root, n)

	// connect edges
	for _, node := range mst {
		nConns := uint(gen.Intn(int(n)))
		perms := randChoice(nConns, int(n))
		for _, idx := range perms {
			if !(*node).HasChild(mst[idx]) {
				(*node).AddChild(mst[idx])
			}
		}
	}
}

// =============================================
//                    Private
// =============================================

// build a minimum spanning tree of n nodes connected to root
func randMinSpanTree(root *TreeNode, n uint) []*TreeNode {
	src := []*TreeNode{root}
	var i uint = 0
	for i < n {
		parent := src[0]
		if len(src) > 1 {
			parent = src[gen.Intn(len(src))]
		}

		child := (*parent).NewChild()
		if child != nil {
			src = append(src, child)
			i++
		}
	}
	return src
}

// randomly choose c elements from array 1 to n
func randChoice(c uint, n int) []int {
	list := rand.Perm(n)
	return list[:c]
}
