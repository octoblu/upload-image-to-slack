package slack

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const slackURL = "https://slack.com/api/files.upload"

// Slack can upload images to slack
type Slack interface {
	// Upload takes file content and upload it up to slack
	Upload(content *bufio.Reader) error
}

type httpSlack struct {
	channel, token string
}

type slackResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

// New constructs a new Slack instance
func New(channel, token string) Slack {
	return &httpSlack{channel, token}
}

func (slack *httpSlack) Upload(content *bufio.Reader) error {
	request, err := slack.createRequest(content)
	if err != nil {
		return err
	}

	return slack.doRequest(request)
}

func (slack *httpSlack) createRequest(content *bufio.Reader) (*http.Request, error) {
	body, contentType, err := slack.createRequestBody(content)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", slackURL, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", contentType)

	return request, nil
}

func (slack *httpSlack) createRequestBody(content *bufio.Reader) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "foo.txt")
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(part, content)
	if err != nil {
		return nil, "", err
	}

	err = writer.WriteField("filename", "foo.txt")
	if err != nil {
		return nil, "", err
	}

	err = writer.WriteField("token", slack.token)
	if err != nil {
		return nil, "", err
	}

	err = writer.WriteField("channels", slack.channel)
	if err != nil {
		return nil, "", err
	}

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func (slack *httpSlack) doRequest(request *http.Request) error {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	return slack.validateResponse(response)
}

func (slack *httpSlack) validateResponse(response *http.Response) error {
	if response.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("Received non 200 from slack: %v, %v", response.StatusCode, string(bodyBytes))
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	parsedBody := &slackResponse{}
	err = json.Unmarshal(body, parsedBody)
	if err != nil {
		return err
	}

	if !parsedBody.Ok {
		return fmt.Errorf("Slack Error: %v", parsedBody.Error)
	}

	return nil
}
