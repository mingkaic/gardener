package gardener

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

// =============================================
//                    Declarations
// =============================================

// PageNode ...
// Is the pointer content of SiteNode
// to persist values after frequent value conversion
type PageNode struct {
	Depth                        uint
	Hostname, LinkPath, FullLink string
	Page                         *HTMLNode
	Refs                         []TreeNode
	BackLinks                    []*PageNode
}

// SiteInfo ...
// Aggregates Site data for generated site
// Serves as the Ad Hoc information (expected output) when testing web crawlers
type SiteInfo struct {
	Pages    PageMap
	MaxDepth uint
	HostList []string
}

// SiteNode ...
// Is a TreeNode implementation and
// an atomic element of a site
type SiteNode struct {
	*PageNode
	Info *SiteInfo
}

// PageMap ...
// Associates links to a single HTML page
type PageMap map[string]*PageNode

// =============================================
//                    Globals
// =============================================

const protocol = "http"

// =============================================
//                    Public
// =============================================

//// Gardener Extension

// GeneratePage ...
// Randomly generates a Website graph
func (gardener *Gardener) GenerateSite(nSites uint) *SiteNode {
	hostname := protocol + "://" + uuid.New().String() + ".com"
	linkPath := uuid.New().String()
	fullLink := hostname + "/" + linkPath
	content := &PageNode{Hostname: hostname, LinkPath: linkPath, FullLink: fullLink}
	site := &SiteNode{
		content,
		&SiteInfo{Pages: PageMap{fullLink: content}, HostList: []string{hostname}},
	}
	var tOrigin TreeNode = site
	gardener.RandGraph(tOrigin, uint(nSites-1))

	// build up Depth, and Page
	visited := make(map[string]struct{})
	var siteTraversal func(TreeNode, uint) string
	siteTraversal = func(node TreeNode, depth uint) string {
		page := node.(*SiteNode)
		if site.Info.MaxDepth < depth {
			site.Info.MaxDepth = depth
		}
		if _, ok := visited[page.FullLink]; !ok {
			visited[page.FullLink] = struct{}{}
			page.Depth = depth

			links := make(map[string]struct{})
			for _, ref := range page.Refs {
				links[siteTraversal(ref, depth+1)] = struct{}{}
			}
			// number of elements on a page should be significantly higher than number of potential links
			nElems := nSites*uint(2) + uint(gardener.Intn(91))
			page.Page = gardener.GeneratePage(nElems, links)
			if page.Page == nil {
				panic("Generated nil page")
			}
		}
		return page.FullLink
	}
	siteTraversal(tOrigin, 0)

	return site
}

//// Members for SiteNode

// NewChild ...
// gardener a new Site and add as reference
func (gardener *SiteNode) NewChild(gen *rand.Rand) TreeNode {
	var hostname string
	nHosts := uint(len(gardener.Info.HostList))
	if nHosts == 0 {
		hostname = protocol + "://" + uuid.New().String() + ".com"
		gardener.Info.HostList = []string{hostname}
	} else {
		idx := uint(math.Abs(gen.NormFloat64() * float64(nHosts) / 2))
		if idx >= nHosts {
			// add new hostname
			hostname = protocol + "://" + uuid.New().String() + ".com"
			gardener.Info.HostList = append(gardener.Info.HostList, hostname)
		} else {
			// use existing hostname
			hostname = gardener.Info.HostList[idx]
		}
	}

	linkPath := uuid.New().String()
	s := &SiteNode{
		&PageNode{
			Hostname:  hostname,
			LinkPath:  linkPath,
			FullLink:  hostname + "/" + linkPath,
			BackLinks: []*PageNode{gardener.PageNode},
		},
		gardener.Info,
	}
	gardener.Info.Pages[s.FullLink] = s.PageNode

	var out TreeNode = s
	gardener.AddChild(out)
	return out
}

// AddChild ...
// Add an existing SiteNode as reference
func (gardener *SiteNode) AddChild(child TreeNode) {
	sChild := child.(*SiteNode)
	sChild.BackLinks = append(sChild.BackLinks, gardener.PageNode)
	gardener.Refs = append(gardener.Refs, child)
}

// HasChild ...
// Check if gardener already have SiteNode as a reference
func (gardener *SiteNode) HasChild(child TreeNode) bool {
	has := false
	n := len(gardener.Refs)
	sChild := child.(*SiteNode)
	for i := 0; i < n && !has; i++ {
		has = has || gardener.Refs[i].(*SiteNode).PageNode == sChild.PageNode
	}
	return has
}
