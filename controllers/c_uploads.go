package controllers

import (
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/dimaskiddo/frame-go/utils"
)

// Function to Upload a File
func AddUpload(w http.ResponseWriter, r *http.Request) {
	// Limit Body Size with 1 MiB Margin
	r.Body = http.MaxBytesReader(w, r.Body, (utils.Config.GetInt64("SERVER_UPLOAD_LIMIT")+1)*int64(math.Pow(1024, 2)))

	// Get File Content from Multipart Data
	fileUploadContent, fileUploadHeader, err := r.FormFile("fileUpload")
	if err == nil {
		fileUploadName := fileUploadHeader.Filename

		// Create Uploaded File
		fileUploadWrite, err := os.OpenFile(utils.Config.GetString("SERVER_UPLOAD_PATH")+"/"+fileUploadName, os.O_WRONLY|os.O_CREATE, 0666)
		if err == nil {
			var response utils.Response

			// Copy Uploaded File Data from Multipart Data
			io.Copy(fileUploadWrite, fileUploadContent)

			// Close File Handle
			fileUploadContent.Close()
			fileUploadWrite.Close()

			// If Storage Driver Defined Then Re-upload It to Custom Storage
			if len(strings.ToLower(utils.Config.GetString("STORAGE_DRIVER"))) != 0 {
				switch strings.ToLower(utils.Config.GetString("STORAGE_DRIVER")) {
				case "aws", "minio":
					// Try to Upload File to Custom Storage
					err := utils.StoreS3UploadFile(utils.Config.GetString("SERVER_UPLOAD_PATH")+"/"+fileUploadName, r.Header.Get("Content-Type"))
					if err == nil {
						// Set Response Data
						response.Status = true
						response.Code = http.StatusOK
						response.Message = "success"

						// Write Response Data to HTTP
						utils.ResponseWrite(w, response.Code, response)
					} else {
						utils.ResponseInternalError(w, err.Error())
						log.Println(err.Error())
					}

					// Remove Uploaded File
					err := os.Remove(utils.Config.GetString("SERVER_UPLOAD_PATH") + "/" + fileUploadName)
					if err != nil {
						log.Println(err)
					}

				default:
					utils.ResponseInternalError(w, "Invalid storage driver type")
					log.Println("Invalid storage driver type")
				}
			} else {
				// Set Response Data
				response.Status = true
				response.Code = http.StatusOK
				response.Message = "success"

				// Write Response Data to HTTP
				utils.ResponseWrite(w, response.Code, response)
			}
		} else {
			utils.ResponseInternalError(w, err.Error())
			log.Println(err.Error())
		}
	} else {
		utils.ResponseInternalError(w, err.Error())
		log.Println(err.Error())
	}
}
