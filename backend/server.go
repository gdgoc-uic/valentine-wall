package main

import (
	"bytes"
	"context"
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
	"google.golang.org/api/option"

	goNanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"github.com/dghubble/oauth1"

	sq "github.com/Masterminds/squirrel"
	"github.com/mailgun/mailgun-go/v4"
)

var messagesPaginator = Paginator{
	OrderKey: "id",
}

type Gift struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Label string `json:"label"`
}

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
	GiftID      *int      `db:"gift_id" json:"gift_id" validate:"omitempty,numeric"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// TODO: Add mail template
func (msg Message) Message(mg *mailgun.MailgunImpl, toRecipientEmail string) *mailgun.Message {
	return mg.NewMessage(
		fmt.Sprintf("Mr. Kupido <mailgun@%s>", mailgunDomain),
		"You've got a message!",
		msg.Content,
		toRecipientEmail,
	)
}

func (msg Message) PendingMessageID() string {
	return msg.ID
}

type RawMessage struct {
	Message
	UID string `db:"submitter_user_id" json:"uid" validate:"required"`
}

type MessageReply struct {
	MessageID string    `db:"message_id" json:"message_id"`
	Content   string    `db:"content" json:"content" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TODO: Add mail template
// TODO: include emojis to mails
func (mr MessageReply) Message(mg *mailgun.MailgunImpl, toRecipientEmail string) *mailgun.Message {
	return mg.NewMessage(
		fmt.Sprintf("Mr. Kupido <mailgun@%s>", mailgunDomain),
		"Your message has received a reply!",
		mr.Content,
		toRecipientEmail,
	)
}

func (mr MessageReply) PendingMessageID() string {
	return mr.MessageID
}

type ResponseError struct {
	StatusCode int    `json:"-"`
	WError     error  `json:"-"`
	Message    string `json:"error_message"`
}

func (re *ResponseError) Error() string {
	if re.WError != nil {
		return re.WError.Error()
	} else if len(re.Message) == 0 {
		return http.StatusText(re.StatusCode)
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

func getAuthToken(r *http.Request, firebaseApp *firebase.App) (*auth.Token, *auth.Client, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, nil, &ResponseError{
			StatusCode: http.StatusForbidden,
		}
	}

	authClient, err := firebaseApp.Auth(r.Context())
	if err != nil {
		return nil, nil, err
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := authClient.VerifyIDToken(r.Context(), idToken)
	if err != nil {
		return nil, nil, &ResponseError{
			WError:     err,
			StatusCode: http.StatusForbidden,
		}
	}

	return token, authClient, nil
}

func verifyUser(firebaseApp *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token, authClient, err := getAuthToken(r, firebaseApp)
			if err != nil {
				return err
			}

			ctx := context.WithValue(r.Context(), "authToken", token)
			ctxWithClient := context.WithValue(ctx, "authClient", authClient)
			next.ServeHTTP(rw, r.WithContext(ctxWithClient))
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

func wrapValidationError(rw http.ResponseWriter, err error) error {
	if validatorErrors, vOk := err.(goValidator.ValidationErrors); vOk {
		errs := []map[string]interface{}{}
		for _, er := range validatorErrors {
			errs = append(errs, map[string]interface{}{
				"field":   er.Field(),
				"message": er.Error(),
			})
		}
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return jsonEncode(rw, map[string]interface{}{
			"error_message": "there were errors when submitting your data.",
			"errors":        errs,
		})
	}
	return err
}

func isConnectedTo(conns []UserConnection, provider string) (int, bool) {
	for i, c := range conns {
		if c.Provider == provider {
			return i, true
		}
	}
	return -1, false
}

func replyViaTwitter(twitterUserConnection UserConnection, message RawMessage, reply MessageReply) error {
	if twitterUserConnection.Provider != "twitter" {
		return fmt.Errorf("invalid provider: expected twitter, got %s", twitterUserConnection.Provider)
	}

	// commence posting process
	twClient := twitterOauth1Config.Client(oauth1.NoContext, twitterUserConnection.ToOauth1Token())

	// upload image first
	uploadImgBody := &bytes.Buffer{}
	uploadImgData := multipart.NewWriter(uploadImgBody)
	mw, _ := uploadImgData.CreateFormFile("media", "msg.png")

	imageData := &bytes.Buffer{}
	if err := generateImagePNG(imageData, imageTypeTwitter, message.Message); err != nil {
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
		"text": reply.Content,
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

	return nil
}

type ReplyFunc func(UserConnection, RawMessage, MessageReply) error

func replyViaEmail(mg *mailgun.MailgunImpl, db *sqlx.DB, authClient *auth.Client) ReplyFunc {
	return func(emailUserConnection UserConnection, message RawMessage, reply MessageReply) error {
		if emailUserConnection.Provider != "email" {
			return fmt.Errorf("invalid provider: expected email, got %s", emailUserConnection.Provider)
		}

		// get sender email
		senderEmail, err := getUserEmailByUID(authClient, message.UID)
		if err != nil {
			return err
		}

		// TODO: send reply later (?)
		if _, _, err := sendEmail(mg, reply, senderEmail); err != nil {
			return err
		}

		return nil
	}
}

func checkProfanity(content string) error {
	if profanityDetector.IsProfane(content) {
		return &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Your submission contains inappropriate content.",
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

func getUserConnections(db *sqlx.DB, uid string) []UserConnection {
	connections := []UserConnection{}
	if err := db.Select(&connections, "SELECT * FROM user_connections WHERE user_id = ?", uid); err != nil {
		log.Println(err)
		// return err
	}
	return connections
}

func getUserEmailByUID(authClient *auth.Client, uid string) (string, error) {
	gotUser, err := authClient.GetUser(context.Background(), uid)
	if err != nil {
		return "", err
	}

	return gotUser.Email, nil
}

func getUserBySID(db *sqlx.DB, authClient *auth.Client, sid string) (*auth.UserRecord, error) {
	var associatedData struct {
		UID string `db:"user_id"`
	}
	if err := db.Get(&associatedData, "SELECT user_id FROM associated_ids WHERE associated_id = ?", sid); err != nil {
		return nil, err
	}
	return authClient.GetUser(context.Background(), associatedData.UID)
}

func main() {
	// TODO:
	store := sessions.NewCookieStore([]byte("TEST_123"))
	store.Options.SameSite = http.SameSiteDefaultMode
	store.Options.HttpOnly = true

	// mailgun
	mg := mailgun.NewMailgun(mailgunDomain, mailgunApiKey)

	// firebase
	opt := option.WithCredentialsFile(gAppCredPath)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	db := initializeDb()
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(recoverer)
	r.Use(middleware.CleanPath)

	// enable cors only on development or when frontend is not the same as base
	if targetEnv == "development" || frontendUrl != baseUrl {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{frontendUrl, baseUrl},
			AllowCredentials: true,
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			Debug:            targetEnv == "development",
		}))
	}

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
	appVerifyUser := verifyUser(firebaseApp)

	r.Get("/gifts", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		return jsonEncode(rw, giftList)
	}))

	getMessagesHandler := wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
		recipientId := chi.URLParam(rr, "recipientId")
		pg := rr.Context().Value("paginator").(Paginator)

		mainQuery, hasQuery := rr.Context().Value("selectQuery").(sq.SelectBuilder)
		if !hasQuery {
			mainQuery = sq.Select()
		}
		mainQuery = mainQuery.Columns("id", "recipient_id", "content", "has_replied", "gift_id", "created_at", "updated_at")
		if len(recipientId) != 0 {
			mainQuery = mainQuery.Where(sq.Eq{"recipient_id": recipientId})
		}

		resp, err := pg.Load(db, "messages", mainQuery, func(r *sqlx.Rows) (interface{}, error) {
			msg := Message{}
			if err := r.StructScan(&msg); err != nil {
				return nil, err
			}
			return msg, nil
		})

		if err != nil {
			return err
		} else if len(recipientId) != 0 && len(resp.Data) == 0 {
			r.NotFoundHandler().ServeHTTP(rw, rr)
			return nil
		}

		return jsonEncode(rw, resp)
	})

	r.With(pagination(messagesPaginator)).Get("/messages", getMessagesHandler)

	customMsgQueryFilters := customSelectFilters(map[string]FilterFunc{
		"has_gift": func(r *http.Request, queryVal string, sb *sq.SelectBuilder) error {
			token, _, _ := getAuthToken(r, firebaseApp)
			if token != nil && (queryVal == "true" || queryVal == "1") {
				*sb = (*sb).Where("gift_id IS NOT NULL")
				return nil
			} else if queryVal != "2" {
				*sb = (*sb).Where("gift_id IS NULL")
			}
			return nil
		},
	})

	r.With(customMsgQueryFilters, pagination(messagesPaginator)).Get("/messages/{recipientId}", getMessagesHandler)
	r.With(jsonOnly, appVerifyUser).
		Post("/messages", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := r.Context().Value("authToken").(*auth.Token)
			authClient := r.Context().Value("authClient").(*auth.Client)

			var submittedMsg RawMessage
			if err := json.NewDecoder(r.Body).Decode(&submittedMsg); err != nil {
				return err
			}

			if err := checkProfanity(submittedMsg.Content); err != nil {
				return err
			}

			if token.UID != submittedMsg.UID {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
				}
			}

			submittedMsg.CreatedAt = time.Now()
			// set to null if none
			if submittedMsg.GiftID != nil && *submittedMsg.GiftID == -1 {
				submittedMsg.GiftID = nil
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

			// validate
			if err := validator.Struct(&submittedMsg); err != nil {
				return wrapValidationError(rw, err)
			}

			// generate ID
			id, err := goNanoid.New()
			if err != nil {
				return err
			}

			submittedMsg.ID = id
			res, err := db.NamedExec("INSERT INTO messages (id, recipient_id, content, gift_id, submitter_user_id) VALUES (:id, :recipient_id, :content, :gift_id, :submitter_user_id)", &submittedMsg)
			if err != nil {
				return err
			}

			if err := wrapSqlResult(res); err != nil {
				return err
			}

			// send email to recipient if available
			recipientUser, err := getUserBySID(db, authClient, submittedMsg.RecipientID)
			// ignore the errors, just pass through
			if err != nil {
				log.Println(err)
			}

			// send the mail within 10 minutes.
			defer addEmailSendEntry("send", mg, submittedMsg.Message, recipientUser.Email)
			return jsonEncode(rw, map[string]interface{}{
				"message": "Message created successfully",
				"path":    fmt.Sprintf("/messages/%s/%s", submittedMsg.RecipientID, submittedMsg.ID),
			})
		}))

	getRawMessage := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
			recipientId := chi.URLParam(rr, "recipientId")
			messageId := chi.URLParam(rr, "messageId")
			message := RawMessage{}
			if err := db.Get(&message, "SELECT * FROM messages WHERE id = ? AND recipient_id = ?", messageId, recipientId); err != nil {
				log.Println(err)
				r.NotFoundHandler().ServeHTTP(rw, rr)
				return
			}

			newCtx := context.WithValue(rr.Context(), "gotMessage", message)
			next.ServeHTTP(rw, rr.WithContext(newCtx))
			return
		})
	}

	r.With(getRawMessage).Get("/messages/{recipientId}/{messageId}", wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
		message := rr.Context().Value("gotMessage").(RawMessage)

		// generate image if ?image query
		if rr.URL.Query().Has("image") {
			return generateImagePNG(rw, imageTypeTwitter, message.Message)
		}

		isUserSenderOrReceiver := false
		reply := &MessageReply{}
		if token, authClient, tErr := getAuthToken(rr, firebaseApp); token != nil {
			gotRecipientUser, _ := getUserBySID(db, authClient, message.RecipientID)
			if token.UID == message.UID || (gotRecipientUser != nil && token.UID == gotRecipientUser.UID) {
				isUserSenderOrReceiver = true
			}
		} else if tErr != nil {
			log.Println(tErr.Error())
		}

		// get reply if possible
		if isUserSenderOrReceiver {
			if message.HasReplied {
				// ignore error
				if err := db.Get(reply, "SELECT * FROM message_replies WHERE message_id = ?", message.ID); err != nil {
					log.Println(err)
				}
			}
		} else if message.GiftID != nil {
			// make notes with gifts limited to sender and receivers only
			return &ResponseError{
				StatusCode: http.StatusForbidden,
			}
		}

		return jsonEncode(rw, map[string]interface{}{
			"message": message,
			"reply":   reply,
		})
	}))

	r.With(appVerifyUser, getRawMessage).Delete("/messages/{recipientId}/{messageId}", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := r.Context().Value("authToken").(*auth.Token)
		message := r.Context().Value("gotMessage").(RawMessage)
		if message.UID != token.UID {
			return &ResponseError{
				StatusCode: http.StatusForbidden,
			}
		}

		res, err := db.Exec("DELETE FROM messages WHERE id = ?", message.ID)
		if err != nil {
			return err
		} else if err := wrapSqlResult(res); err != nil {
			return err
		}

		// cancel pending email job
		jobId := fmt.Sprintf("send_%s", message.PendingMessageID())
		if pendingMailAfterFunc, hasPending := pendingEmailMessages.Load(jobId); hasPending {
			if pendingTimer, ok := pendingMailAfterFunc.(*time.Timer); ok {
				pendingTimer.Stop()
				pendingEmailMessages.Delete(jobId)
			}
		}

		return jsonEncode(rw, map[string]string{
			"message": "message deleted successfully",
		})
	}))

	r.With(jsonOnly, appVerifyUser, getRawMessage).
		Post("/messages/{recipientId}/{messageId}/reply", wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
			// retrieve message
			message := rr.Context().Value("gotMessage").(RawMessage)

			// retrieve token
			token := rr.Context().Value("authToken").(*auth.Token)
			authClient := rr.Context().Value("authClient").(*auth.Client)

			// retrieve connections
			connections := getUserConnections(db, token.UID)
			noConnectionErr := &ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "must be connected either to e-mail or twitter",
			}

			if len(connections) == 0 {
				return noConnectionErr
			}

			// decode reply payload
			var reply MessageReply
			if err := json.NewDecoder(rr.Body).Decode(&reply); err != nil {
				return err
			} else if err := checkProfanity(reply.Content); err != nil {
				return err
			} else if err := validator.Struct(&reply); err != nil {
				return wrapValidationError(rw, err)
			}

			reply.MessageID = message.ID
			if twitterIdx, hasTwitter := isConnectedTo(connections, "twitter"); hasTwitter {
				if err := replyViaTwitter(connections[twitterIdx], message, reply); err != nil {
					return err
				}
			} else if emailIdx, hasEmail := isConnectedTo(connections, "email"); hasEmail {
				if err := replyViaEmail(mg, db, authClient)(connections[emailIdx], message, reply); err != nil {
					return err
				}
			} else {
				// just to be sure
				return noConnectionErr
			}

			updateRes, err := db.NamedExec("INSERT INTO message_replies (message_id, content) VALUES (:message_id, :content)", &reply)
			if err != nil {
				return err
			} else if err := wrapSqlResult(updateRes); err != nil {
				return err
			}

			res, err := db.Exec("UPDATE messages SET has_replied = true WHERE id = ?", message.ID)
			if err != nil {
				return err
			} else if err := wrapSqlResult(res); err != nil {
				return err
			}

			return jsonEncode(rw, map[string]string{
				"message": "reply success",
			})
		}))

	r.With(appVerifyUser).
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

	r.With(appVerifyUser).
		Post("/user/login_callback", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := r.Context().Value("authToken").(*auth.Token)
			var associatedData struct {
				AssociatedID string `db:"associated_id" json:"associated_id"`
			}

			if err := db.Get(&associatedData, "SELECT associated_id FROM associated_ids WHERE user_id = ?", token.UID); err != nil {
				log.Println(err)
				// return err
			}

			userConnections := getUserConnections(db, token.UID)
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

	r.With(jsonOnly, appVerifyUser).
		Post("/user/setup", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := r.Context().Value("authToken").(*auth.Token)
			authClient := r.Context().Value("authClient").(*auth.Client)

			var submittedData struct {
				UID          string `db:"user_id" json:"-"`
				AssociatedID string `db:"associated_id" json:"associated_id"`
				TermsAgreed  bool   `db:"terms_agreed" json:"terms_agreed"`
			}
			if err := json.NewDecoder(r.Body).Decode(&submittedData); err != nil {
				return err
			}

			existingAssoc := struct {
				AssociatedID string `db:"associated_id"`
			}{}
			if err := db.Get(&existingAssoc, "SELECT associated_id FROM associated_ids WHERE user_id = ? OR associated_id = ?", token.UID, submittedData.AssociatedID); err == nil {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					Message:    "You have already registered.",
				}
			}

			if !submittedData.TermsAgreed {
				// delete user
				if err := authClient.DeleteUser(context.Background(), token.UID); err != nil {
					log.Println(err)
				}

				return &ResponseError{
					StatusCode: http.StatusForbidden,
					Message:    "Access to the service is denied.",
				}
			}

			submittedData.UID = token.UID
			res, err := db.NamedExec("INSERT INTO associated_ids (user_id, associated_id) VALUES (:user_id, :associated_id, :terms_agreed)", &submittedData)
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
				"message":       "ID was connected to user successfully.",
				"associated_id": submittedData.AssociatedID,
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

	r.With(appVerifyUser).Get("/user/connect_email", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := r.Context().Value("authToken").(*auth.Token)
		authClient := r.Context().Value("authClient").(*auth.Client)
		userEmail, err := getUserEmailByUID(authClient, token.UID)
		if err != nil {
			return err
		}

		newEmailConnection := UserConnection{
			UserID:      token.UID,
			Provider:    "email",
			Token:       userEmail,
			TokenSecret: "",
		}

		res, err := db.NamedExec("INSERT INTO user_connections (user_id, provider, token, token_secret) VALUES (:user_id, :provider, :token, :token_secret)", &newEmailConnection)
		if err != nil {
			return err
		}

		if err := wrapSqlResult(res, "Unable to connect e-mail."); err != nil {
			return err
		}

		connections := getUserConnections(db, token.UID)
		return jsonEncode(rw, map[string]interface{}{
			"message":          "e-mail connected successfully",
			"user_connections": connections,
		})
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
