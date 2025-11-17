package seed

import (
	"context"
	"log"

	"github.com/odilonjk/golang-examples/rag/openai"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

func Generate(ctx context.Context, client *weaviate.Client) error {
	rubiconTrail := "The Rubicon Trail is a 22-mile-long route, part road and part 4x4 trail, located in the Sierra Nevada of the western United States, due west of Lake Tahoe and about 80 miles (130 km) east of Sacramento."
	rubiconEmbeddings, err := openai.GetEmbedding(rubiconTrail)
	if err != nil {
		return err
	}
	creator := client.Data().Creator()
	res, err := creator.WithClassName("Trail").
		WithProperties(map[string]interface{}{
			"name":        "Rubicon Trail",
			"region":      "Sierra Nevada, California",
			"description": rubiconTrail,
		}).
		WithVector(rubiconEmbeddings).Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("Embeddings persisted: %v", res)

	honeymoonPackTrail := "Honeymoon Pack Trail in St. George, Utah, is a historically significant landmark for many North American southwesterners. In the late 1800s, it was common for many newlywed Mormon settlers living in Arizona to travel to St. George Temple to formalize and celebrate their marriage vows. This pilgrimage was so frequent that the route they rode their wagons on was eventually named the Honeymoon Trail. Today, you can explore this historic trail on a 14-mile out-and-back route. Motorists and pedestrians, especially bird watchers and hikers, share this route. Therefore, it’s a great environment to practice side-by-side driving at slower speeds."
	honeymoonEmbeddings, err := openai.GetEmbedding(honeymoonPackTrail)
	if err != nil {
		return err
	}
	creator = client.Data().Creator()
	res, err = creator.WithClassName("Trail").
		WithProperties(map[string]interface{}{
			"name":        "Honeymoon Pack Trail",
			"region":      "St. George, Utah",
			"description": honeymoonPackTrail,
		}).
		WithVector(honeymoonEmbeddings).Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("Embeddings persisted: %v", res)

	return nil
}
