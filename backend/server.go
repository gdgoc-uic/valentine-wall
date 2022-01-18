package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	goValidator "github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"

	"github.com/jmoiron/sqlx"
	goNanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"github.com/dghubble/oauth1"
)

var validator = goValidator.New()

type UserConnection struct {
	UserID      string `db:"user_id" json:"-"`
	Provider    string `db:"provider" json:"provider"`
	Token       string `db:"token" json:"-"`
	TokenSecret string `db:"token_secret" json:"-"`
}

func (uc UserConnection) ToOauth1Token() *oauth1.Token {
	return &oauth1.Token{
		Token:       uc.Token,
		TokenSecret: uc.TokenSecret,
	}
}

type Message struct {
	ID          string    `db:"id" json:"id"`
	RecipientID string    `db:"recipient_id" json:"recipient_id" validate:"required,len=12"`
	Content     string    `db:"content" json:"content" validate:"required,max=240"`
	HasReplied  bool      `db:"has_replied" json:"has_replied"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type SubmittedMessage struct {
	Message
	UID string `db:"submitter_user_id" json:"uid" validate:"required"`
}

type ResponseError struct {
	StatusCode int    `json:"-"`
	WError     error  `json:"-"`
	Message    string `json:"error_message"`
}

func (re *ResponseError) Error() string {
	if re.WError != nil {
		return re.WError.Error()
	} else {
		return re.Message
	}
}

func jsonEncode(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(data)
}

func htmlEncode(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", "text/html")
	if _, err := rw.Write([]byte("<html><body>")); err != nil {
		return err
	}
	if respErr, ok := data.(*ResponseError); ok {
		if _, err := rw.Write([]byte("<p>" + respErr.Message + "</p>")); err != nil {
			return err
		}
	} else if _, err := rw.Write([]byte(fmt.Sprintf("%s", data))); err != nil {
		return err
	}
	if _, err := rw.Write([]byte("</body></html>")); err != nil {
		return err
	}
	return nil
}

func writeRespError(rw http.ResponseWriter, respErr *ResponseError, encoder ...func(http.ResponseWriter, interface{}) error) {
	if len(respErr.Message) == 0 {
		respErr.Message = http.StatusText(respErr.StatusCode)
	}
	log.Println(respErr.Error())
	rw.WriteHeader(respErr.StatusCode)
	if len(encoder) == 0 {
		jsonEncode(rw, respErr)
	} else {
		encoder[0](rw, respErr)
	}
}

func wrapHandler(handler func(http.ResponseWriter, *http.Request) error, encoder ...func(http.ResponseWriter, interface{}) error) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := handler(rw, r); err != nil {
			if respErr, ok := err.(*ResponseError); ok {
				writeRespError(rw, respErr, encoder...)
			} else {
				writeRespError(rw, &ResponseError{
					StatusCode: http.StatusInternalServerError,
					WError:     err,
					Message:    "Something went wrong.",
				}, encoder...)
			}
		}
	})
}

func verifyUser(firebaseApp *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			// validate user
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
				return &ResponseError{
					StatusCode: http.StatusForbidden,
				}
			}

			authClient, err := firebaseApp.Auth(r.Context())
			if err != nil {
				return err
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := authClient.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				return &ResponseError{
					WError:     err,
					StatusCode: http.StatusForbidden,
				}
			}

			ctx := context.WithValue(r.Context(), "authToken", token)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return nil
		})
	}
}

func wrapRequestError(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("resp is nil")
	}
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		log.Println(err2.Error())
	}
	fmt.Printf("%s: body read - %d\n", resp.Request.URL, len(bodyBytes))
	return fmt.Errorf("%s: got status %d, headers: %s, body: %s", resp.Request.URL, resp.StatusCode, resp.Request.Header, string(bodyBytes))
}

func wrapSqlResult(res sql.Result, customErrorMessage ...string) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		errMessage := "Unable to process your submission. Please try again."
		if len(customErrorMessage) != 0 {
			errMessage = customErrorMessage[0]
		}
		return &ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errMessage,
		}
	}
	return nil
}

// func recoverer(next http.Handler) http.Handler {
// 	return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
// 		if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
// 			logEntry := middleware.GetLogEntry(r)
// 			if logEntry != nil {
// 				logEntry.Panic(rvr, debug.Stack())
// 			} else {
// 				middleware.PrintPrettyStack(rvr)
// 			}

// 			return &ResponseError{
// 				StatusCode: http.StatusInternalServerError,
// 			}
// 		}

// 		next.ServeHTTP(w, r)
// 		return nil
// 	})
// }

func main() {
	// TODO:
	store := sessions.NewCookieStore([]byte("TEST_123"))
	store.Options.SameSite = http.SameSiteDefaultMode
	store.Options.HttpOnly = true

	firebaseApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqlx.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(recoverer)
	r.Use(middleware.CleanPath)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendUrl, baseUrl},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		Debug:            targetEnv == "development",
	}))

	r.NotFound(wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		return &ResponseError{
			StatusCode: http.StatusNotFound,
		}
	}))

	r.MethodNotAllowed(wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		return &ResponseError{
			StatusCode: http.StatusMethodNotAllowed,
		}
	}))

	jsonOnly := middleware.AllowContentType("application/json")

	r.Get("/messages", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		messages := []Message{}
		if err := db.Select(&messages, "SELECT id, recipient_id, content, has_replied, created_at, updated_at FROM messages"); err != nil {
			return err
		}
		return json.NewEncoder(rw).Encode(map[string]interface{}{
			"messages": messages,
		})
	}))

	r.With(jsonOnly, verifyUser(firebaseApp)).
		Post("/messages", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := r.Context().Value("authToken").(*auth.Token)

			var submittedMsg SubmittedMessage
			if err := json.NewDecoder(r.Body).Decode(&submittedMsg); err != nil {
				return err
			}

			submittedMsg.CreatedAt = time.Now()
			if token.UID != submittedMsg.UID {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
				}
			}

			// make lastpostinfo an array in order to avoid false positive error
			// when user posts for the first time
			lastPostInfos := []Message{}
			if err := db.Select(&lastPostInfos, "SELECT recipient_id, content, created_at FROM messages WHERE submitter_user_id = ? ORDER BY datetime(created_at) DESC LIMIT 1", submittedMsg.UID); err != nil {
				return err
			}

			if len(lastPostInfos) != 0 {
				lastPostInfo := lastPostInfos[0]
				if submittedMsg.RecipientID == lastPostInfo.RecipientID && submittedMsg.Content == lastPostInfo.Content {
					return &ResponseError{
						StatusCode: http.StatusBadRequest,
						Message:    "You have posted a similar message to a similar recipient.",
					}
				} else if submittedMsg.CreatedAt.Sub(lastPostInfo.CreatedAt) < 10*time.Minute {
					return &ResponseError{
						StatusCode: http.StatusTooManyRequests,
						Message:    "You are being limited to post every 10 minutes.",
					}
				}
			}

			if err := validator.Struct(&submittedMsg); err != nil {
				if validatorErrors, vOk := err.(goValidator.ValidationErrors); vOk {
					errs := []string{}
					for _, er := range validatorErrors {
						errs = append(errs, er.Error())
					}

					rw.WriteHeader(http.StatusUnprocessableEntity)
					return jsonEncode(rw, map[string]interface{}{
						"errors": errs,
					})
				}

				return err
			}

			// generate ID
			id, err := goNanoid.New()
			if err != nil {
				return err
			}

			submittedMsg.ID = id
			res, err := db.NamedExec("INSERT INTO messages (id, recipient_id, content, submitter_user_id) VALUES (:id, :recipient_id, :content, :submitter_user_id)", &submittedMsg)
			if err != nil {
				return err
			}

			if err := wrapSqlResult(res); err != nil {
				return err
			}

			return jsonEncode(rw, map[string]interface{}{
				"message": "Message created successfully",
				"path":    fmt.Sprintf("/messages/%s/%s", submittedMsg.RecipientID, submittedMsg.ID),
			})
		}))

	r.Get("/messages/{recipientId}", wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
		recipientId := chi.URLParam(rr, "recipientId")
		messages := []Message{}
		if err := db.Select(&messages, "SELECT id, recipient_id, content, has_replied, created_at, updated_at FROM messages WHERE recipient_id = ?", recipientId); err != nil {
			return err
		}

		if len(messages) == 0 {
			r.NotFoundHandler().ServeHTTP(rw, rr)
			return nil
		}

		rw.WriteHeader(http.StatusOK)
		return jsonEncode(rw, map[string]interface{}{
			"count":    len(messages),
			"messages": messages,
		})
	}))

	getMessage := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
			recipientId := chi.URLParam(rr, "recipientId")
			messageId := chi.URLParam(rr, "messageId")
			message := Message{}
			if err := db.Get(&message, "SELECT id, recipient_id, content, has_replied, created_at, updated_at FROM messages WHERE id = ? AND recipient_id = ?", messageId, recipientId); err != nil {
				log.Println(err)
				r.NotFoundHandler().ServeHTTP(rw, rr)
				return
			}

			newCtx := context.WithValue(rr.Context(), "gotMessage", message)
			next.ServeHTTP(rw, rr.WithContext(newCtx))
			return
		})
	}

	r.With(getMessage).Get("/messages/{recipientId}/{messageId}", wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
		message := rr.Context().Value("gotMessage").(Message)
		return jsonEncode(rw, message)
	}))

	r.With(jsonOnly, verifyUser(firebaseApp), getMessage).
		Post("/messages/{recipientId}/{messageId}/reply", wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
			// retrieve message
			message := rr.Context().Value("gotMessage").(Message)

			// retrieve token
			token := rr.Context().Value("authToken").(*auth.Token)

			// retrieve twitter connection
			var twitterUserConnection UserConnection
			if err := db.Get(&twitterUserConnection, "SELECT * FROM user_connections WHERE user_id = ? AND provider = ?", token.UID, "twitter"); err != nil {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					WError:     err,
					Message:    "no twitter connection found for user",
				}
			}

			// decode reply payload
			var payload struct {
				Content string `json:"content"`
			}
			if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
				return err
			}

			// commence posting process
			twClient := twitterOauth1Config.Client(oauth1.NoContext, twitterUserConnection.ToOauth1Token())

			// upload image first
			uploadImgBody := &bytes.Buffer{}
			uploadImgData := multipart.NewWriter(uploadImgBody)
			mw, _ := uploadImgData.CreateFormFile("media", "msg.png")

			imageData := &bytes.Buffer{}
			if err := generateImagePNG(imageData, imageTypeTwitter, message); err != nil {
				return err
			}

			_, err := io.Copy(mw, imageData)
			if err != nil {
				return err
			}

			uploadImgData.Close()
			uploadImageResponse, err := twClient.Post("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image", uploadImgData.FormDataContentType(), bytes.NewReader(uploadImgBody.Bytes()))
			if err != nil {
				return err
			}

			defer uploadImageResponse.Body.Close()
			if uploadImageResponse.StatusCode != http.StatusOK {
				return wrapRequestError(uploadImageResponse)
			}

			var uploadImageRespPayload struct {
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

			if err := json.NewDecoder(uploadImageResponse.Body).Decode(&uploadImageRespPayload); err != nil {
				return err
			}

			// tweet
			bodyBuf := &bytes.Buffer{}
			if err := json.NewEncoder(bodyBuf).Encode(map[string]interface{}{
				"text": payload.Content,
				"media": map[string]interface{}{
					"media_ids": []string{uploadImageRespPayload.MediaIDString},
				},
			}); err != nil {
				return &ResponseError{
					StatusCode: http.StatusUnprocessableEntity,
					WError:     err,
				}
			}

			createTweetResp, err := twClient.Post("https://api.twitter.com/2/tweets", "application/json", bytes.NewReader(bodyBuf.Bytes()))
			if err != nil {
				return err
			}
			defer createTweetResp.Body.Close()
			if createTweetResp.StatusCode < 200 || createTweetResp.StatusCode > 299 {
				return wrapRequestError(createTweetResp)
			}

			res, err := db.Exec("UPDATE messages SET has_replied = true WHERE id = ?", message.ID)
			if err != nil {
				return err
			}

			if err := wrapSqlResult(res); err != nil {
				return err
			}

			return jsonEncode(rw, map[string]string{
				"message": "reply success",
			})
		}))

	r.With(verifyUser(firebaseApp)).
		Post("/user/logout_callback", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			session, err := store.Get(r, sessionName)
			if err != nil {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					WError:     err,
				}
			}

			session.Options.MaxAge = -1
			if err := session.Save(r, rw); err != nil {
				return err
			}

			return jsonEncode(rw, map[string]string{
				"message": "logout success",
			})
		}))

	r.With(verifyUser(firebaseApp)).
		Post("/user/login_callback", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := r.Context().Value("authToken").(*auth.Token)
			var associatedData struct {
				AssociatedID string `db:"associated_id" json:"associated_id"`
			}

			if err := db.Get(&associatedData, "SELECT associated_id FROM associated_ids WHERE user_id = ?", token.UID); err != nil {
				log.Println(err)
				// return err
			}

			userConnections := []UserConnection{}
			if err := db.Select(&userConnections, "SELECT * FROM user_connections WHERE user_id = ?", token.UID); err != nil {
				log.Println(err)
				// return err
			}

			session, _ := store.Get(r, sessionName)
			session.Values["uid"] = token.UID
			if err := session.Save(r, rw); err != nil {
				return &ResponseError{
					StatusCode: http.StatusUnprocessableEntity,
					WError:     err,
				}
			}

			return jsonEncode(rw, map[string]interface{}{
				"associated_id":    associatedData.AssociatedID,
				"user_connections": userConnections,
			})
		}))

	r.With(jsonOnly, verifyUser(firebaseApp)).
		Post("/user/connect_id", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := r.Context().Value("authToken").(*auth.Token)
			var submittedData struct {
				AssociatedID string `db:"associated_id" json:"associated_id"`
			}

			if err := json.NewDecoder(r.Body).Decode(&submittedData); err != nil {
				return err
			}

			res, err := db.Exec("INSERT INTO associated_ids (user_id, associated_id) VALUES (?, ?)", token.UID, submittedData.AssociatedID)
			if err != nil {
				return &ResponseError{
					WError:     err,
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "Failure to connect ID to user. Please try again.",
				}
			}

			if err := wrapSqlResult(res, "Failure to connect ID to user. Please try again"); err != nil {
				return err
			}

			return jsonEncode(rw, map[string]string{
				"message": "ID was connected to user successfully.",
			})
		}))

	r.Get("/user/connect_twitter", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		requestToken, _, err := twitterOauth1Config.RequestToken()
		if err != nil {
			return err
		}

		authUrl, err := twitterOauth1Config.AuthorizationURL(requestToken)
		if err != nil {
			return err
		}

		http.Redirect(rw, r, authUrl.String(), http.StatusFound)
		return nil
	}))

	r.Get("/user/twitter_callback", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		// get session cookie
		session, err := store.Get(r, sessionName)
		if err != nil {
			return &ResponseError{
				StatusCode: http.StatusUnauthorized,
				WError:     err,
			}
		}

		uid, ok := session.Values["uid"].(string)
		if !ok {
			return &ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "uid is missing from session",
			}
		}

		requestToken := r.FormValue("oauth_token")
		verifier := r.FormValue("oauth_verifier")
		accessToken, accessSecret, err := twitterOauth1Config.AccessToken(requestToken, "", verifier)
		if err != nil {
			return err
		}

		token := oauth1.NewToken(accessToken, accessSecret)
		newTwitterConnection := UserConnection{
			UserID:      uid,
			Provider:    "twitter",
			Token:       token.Token,
			TokenSecret: token.TokenSecret,
		}

		res, err := db.NamedExec("INSERT INTO user_connections (user_id, provider, token, token_secret) VALUES (:user_id, :provider, :token, :token_secret)", &newTwitterConnection)
		if err != nil {
			return err
		}

		if err := wrapSqlResult(res, "Unable to process Twitter login"); err != nil {
			return err
		}

		scriptJs := `<script type="text/javascript">window.opener.postMessage({message:'twitter connect success', user_connections:[{provider:'twitter'}]}, '` + frontendUrl + `')</script>`
		return htmlEncode(rw, "<p>success</p>"+scriptJs)
	}, htmlEncode))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), r); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Server opened on http://localhost:%d\n", serverPort)
}
