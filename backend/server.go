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
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/chromedp/chromedp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/newrelic/go-agent/v3/newrelic"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go/v4"
	goNanoid "github.com/matoous/go-nanoid/v2"

	"github.com/dghubble/oauth1"

	sq "github.com/Masterminds/squirrel"
	"github.com/hako/durafmt"

	poClient "github.com/nedpals/valentine-wall/postal_office/client"

	"github.com/patrickmn/go-cache"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type gotMessageKey struct{}

var messagesPaginator = &Paginator{
	OrderKey: "id",
}

type CollegeDepartment struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Gift struct {
	ID    int     `json:"id"`
	UID   string  `json:"uid"`
	Label string  `json:"label"`
	Price float32 `json:"price"`
}

type Gifts []Gift

func (gs Gifts) GetPriceByID(id int) float32 {
	for _, g := range gs {
		if g.ID == id {
			return g.Price
		}
	}
	return 0
}

type ReturnedStringID struct {
	ID string `db:"id"`
}

type RecipientStats struct {
	RecipientID       string `db:"recipient_id" json:"recipient_id"`
	MessagesCount     int    `db:"messages_count" json:"messages_count"`
	GiftMessagesCount int    `db:"gift_messages_count" json:"gift_messages_count"`
}

type RecipientStats2 struct {
	RecipientID string  `db:"recipient_id" json:"recipient_id"`
	Department  string  `db:"department" json:"department"`
	Sex         string  `db:"sex" json:"sex"`
	TotalCoins  float32 `db:"-" json:"total_coins"`
}

type Recipients []*RecipientStats2

func (a Recipients) Len() int           { return len(a) }
func (a Recipients) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Recipients) Less(i, j int) bool { return a[i].TotalCoins > a[j].TotalCoins }

func (a Recipients) BySex(sex string) Recipients {
	if len(sex) == 0 || sex == "all" {
		return a
	}
	res := make(Recipients, 0, len(a))
	for _, v := range a {
		if v.Sex == sex {
			res = append(res, v)
		}
	}
	return res
}

func fetchRecipientRankings(db *sqlx.DB) (Recipients, error) {
	recipientsMap := map[string]*RecipientStats2{}

	// get all associated_users data
	rows, err := db.Query("SELECT associated_id, department, sex FROM associated_users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var associatedId string
		var collegeDepartment string
		var sex string
		rows.Scan(&associatedId, &collegeDepartment, &sex)
		if len(associatedId) == 0 {
			continue
		} else if _, exists := recipientsMap[associatedId]; !exists {
			recipientsMap[associatedId] = &RecipientStats2{
				RecipientID: associatedId,
			}
		}
		recipientsMap[associatedId].Department = collegeDepartment
		recipientsMap[associatedId].Sex = sex
	}

	// get gift id and it's associated recipient
	giftMessagesCountQuerySQL := "SELECT recipient_id, gift_id FROM messages LEFT JOIN message_gifts mg on mg.message_id = messages.id"
	recipientsWithGiftsRows, err := db.Query(giftMessagesCountQuerySQL)
	if err != nil {
		return nil, err
	}
	defer recipientsWithGiftsRows.Close()

	for recipientsWithGiftsRows.Next() {
		var recipientId string
		giftId := 0

		recipientsWithGiftsRows.Scan(&recipientId, &giftId)
		if len(recipientId) == 0 {
			continue
		} else if _, exists := recipientsMap[recipientId]; !exists {
			recipientsMap[recipientId] = &RecipientStats2{
				RecipientID: recipientId,
				Sex:         "unknown",
				Department:  "Unknown",
			}
		}

		if giftId > 0 {
			recipientsMap[recipientId].TotalCoins += giftList[giftId-1].Price
		}

		recipientsMap[recipientId].TotalCoins += sendPrice
	}

	// sort and get results
	results := make(Recipients, 0, len(recipientsMap))
	for _, val := range recipientsMap {
		results = append(results, val)
	}

	sort.Sort(results)
	return results, nil
}

type AssociatedUser struct {
	UserID       string     `db:"user_id" json:"user_id"`
	AssociatedID string     `db:"associated_id" json:"associated_id" validate:"required,numeric"`
	TermsAgreed  bool       `db:"terms_agreed" json:"terms_agreed"`
	Department   string     `db:"department" json:"department" validate:"required"`
	Sex          string     `db:"sex" json:"sex" validate:"required"`
	LastActiveAt *time.Time `db:"last_active_at" json:"-"`
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

type MessageSearchEntry struct {
	MessageID   string    `db:"id" json:"-"`
	RecipientID string    `db:"recipient_id" json:"recipient_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	HasGifts    bool      `db:"has_gifts" json:"has_gifts"`
}

type Message struct {
	ID          string `db:"id" json:"id"`
	RecipientID string `db:"recipient_id" json:"recipient_id" validate:"required,min=6,max=12,numeric"`
	Content     string `db:"content" json:"content" validate:"required,max=240"`
	HasReplied  bool   `db:"has_replied" json:"has_replied"`
	HasGifts    bool   `db:"has_gifts" json:"has_gifts"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

var recentMessagesChan chan Message
var messagesCol = []string{"id", "recipient_id", "content", "has_replied", "has_gifts", "created_at", "updated_at"}

type RawMessage struct {
	Message
	UID       string     `db:"submitter_user_id" json:"uid" validate:"required"`
	GiftIDs   []int      `json:"gift_ids,omitempty" validate:"omitempty,max=3"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
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

func setupMessagesSearchIndex() (bleve.Index, error) {
	idxMapping := bleve.NewIndexMapping()
	idxMapping.DefaultAnalyzer = "en"
	mapping := bleve.NewDocumentMapping()

	mapping.AddFieldMappingsAt("id", bleve.NewTextFieldMapping())
	mapping.AddFieldMappingsAt("recipient_id", bleve.NewTextFieldMapping())
	mapping.AddFieldMappingsAt("has_gifts", bleve.NewBooleanFieldMapping())
	mapping.AddFieldMappingsAt("created_at", bleve.NewDateTimeFieldMapping())
	mapping.AddFieldMappingsAt("updated_at", bleve.NewDateTimeFieldMapping())
	idxMapping.AddDocumentMapping("message", mapping)

	return NewSearch("messages", idxMapping)
}

func importExistingMessages(db *sqlx.DB, index bleve.Index) error {
	batch := index.NewBatch()
	sqlQuery, _, _ := psql.
		Select("id", "recipient_id", "created_at", "updated_at", "has_gifts").
		Where(sq.Eq{"deleted_at": nil}).
		From("messages").
		ToSql()

	rows, err := db.Queryx(sqlQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var result MessageSearchEntry
		if err := rows.StructScan(&result); err != nil {
			return err
		} else if err := batch.Index(result.MessageID, result); err != nil {
			return err
		}
	}

	return index.Batch(batch)
}

func main() {
	var chromeCtx context.Context
	var chromeCancel context.CancelFunc
	htmlTemplates := &htmlTemplate.Template{}
	recentMessagesChan = make(chan Message, 10)
	defer close(recentMessagesChan)

	// cache
	cacher := cache.New(1*time.Hour, 1*time.Minute)

	// image renderer
	imageRenderer := &ImageRenderer{
		CacheStore: cacher,
	}

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
		var err error
		if htmlTemplates, err = htmlTemplate.ParseGlob("./templates/html/*.html.tpl"); err != nil {
			log.Panicln(err)
		} else {
			log.Printf("%d html templates have been loaded\n", len(htmlTemplates.Templates()))
		}

		imageRenderer.ChromeCtx = chromeCtx
		imageRenderer.Template = htmlTemplates.Lookup("message_image.html.tpl")
	}

	// postal client
	log.Printf("connecting postal service via %s...\n", postalOfficeAddress)
	postalOfficeClient, err := poClient.DialHTTP(postalOfficeAddress)
	if err != nil {
		log.Println("dialing:", err)
	}

	// load email templates
	log.Println("loading email templates...")
	rawEmailTemplates := template.Must(template.ParseGlob("./templates/mail/*.txt.tpl"))
	log.Printf("%d email templates have been loaded\n", len(rawEmailTemplates.Templates()))
	emailTemplates := map[string]*TemplatedMailSender{
		"reply":   newTemplatedMailSender(rawEmailTemplates.Lookup("reply.txt.tpl"), "Mr. Kupido", "Your message has received a reply!", 10*time.Second),
		"message": newTemplatedMailSender(rawEmailTemplates.Lookup("message.txt.tpl"), "Mr. Kupido", "You received a new message!", 10*time.Second),
		"welcome": newTemplatedMailSender(rawEmailTemplates.Lookup("welcome.txt.tpl"), "Mr. Kupido", "Welcome to UIC Valentine Wall 2022!", 10*time.Second),
	}

	// email verification
	log.Println("compiling email regex...")
	emailRegex, err := regexp.Compile(`\A[a-z]+_([0-9]+)@uic.edu.ph\z`)
	if err != nil {
		log.Panicln(err)
	}

	// TODO:
	log.Println("setting up sessions...")
	store := sessions.NewCookieStore([]byte("TEST_123"))
	store.Options.SameSite = http.SameSiteDefaultMode
	store.Options.HttpOnly = true

	// firebase
	log.Println("connecting to firebase admin api...")
	opt := option.WithCredentialsFile(gAppCredPath)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Panicln(err)
	}

	// database
	log.Println("initializing database...")
	db := initializeDb()
	defer db.Close()

	// search engine
	log.Println("indexing existing messages to search database...")
	messagesSearchIndex, err := setupMessagesSearchIndex()
	if err != nil {
		log.Panicln(err)
	} else if err := importExistingMessages(db, messagesSearchIndex); err != nil {
		log.Panicln(err)
	}

	// virtual currency system (aka virtual bank)
	b := &VirtualBank{DB: db}
	if err := b.AddInitialAmountToExistingAccounts(firebaseApp); err != nil {
		log.Println(err)
	}

	appCheckBalance := checkBalance(b)

	// invitation system
	invSys := &InvitationSystem{
		DB:         db,
		CookieName: invitationCookieName,
	}

	// middlewares
	jsonOnly := middleware.AllowContentType("application/json")
	appVerifyUser := verifyUser(firebaseApp)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(recoverer)
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
				http.MethodPatch,
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

	rankingPaginator := &Paginator{}
	r.With(pagination(rankingPaginator), customFilters(map[string]FilterFunc{
		"sex": func(r *http.Request, c context.Context, f Filter) error {
			if !f.Exists {
				return nil
			} else if f.Value != "male" && f.Value != "female" && f.Value != "unknown" && f.Value != "all" {
				return &ResponseError{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "sex query should be either male, female, unknown, or all",
				}
			}
			return nil
		},
	})).Get("/rankings", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		pg := getPaginatorFromReq(r)
		cachedRankings, hasRankingsCached := cacher.Get("rankings")
		var results Recipients
		if !hasRankingsCached {
			results, err = fetchRecipientRankings(db)
			if err != nil {
				return err
			}

			// timeToCache := 30 * time.Minute
			// if targetEnv == "development" {
			// timeToCache := 30 * time.Second
			timeToCache := 5 * time.Second
			// }

			cacher.Set("rankings", results, timeToCache)
		} else if recs, ok := cachedRankings.(Recipients); ok {
			results = recs
		}

		if r.URL.Query().Has("sex") {
			filterBySex := r.URL.Query().Get("sex")
			results = results.BySex(filterBySex)
		}

		resp, err := pg.Load(&ArrayPaginatorSource{results})
		if err != nil {
			return err
		}

		return jsonEncode(rw, resp)
	}))

	getMessagesHandler := wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
		recipientId := chi.URLParam(rr, "recipientId")
		pg := getPaginatorFromReq(rr)
		searchRequest := rr.Context().Value(searchRequestKey{}).(*bleve.SearchRequest)
		if len(recipientId) != 0 {
			matchRecipientQuery := bleve.NewTermQuery(recipientId)
			matchRecipientQuery.SetField("recipient_id")
			searchRequest.Query.(*query.ConjunctionQuery).AddQuery(matchRecipientQuery)
		} else if conjQ, isConjQuery := searchRequest.Query.(*query.ConjunctionQuery); isConjQuery && len(conjQ.Conjuncts) == 0 {
			searchRequest.Query = bleve.NewMatchAllQuery()
		}
		resp, err := pg.Load(&PipePaginatorSource{
			Source: &BlevePaginatorSource{
				Index:         messagesSearchIndex,
				SearchRequest: searchRequest,
			},
			PipeFunc: func(inputs []interface{}) ([]interface{}, error) {
				ids := sq.Or{}
				for _, res := range inputs {
					docMatch, ok := res.(*search.DocumentMatch)
					if !ok {
						return nil, fmt.Errorf("not a document match")
					}
					ids = append(ids, sq.Eq{"id": docMatch.ID})
				}
				querySql, queryArgs, err := psql.
					Select(messagesCol...).
					From("messages").Where(sq.And{ids, sq.Eq{"deleted_at": nil}}).
					OrderBy(fmt.Sprintf("%s %s", pg.OrderKey, pg.Order)).ToSql()
				if err != nil {
					return nil, err
				}

				rows, err := db.Queryx(querySql, queryArgs...)
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				results := []interface{}{}
				for rows.Next() {
					msg := Message{}
					if err := rows.StructScan(&msg); err != nil {
						return nil, err
					}
					results = append(results, msg)
				}
				return results, nil
			},
		})
		if err != nil {
			return err
		} else if resp.Page > resp.PageCount && resp.Len == 0 {
			r.NotFoundHandler().ServeHTTP(rw, rr)
			return nil
		}
		return jsonEncode(rw, resp)
	})

	customMsgQueryFilters := customFilters(map[string]FilterFunc{
		"has_gift": func(r *http.Request, ctx context.Context, filter Filter) error {
			searchReq := ctx.Value(searchRequestKey{}).(*bleve.SearchRequest)
			hasGiftQuery := bleve.NewBoolFieldQuery(false)
			hasGiftQuery.SetField("has_gifts")

			// TODO: disable_restricted_access_to_gift_messages
			token, _, err := getAuthToken(r, firebaseApp)
			if filter.Value == "1" || filter.Value == "2" {
				if token == nil {
					return &ResponseError{
						WError:     err,
						StatusCode: http.StatusForbidden,
					}
				}
				recipientId := chi.URLParam(r, "recipientId")
				associatedUser, err := getAssociatedUserBy(db, sq.Eq{"user_id": token.UID})
				if err != nil {
					return &ResponseError{
						WError:     err,
						StatusCode: http.StatusForbidden,
					}
				}

				if err == nil && associatedUser.AssociatedID == recipientId {
					if filter.Value == "1" {
						hasGiftQuery.Bool = true
					} else if filter.Value == "2" {
						// leave as is
						return nil
					}
				}
			}

			searchReq.Query.(*query.ConjunctionQuery).AddQuery(hasGiftQuery)
			return nil
		},
	})

	r.Get("/recent-messages", func(rw http.ResponseWriter, r *http.Request) {
		tx := newrelic.FromContext(r.Context())
		tx.Ignore()

		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Connection", "keep-alive")

		existingEntriesChan := make(chan Message, 10)

		go func() {
			mSql, mArgs, _ := psql.Select(messagesCol...).
				From("messages").Limit(12).OrderBy("created_at ASC").
				Where(sq.And{sq.Eq{"deleted_at": nil, "has_gifts": false}, sq.LtOrEq{"created_at": time.Now()}}).ToSql()

			rows, err := db.Queryx(mSql, mArgs...)
			if err != nil {
				log.Println(err)
				return
			}
			defer rows.Close()

			for rows.Next() {
				msg := Message{}
				if err := rows.StructScan(&msg); err != nil {
					log.Println(err)
				}
				existingEntriesChan <- msg
			}
		}()
		defer close(existingEntriesChan)

		for {
			select {
			case entry := <-existingEntriesChan:
				encodeDataSSE(rw, entry)
			case entry2 := <-recentMessagesChan:
				encodeDataSSE(rw, entry2)
			case <-r.Context().Done():
				return
			}
		}
	})

	messageListMiddlewares := []func(http.Handler) http.Handler{injectSearchQuery, customMsgQueryFilters, pagination(messagesPaginator)}
	r.With(messageListMiddlewares...).Get("/messages", getMessagesHandler)
	r.With(messageListMiddlewares...).Get("/messages/{recipientId}", getMessagesHandler)
	r.Get("/messages/{recipientId}/stats", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		recipientId := chi.URLParam(r, "recipientId")
		stats, err := getRecipientStatsBySID(messagesSearchIndex, recipientId)
		if err != nil {
			return err
		}
		return jsonEncode(rw, stats)
	}))
	r.With(jsonOnly, appVerifyUser, appCheckBalance).
		Post("/messages", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := getAuthTokenByReq(r)
			authClient := getAuthClientByReq(r)

			var submittedMsg RawMessage
			if err := json.NewDecoder(r.Body).Decode(&submittedMsg); err != nil {
				return err
			}

			if token.UID != submittedMsg.UID {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
				}
			} else if err := checkProfanity(submittedMsg.Content); err != nil {
				return err
			}

			submittedMsg.CreatedAt = time.Now()

			// make lastpostinfo an array in order to avoid false positive error
			// when user posts for the first time
			lastPostInfos := []Message{}
			if err := db.Select(&lastPostInfos, "SELECT recipient_id, content, created_at FROM messages WHERE submitter_user_id = $1 AND deleted_at IS NULL ORDER BY created_at ASC LIMIT 1", submittedMsg.UID); err != nil {
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

			if len(submittedMsg.GiftIDs) != 0 {
				submittedMsg.HasGifts = true
			}

			hasMoney := false
			for _, giftId := range submittedMsg.GiftIDs {
				giftPrice := giftList.GetPriceByID(giftId)
				sendPrice += giftPrice

				if giftId == moneyGift.ID {
					hasMoney = true
				}
			}

			var currentBalance float32
			recipientUser, gotUserErr := getUserBySID(db, authClient, submittedMsg.RecipientID)
			if gotUserErr != nil {
				passivePrintError(err)
			}

			if err := Transact(db, func(tx *sqlx.Tx) error {
				// transact first before proceeding
				gotCurrentBalance, err := b.DeductBalanceTo(
					token.UID,
					sendPrice,
					fmt.Sprintf("Send message to %s", submittedMsg.RecipientID),
					tx,
				)
				if err != nil {
					return err
				}

				currentBalance = gotCurrentBalance

				// insert message
				if rows, err := tx.NamedQuery(
					"INSERT INTO messages (recipient_id, content, submitter_user_id, has_gifts) VALUES (:recipient_id, :content, :submitter_user_id, :has_gifts) RETURNING id",
					&submittedMsg,
				); err != nil {
					return err
				} else if id, err := wrapSqlRowsAfterInsert(rows); err != nil {
					return err
				} else {
					submittedMsg.ID = id
				}

				for _, giftId := range submittedMsg.GiftIDs {
					if res, err := tx.Exec("INSERT INTO message_gifts (message_id, gift_id) VALUES ($1, $2)", submittedMsg.ID, giftId); err != nil {
						return err
					} else if err := wrapSqlResult(res); err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				return err
			}

			if err := Transact(db, func(tx *sqlx.Tx) error {
				// give the money to the person
				if hasMoney && gotUserErr == nil {
					if _, err := b.AddBalanceTo(
						recipientUser.UID,
						moneyGift.Price,
						fmt.Sprintf("Money gift from message %s", submittedMsg.ID),
						tx,
					); err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				passivePrintError(err)
			}

			// save to search engine
			passivePrintError(UpsertEntry(messagesSearchIndex, submittedMsg.ID, MessageSearchEntry{
				MessageID:   submittedMsg.ID,
				RecipientID: submittedMsg.RecipientID,
				CreatedAt:   submittedMsg.CreatedAt,
				UpdatedAt:   submittedMsg.UpdatedAt,
				HasGifts:    submittedMsg.HasGifts,
			}))

			// send email to recipient if available
			// ignore the errors, just pass through
			if gotUserErr == nil {
				notifier := &EmailNotifier{
					Client:   postalOfficeClient,
					Template: emailTemplates["message"],
					Context: &MailSenderContext{
						Email:       recipientUser.Email,
						RecipientID: submittedMsg.RecipientID,
						MessageID:   submittedMsg.ID,
						FrontendURL: frontendUrl,
					},
				}

				// send the mail within n minutes.
				passivePrintError(notifier.Notify())
			}

			if !submittedMsg.HasGifts {
				recentMessagesChan <- submittedMsg.Message
			}

			return jsonEncode(rw, map[string]interface{}{
				"message":         "Message created successfully",
				"current_balance": currentBalance,
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
			if err := db.Get(&message, "SELECT * FROM messages WHERE id = $1 AND recipient_id = $2 AND deleted_at IS NULL", messageId, recipientId); err != nil {
				passivePrintError(err)
				r.NotFoundHandler().ServeHTTP(rw, rr)
				return
			}

			if err := db.Select(&message.GiftIDs, "SELECT gift_id FROM message_gifts WHERE message_id = $1", messageId); err != nil {
				passivePrintError(err)
			}

			newCtx := context.WithValue(rr.Context(), gotMessageKey{}, message)
			next.ServeHTTP(rw, rr.WithContext(newCtx))
		})
	}

	r.With(getRawMessage).Get("/messages/{recipientId}/{messageId}", wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
		message := rr.Context().Value(gotMessageKey{}).(RawMessage)

		// generate image if ?image query
		if rr.URL.Query().Has("image") {
			buf, err := imageRenderer.Render(imageTypeTwitter, message.Message)
			if err != nil {
				return err
			}

			rw.Header().Set("Content-Type", "image/png")
			rw.Write(buf)
			return nil
		}

		isDeletable := false
		isUserSenderOrReceiver := false
		reply := &MessageReply{}
		if token, authClient, tErr := getAuthToken(rr, firebaseApp); token != nil {
			gotRecipientUser, _ := getUserBySID(db, authClient, message.RecipientID)
			if token.UID == message.UID || (gotRecipientUser != nil && token.UID == gotRecipientUser.UID) {
				isUserSenderOrReceiver = true
			}

			if token.UID == message.UID {
				isDeletable = true
			}
		} else if tErr != nil {
			log.Println(tErr.Error())
		}

		respPayload := map[string]interface{}{
			"is_deletable": isDeletable,
			"message":      message,
		}

		// get reply if possible
		if isUserSenderOrReceiver {
			if message.HasReplied {
				// ignore error
				if err := db.Get(reply, "SELECT * FROM message_replies WHERE message_id = $1", message.ID); err != nil {
					log.Println(err)
					respPayload["reply"] = nil
				} else {
					respPayload["reply"] = reply
				}
			}
		} else if message.HasGifts {
			// make notes with gifts limited to sender and receivers only
			// TODO: disable_restricted_access_to_gift_messages
			return &ResponseError{
				StatusCode: http.StatusForbidden,
			}
		}

		return jsonEncode(rw, respPayload)
	}))

	r.With(appVerifyUser, getRawMessage).Delete("/messages/{recipientId}/{messageId}", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := getAuthTokenByReq(r)
		message := r.Context().Value(gotMessageKey{}).(RawMessage)
		// timeToSend := emailTemplates["message"].TimeToSend()
		//  || time.Since(message.CreatedAt) >= timeToSend

		if message.UID != token.UID {
			return &ResponseError{
				StatusCode: http.StatusForbidden,
			}
		}

		res, err := db.Exec("UPDATE messages SET deleted_at = $1 WHERE id = $2", time.Now(), message.ID)
		if err != nil {
			return err
		} else if err := wrapSqlResult(res); err != nil {
			return err
		}

		// cancel send job if possible
		if err := postalOfficeClient.CancelJobByUID(message.ID); err != nil {
			log.Println(err)
		}

		if err := DeleteEntry(messagesSearchIndex, message.ID); err != nil {
			log.Println(err)
		}

		cacher.Delete(fmt.Sprintf("image/%s", message.ID))
		return jsonEncode(rw, map[string]string{
			"message": "message deleted successfully",
		})
	}))

	r.With(jsonOnly, appVerifyUser, getRawMessage).
		Post("/messages/{recipientId}/{messageId}/reply", wrapHandler(func(rw http.ResponseWriter, rr *http.Request) error {
			// retrieve message
			message := rr.Context().Value(gotMessageKey{}).(RawMessage)

			// retrieve token
			token := getAuthTokenByReq(rr)
			authClient := getAuthClientByReq(rr)

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
			var submittedData struct {
				Reply   MessageReply `json:"reply"`
				Options struct {
					PostToTwitter bool `json:"post_to_twitter"`
					PostToEmail   bool `json:"post_to_email"`
				} `json:"options"`
			}

			if err := json.NewDecoder(rr.Body).Decode(&submittedData); err != nil {
				return err
			}

			reply := submittedData.Reply
			if err := checkProfanity(reply.Content); err != nil {
				return err
			} else if err := validator.Struct(&reply); err != nil {
				return wrapValidationError(rw, err)
			}

			reply.MessageID = message.ID
			var currentBalance float32

			Transact(db, func(tx *sqlx.Tx) error {
				// transact first before proceeding
				if gotCurrentBalance, err := b.DeductBalanceTo(
					token.UID,
					150.0,
					fmt.Sprintf("Reply message to %s", reply.MessageID),
					tx,
				); err != nil {
					return err
				} else {
					currentBalance = gotCurrentBalance
				}

				if updateRes, err := tx.NamedExec("INSERT INTO message_replies (message_id, content) VALUES (:message_id, :content)", &reply); err != nil {
					return err
				} else if err := wrapSqlResult(updateRes); err != nil {
					return err
				}

				if res, err := tx.Exec("UPDATE messages SET has_replied = true WHERE id = $1", message.ID); err != nil {
					return err
				} else if err := wrapSqlResult(res); err != nil {
					return err
				}

				return nil
			})

			notifier := MultiNotifier{[]Notifier{}}
			if twitterIdx, hasTwitter := isConnectedTo(connections, "twitter"); submittedData.Options.PostToTwitter && hasTwitter {
				imageData, err := imageRenderer.Render(imageTypeTwitter, message.Message)
				if err == nil {
					notifier.Add(&TwitterNotifier{
						Connection:  connections[twitterIdx],
						ImageData:   bytes.NewReader(imageData),
						TextContent: reply.Content,
						Hashtag:     twHashTag,
						Link:        fmt.Sprintf("%s/wall/%s/%s", frontendUrl, message.RecipientID, message.ID),
					})
				} else {
					log.Println(err)
				}
			}

			if _, hasEmail := isConnectedTo(connections, "email"); submittedData.Options.PostToEmail && hasEmail {
				// get sender email
				senderEmail, err := getUserEmailByUID(authClient, message.UID)
				if err == nil {
					notifier.Add(&EmailNotifier{
						Client:   postalOfficeClient,
						Template: emailTemplates["reply"],
						Context: &MailSenderContext{
							Email:       senderEmail,
							RecipientID: message.RecipientID,
							MessageID:   message.ID,
							FrontendURL: frontendUrl,
						},
					})
				} else {
					log.Println(err)
				}
			}

			// dont let notifier errors stop the process.
			passivePrintError(notifier.Notify())

			return jsonEncode(rw, map[string]interface{}{
				"message":         "reply success",
				"current_balance": currentBalance,
			})
		}))

	r.With(jsonOnly, appVerifyUser, getRawMessage).
		Delete("/messages/{recipientId}/{messageId}/reply", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			messageId := chi.URLParam(r, "messageId")
			if err := Transact(db, func(tx *sqlx.Tx) error {
				if res, err := tx.Exec("DELETE FROM message_replies WHERE message_id = $1", messageId); err != nil {
					return err
				} else if err := wrapSqlResult(res); err != nil {
					return err
				}

				if res, err := tx.Exec("UPDATE messages SET has_replied = false WHERE id = $1", messageId); err != nil {
					return err
				} else if err := wrapSqlResult(res); err != nil {
					return err
				}

				return nil
			}); err != nil {
				return err
			}

			return jsonEncode(rw, map[string]string{"message": "reply deleted successfully."})
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

	r.With(jsonOnly, appVerifyUser).
		Patch("/user/info", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := getAuthTokenByReq(r)
			updatedAssocInfo := &AssociatedUser{}

			if err := json.NewDecoder(r.Body).Decode(updatedAssocInfo); err != nil {
				return err
			}

			if res, err := db.Exec(
				"UPDATE associated_users SET department = $1, sex = $2 WHERE user_id = $3",
				updatedAssocInfo.Department,
				updatedAssocInfo.Sex,
				token.UID,
			); err != nil {
				return err
			} else if err := wrapSqlResult(res); err != nil {
				return err
			}
			return jsonEncode(rw, map[string]string{
				"message": "user details updated successfully",
			})
		}))

	r.With(appVerifyUser).Patch("/user/session", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := getAuthTokenByReq(r)

		// if user is already registered, get idle time
		if gotAssociatedUser, err := getAssociatedUserBy(db, sq.Eq{"user_id": token.UID}); err == nil {
			if gotAssociatedUser.LastActiveAt != nil {
				// convert idle time to virtual money
				receipt, err := b.ConvertIdleTime(token.UID, *gotAssociatedUser.LastActiveAt)
				if receipt != nil {
					// TODO: notify user
				} else if err != nil {
					passivePrintError(err)
				}
			}

			passivePrintError(updateUserLastActive(db, token.UID))
		}

		return jsonEncode(rw, map[string]string{"message": "ok"})
	}))

	r.With(appVerifyUser).
		Post("/user/session", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := getAuthTokenByReq(r)
			session, _ := store.Get(r, sessionName)
			session.Values["uid"] = token.UID
			if err := session.Save(r, rw); err != nil {
				return &ResponseError{
					StatusCode: http.StatusUnprocessableEntity,
					WError:     err,
				}
			}
			return jsonEncode(rw, map[string]string{
				"message": "ok",
			})
		}))

	r.Get("/invite/{invitationCode}", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		_, _, err := getAuthToken(r, firebaseApp)
		if err == nil {
			return &ResponseError{
				StatusCode: http.StatusForbidden,
				Message:    "Invitation links are only available for unregistered users.",
			}
		}

		invitationCode := chi.URLParam(r, "invitationCode")

		// verify invitation
		gotInvitation, err := invSys.VerifyInvitationCode(invitationCode)
		if err != nil {
			return err
		}

		// inject invitation referral
		invSys.InjectToRequest(rw, gotInvitation)
		return jsonEncode(rw, map[string]string{
			"message":         "ok",
			"invitation_code": gotInvitation.ID,
		})
	}))

	r.With(appVerifyUser).Get("/user/invitations", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := getAuthTokenByReq(r)
		invitations, err := invSys.GetInvitationsByUID(token.UID)
		if err != nil {
			return err
		}
		return jsonEncode(rw, invitations)
	}))

	r.With(appVerifyUser).Post("/user/invitations/generate", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := getAuthTokenByReq(r)
		if err := invSys.CheckEligibilityByUID(token.UID); err != nil {
			return err
		}

		r.ParseForm()
		rawMaxUsers := r.PostFormValue("max_users")
		if len(rawMaxUsers) == 0 {
			rawMaxUsers = "1"
		}

		maxUsers, err := strconv.Atoi(rawMaxUsers)
		if err != nil {
			return &ResponseError{
				StatusCode: http.StatusBadRequest,
				WError:     err,
				Message:    "Invalid value for max_users. Must be a number",
			}
		}

		newInvitationId, err := invSys.Generate(token.UID, maxUsers)
		if err != nil {
			return err
		}

		return jsonEncode(rw, map[string]interface{}{
			"message":         "Invitation link created successfully.",
			"invitation_code": newInvitationId,
			"expires_in_hrs":  defaultInvitationExpiration / time.Hour,
			"max_users":       maxUsers,
			"reward_coins":    int(invitationMoney),
		})
	}))

	r.With(appVerifyUser, pagination(&Paginator{
		OrderKey:  "created_at",
		TableName: virtualTransactionsTableName,
	})).
		Get("/user/transactions", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := getAuthTokenByReq(r)
			pg := getPaginatorFromReq(r)
			query := psql.Select().From(virtualTransactionsTableName).Where(sq.Eq{"user_id": token.UID})
			resp, err := pg.Load(&DatabasePaginatorSource{
				DB:        db,
				BaseQuery: query,
				DataQuery: query.Columns("*"),
				Converter: func(r *sqlx.Rows) (interface{}, error) {
					vtx := VirtualTransaction{}
					if err := r.StructScan(&vtx); err != nil {
						return nil, err
					}
					return vtx, nil
				},
			})
			if err != nil {
				return err
			}
			return jsonEncode(rw, resp)
		}))

	r.With(jsonOnly, appVerifyUser).Post("/user/cheque/deposit", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		var parsedData struct {
			ChequeID string `json:"cheque_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&parsedData); err != nil {
			return err
		}

		token := getAuthTokenByReq(r)
		if err := b.DepositChequeByID(parsedData.ChequeID, token.UID); err != nil {
			return err
		}

		return jsonEncode(rw, map[string]string{
			"message": "cheque deposited successfully",
		})
	}))

	r.With(appVerifyUser).
		Get("/user/info", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := getAuthTokenByReq(r)
			vWallet, err := b.GetWalletByUID(token.UID)
			passivePrintError(err)

			associatedData, err := getAssociatedUserBy(db, sq.Eq{"user_id": token.UID})
			if err != nil {
				log.Println(err)
				// return err
			}

			if associatedData == nil {
				associatedData = &AssociatedUser{}
			}

			userConnections := getUserConnections(db, token.UID)
			return jsonEncode(rw, map[string]interface{}{
				"associated_id":    associatedData.AssociatedID,
				"department":       associatedData.Department,
				"sex":              associatedData.Sex,
				"user_connections": userConnections,
				"wallet":           vWallet,
			})
		}))

	r.With(jsonOnly, appVerifyUser).
		Post("/user/setup", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			shouldDenyService := false
			token := getAuthTokenByReq(r)
			authClient := getAuthClientByReq(r)
			submittedData := AssociatedUser{}
			if err := json.NewDecoder(r.Body).Decode(&submittedData); err != nil {
				return err
			}

			if err := validator.Struct(&submittedData); err != nil {
				return wrapValidationError(rw, err)
			}

			if !submittedData.TermsAgreed {
				shouldDenyService = true
			} else if userEmail, err := getUserEmailByUID(authClient, token.UID); err == nil {
				matches := emailRegex.FindAllStringSubmatch(userEmail, -1)
				// deny service if no matching ID found when scanning the email through regex
				if (matches == nil || len(matches) == 0) || (len(matches[0]) < 2 || len(matches[0][1]) == 0) {
					shouldDenyService = true
				} else if gotId := matches[0][1]; gotId != submittedData.AssociatedID {
					shouldDenyService = true
				}
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
			if err := Transact(db, func(tx *sqlx.Tx) error {
				if res, err := tx.NamedExec(
					"INSERT INTO associated_users (user_id, associated_id, terms_agreed, sex, department) VALUES (:user_id, :associated_id, :terms_agreed, :sex, :department)",
					&submittedData,
				); err != nil {
					return &ResponseError{
						WError:     err,
						StatusCode: http.StatusUnprocessableEntity,
						Message:    "Failed to connect ID to user. Please try again.",
					}
				} else if err := wrapSqlResult(res, "Failed to connect ID to user. Please try again"); err != nil {
					return err
				}

				// initialize virtual wallet
				if _, err := b.AddInitialBalanceTo(token.UID, tx); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return &ResponseError{
					WError:     err,
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "Something went wrong while saving your account. Please try again.",
				}
			}

			// verify and save invitation
			if err := invSys.UseInvitationFromReq(rw, r, token.UID, b); err != nil {
				log.Println(err)
			}

			// generate welcome message
			if userEmail, err := getUserEmailByUID(authClient, token.UID); err == nil {
				emailId, _ := goNanoid.New()
				associatedUser, err := getAssociatedUserByEmail(db, authClient, userEmail)
				if err != nil {
					log.Println(err)
				}

				stats, err := getRecipientStatsBySID(messagesSearchIndex, associatedUser.AssociatedID)
				if err != nil {
					log.Println(err)
				}

				sender := emailTemplates["welcome"].With(struct {
					Email string
					Stats *RecipientStats
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

	// r.With(appVerifyUser).Get("/user/share_callback", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
	// 	queries := r.URL.Query()
	// 	action := queries.Get("action")
	// 	messageId := queries.Get("message_id")
	// 	sharerUserId := queries.Get("sharer_user_id")

	// }))

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
	}, htmlEncoder))

	r.With(appVerifyUser).Post("/user/connect_email", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := getAuthTokenByReq(r)
		authClient := getAuthClientByReq(r)
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

		res, err := db.NamedExec("INSERT INTO user_connections_new (user_id, provider, token, token_secret) VALUES (:user_id, :provider, :token, :token_secret)", &newEmailConnection)
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

		res, err := db.NamedExec("INSERT INTO user_connections_new (user_id, provider, token, token_secret) VALUES (:user_id, :provider, :token, :token_secret)", &newTwitterConnection)
		if err != nil {
			return err
		}

		if err := wrapSqlResult(res, "Unable to process Twitter login"); err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(getUserConnections(db, uid)); err != nil {
			passivePrintError(err)
			buf.WriteString("[]")
		}

		scriptJs := fmt.Sprintf(
			`<p>success</p>
<script type="text/javascript">
window.opener.postMessage({message:'twitter connect success',user_connections:%s}, '%s')
</script>`,
			buf.String(),
			frontendUrl,
		)
		return htmlEncode(rw, scriptJs)
	}, htmlEncoder))

	r.With(appVerifyUser).Delete("/user/connections/{connectionName}", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		token := getAuthTokenByReq(r)
		connectionName := chi.URLParam(r, "connectionName")

		// retrieve connections
		connections := getUserConnections(db, token.UID)
		if _, connected := isConnectedTo(connections, connectionName); !connected {
			return &ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "not connected to provider",
			}
		}

		if res, err := db.Exec("DELETE FROM user_connections_new WHERE user_id = $1 AND provider = $2", token.UID, connectionName); err != nil {
			return err
		} else if err := wrapSqlResult(res); err != nil {
			return err
		}

		return jsonEncode(rw, map[string]string{
			"message": "user third-party disconnection success",
		})
	}))

	r.With(jsonOnly, appVerifyUser).Post("/user/delete", wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		authClient := getAuthClientByReq(r)
		token := getAuthTokenByReq(r)

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

		// delete from associated_ids and user_connections_new
		if err := Transact(db, func(tx *sqlx.Tx) error {
			for _, tableName := range []string{"user_connections_new", "virtual_wallets", "associated_users"} {
				deleteSql, deleteArgs, _ := psql.Delete(tableName).Where(sq.Eq{"user_id": token.UID}).ToSql()
				if res, err := tx.Exec(deleteSql, deleteArgs...); err != nil {
					if err == sql.ErrNoRows {
						continue
					} else {
						return err
					}
				} else if err := wrapSqlResult(res); err != nil {
					passivePrintError(err)
					continue
				}
			}
			return nil
		}); err != nil {
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
		log.Panicln(err)
	}
}
