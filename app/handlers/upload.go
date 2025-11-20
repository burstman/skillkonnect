package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anthdm/superkit/kit"
)

// AdminUploadFile handles file uploads
// @Summary Upload a file
// @Description Upload a file to the server. Allowed types: .png, .jpg, .jpeg, .webp, .pdf, .docx
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param Authorization header string false "Bearer token in format 'Bearer <token>'"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/upload [post]
func AdminUploadFile(kit *kit.Kit) error {
	err := kit.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_form",
		})
	}

	file, handler, err := kit.Request.FormFile("file")
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "file_missing",
		})
	}
	defer file.Close()

	// Validate file type
	allowedExt := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".webp": true,
		".pdf":  true,
		".docx": true,
	}

	ext := filepath.Ext(handler.Filename)
	if !allowedExt[ext] {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid_type",
		})
	}

	// Generate safe unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	filepath := filepath.Join("public/uploads", filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "file_save_error",
		})
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "file_save_failed",
		})
	}

	// Return public URL
	url := fmt.Sprintf("/uploads/%s", filename)

	return kit.JSON(http.StatusOK, map[string]string{
		"filename": filename,
		"url":      url,
	})
}
