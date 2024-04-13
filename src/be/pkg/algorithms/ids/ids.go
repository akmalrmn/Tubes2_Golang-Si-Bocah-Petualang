package ids

import (
    "be/pkg/scraper"
    "be/pkg/tree"
)

func IterativeDeepeningSearch(start, end string) []string {
    root := tree.NewNode(start)
    for depth := 0; ; depth++ {
        visited := make(map[string]bool)
        result := depthLimitedSearch(root, end, depth, visited)
        if result != nil {
            var path []string
            for node := result; node != nil; node = node.Parent {
                path = append([]string{node.Value}, path...)
            }
            return path
        }
    }
}

func depthLimitedSearch(node *tree.Node, end string, depth int, visited map[string]bool) *tree.Node {
    if node.Value == end {
        return node
    }
    if depth <= 0 {
        return nil
    }
    visited[node.Value] = true
    links, _ := scraper.ExtractLinks("https://en.wikipedia.org" + node.Value)
    for _, link := range links {
        child := tree.NewNode(link)
        node.AddChild(child)
        if !visited[child.Value] {
            result := depthLimitedSearch(child, end, depth-1, visited)
            if result != nil {
                return result
            }
        }
    }
    return nil
}