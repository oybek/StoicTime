package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/sashabaranov/go-openai"
	"oybek.io/kerege/config"
	"oybek.io/kerege/openaiw"
)

//go:embed assets/taxcode.txt
var taxCode string

func splitByArticle(text string) []string {
	re := regexp.MustCompile(`(?m)^Статья \d+\.`)
	matches := re.FindAllStringIndex(text, -1)

	blocks := make([]string, len(matches))
	for i, match := range matches {
		start := match[0]
		// If it's the last match, the block ends at the end of the string.
		end := len(text)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}

		block := strings.TrimSpace(text[start:end])
		blocks = append(blocks, block)
	}

	return blocks
}

func main() {
	ctx := context.Background()

	blocks := splitByArticle(taxCode)
	slices.SortFunc(blocks, func(a, b string) int {
		return len(b) - len(a)
	})
	log.Printf("Кол-во слов: %d\n%s", len(strings.Fields(blocks[1])), blocks[1])

	theConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error create config: %v", err)
	}

	//
	var openAIWrapper *openaiw.OpenAIWrapper
	{
		client := openai.NewClient(theConfig.OpenAIToken)
		openAIWrapper = openaiw.NewOpenAIWrapper(client)
	}

	pgConn, err := pgx.Connect(ctx, theConfig.PGDSN)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer pgConn.Close(ctx)

	embedding, err := openAIWrapper.GetEmbedding(ctx, blocks[1])

	fmt.Printf("error: %#v\n", err)
	fmt.Printf("embedding: %#v\n", embedding)
}
