package Result

import (
	"log"
	"strings"
)

var (
	ArticleCount = 1
)

type Node struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Relation struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

type Graph struct {
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"links"`
}

type Result struct {
	Time                int64 `json:"time"`
	TotalArticleChecked int   `json:"total_article_checked"`
	TotalArtcleVisit    int   `json:"total_article_visit"`
	Traph               Graph `json:"graph"`
}

func (g *Graph) AddNode(node Node) {
	g.Nodes = append(g.Nodes, node)
}

func (g *Graph) AddRelation(source, target int) {
	g.Relations = append(g.Relations, Relation{
		Source: source,
		Target: target,
	})
}

func (g *Graph) FindNode(s string) int {
	for _, v := range g.Nodes {
		if v.Title == s {
			return v.Id
		}
	}
	return -1
}

func (g *Graph) GenerateGraph(input [][]string) {

	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[i]); j++ {
			input[i][j] = trimPrefix(input[i][j])
		}
	}

	g.GenerateNodeBFS(input)
	g.GenerateRelationBFS(input)
}

func (g *Graph) GetNodesCount() int {
	return len(g.Nodes)
}

func (g *Graph) GenerateNodeBFS(input [][]string) {
	ArticleCount = 0
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[i]); j++ {
			if g.FindNode(input[i][j]) == -1 {
				g.AddNode(Node{
					Id:    ArticleCount,
					Title: input[i][j],
				})
				ArticleCount++
			}
		}
	}
}

func trimPrefix(input string) string {
	// Delete https://en.wikipedia.org/wiki/
	input = input[30:]
	// Delete _ from the title
	return strings.Replace(input, "_", " ", -1)
}

func (g *Graph) GenerateRelationBFS(input [][]string) {
	log.Println("Generating relations")
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[i])-1; j++ {

			source := g.FindNode(input[i][j])
			target := g.FindNode(input[i][j+1])

			if source == -1 {
				log.Println("Error: Source not found in the graph", input[i][j])
				continue
			}

			if target == -1 {
				log.Println("Error: Target not found in the graph", input[i][j+1])
				continue
			}

			log.Println("Source:", source, "Target:", target)
			g.AddRelation(source, target)
			log.Println("Relation added:", source, "->", target)
		}
	}
}
