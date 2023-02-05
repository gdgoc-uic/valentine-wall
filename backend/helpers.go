package main

import (
	"net/http"
)

// func wrapValidationError(rw http.ResponseWriter, err error) error {
// 	if validatorErrors, vOk := err.(goValidator.ValidationErrors); vOk {
// 		errs := []map[string]interface{}{}
// 		for _, er := range validatorErrors {
// 			errs = append(errs, map[string]interface{}{
// 				"field":   er.Field(),
// 				"message": er.Error(),
// 			})
// 		}

// 		jsonEncoder.WriteHeader(rw, http.StatusUnprocessableEntity)
// 		return jsonEncoder.Write(rw, map[string]interface{}{
// 			"error_message": "there were errors when submitting your data.",
// 			"errors":        errs,
// 		})
// 	}
// 	return err
// }

func checkProfanity(content string) *ResponseError {
	if profanityDetector.IsProfane(content) {
		return &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Your submission contains inappropriate content.",
		}
	}

	return nil
}

// func getUserConnections(db *sqlx.DB, uid string) []UserConnection {
// 	connections := []UserConnection{}
// 	if err := db.Select(&connections, "SELECT * FROM user_connections_new WHERE user_id = $1", uid); err != nil {
// 		log.Println(err)
// 		// return err
// 	}
// 	return connections
// }

// func getUserEmailByUID(authClient *auth.Client, uid string) (string, error) {
// 	gotUser, err := authClient.GetUser(context.Background(), uid)
// 	if err != nil {
// 		return "", err
// 	}

// 	return gotUser.Email, nil
// }

// func getUserBySID(db *sqlx.DB, authClient *auth.Client, sid string) (*auth.UserRecord, error) {
// 	associatedData, err := getAssociatedUserBy(db, sq.Eq{"associated_id": sid})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return authClient.GetUser(context.Background(), associatedData.UserID)
// }

// func getAssociatedUserByEmail(db *sqlx.DB, authClient *auth.Client, email string) (*AssociatedUser, error) {
// 	gotUser, err := authClient.GetUserByEmail(context.Background(), email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return getAssociatedUserBy(db, sq.Eq{"user_id": gotUser.UID})
// }

// func getAssociatedUserBy(db *sqlx.DB, pred Predicate) (*AssociatedUser, error) {
// 	associatedData := &AssociatedUser{}
// 	sql, args, err := psql.Select("*").From("associated_users").Where(pred).ToSql()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := db.Get(associatedData, sql, args...); err != nil {
// 		return nil, err
// 	}
// 	return associatedData, nil
// }

// func updateUserLastActive(db *sqlx.DB, uid string) error {
// 	return Transact(db, func(tx *sqlx.Tx) error {
// 		if res, err := tx.Exec(
// 			"UPDATE associated_users SET last_active_at = $1 WHERE user_id = $2",
// 			time.Now(), uid,
// 		); err != nil {
// 			return err
// 		} else if err := wrapSqlResult(res); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
