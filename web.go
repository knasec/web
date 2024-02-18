package web

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

func Login(url string, username string, password string) Responce {
	var answer_body []byte
	req, err := http.NewRequest(http.MethodGet, url, nil)

}

func Get(url string, authType string, token string, expectedCode int) (Responce, error) {
	var http_response Responce
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return http_response, err
	}

	req.Header = authheader(authType, token)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http_response, err
	}

	http_response.Body, err = io.ReadAll(resp.Body)
	http_response.Status = resp.StatusCode
	if err != nil {
		return http_response, err
	}

	if resp.StatusCode != expectedCode {
		return http_response, errors.New(string(http_response.Body))
	}
	return http_response, nil
}

func Post(url string, authType string, token string, body []byte, expectedCode int, options ...string) (Responce, error) {
	var http_response Responce
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(body))
	if err != nil {
		return http_response, err
	}

	req.Header = authheader(authType, token)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return http_response, err
	}

	http_response.Body, err = io.ReadAll(resp.Body)
	http_response.Status = resp.StatusCode
	if err != nil {
		return http_response, err
	}

	if resp.StatusCode != expectedCode {
		return http_response, errors.New(string(http_response.Body))
	}
	return http_response, nil
}

func Upload(url string, authType string, token string, body []byte, expectedCode int, filename string) (Responce, error) {
	var http_response Responce
	upload_body := &bytes.Buffer{}
	upload_writer := multipart.NewWriter(upload_body)
	upload_file, err := upload_writer.CreateFormFile("file", filename)
	if err != nil {

		return http_response, err
	}
	defer upload_writer.Close()
	_, err = io.Copy(upload_file, bytes.NewBuffer(body))
	if err != nil {
		return http_response, err
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(upload_body.Bytes()))
	if err != nil {
		return http_response, err
	}
	req.Header = authheader(authType, token)
	req.Header.Add("Content-Type", upload_writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return http_response, err
	}

	http_response.Body, err = io.ReadAll(resp.Body)
	http_response.Status = resp.StatusCode
	if err != nil {
		return http_response, err
	}

	if resp.StatusCode != expectedCode {
		return http_response, errors.New(string(http_response.Body))
	}
	return http_response, nil

}

func authheader(authType string, token string, options ...string) http.Header {
	switch authType {
	// Default no authorazation
	default:
		return http.Header{}
	// Basic authorization
	case "Basic":
		return http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Basic " + token},
		}
	case "Bearer":
		return http.Header{
			"Authorization": {"Bearer " + token},
			"Content-Type":  {"application/json"},
		}
	case "Jira":
		return http.Header{
			"Authorization":     {"Bearer " + token},
			"Content-Type":      {"application/json"},
			"X-Atlassian-Token": {"no-check"},
		}
	}
}
