package cvs

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	rag "github.com/izzy-Ti/PromptRecruit/internals/Rag"
	"github.com/izzy-Ti/PromptRecruit/internals/Utils"
	"github.com/izzy-Ti/PromptRecruit/internals/config"
	"github.com/ledongthuc/pdf"
	"github.com/pgvector/pgvector-go"
)

type NewHandler struct {
	svc *CVservice
}

func NewCVHnadler(svc *CVservice) *NewHandler {
	return &NewHandler{svc: svc}
}

func (s *NewHandler) CVUploader(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil || file == nil {
		Utils.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "No file uploaded",
		})
		return
	}
	user, _ := r.Context().Value("user").(*models.User)
	if user == nil {
		Utils.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "No user given",
		})
		return
	}
	defer file.Close()

	fileUrl, err := config.UploaderCloudinary(file, header.Filename)
	if err != nil {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var pdfText strings.Builder

	reader, err := pdf.NewReader(file, header.Size)
	if err != nil {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	totalPage := reader.NumPage()

	for i := 1; i <= totalPage; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, _ := page.GetPlainText(nil)
		pdfText.WriteString(content)
	}
	var cvs []models.Cvs
	chunks := rag.ChunkText(pdfText.String(), 500)
	vecs, err := rag.EmbedText(pdfText.String())
	if err != nil {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	for i, vec := range vecs {
		cvs = append(cvs, models.Cvs{
			Content:   chunks[i],
			Vector:    pgvector.NewVector(vec),
			SourceURL: fileUrl,
			Uploadby:  user.ID,
		})
	}
	for _, cv := range cvs {
		_ = s.svc.repo.db.Create(&cv)
	}
	Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
		"success": true,
		"message": "CV uploaded successfully",
	})
}
func (s *NewHandler) Application(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId, ok := vars["jobId"]
	jobIdin, err := strconv.Atoi(jobId)
	user, _ := r.Context().Value("user").(*models.User)
	score, ok, err := s.svc.ApplicationService(user.ID, uint(jobIdin))
	if !ok {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	Utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"score":   score,
		"message": "Application submitted successfully",
	})
}
func (s *NewHandler) JobPost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	Utils.ParseJSON(r, &req)
	user, _ := r.Context().Value("user").(*models.User)
	ok, err := s.svc.jobAddService(req.Title, req.Content, user.ID)
	if !ok {
		Utils.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	Utils.WriteJson(w, http.StatusAccepted, map[string]interface{}{
		"success": true,
		"message": "job post successfully",
	})
}
