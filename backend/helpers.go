package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	sq "github.com/Masterminds/squirrel"
	goValidator "github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

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
	associatedData, err := getAssociatedUserBy(db, sq.Eq{"associated_id": sid})
	if err != nil {
		return nil, err
	}
	return authClient.GetUser(context.Background(), associatedData.UserID)
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

func getMessageStatsBySID(db *sqlx.DB, sid string) (*MessageStats, error) {
	stats := &MessageStats{}
	eqId := sq.Eq{"recipient_id": sid}
	joinStmt := "message_gifts on message_gifts.message_id = messages.id"
	baseQuery := sq.Select("count(*)")
	giftMessageCountQuery := sq.Select("id", "recipient_id").Distinct().From("messages").InnerJoin(joinStmt).Where(eqId)
	giftMessagesCountQuery2 := baseQuery.FromSelect(giftMessageCountQuery, "msg").GroupBy("recipient_id")
	nonGiftMessagesCountSql, ngmArgs, err := baseQuery.From("messages").LeftJoin(joinStmt).Where("message_gifts.gift_id IS NULL").Where(eqId).ToSql()
	if err != nil {
		return nil, err
	}

	statSql, args, err := giftMessagesCountQuery2.Suffix("UNION ALL "+nonGiftMessagesCountSql, ngmArgs...).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(statSql, args...)
	if err != nil {
		return nil, err
	}

	rows.Next()
	if err := rows.Scan(&stats.GiftMessages); err != nil {
		log.Println(err)
	}

	rows.Next()
	if err := rows.Scan(&stats.Messages); err != nil {
		log.Println(err)
	}
	return stats, nil
}
