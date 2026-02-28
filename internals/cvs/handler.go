package cvs

import (
	"net/http"
	"strings"

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

func CVUploader(w http.ResponseWriter, r *http.Request) {
	var req struct {
		JobId string
	}
	Utils.ParseJSON(r, &req)
	file, header, err := r.FormFile("file")
	user, _ := r.Context().Value("user").(*models.User)
	if err != nil {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err,
		})
		return
	}
	defer file.Close()

	fileUrl, err := config.UploaderCloudinary(file, header.Filename)
	if err != nil {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err,
		})
		return
	}

	var pdfText strings.Builder

	reader, err := pdf.NewReader(file, header.Size)
	if err != nil {
		Utils.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err,
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
	for i, vec := range vecs {
		cvs = append(cvs, models.Cvs{
			Content:   chunks[i],
			Vector:    pgvector.NewVector(vec),
			SourceURL: fileUrl,
			Uploadby:  user.ID,
		})
	}
}
func Application(w http.ResponseWriter, r *http.Request) {

}
