package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func UserScore(job, userCv string) (int, error) {
	body := map[string]interface{}{
		"model": "moonshotai/kimi-k2-instruct-0905",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `You are an AI job recruiter.

			Compare the candidate CV with the job description using this exact scoring rubric.

			Scoring rules:
			- skills_match: integer from 0 to 30
			- experience_match: integer from 0 to 25
			- project_relevance: integer from 0 to 20
			- tools_alignment: integer from 0 to 15
			- education_fit: integer from 0 to 10

			Important rules:
			- Be strict and consistent
			- Use only integers
			- Do not guess extra info not found in the CV
			- Return valid JSON only
			- Do not return explanations

			The final score will be calculated as:
			skills_match + experience_match + project_relevance + tools_alignment + education_fit`,
			},
			{
				"role":    "user",
				"content": "job:\n\n" + job + "\n\n" + "cv: \n" + userCv,
			},
		},
		"temperature": 0,
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name": "cv_scoring",
				"schema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"skills_match": map[string]interface{}{
							"type": "integer",
						},
						"experience_match": map[string]interface{}{
							"type": "integer",
						},
						"project_relevance": map[string]interface{}{
							"type": "integer",
						},
						"tools_alignment": map[string]interface{}{
							"type": "integer",
						},
						"education_fit": map[string]interface{}{
							"type": "integer",
						},
					},
					"required": []string{
						"skills_match",
						"experience_match",
						"project_relevance",
						"tools_alignment",
						"education_fit",
					},
					"additionalProperties": false,
				},
			},
		},
	}
	json_body, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %v", err)
	}
	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(json_body),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return 0, fmt.Errorf("Groq API error %d: %s", res.StatusCode, string(body))
	}
	type CvScoreResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	var groqResp CvScoreResponse
	if err := json.NewDecoder(res.Body).Decode(&groqResp); err != nil {
		return 0, fmt.Errorf("failed to decode Groq response: %v", err)
	}
	if len(groqResp.Choices) == 0 {
		return 0, fmt.Errorf("no choices returned from Groq")
	}
	content := groqResp.Choices[0].Message.Content
	fmt.Println("Groq returned:", content)
	var scoreStruct struct {
		SkillsMatch      int `json:"skills_match"`
		ExperienceMatch  int `json:"experience_match"`
		ProjectRelevance int `json:"project_relevance"`
		ToolsAlignment   int `json:"tools_alignment"`
		EducationFit     int `json:"education_fit"`
	}

	if err := json.Unmarshal([]byte(content), &scoreStruct); err != nil {
		return 0, fmt.Errorf("failed to unmarshal into decision struct: %v\nRaw content: %s", err, content)
	}
	total := scoreStruct.SkillsMatch +
		scoreStruct.ExperienceMatch +
		scoreStruct.ProjectRelevance +
		scoreStruct.ToolsAlignment +
		scoreStruct.EducationFit
	fmt.Print(total)
	return total, nil
}
