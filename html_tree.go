//// file: html_tree.go

package gardener

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"gopkg.in/eapache/queue.v1" // this queue isn't threadsafe
)

// =============================================
//                    Declarations
// =============================================

// NodeInfo ...
// Is the pointer content of HTMLNode
// to persist values after frequent value conversion
type NodeInfo struct {
	Pos      uint
	Tag      string
	Children []*TreeNode
	Attrs    map[string][]string
}

// PageInfo ...
// Aggregates DOM data for generated pages
// During generation, PageInfo is shared for every node
// It serves as the Ad Hoc information (expected output) when testing web scrapers
type PageInfo struct {
	Tags  NodeMap
	Attrs NodeMap

	// parameters used during DOM construction
	links      map[string]struct{}
	nRemaining int
}

// HTMLNode ...
// Is a TreeNode implementation and
// an atomic element of a page
type HTMLNode struct {
	*NodeInfo
	Info *PageInfo
}

// NodeMap ...
// Maps a string key to many Nodes
type NodeMap map[string][]*HTMLNode

// =============================================
//                    Globals
// =============================================

var secContent = []string{
	"h1", "h2", "h3", "h4", "h5", "h6",
	"article", "section", "footer", "header", "nav",
}

var textContent = []string{"div", "hr", "li", "main", "p", "ul"}

var content = []string{"a", "img", "span", "audio", "video", "source"}

var tagPool = map[string][]string{
	"body": append(append(secContent, textContent...), content...),

	"h1":      textContent,
	"h2":      textContent,
	"h3":      textContent,
	"h4":      textContent,
	"h5":      textContent,
	"h6":      textContent,
	"article": append(content, textContent...),
	"section": append(content, textContent...),
	"footer":  append(content, textContent...),
	"header":  append(content, textContent...),
	"nav":     append(content, textContent...),

	"main": append(content, textContent...),
	"div":  append(content, textContent...),
	"ul":   append(content, "li"),
	"li":   content,
	"hr":   {},
	"p":    {},

	"a":      {"img", "span", "audio", "video", "source"},
	"audio":  {"source"},
	"video":  {"source"},
	"source": {},
	"span":   {},
	"img":    {},
}

var commonAttrs = []string{"class", "id"}

var attrPool = map[string][]string{
	"head":  commonAttrs,
	"body":  commonAttrs,
	"title": commonAttrs,

	"h1":      commonAttrs,
	"h2":      commonAttrs,
	"h3":      commonAttrs,
	"h4":      commonAttrs,
	"h5":      commonAttrs,
	"h6":      commonAttrs,
	"article": commonAttrs,
	"section": commonAttrs,
	"footer":  commonAttrs,
	"header":  commonAttrs,
	"nav":     commonAttrs,

	"main": commonAttrs,
	"div":  commonAttrs,
	"ul":   commonAttrs,
	"li":   append(commonAttrs, "value"),
	"hr":   commonAttrs,
	"p":    commonAttrs,

	"a":      append(commonAttrs, "href"),
	"audio":  append(commonAttrs, "controls"),
	"img":    append(commonAttrs, "src"),
	"source": append(commonAttrs, "src"),
	"span":   commonAttrs,
	"video":  append(commonAttrs, "controls"),
}

// =============================================
//                    Public
// =============================================

//// Gardener Extension

// GeneratePage ...
// Randomly generates a DOM structure
// guaranteeing it contains nElems elements and input links
func (this Gardener) GeneratePage(nElems uint, links map[string]struct{}) *HTMLNode {
	info := &PageInfo{make(NodeMap), make(NodeMap), links, int(nElems - 4)}
	title := HTMLNode{
		&NodeInfo{Tag: "title", Attrs: map[string][]string{}},
		info}
	var tTitle TreeNode = title
	head := HTMLNode{
		&NodeInfo{Tag: "head", Attrs: make(map[string][]string), Children: []*TreeNode{&tTitle}},
		info}
	body := HTMLNode{
		&NodeInfo{Tag: "body", Attrs: make(map[string][]string)},
		info}
	var tHead TreeNode = head
	var tBody TreeNode = body
	html := HTMLNode{
		&NodeInfo{Tag: "html", Attrs: make(map[string][]string), Children: []*TreeNode{&tHead, &tBody}},
		info}
	var tHtml TreeNode = html
	root := &HTMLNode{
		&NodeInfo{Attrs: make(map[string][]string), Children: []*TreeNode{&tHtml}},
		info}

	this.RandTree(&tBody, nElems-4)

	var pos uint = 1
	q := queue.New()
	q.Add(&tHtml)
	for q.Length() > 0 {
		var node = q.Remove()
		tPtr := node.(*TreeNode)
		cVal := (*tPtr).(HTMLNode)
		cVal.Pos = pos
		pos++
		for _, child := range cVal.Children {
			q.Add(child)
		}
	}

	return root
}

//// Public Utility

// ToHTML ...
// Obtain the HTML string of input HTML tree
func ToHTML(node *HTMLNode) string {
	if node == nil {
		panic("printing nil HTMLNode")
	}
	result := ""
	if len(node.Tag) > 0 {
		result += "<" + node.Tag
		for key, val := range node.Attrs {
			if len(key) > 0 {
				result += fmt.Sprintf(" %s=\"%s\"", key, val[0])
			}
		}
		result += ">"
		if attr, ok := node.Attrs[""]; ok {
			result += attr[0]
		}
	}
	for _, child := range node.Children {
		hChild := (*child).(HTMLNode)
		result += ToHTML(&hChild)
	}
	if len(node.Tag) > 0 {
		result += "</" + node.Tag + ">"
	}
	return result
}

//// Methods for HTMLNode

// NewChild ...
// Make a new HTMLNode and add as child
func (this HTMLNode) NewChild(gen *rand.Rand) *TreeNode {
	// check whether this node supports children
	potentialTags, ok := tagPool[this.Tag]
	if !ok || len(potentialTags) == 0 {
		return nil
	}

	if this.Info == nil {
		panic("NewChild should never have nil Info")
	}

	s := HTMLNode{
		&NodeInfo{Attrs: make(map[string][]string)},
		this.Info,
	}

	// if the number of remaining elements to fill is less than the link set use a tag instead
	if this.Info.links != nil && this.Info.nRemaining <= len(this.Info.links) {
		s.Tag = "a"
	} else {
		// determine likely tags given this parent
		s.Tag = potentialTags[gen.Intn(len(potentialTags))]
	}
	this.Info.Tags[s.Tag] = append(this.Info.Tags[s.Tag], &s)

	var potentialAttrs = attrPool[s.Tag]
	for _, attr := range potentialAttrs {
		if attr == "href" { // always assign href
			links := this.Info.links
			if links == nil || len(links) == 0 {
				s.Attrs[attr] = []string{"#"}
			} else {
				var selLink string
				for link := range links {
					selLink = link
					break
				}
				s.Attrs[attr] = []string{selLink}
				delete(links, selLink)
			}
			this.Info.Attrs[attr] = append(this.Info.Attrs[attr], &s)
		} else if gen.Intn(2) == 1 {
			s.Attrs[attr] = []string{uuid.New().String()}
			this.Info.Attrs[attr] = append(this.Info.Attrs[attr], &s)
		}
	}
	this.Info.nRemaining-- // once this reaches 0, every element generated is an a element

	var out TreeNode = s
	this.AddChild(&out)
	return &out
}

// AddChild ...
// Add an existing HTMLNode as child
func (this HTMLNode) AddChild(child *TreeNode) {
	this.Children = append(this.Children, child)
}

// HasChild ...
// Check if this already have HTMLNode child
func (this HTMLNode) HasChild(child *TreeNode) bool {
	has := false
	n := len(this.Children)
	hChild := (*child).(HTMLNode)
	for i := 0; i < n && !has; i++ {
		hThis := (*this.Children[i]).(HTMLNode)
		has = has || hThis.NodeInfo == hChild.NodeInfo
	}
	return has
}
