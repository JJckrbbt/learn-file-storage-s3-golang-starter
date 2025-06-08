package main

import (
	"fmt"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	const maxMemory int64 = 10 << 20

	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		return
	}

	file, fileheader, err := r.FormFile("thumbnail")
	if err != nil {
		return
	}

	ext_map := map[string]string{
		"image/apng": "apng",
		"image/jpeg": "jpg",
		"image/gif":  "gif",
		"image/png":  "png",
		"image/svg":  "svg",
	}
	content_type := fileheader.Header.Get("Content-Type")

	file_ext := ext_map[content_type]

	uploadname := fmt.Sprintf("assets/%s.%s", videoID, file_ext)

	upload, err := os.Create(uploadname)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "File not created", err)
		return
	}

	if _, err := io.Copy(upload, file); err != nil {
		log.Fatal(err)
	}

	vdmeta, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Video Not Found", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/%s", cfg.port, uploadname)

	vdmeta.ThumbnailURL = &url

	err = cfg.db.UpdateVideo(vdmeta)
	if err != nil {
		respondWithError(w, http.StatusNotModified, "Update incomplete", err)
		return
	}

	respondWithJSON(w, http.StatusOK, vdmeta)

}
