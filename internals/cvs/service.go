package cvs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CVservice struct {
	repo *CvRepo
}

func NewUserService(repo *CvRepo) *CVservice {
	return &CVservice{repo: repo}
}

func CvUploadSvc() {

}

func (s *CVservice) ApplicationService(userId, JobId uint) (bool, error) {

	return true, nil
}
func (r *CvRepo) UserScore(job, userCv string) (int, error) {
	body := map[string]interface{}{
		"model": "moonshotai/kimi-k2-instruct-0905",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `You are an AI job recruiter. 

	Compare the candidate's CV to the job description and return a single numeric score that represents how well the candidate matches the job. 

	- Score range: 0.0000000 (worst match) to 100.0000000 (perfect match)  
	- Always return exactly 7 decimal places.  
	- Do not include any text, only the numeric score.  
	- Evaluate skills, experience, and relevance of the CV content to the job description.`,
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
						"score": map[string]interface{}{
							"type":        "string",
							"description": "Numeric score from 0.0000000 to 100.0000000 representing CV-job match with 7 decimal places.",
						},
					},
					"required":             []string{"score"},
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
		Message []struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	var groqResp CvScoreResponse
	if err := json.NewDecoder(res.Body).Decode(&groqResp); err != nil {
		return 0, fmt.Errorf("failed to decode Groq response: %v", err)
	}
	content := groqResp.Message[0].Content
	fmt.Println("Groq returned:", content)
	var scoreStruct struct {
		Score int `json:"score"`
	}

	if err := json.Unmarshal([]byte(content), &scoreStruct); err != nil {
		return 0, fmt.Errorf("failed to unmarshal into decision struct: %v\nRaw content: %s", err, content)
	}
	score := scoreStruct.Score

	return score, nil
}
