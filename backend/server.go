package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	htmlTemplate "html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	goNanoid "github.com/matoous/go-nanoid/v2"

	"github.com/dghubble/oauth1"

	sq "github.com/Masterminds/squirrel"
	"github.com/hako/durafmt"

	poClient "github.com/nedpals/valentine-wall/postal_office/client"
	poTypes "github.com/nedpals/valentine-wall/postal_office/types"
)

var messagesPaginator = Paginator{
	OrderKey: "id",
}

type CollegeDepartment struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Gift struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Label string `json:"label"`
}

type MessageStats struct {
	Messages     int `json:"messages"`
	GiftMessages int `json:"gift_messages"`
}

type AssociatedUser struct {
	UserID       string `db:"user_id" json:"user_id"`
	AssociatedID string `db:"associated_id" json:"associated_id" validate:"required,numeric"`
	TermsAgreed  bool   `db:"terms_agreed" json:"terms_agreed" validate:"required"`
	Department   string `db:"department" json:"department" validate:"required"`
	Gender       string `db:"gender" json:"gender" validate:"required"`
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
	ID          string `db:"id" json:"id"`
	RecipientID string `db:"recipient_id" json:"recipient_id" validate:"required,min=6,max=12,numeric"`
	Content     string `db:"content" json:"content" validate:"required,max=240"`
	HasReplied  bool   `db:"has_replied" json:"has_replied"`
	GiftIDs     []int  `json:"gift_ids" validate:"omitempty,max=3"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
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

func main() {
	var chromeCtx context.Context
	var chromeCancel context.CancelFunc
	htmlTemplates := &htmlTemplate.Template{}

	// chrome/browser-based image rendering specific code
	if len(chromeDevtoolsURL) != 0 {
		// launch chrome instance
		log.Printf("connecting chrome via: %s\n", chromeDevtoolsURL)
		remoteChromeCtx, remoteCtxCancel := chromedp.NewRemoteAllocator(context.Background(), chromeDevtoolsURL)
		defer remoteCtxCancel()

		chromeCtx, chromeCancel = chromedp.NewContext(remoteChromeCtx)
		defer chromeCancel()

		// load template
		log.Println("loading image template...")
		htmlTemplate.Must(htmlTemplates.New("image").ParseFiles("./templates/html/message_image.html.tpl"))
	}

	// postal client
	log.Printf("connecting postal service via %s...\n", postalOfficeAddress)
	postalOfficeClient, err := poClient.DialHTTP(postalOfficeAddress)
	if err != nil {
		log.Println("dialing:", err)
	}

	// load email templates
	rawEmailTemplates := template.Must(template.ParseGlob("./templates/mail/*.txt.tpl"))
	emailTemplates := map[string]*TemplatedMailSender{
		"reply":   newTemplatedMailSender(rawEmailTemplates.Lookup("reply.txt.tpl"), "Mr. Kupido", "Your message has received a reply!", 10*time.Second),
		"message": newTemplatedMailSender(rawEmailTemplates.Lookup("message.txt.tpl"), "Mr. Kupido", "You received a new message!", poTypes.DefaultEmailSendExp),
		"welcome": newTemplatedMailSender(rawEmailTemplates.Lookup("welcome.txt.tpl"), "Mr. Kupido", "Welcome to UIC Valentine Wall 2021!", 10*time.Second),
	}

	// email verification
	log.Println("compiling email regex...")
	emailRegex, err := regexp.Compile(`\A[a-z]+_([0-9]+)@uic.edu.ph\z`)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO:
	log.Println("setting up sessions...")
	store := sessions.NewCookieStore([]byte("TEST_123"))
	store.Options.SameSite = http.SameSiteDefaultMode
	store.Options.HttpOnly = true

	// firebase
	log.Println("connect firebase admin api...")
	opt := option.WithCredentialsFile(gAppCredPath)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("initializing database...")
	db := initializeDb()
	defer db.Close()

	// middlewares
	jsonOnly := middleware.AllowContentType("application/json")
	appVerifyUser := verifyUser(firebaseApp)

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
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
			},
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

	r.Get("/departments", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		return jsonEncode(rw, collegeDepartments)
	}))

	r.Get("/gifts", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		return jsonEncode(rw, giftList)
	}))

	r.Get("/rankings", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		limit := 0
		if gotLimit, exists := r.URL.Query()["limit"]; exists || len(gotLimit) != 0 {
			var err error
			if limit, err = strconv.Atoi(gotLimit[0]); err != nil {
				return err
			} else if limit < 5 {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					Message:    "invalid limit value", // TODO
				}
			}
		}

		// TODO: cache results
		// get number of gift message for each recipient
		giftInnerJoinSubquery := sq.Select("id", "recipient_id").Distinct().From("messages").InnerJoin("message_gifts on message_gifts.message_id = messages.id")
		giftMessagesCountSubquerySQL, gSqArgs, err := sq.Select("recipient_id", "count(*) gift_messages_count").FromSelect(giftInnerJoinSubquery, "msg").GroupBy("recipient_id").ToSql()
		if err != nil {
			return err
		}

		// get number of ordinary messages for each recipient
		messagesCountSubQuery := sq.Select("messages.recipient_id", "count(*) messages_count", "gift_messages_rankings.gift_messages_count").
			From("messages").LeftJoin(fmt.Sprintf("(%s) gift_messages_rankings on gift_messages_rankings.recipient_id = messages.recipient_id", giftMessagesCountSubquerySQL), gSqArgs...).
			GroupBy("messages.recipient_id")

		// get all
		rankingsResults := []struct {
			RecipientID       string `db:"recipient_id" json:"recipient_id"`
			Department        string `db:"department" json:"department"`
			GiftMessagesCount string `db:"gift_messages_count" json:"gift_messages_count"`
			MessagesCount     string `db:"messages_count" json:"messages_count"`
		}{}

		rankingsQuery := sq.Select("recipient_id", "associated_ids.department", "ifnull(rankings.gift_messages_count, 0) gift_messages_count", "ifnull(rankings.messages_count, 0) messages_count").
			FromSelect(messagesCountSubQuery, "rankings").InnerJoin("associated_ids on associated_ids.associated_id = rankings.recipient_id").
			OrderBy("rankings.gift_messages_count desc", "rankings.messages_count desc")

		if limit != 0 {
			rankingsQuery = rankingsQuery.Limit(uint64(limit))
		}

		rankingsQuerySQL, rqArgs, err := rankingsQuery.ToSql()
		if err != nil {
			return err
		} else if err := db.Select(&rankingsResults, rankingsQuerySQL, rqArgs...); err != nil {
			return err
		}

		return jsonEncode(rw, rankingsResults)
	}))

	getMessagesHandler := wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
		recipientId := chi.URLParam(rr, "recipientId")
		pg := rr.Context().Value("paginator").(Paginator)

		baseQuery, okQuery := rr.Context().Value("selectQuery").(sq.SelectBuilder)
		if okQuery {
			baseQuery = baseQuery.From("messages")
		} else {
			baseQuery = sq.Select().From("messages")
		}

		if len(recipientId) != 0 {
			baseQuery = baseQuery.Where(sq.Eq{"recipient_id": recipientId})
		}

		dataQuery := baseQuery.Columns("id", "recipient_id", "content", "has_replied", "created_at", "updated_at")
		resp, err := pg.Load(db, baseQuery, dataQuery, func(r *sqlx.Rows) (interface{}, error) {
			msg := Message{}
			if err := r.StructScan(&msg); err != nil {
				return nil, err
			}
			return msg, nil
		})

		if err != nil {
			return err
		} else if resp.Page > resp.PageCount && len(resp.Data) == 0 {
			r.NotFoundHandler().ServeHTTP(rw, rr)
			return nil
		}

		return jsonEncode(rw, resp)
	})

	customMsgQueryFilters := customSelectFilters(map[string]FilterFunc{
		"has_gift": func(r *http.Request, queryVal string, sb *sq.SelectBuilder) error {
			// TODO: disable_restricted_access_to_gift_messages
			// token, _, err := getAuthToken(r, firebaseApp)
			switch queryVal {
			case "1", "2":
				// if token == nil {
				// 	return &ResponseError{
				// 		WError:     err,
				// 		StatusCode: http.StatusForbidden,
				// 	}
				// }
				// recipientId := chi.URLParam(r, "recipientId")
				// if associatedUser, err := getAssociatedUserBy(db, sq.Eq{"user_id": token.UID}); err == nil && associatedUser.AssociatedID == recipientId {
				if queryVal == "1" {
					*sb = (*sb).Distinct().InnerJoin("message_gifts on message_gifts.message_id = messages.id")
				} else if queryVal == "2" {
					// leave as is
				}
				return nil
				// }
				// if err != nil {
				// 	log.Println(err)
				// }
				// return &ResponseError{
				// 	StatusCode: http.StatusForbidden,
				// }
			default:
				*sb = (*sb).LeftJoin("message_gifts on message_gifts.message_id = messages.id").Where("message_gifts.gift_id IS NULL")
			}
			return nil
		},
	})

	r.With(customMsgQueryFilters, pagination(messagesPaginator)).Get("/messages", getMessagesHandler)
	r.With(customMsgQueryFilters, pagination(messagesPaginator)).Get("/messages/{recipientId}", getMessagesHandler)
	r.Get("/messages/{recipientId}/stats", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		recipientId := chi.URLParam(r, "recipientId")
		stats, err := getMessageStatsBySID(db, recipientId)
		if err != nil {
			return err
		}
		return jsonEncode(rw, stats)
	}))
	r.With(jsonOnly, appVerifyUser).
		Post("/messages", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := r.Context().Value("authToken").(*auth.Token)
			authClient := r.Context().Value("authClient").(*auth.Client)

			var submittedMsg RawMessage
			if err := json.NewDecoder(r.Body).Decode(&submittedMsg); err != nil {
				return err
			} else if err := checkProfanity(submittedMsg.Content); err != nil {
				return err
			} else if token.UID != submittedMsg.UID {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
				}
			}

			submittedMsg.CreatedAt = time.Now()

			// make lastpostinfo an array in order to avoid false positive error
			// when user posts for the first time
			lastPostInfos := []Message{}
			if err := db.Select(&lastPostInfos, "SELECT recipient_id, content, created_at FROM messages WHERE submitter_user_id = ? ORDER BY datetime(created_at) DESC LIMIT 1", submittedMsg.UID); err != nil {
				return err
			}

			if len(lastPostInfos) != 0 {
				lastPostInfo := lastPostInfos[0]
				timeToSend := emailTemplates["message"].TimeToSend()
				if submittedMsg.RecipientID == lastPostInfo.RecipientID && submittedMsg.Content == lastPostInfo.Content {
					return &ResponseError{
						StatusCode: http.StatusBadRequest,
						Message:    "You have posted a similar message to a similar recipient.",
					}
				} else if diff := submittedMsg.CreatedAt.Sub(lastPostInfo.CreatedAt); diff < timeToSend {
					fmtDuration := durafmt.Parse(timeToSend - diff).LimitFirstN(1)
					return &ResponseError{
						StatusCode: http.StatusTooManyRequests,
						Message:    fmt.Sprintf("You have %s left before you can post again.", fmtDuration),
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
			tx, err := db.BeginTxx(r.Context(), &sql.TxOptions{})
			if err != nil {
				return err
			}

			res, err := tx.NamedExec("INSERT INTO messages (id, recipient_id, content, submitter_user_id) VALUES (:id, :recipient_id, :content, :submitter_user_id)", &submittedMsg)
			if err != nil {
				log.Println(tx.Rollback())
				return err
			}

			if err := wrapSqlResult(res); err != nil {
				log.Println(tx.Rollback())
				return err
			}

			for _, giftId := range submittedMsg.GiftIDs {
				if res, err := tx.Exec("INSERT INTO message_gifts (message_id, gift_id) VALUES (?, ?)", submittedMsg.ID, giftId); err != nil {
					log.Println(tx.Rollback())
					return err
				} else if err := wrapSqlResult(res); err != nil {
					log.Println(tx.Rollback())
					return err
				}
			}

			if err := tx.Commit(); err != nil {
				return err
			}

			// send email to recipient if available
			recipientUser, err := getUserBySID(db, authClient, submittedMsg.RecipientID)
			// ignore the errors, just pass through
			if err != nil {
				log.Println(err)
			} else {
				// send the mail within n minutes.
				sender := emailTemplates["message"].With(submittedMsg.Message)
				if _, err := newSendJob(postalOfficeClient, sender, recipientUser.Email, submittedMsg.ID); err != nil {
					log.Println(err)
				}
			}

			return jsonEncode(rw, map[string]interface{}{
				"message": "Message created successfully",
				"route": map[string]interface{}{
					"name":   "message-page",
					"params": map[string]string{"recipientId": submittedMsg.RecipientID, "messageId": submittedMsg.ID},
					"query":  map[string]string{"from": "send_message_modal"},
				},
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

			if err := db.Select(&message.GiftIDs, "SELECT gift_id FROM message_gifts WHERE message_id = ?", messageId); err != nil {
				log.Println(err)
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
			if len(chromeDevtoolsURL) != 0 {
				return generateImagePNGChrome(rw, chromeCtx, htmlTemplates.Lookup("image"), message.Message)
			} else {
				return generateImagePNG(rw, imageTypeTwitter, message.Message)
			}
		}

		isDeletable := false
		isUserSenderOrReceiver := false
		reply := &MessageReply{}
		if token, authClient, tErr := getAuthToken(rr, firebaseApp); token != nil {
			gotRecipientUser, _ := getUserBySID(db, authClient, message.RecipientID)
			if token.UID == message.UID || (gotRecipientUser != nil && token.UID == gotRecipientUser.UID) {
				isUserSenderOrReceiver = true
			}

			timeToSend := emailTemplates["message"].TimeToSend()
			if token.UID == message.UID && time.Now().Sub(message.CreatedAt) < timeToSend {
				isDeletable = true
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
		} else if len(message.GiftIDs) != 0 {
			// make notes with gifts limited to sender and receivers only
			// TODO: disable_restricted_access_to_gift_messages
			// return &ResponseError{
			// 	StatusCode: http.StatusForbidden,
			// }
		}

		return jsonEncode(rw, map[string]interface{}{
			"is_deletable": isDeletable,
			"message":      message,
			"reply":        reply,
		})
	}))

	r.With(appVerifyUser, getRawMessage).Delete("/messages/{recipientId}/{messageId}", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := r.Context().Value("authToken").(*auth.Token)
		message := r.Context().Value("gotMessage").(RawMessage)
		timeToSend := emailTemplates["message"].TimeToSend()

		if message.UID != token.UID || time.Now().Sub(message.CreatedAt) >= timeToSend {
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

		// cancel send job if possible
		if err := postalOfficeClient.CancelJobByUID(message.ID); err != nil {
			log.Println(err)
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
			var notifier Notifier
			if twitterIdx, hasTwitter := isConnectedTo(connections, "twitter"); hasTwitter {
				imageData := &bytes.Buffer{}
				if err := generateImagePNG(imageData, imageTypeTwitter, message.Message); err != nil {
					return err
				}

				notifier = &TwitterNotifier{
					Connection:  connections[twitterIdx],
					ImageData:   imageData,
					TextContent: reply.Content,
				}
			} else if _, hasEmail := isConnectedTo(connections, "email"); hasEmail {
				// get sender email
				senderEmail, err := getUserEmailByUID(authClient, message.UID)
				if err != nil {
					return err
				}

				notifier = &EmailNotifier{
					PostalOfficeClient: postalOfficeClient,
					Template:           emailTemplates["reply"].With(reply),
					RecipientEmail:     senderEmail,
					MessageID:          reply.MessageID,
				}
			} else {
				// just to be sure
				return noConnectionErr
			}

			if err := notifier.Notify(); err != nil {
				return err
			} else if updateRes, err := db.NamedExec("INSERT INTO message_replies (message_id, content) VALUES (:message_id, :content)", &reply); err != nil {
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
			associatedData, err := getAssociatedUserBy(db, sq.Eq{"user_id": token.UID})
			if err != nil {
				log.Println(err)
				// return err
			}

			if associatedData == nil {
				associatedData = &AssociatedUser{}
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
			shouldDenyService := false
			token := r.Context().Value("authToken").(*auth.Token)
			authClient := r.Context().Value("authClient").(*auth.Client)
			submittedData := AssociatedUser{}
			if err := json.NewDecoder(r.Body).Decode(&submittedData); err != nil {
				return err
			}

			if err := validator.Struct(&submittedData); err != nil {
				return wrapValidationError(rw, err)
			}

			if userEmail, err := getUserEmailByUID(authClient, token.UID); err == nil {
				matches := emailRegex.FindAllStringSubmatch(userEmail, -1)
				// deny service if no matching ID found when scanning the email through regex
				if (matches == nil || len(matches) == 0) || (len(matches[0]) < 2 || len(matches[0][1]) == 0) {
					shouldDenyService = true
				} else if gotId := matches[0][1]; gotId != submittedData.AssociatedID {
					shouldDenyService = true
				}
			} else if !submittedData.TermsAgreed {
				shouldDenyService = true
			}

			if shouldDenyService {
				if targetEnv == "production" {
					// delete user
					if err := authClient.DeleteUser(context.Background(), token.UID); err != nil {
						log.Println(err)
					}
				}

				return &ResponseError{
					StatusCode: http.StatusForbidden,
					Message:    "Access to the service is denied.",
				}
			}

			if _, err := getAssociatedUserBy(db, sq.Or{sq.Eq{"user_id": token.UID, "associated_id": submittedData.AssociatedID}}); err == nil {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					Message:    "You have already registered.",
				}
			}

			submittedData.UserID = token.UID
			res, err := db.NamedExec("INSERT INTO associated_ids (user_id, associated_id, terms_agreed) VALUES (:user_id, :associated_id, :terms_agreed)", &submittedData)
			if err != nil {
				return &ResponseError{
					WError:     err,
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "Failed to connect ID to user. Please try again.",
				}
			}

			if err := wrapSqlResult(res, "Failed to connect ID to user. Please try again"); err != nil {
				return err
			}

			// generate welcome message
			if userEmail, err := getUserEmailByUID(authClient, token.UID); err == nil {
				emailId, _ := goNanoid.New()
				associatedUser, err := getAssociatedUserByEmail(db, authClient, userEmail)
				if err != nil {
					log.Println(err)
				}

				stats, err := getMessageStatsBySID(db, associatedUser.AssociatedID)
				if err != nil {
					log.Println(err)
				}

				sender := emailTemplates["welcome"].With(struct {
					Email string
					Stats *MessageStats
				}{
					Email: userEmail,
					Stats: stats,
				})
				if _, err := newSendJob(postalOfficeClient, sender, userEmail, emailId); err != nil {
					log.Println(err)
				}
			} else {
				log.Println(err)
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

	r.With(appVerifyUser).Get("/user/delete", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		authClient := r.Context().Value("authClient").(*auth.Client)
		token := r.Context().Value("authToken").(*auth.Token)

		var confirmationData struct {
			InputSID string `json:"input_sid"`
			InputUID string `json:"input_uid"`
		}

		if err := json.NewDecoder(r.Body).Decode(&confirmationData); err != nil {
			return err
		} else if confirmationData.InputUID != token.UID {
			return &ResponseError{
				StatusCode: http.StatusForbidden,
				WError:     fmt.Errorf("input uid mismatched"),
				Message:    "Unable to delete account.",
			}
		}

		gotAssociatedUser, err := getAssociatedUserBy(db, sq.Eq{"associated_id": confirmationData.InputSID})
		if err != nil {
			return &ResponseError{
				StatusCode: http.StatusForbidden,
				WError:     err,
				Message:    "Unable to delete account.",
			}
		} else if gotAssociatedUser.UserID != confirmationData.InputUID {
			return &ResponseError{
				StatusCode: http.StatusForbidden,
				WError:     fmt.Errorf("input uid mismatched"),
				Message:    "Unable to delete account.",
			}
		}

		// delete from associated_ids and user_connections
		tx, err := db.BeginTxx(r.Context(), &sql.TxOptions{})
		if err != nil {
			return err
		}

		for _, tableName := range []string{"user_connections", "associated_ids"} {
			if deleteSql, deleteArgs, err := sq.Delete(tableName).Where(sq.Eq{"user_id": token.UID}).ToSql(); err != nil {
				return err
			} else if res, err := tx.Exec(deleteSql, deleteArgs...); err != nil {
				return err
			} else if err := wrapSqlResult(res); err != nil {
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
			}
			return err
		} else if err := authClient.DeleteUser(r.Context(), token.UID); err != nil {
			return err
		}

		return jsonEncode(rw, map[string]interface{}{
			"message": "user deleted successfully",
		})
	}))

	log.Printf("Server opened on http://localhost:%d\n", serverPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), r); err != nil {
		log.Fatalln(err)
	}
}
