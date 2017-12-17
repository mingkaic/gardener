//// file: html_tree.go

package gardener

import (
	"fmt"
	"gopkg.in/eapache/queue.v1"
)

//// ====== Structures ======

type HTMLInfo struct {
	Pos uint
	Tag string
	Children []*TreeNode
	Attrs map[string][]string
}

type MapInfo struct {
	Tags nodeMap
	Attrs nodeMap
}

type HTMLNode struct {
	*HTMLInfo
	TreeInfo *MapInfo
}

type nodeMap map[string][]*HTMLNode

//// ====== Globals ======

var secContent = []string {
	"h1", "h2", "h3", "h4", "h5", "h6",
	"article", "section", "footer", "header", "nav",
}

var textContent = []string {"div", "hr", "li", "main", "p", "ul"}

var content = []string {"a", "img", "span", "audio", "video", "source"}

var tagPool = map[string][]string {
	"body": append(append(secContent, textContent...), content...),

	"h1": textContent,
	"h2": textContent,
	"h3": textContent,
	"h4": textContent,
	"h5": textContent,
	"h6": textContent,
	"article": append(content, textContent...),
	"section": append(content, textContent...),
	"footer": append(content, textContent...),
	"header": append(content, textContent...),
	"nav": append(content, textContent...),

	"main": append(content, textContent...),
	"div": append(content, textContent...),
	"ul": append(content, "li"),
	"li": content,
	"hr": {},
	"p": {},

	"a": {"img", "span", "audio", "video", "source"},
	"audio": {"source"},
	"video": {"source"},
	"source": {},
	"span": {},
	"img": {},
}

var attrPool = map[string][]string {
	"head": {},
	"body": {},
	"title": {},

	"h1": {},
	"h2": {},
	"h3": {},
	"h4": {},
	"h5": {},
	"h6": {},
	"article": {},
	"section": {},
	"footer": {},
	"header": {},
	"nav": {},

	"main": {},
	"div": {},
	"ul": {},
	"li": {},
	"hr": {},
	"p": {},

	"a": {},
	"audio": {},
	"video": {},
	"source": {},
	"span": {},
	"img": {},
}

//// ====== Public ======

//// Member for HTMLNode

func (this HTMLNode) NewChild() *TreeNode {
	// check whether this node supports children
	potentialTags, ok := tagPool[this.Tag]
	if !ok || len(potentialTags) == 0 {
		return nil
	}

	if this.TreeInfo == nil {
		this.TreeInfo = &MapInfo{make(nodeMap), make(nodeMap)}
		if len(this.Tag) > 0 {
			this.TreeInfo.Tags[this.Tag] = append(this.TreeInfo.Tags[this.Tag], &this)
		}
		for attr, _ := range this.Attrs {
			this.TreeInfo.Attrs[attr] = append(this.TreeInfo.Attrs[attr], &this)
		}
	}

	s := HTMLNode{
		&HTMLInfo{Attrs: make(map[string][]string)},
		this.TreeInfo,
	}

	// determine likely tags given this parent
	s.Tag = potentialTags[gen.Intn(len(potentialTags))]
	this.TreeInfo.Tags[s.Tag] = append(this.TreeInfo.Tags[s.Tag], &s)

	var potentialAttrs = attrPool[s.Tag]
	for _, attr := range potentialAttrs {
		if gen.Intn(2) == 1 {
			s.Attrs[attr] = append(s.Attrs[attr], RandString(17))
			this.TreeInfo.Attrs[attr] = append(this.TreeInfo.Attrs[attr], &s)
		}
	}

	var out TreeNode = s
	this.AddChild(&out)
	return &out
}

func (this HTMLNode) AddChild(child *TreeNode) {
	this.Children = append(this.Children, child)
}

func (this HTMLNode) HasChild(child* TreeNode) bool {
	has := false
	n := len(this.Children)
	for i := 0; i < n && !has; i++ {
		has = has || this.Children[i] == child
	}
	return has
}

//// Core Functions

func GenerateSite() *HTMLNode {
	info := &MapInfo{make(nodeMap), make(nodeMap)}
	title := HTMLNode{
		&HTMLInfo{Tag: "title", Attrs: map[string][]string{}},
		info}
	var tTitle TreeNode = title
	head := HTMLNode{
		&HTMLInfo{Tag: "head", Attrs: make(map[string][]string), Children: []*TreeNode{&tTitle}},
		info}
	body := HTMLNode{
		&HTMLInfo{Tag: "body", Attrs: make(map[string][]string)},
		info}
	var tHead TreeNode = head
	var tBody TreeNode = body
	html := HTMLNode{
		&HTMLInfo{Tag: "html", Attrs: make(map[string][]string), Children: []*TreeNode{&tHead, &tBody}},
		info}
	var tHtml TreeNode = html
	root := &HTMLNode{
		&HTMLInfo{Attrs: make(map[string][]string), Children: []*TreeNode{&tHtml}},
		info}

	RandTree(&tBody, 96)

	var pos uint = 1
	q := queue.New()
	q.Add(&tHtml)
	for q.Length() > 0 {
		var node  = q.Peek()
		q.Remove()
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

func ToHTML(node *HTMLNode) string {
	result := "<" + node.Tag
	for key, val := range node.Attrs {
		if len(key) > 0 {
			result += fmt.Sprintf(" %s=\"%s\"", key, val[0])
		}
	}
	result += ">"
	if attr, ok := node.Attrs[""]; ok {
		result += attr[0]
	}
	for _, child := range node.Children {
		hChild := (*child).(HTMLNode)
		result += ToHTML(&hChild)
	}
	result += "</" + node.Tag + ">"
	return result
}
