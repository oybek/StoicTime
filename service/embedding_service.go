package service

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"oybek.io/kerege/model"
)

const maxWordNum = 2000

type EmbeddingService struct {
	documentStorage     DocumentStorage
	embeddingCalculator EmbeddingCalculator
}

func NewEmbeddingService(
	documentStorage DocumentStorage,
	embeddingCalculator EmbeddingCalculator,
) *EmbeddingService {
	return &EmbeddingService{
		documentStorage:     documentStorage,
		embeddingCalculator: embeddingCalculator,
	}
}

func (es *EmbeddingService) GetAndStoreEmbeddings(ctx context.Context, text string) error {
	blocks := splitByArticle(text)

	excerpt := func(s string) string {
		return string([]rune(s)[:min(len(s), 16)])
	}

	for i, block := range blocks {
		time.Sleep(2 * time.Second)

		if len(strings.Fields(block)) > maxWordNum {
			log.Printf("block %d: %s... - skip ❌", i, excerpt(block))
			continue
		}

		embedding, err := es.embeddingCalculator.GetEmbedding(ctx, block)
		if err != nil {
			log.Printf("block %d: %s... - error get embedding: %v ❌", i, excerpt(block), err)
			continue
		}

		document := model.Document{
			Content:   block,
			Embedding: embedding,
		}
		err = es.documentStorage.Store(ctx, &document)
		if err != nil {
			log.Printf("block %d: %s... - error store: %v ❌", i, excerpt(block), err)
			continue
		}

		log.Printf("block %d: %s... - ok ✅", i, excerpt(block))
	}

	return nil
}

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

		blocks[i] = cleanString(text[start:end])
	}

	return blocks
}

func cleanString(s string) string {
	// 1. Replace all newlines with spaces
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ") // for Windows \r\n cases

	// 2. Collapse multiple spaces into one
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	// 3. Optionally trim spaces at start and end
	s = strings.TrimSpace(s)

	return s
}
