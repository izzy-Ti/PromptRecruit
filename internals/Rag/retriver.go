package rag

import (
	"context"
	"fmt"

	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	"gorm.io/gorm"
)

func Retriver(db *gorm.DB, query string) ([]string, error) {
	ctx := context.Background()
	queryVec, err := VectorizeText(query)
	if err != nil {
		return nil, err
	}
	var docs []models.Cvs
	err = db.WithContext(ctx).
		Model(&models.Cvs{}).
		Select("content").
		Order(gorm.Expr("embedding <-> ?", queryVec)).
		Limit(3).
		Find(&docs).Error
	if err != nil {
		return nil, err
	}

	var results []string
	for _, d := range docs {
		results = append(results, d.Content)
	}

	fmt.Println("Retriever found documents:", len(results))
	return results, nil
}
