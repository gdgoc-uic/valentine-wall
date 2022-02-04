package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/dghubble/oauth1"
	poClient "github.com/nedpals/valentine-wall/postal_office/client"
)

type Notifier interface {
	Notify() error
}

var twitterUploadImageUrl = "https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image"
var twitterTweetUrl = "https://api.twitter.com/2/tweets"

func assertConnectionProvider(userConn UserConnection, expectedProvider string) error {
	if userConn.Provider != expectedProvider {
		return fmt.Errorf("invalid provider: expected %s, got %s", expectedProvider, userConn.Provider)
	}
	return nil
}

type TwUploadImageResponse struct {
	MediaID          int    `json:"media_id"`
	MediaIDString    string `json:"media_id_string"`
	MediaKey         string `json:"media_key"`
	Size             int    `json:"size"`
	ExpiresAfterSecs int    `json:"expires_after_secs"`
	Image            struct {
		ImageType string `json:"image_type"`
		Width     int    `json:"w"`
		Height    int    `json:"h"`
	} `json:"image"`
}

type TwitterNotifier struct {
	Connection  UserConnection
	ImageData   io.Reader
	TextContent string
}

// TODO: add hashtag
func (tw *TwitterNotifier) Notify() error {
	if err := assertConnectionProvider(tw.Connection, "twitter"); err != nil {
		return err
	}

	tweetRespBody := map[string]interface{}{
		"text": tw.TextContent,
	}

	// commence posting process
	twClient := twitterOauth1Config.Client(oauth1.NoContext, tw.Connection.ToOauth1Token())

	// upload image first if available
	if tw.ImageData != nil {
		uploadImgBody := &bytes.Buffer{}
		uploadImgData := multipart.NewWriter(uploadImgBody)
		mw, _ := uploadImgData.CreateFormFile("media", "msg.png")

		_, err := io.Copy(mw, tw.ImageData)
		if err != nil {
			return err
		}

		uploadImgData.Close()
		uploadImageResponse, err := twClient.Post(twitterUploadImageUrl, uploadImgData.FormDataContentType(), bytes.NewReader(uploadImgBody.Bytes()))
		if err != nil {
			return err
		}

		defer uploadImageResponse.Body.Close()
		if uploadImageResponse.StatusCode != http.StatusOK {
			return wrapRequestError(uploadImageResponse)
		}

		uploadImageRespPayload := &TwUploadImageResponse{}
		if err := json.NewDecoder(uploadImageResponse.Body).Decode(uploadImageRespPayload); err != nil {
			return err
		}

		tweetRespBody["media"] = map[string]interface{}{
			"media_ids": []string{uploadImageRespPayload.MediaIDString},
		}
	}

	// tweet
	bodyBuf := &bytes.Buffer{}
	if err := json.NewEncoder(bodyBuf).Encode(tweetRespBody); err != nil {
		return &ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			WError:     err,
		}
	}

	createTweetResp, err := twClient.Post(twitterTweetUrl, "application/json", bytes.NewReader(bodyBuf.Bytes()))
	if err != nil {
		return err
	}
	defer createTweetResp.Body.Close()

	if createTweetResp.StatusCode < 200 || createTweetResp.StatusCode > 299 {
		return wrapRequestError(createTweetResp)
	}

	return nil
}

type EmailNotifier struct {
	Client   *poClient.Client
	Template *TemplatedMailSender
	Context  *MailSenderContext
}

func (em *EmailNotifier) Notify() error {
	_, err := newSendJob(
		em.Client,
		em.Template.With(em.Context),
		em.Context.Email,
		em.Context.MessageID,
	)
	return err
}
