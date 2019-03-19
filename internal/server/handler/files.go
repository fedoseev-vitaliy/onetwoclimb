package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/onetwoclimb/internal/server/restapi/operations"
)

type FileType int

const (
	FILES_PNG FileType = iota
	FILES_JPEG
)

const (
	FILE_EXTENSION_PNG  = "png"
	FILE_EXTENSION_JPEG = "jpeg"
)

func (h *Handler) PostUploadFile(params operations.UploadFileParams) middleware.Responder {
	defer params.File.Close()

	lr := &io.LimitedReader{N: int64(h.config.MaxFileSize), R: params.File}
	ib, err := ioutil.ReadAll(lr)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to read image")
		return operations.NewUploadFileInternalServerError()
	}

	var fileName string
	ct := http.DetectContentType(ib)
	switch {
	case strings.Contains(ct, "image/png"):
		img, err := png.Decode(bytes.NewReader(ib))
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to decode png file")
			return operations.NewUploadFileInternalServerError()
		}

		// Prepare parent image where we want to position child image.
		target := image.NewRGBA(img.Bounds())
		// Draw child image.
		draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

		fileName, err = generateFileName(FILES_PNG)
		if err != nil {
			l.WithError(err).Error("faile to generate filename")
		}

		filePath := fmt.Sprintf("%s/%s.png", h.config.FilesDst, fileName)
		f, err := os.Create(filePath)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to create png file")
			return operations.NewUploadFileInternalServerError()
		}

		err = png.Encode(f, target)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to save file")
			f.Close()
			if err := os.Remove(filePath); err != nil {
				l.WithError(errors.WithStack(err)).Errorf("failed to delete file by path: %s", filePath)
			}
			return operations.NewUploadFileInternalServerError()
		}
	case strings.Contains(ct, "image/jpeg"):
		img, err := jpeg.Decode(bytes.NewReader(ib))
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to decode png file")
			return operations.NewUploadFileInternalServerError()
		}
		// Prepare parent image where we want to position child image.
		target := image.NewRGBA(img.Bounds())
		// Draw child image.
		draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

		fileName, err = generateFileName(FILES_JPEG)
		if err != nil {
			l.WithError(err).Error("faile to generate filename")
		}

		filePath := fmt.Sprintf("%s/%s.jpeg", h.config.FilesDst, fileName)
		f, err := os.Create(filePath)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to create png file")
			return operations.NewUploadFileInternalServerError()
		}

		err = jpeg.Encode(f, target, nil)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to save file")
			f.Close()
			if err := os.Remove(filePath); err != nil {
				l.WithError(errors.WithStack(err)).Errorf("failed to delete file by path: %s", filePath)
			}
			return operations.NewUploadFileInternalServerError()
		}
	default:
		l.WithError(errors.New("unsupported content type")).Errorf("content type:%s", ct)
		return operations.NewUploadFileInternalServerError()
	}

	return operations.NewUploadFileOK().WithPayload(&operations.UploadFileOKBody{ID: &fileName})
}

func generateFileName(fileType FileType) (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return fmt.Sprintf("%s-%d", uuid.String(), fileType), nil
}

func getFileType(fileName string) (string, error) {
	t, err := strconv.Atoi(string(fileName[len(fileName)-1]))
	if err != nil {
		return "", nil
	}

	switch FileType(t) {
	case FILES_PNG:
		return "png", nil
	case FILES_JPEG:
		return FILE_EXTENSION_PNG, nil
	default:
		return "", errors.New(fmt.Sprintf("unsupported file type :%d", t))
	}
}

func getFileName(id string) (string, error) {
	ft, err := getFileType(id)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if ft == "" {
		return id, nil
	}

	return fmt.Sprintf("%s.%s", id, ft), nil
}

// todo implement download
func (h *Handler) GetDownloadFile(params operations.DownloadFileParams) middleware.Responder {
	if params.ID == "" {
		return operations.NewDownloadFileBadRequest()
	}

	fileName, err := getFileName(params.ID)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to get file name")
		return operations.NewDownloadFileBadRequest()
	}

	filePath := filepath.Join(h.config.FilesDst, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		return operations.NewDownloadFileNotFound()
	}

	f, err := os.Open(filePath)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to open file")
		return operations.NewUploadFileInternalServerError()
	}
	defer f.Close()
	fb, err := ioutil.ReadAll(f)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to read file")
		return operations.NewUploadFileInternalServerError()
	}
	return operations.NewDownloadFileOK().WithPayload(strfmt.Base64(fb))
}
