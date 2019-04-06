package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/onetwoclimb/internal/server/restapi/operations"
)

type FileType int

const (
	filesPNG FileType = iota
	filesJPEG
)

const (
	fileExtensionPNG  = "png"
	fileExtensionJPEG = "jpeg"
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
		fileName, err = generateFileName(filesPNG)
		if err != nil {
			l.WithError(err).Error("faile to generate filename")
		}

		filePath := fmt.Sprintf("%s/%s.png", h.config.StaticDst, fileName)
		f, err := os.Create(filePath)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to create png file")
			return operations.NewUploadFileInternalServerError()
		}

		if err := EncodePNG(ib, f); err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to save jpeg file")
			f.Close()
			if err := os.Remove(filePath); err != nil {
				l.WithError(errors.WithStack(err)).Errorf("failed to delete file by path: %s", filePath)
			}
			return operations.NewUploadFileInternalServerError()
		}
	case strings.Contains(ct, "image/jpeg"):
		fileName, err = generateFileName(filesJPEG)
		if err != nil {
			l.WithError(err).Error("faile to generate filename")
		}

		filePath := fmt.Sprintf("%s/%s.jpeg", h.config.StaticDst, fileName)
		f, err := os.Create(filePath)
		if err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to create png file")
			return operations.NewUploadFileInternalServerError()
		}

		if err := EncodeJPEG(ib, f); err != nil {
			l.WithError(errors.WithStack(err)).Error("failed to save jpeg file")
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

func (h *Handler) GetDownloadFile(params operations.DownloadFileParams) middleware.Responder {
	if params.ID == "" {
		return operations.NewDownloadFileBadRequest()
	}

	fileName, err := getFileName(params.ID)
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to get file name")
		return operations.NewDownloadFileBadRequest()
	}

	filePath := filepath.Join(h.config.StaticDst, fileName)
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

func generateFileName(fileType FileType) (string, error) {
	return fmt.Sprintf("%s-%d", uuid.NewV4().String(), fileType), nil
}

func getFileType(fileName string) (string, error) {
	t, err := strconv.Atoi(string(fileName[len(fileName)-1]))
	if err != nil {
		return "", nil
	}

	switch FileType(t) {
	case filesPNG:
		return fileExtensionPNG, nil
	case filesJPEG:
		return fileExtensionJPEG, nil
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
