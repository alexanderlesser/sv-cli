package helpers

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexanderlesser/sv-cli/datastore"
	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/alexanderlesser/sv-cli/types"
	"github.com/alexanderlesser/sv-cli/utils/encrypt"
	"github.com/spf13/viper"
)

func GetCSSFiles(p string) ([]types.File, error) {
	// get css path with viper
	cssPath := viper.GetString(constants.CONFIG_CSS_NAME)
	onlyMinified := viper.GetBool(constants.CONFIG_MINIFIED_CSS_NAME)
	// Check if path exists in the path
	assetsPath := filepath.Join(p, cssPath)
	_, err := os.Stat(assetsPath)
	if err != nil {
		return []types.File{}, err
	}

	// List all .css files in the path
	var cssFiles []types.File
	err = filepath.Walk(assetsPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(filePath) == ".css" {

			if onlyMinified {
				fileName := filepath.Base(filePath)
				if strings.Contains(fileName, ".min") {
					file := types.File{
						Name: fileName,
						Path: filePath,
					}
					cssFiles = append(cssFiles, file)
				}
			} else {
				file := types.File{
					Name: filepath.Base(filePath),
					Path: filePath,
				}

				cssFiles = append(cssFiles, file)
			}
		}
		return nil
	})

	if err != nil {
		return []types.File{}, err
	}

	return cssFiles, nil
}

func GetJSFiles(p string) ([]types.File, error) {
	// get js path with viper
	jsPath := viper.GetString(constants.CONFIG_JS_NAME)
	onlyMinified := viper.GetBool(constants.CONFIG_MINIFIED_JS_NAME)
	// Check path exist
	assetsPath := filepath.Join(p, jsPath)
	_, err := os.Stat(assetsPath)
	if err != nil {
		return []types.File{}, err
	}

	// List all .js files in the path
	var jsFiles []types.File
	err = filepath.Walk(assetsPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(filePath) == ".js" {

			if onlyMinified {
				fileName := filepath.Base(filePath)
				if strings.Contains(fileName, ".min") {
					file := types.File{
						Name: fileName,
						Path: filePath,
					}
					jsFiles = append(jsFiles, file)
				}
			} else {
				file := types.File{
					Name: filepath.Base(filePath),
					Path: filePath,
				}

				jsFiles = append(jsFiles, file)
			}
		}
		return nil
	})

	if err != nil {
		return []types.File{}, err
	}

	return jsFiles, nil
}

func GetCSSPath(p string) (string, error) {
	cssPath := viper.GetString(constants.CONFIG_CSS_NAME)
	// onlyMinified := viper.GetBool(constants.CONFIG_MINIFIED_CSS_NAME)
	// Check if path exists in the path
	fullPath := filepath.Join(p, cssPath)
	_, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func GetJSPath(p string) (string, error) {
	jsPath := viper.GetString(constants.CONFIG_JS_NAME)
	// onlyMinified := viper.GetBool(constants.CONFIG_MINIFIED_CSS_NAME)
	// Check if path exists in the path
	fullPath := filepath.Join(p, jsPath)
	_, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func IsCssFile(f types.File) bool {
	minCss := viper.GetBool(constants.CONFIG_MINIFIED_CSS_NAME)

	if strings.HasSuffix(f.Name, ".css") {

		if minCss {
			if strings.Contains(f.Name, ".min") {
				return true
			} else {
				return false
			}
		}

		return true
	}

	return false
}

func DeployFile(record types.Record, file types.File) (types.DeploySuccess, error) {
	username := record.Username
	password, err := encrypt.DecryptPassword(record.Password)

	if err != nil {
		return types.DeploySuccess{}, err
	}

	filePath := file.Path
	routeURL := "https://" + record.Domain + "/rest-api/upload-css/upload"

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return types.DeploySuccess{}, err
	}

	fileName := filepath.Base(filePath)

	// Create the content parameter with base64 encoding
	contentParam := "content=" + base64.StdEncoding.EncodeToString(content) +
		"&fileName=" + fileName

	req, err := http.NewRequest("POST", routeURL, strings.NewReader(contentParam))
	if err != nil {
		return types.DeploySuccess{}, err
	}

	// Add basic authentication headers
	auth := username + ":" + password
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.DeploySuccess{}, err
	}
	defer resp.Body.Close()

	// Handle response
	responseBody := new(strings.Builder)
	_, err = io.Copy(responseBody, resp.Body)
	if err != nil {
		return types.DeploySuccess{}, err
	}

	// Handle response based on status code
	var entry types.Entry
	entry.Name = file.Name
	entry.Time = time.Now().Format("15:04")
	entry.Date = time.Now().Format("2006-01-02")
	entry.ErrorWarning = resp.StatusCode != http.StatusOK // Set to true if status is not OK

	record.Entries = append(record.Entries, entry)

	err = datastore.UpdateRecord(record)

	if err != nil {
		return types.DeploySuccess{}, err
	}

	// Handle response based on status code
	var response types.DeploySuccess

	if resp.StatusCode == http.StatusOK {
		// fmt.Println("Deployed successful")
		response.Success = true
		response.Entry = entry

		return response, nil
	} else {
		// fmt.Println("Deploy failed")
		response.Success = false
		return response, nil
	}
}
