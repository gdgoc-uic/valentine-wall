package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	sq "github.com/Masterminds/squirrel"
	"github.com/blevesearch/bleve"
	goValidator "github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

func getStatusCode(statusCodes ...int) int {
	if len(statusCodes) != 0 {
		return statusCodes[0]
	} else {
		return http.StatusOK
	}
}

type ResponseEncoder interface {
	WriteHeader(rw http.ResponseWriter, statusCode ...int)
	Write(rw http.ResponseWriter, data interface{}) error
}

type JsonEncoder struct{}

func (enc *JsonEncoder) WriteHeader(rw http.ResponseWriter, statusCode ...int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(getStatusCode(statusCode...))
}

func (enc *JsonEncoder) Write(rw http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(rw).Encode(data)
}

type HtmlEncoder struct{}

func (enc *HtmlEncoder) WriteHeader(rw http.ResponseWriter, statusCode ...int) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(getStatusCode(statusCode...))
}

func (enc *HtmlEncoder) Write(rw http.ResponseWriter, data interface{}) error {
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

var jsonEncoder ResponseEncoder = &JsonEncoder{}
var htmlEncoder ResponseEncoder = &HtmlEncoder{}

func jsonEncode(rw http.ResponseWriter, data interface{}) error {
	jsonEncoder.WriteHeader(rw)
	return jsonEncoder.Write(rw, data)
}

func htmlEncode(rw http.ResponseWriter, data interface{}) error {
	htmlEncoder.WriteHeader(rw)
	return htmlEncoder.Write(rw, data)
}

func writeRespError(rw http.ResponseWriter, respErr *ResponseError, encoder ...ResponseEncoder) {
	if len(respErr.Message) == 0 {
		respErr.Message = http.StatusText(respErr.StatusCode)
	}
	log.Println(respErr.Error())
	if len(encoder) == 0 {
		jsonEncoder.WriteHeader(rw, respErr.StatusCode)
		jsonEncoder.Write(rw, respErr)
	} else {
		encoder[0].WriteHeader(rw, respErr.StatusCode)
		encoder[0].Write(rw, respErr)
	}
}

func wrapHandler(handler func(http.ResponseWriter, *http.Request) error, encoder ...ResponseEncoder) http.HandlerFunc {
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

type authTokenKey struct{}
type authClientKey struct{}

func getAuthClientByReq(r *http.Request) *auth.Client {
	cl, ok := r.Context().Value(authClientKey{}).(*auth.Client)
	if !ok {
		return nil
	}
	return cl
}

func getAuthTokenByReq(r *http.Request) *auth.Token {
	token, ok := r.Context().Value(authTokenKey{}).(*auth.Token)
	if !ok {
		return nil
	}
	return token
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

			ctx := context.WithValue(r.Context(), authTokenKey{}, token)
			ctxWithClient := context.WithValue(ctx, authClientKey{}, authClient)
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
	log.Printf("%s: body read - %d\n", resp.Request.URL, len(bodyBytes))
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

		jsonEncoder.WriteHeader(rw, http.StatusUnprocessableEntity)
		return jsonEncoder.Write(rw, map[string]interface{}{
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
	if err := db.Select(&connections, "SELECT * FROM user_connections_new WHERE user_id = ?", uid); err != nil {
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
	associatedData, err := getAssociatedUserBy(db, sq.Eq{"associated_id": sid})
	if err != nil {
		return nil, err
	}
	return authClient.GetUser(context.Background(), associatedData.UserID)
}

func getAssociatedUserByEmail(db *sqlx.DB, authClient *auth.Client, email string) (*AssociatedUser, error) {
	gotUser, err := authClient.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return getAssociatedUserBy(db, sq.Eq{"user_id": gotUser.UID})
}

func getAssociatedUserBy(db *sqlx.DB, pred Predicate) (*AssociatedUser, error) {
	associatedData := &AssociatedUser{}
	sql, args, err := sq.Select("*").From("associated_ids").Where(pred).ToSql()
	if err != nil {
		return nil, err
	}
	if err := db.Get(associatedData, sql, args...); err != nil {
		return nil, err
	}
	return associatedData, nil
}

func updateUserLastActive(db *sqlx.DB, uid string) error {
	if res, err := db.Exec(
		"UPDATE associated_ids SET last_active_at = ? WHERE user_id = ?",
		time.Now(), uid,
	); err != nil {
		return err
	} else if err := wrapSqlResult(res); err != nil {
		return err
	}
	return nil
}

func getRecipientStatsBySID(index bleve.Index, sid string) (*RecipientStats, error) {
	stats := &RecipientStats{RecipientID: sid}
	filter := bleve.NewTermQuery(sid)
	filter.SetField("recipient_id")

	for i, opt := range []bool{true, false} {
		gFilter := bleve.NewBoolFieldQuery(opt)
		gFilter.SetField("has_gifts")
		req := bleve.NewSearchRequest(bleve.NewConjunctionQuery(gFilter, filter))
		req.Fields = []string{}
		res, err := index.Search(req)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			stats.GiftMessagesCount = int(res.Total)
		} else {
			stats.MessagesCount = int(res.Total)
		}
	}

	return stats, nil
}

type FilterFunc func(*http.Request, context.Context, Filter) error

type Filter struct {
	Exists bool
	Value  string
	Name   string
}

func customFilters(filters map[string]FilterFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			query := r.URL.Query()
			ctx := r.Context()
			for targetQueryName, queryFunc := range filters {
				if err := queryFunc(r, ctx, Filter{
					Exists: query.Has(targetQueryName),
					Value:  query.Get(targetQueryName),
					Name:   targetQueryName}); err != nil {
					return err
				}
			}
			next.ServeHTTP(rw, r.WithContext(ctx))
			return nil
		})
	}
}

func passivePrintError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// sse-related
func encodeDataSSE(rw http.ResponseWriter, msg Message) {
	writer := &bytes.Buffer{}
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(msg); err != nil {
		log.Println(err)
		fmt.Fprintf(writer, "null")
	}
	writeResponseDataSSE(rw, writer)
}

func writeResponseDataSSE(rw http.ResponseWriter, buf *bytes.Buffer) {
	fmt.Fprint(rw, "data: ")
	buf.WriteTo(rw)
	fmt.Fprint(rw, "\n\n")
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}
