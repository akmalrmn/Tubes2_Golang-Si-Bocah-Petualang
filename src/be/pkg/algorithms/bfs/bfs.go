package bfs

import (
	"be/pkg/Result"
	"be/pkg/config"
	"be/pkg/scraper"
	"be/pkg/set"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2/queue"
	"log"
	"sync"
	"time"
)

func BFS(starts, ends string, con *config.Config) []byte {

	scraper.ArticleCount = 0

	starts = "/wiki/" + starts
	ends = "/wiki/" + ends

	log.Println("BFS Starts")
	log.Println("Starts:", starts)
	log.Println("Ends:", ends)

	startTime := time.Now()
	parents := sync.Map{}

	queueInput, _ := queue.New(
		con.MaxQueryThread,
		&queue.InMemoryQueueStorage{MaxSize: con.MaxQueueSize},
	)

	queueOutput, _ := queue.New(
		con.MaxQueryThread,
		&queue.InMemoryQueueStorage{MaxSize: con.MaxQueueSize},
	)

	// Add URLs to the first queue
	err := queueInput.AddURL("https://en.wikipedia.org" + starts)
	if err != nil {
		log.Println("Error adding URL to queue:", err)
		return []byte(fmt.Sprintf(`{Error adding URL to queue : %v}`, err))
	}

	var result = set.NewSetOfSlice()

	depth := 0

	for {
		log.Println("Depth:", depth)
		depth++

		startTime2 := time.Now()
		tempResult := scraper.QueueColly(queueInput, queueOutput, starts, ends, con, &parents)
		log.Println("Time elapsed:", time.Since(startTime2).Milliseconds())

		result = result.Union(tempResult)

		if result.Size() > 0 {
			log.Println("Found the target")
			break
		}

		if depth > con.MaxDepth {
			log.Println("Max depth reached")
			break

		} else {
			log.Println("Article checked", scraper.ArticleCount)
			queueInput = queueOutput
			queueOutput, _ = queue.New(
				con.MaxQueryThread,
				&queue.InMemoryQueueStorage{MaxSize: con.MaxQueueSize},
			)
		}
	}
	elapsed := time.Since(startTime).Milliseconds()

	log.Println("Result:")
	for _, v := range result.ToSlice() {
		log.Println(v)
	}

	var graph Result.Graph
	graph.GenerateGraph(result.ToSlice())

	outputResult := Result.Result{
		Traph:               graph,
		Time:                elapsed,
		TotalArticleChecked: scraper.ArticleCount,
		TotalArtcleVisit:    graph.GetNodesCount(),
	}

	jsonOutput, err := json.Marshal(outputResult)

	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return []byte(fmt.Sprintf(`{Error marshalling JSON : %v}`, err))
	}

	return jsonOutput
}
