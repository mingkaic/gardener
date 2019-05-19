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

// Gardener ...
// Is the random generator of graphs and trees
type Gardener struct {
	*rand.Rand
}

// TreeNode ...
// Is the abstract output from the core generation routines
type TreeNode interface {
	NewChild(gen *rand.Rand) TreeNode // create and add node as child
	AddChild(child TreeNode)          // create edge between existing nodes
	HasChild(child TreeNode) bool     // check for existing edge
}

// =============================================
//                    Globals
// =============================================

var tokens = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// =============================================
//                    Public
// =============================================

// New ...
// Creates a new Gardener seeding with current time
func New() *Gardener {
	return &Gardener{rand.New(rand.NewSource(time.Now().Unix()))}
}

// NewSeed ...
// Creates a new Gardener seeding with input
func NewSeed(seed int64) *Gardener {
	return &Gardener{rand.New(rand.NewSource(seed))}
}

// RandTree ...
// Builds an n node subtree below input root
func (gardener *Gardener) RandTree(root TreeNode, n uint) {
	randMinSpanTree(gardener.Rand, root, n)
}

// RandGraph ...
// Builds an n node graph attached to root
func (gardener *Gardener) RandGraph(root TreeNode, n uint) {
	mst := randMinSpanTree(gardener.Rand, root, n)

	// connect edges
	for _, node := range mst {
		nConns := uint(gardener.Intn(int(n)))
		perms := randChoice(nConns, int(n))
		for _, idx := range perms {
			if !node.HasChild(mst[idx]) {
				node.AddChild(mst[idx])
			}
		}
	}
}

// =============================================
//                    Private
// =============================================

// build a minimum spanning tree of n nodes connected to root
func randMinSpanTree(gen *rand.Rand, root TreeNode, n uint) []TreeNode {
	src := []TreeNode{root}
	var i uint
	for i < n {
		parent := src[0]
		if len(src) > 1 {
			parent = src[gen.Intn(len(src))]
		}

		child := parent.NewChild(gen)
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
