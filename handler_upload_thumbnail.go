package main

import (
	"fmt"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
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

	const maxMemory int = 10 << 20

	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		return
	}

	file, fileheader, err := r.FormFile("thumbnail")
	if err != nil {
		return
	}

	content_type := r.header["Content-Type"]

	imageData, err := io.ReadAll(file)
	if err != nil {
		return
	}

	vdmeta, err := cfg.db.GetVideo(videoID)
	if err != nil {
		return
	}

	newThumbnail := thumbnail{
		data:		imageData,
		mediaType:	content_type,
	}

	videoThumbnails[videoID] = newThumbnail


	respondWithJSON(w, http.StatusOK, struct{}{})


}
