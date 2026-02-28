package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
)

var db *gorm.DB

func VectorizeText(txt string) ([]float32, error) {
	body := map[string]interface{}{
		"model": "voyage-3-large",
		"input": []string{txt},
	}
	json_body, _ := json.Marshal(body)
	req, err := http.NewRequest(
		"POST",
		"https://api.voyageai.com/v1/embeddings",
		bytes.NewBuffer(json_body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("VOYAGE_API_KEY"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("voyage error: %s", string(body))
	}
	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	json.NewDecoder(res.Body).Decode(&result)
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	var Vector [][]float32
	for _, d := range result.Data {
		Vector = append(Vector, d.Embedding)
	}
	return result.Data[0].Embedding, nil
}
func ChunkText(text string, size int) []string {
	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += (size - 100) {
		end := i + size

		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		if end == len(runes) {
			break
		}
	}
	return chunks
}
func EmbedText(text string) ([][]float32, error) {
	chunks := ChunkText(text, 500)

	var allVec [][]float32

	for i, chunkContent := range chunks {
		vectorValues, err := VectorizeText(chunkContent)
		if err != nil {
			return nil, fmt.Errorf("failed to embed chunk %d: %v", i, err)
		}
		allVec = append(allVec, vectorValues)
	}
	return allVec, nil
}
