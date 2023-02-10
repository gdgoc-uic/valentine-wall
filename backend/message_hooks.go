package main

import (
	"fmt"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

func expandMessage(dao *daos.Dao, record *models.Record) error {
	errs := dao.ExpandRecord(record, []string{"user", "user.user", "gifts"}, func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
		return dao.FindRecordsByIds(relCollection.Name, relIds)
	})
	if len(errs) != 0 {
		// TODO: add error
		return nil
	}

	// fetch recipient
	if recipient, err := dao.FindFirstRecordByData("user_details", "student_id", record.GetString("recipient")); err == nil {
		dao.ExpandRecord(recipient, []string{"college_departments"}, func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
			return dao.FindRecordsByIds(relCollection.Name, relIds)
		})

		record.MergeExpand(map[string]any{"recipient": recipient})
	}

	return nil
}

func expandMessageReply(dao *daos.Dao, record *models.Record) error {
	errs := dao.ExpandRecord(record, []string{"sender", "message"}, func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
		return dao.FindRecordsByIds(relCollection.Name, relIds)
	})
	if len(errs) != 0 {
		// TODO: add error
		return nil
	}
	return nil
}

func updateRanking(dao *daos.Dao, recipientId string, coinsToAdd float64) error {
	ranking, err := dao.FindFirstRecordByData("rankings", "recipient", recipientId)
	if err != nil {
		collection, err := dao.FindCollectionByNameOrId("rankings")
		if err != nil {
			return err
		}

		ranking = models.NewRecord(collection)
		ranking.Set("recipient", recipientId)
		ranking.Set("college_department", "unknown")
		ranking.Set("sex", "unknown")
		ranking.Set("total_coins", 0)
	}

	// fetch recipient / student id
	recipient, err := dao.FindFirstRecordByData("user_details", "student_id", recipientId)
	if err == nil {
		if ranking.GetString("college_department") == "unknown" && ranking.GetString("sex") == "unknown" {
			cDept, _ := dao.FindRecordById("college_departments", recipient.GetString("college_department"))
			ranking.Set("college_department", cDept.Id)
			ranking.Set("sex", recipient.GetString("sex"))
		}
	}

	ranking.Set("total_coins", ranking.GetFloat("total_coins")+coinsToAdd)
	return dao.SaveRecord(ranking)
}

func computeGiftCost(record *models.Record) (totalAmount float64, remittableAmount float64) {
	if giftsList, giftsListExists := record.Expand()["gifts"]; giftsListExists {
		if gifts, gOk := giftsList.([]*models.Record); gOk {
			for _, msgGift := range gifts {
				totalAmount += msgGift.GetFloat("price")

				if msgGift.GetBool("is_remittable") {
					remittableAmount += msgGift.GetFloat("price")
				}
			}
		}
	}

	return totalAmount, remittableAmount
}

func onBeforeAddMessage(dao *daos.Dao, e *core.RecordCreateEvent) error {
	// to avoid spams
	if r, err := dao.FindRecordsByExpr(
		e.Record.Collection().Name,
		dbx.HashExp{
			"content":   e.Record.GetString("content"),
			"recipient": e.Record.GetString("recipient"),
		},
	); err == nil && len(r) != 0 {
		return apis.NewBadRequestError(
			"You have posted a similar message to a similar recipient.", nil)
	}

	// check profanity content
	if err := checkProfanity(e.Record.GetString("content")); err != nil {
		return err.ToApiError()
	}

	return nil
}

func onAddMessage(dao *daos.Dao, e *core.RecordCreateEvent) error {
	expandMessage(dao, e.Record)

	user := e.Record.Expand()["user"].(*models.Record)
	recipient, isRecipientAccessible := e.Record.Expand()["recipient"].(*models.Record)
	totalAmount, remittableAmount := computeGiftCost(e.Record)

	wallet, err := getWalletByUserId(dao, user.GetString("user"))

	if err != nil {
		return err
	} else if wallet.GetFloat("amount") < remittableAmount {
		return (&ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "You have insufficient funds to send virtual remittable gifts to recipient.",
		}).ToApiError()
	}

	studentId := e.Record.GetString("recipient")
	if err := updateRanking(dao, studentId, totalAmount+sendPrice); err != nil {
		// TODO: add error
		// return err
		return nil
	}

	if err := createTransaction(dao, wallet.Id, -sendPrice, fmt.Sprintf("Send message to %s", studentId)); err != nil {
		return err
	}

	if totalAmount != 0 {
		if err := createTransaction(dao,
			wallet.Id, -totalAmount,
			fmt.Sprintf("Sent virtual gifts for %s", studentId)); err != nil {
			// TODO: return err
		}
	}

	if isRecipientAccessible && remittableAmount != 0 {
		if err := createTransactionFromUser(dao, recipient.Id,
			remittableAmount, fmt.Sprintf("Gift message from message %s", e.Record.Id)); err != nil {
			// TODO: return err
		}
	}

	// send email
	emailTemplates.message.With(map[string]any{
		"Name":       recipient.GetString("name"),
		"MessageURL": "https://test",
	})

	// update last active at
	user.Set("last_active", types.DateTime{})
	dao.SaveRecord(user) //TODO: error

	return nil
}

func onRemoveMessage(dao *daos.Dao, e *core.RecordDeleteEvent) error {
	expandMessage(dao, e.Record)

	// deduct total_cost
	ranking, err := dao.FindFirstRecordByData("rankings", "recipient_id", e.Record.GetString("recipient"))
	if err != nil {
		// TODO: add error
		return nil
	}

	totalAmount, _ := computeGiftCost(e.Record)
	ranking.Set("total_coins", ranking.GetFloat("total_coins")-totalAmount)
	dao.SaveRecord(ranking) // TODO: add error

	return nil
}

func onAddMessageReply(dao *daos.Dao, e *core.RecordCreateEvent) error {
	// check profanity content
	if err := checkProfanity(e.Record.GetString("content")); err != nil {
		return err.ToApiError()
	}

	expandMessageReply(dao, e.Record)
	user := e.Record.Expand()["sender"].(*models.Record)

	// update last active at
	if user != nil {
		user.Set("last_active", types.DateTime{})
		dao.SaveRecord(user) //TODO: error
	}

	msg, msgOk := e.Record.Expand()["message"].(*models.Record)
	if msgOk {
		msg.Set("replies_count", msg.GetInt("replies_count")+1)
		passivePrintError(dao.SaveRecord(msg))
	}

	return createTransactionFromUser(dao, e.Record.GetString("sender"), sendPrice, fmt.Sprintf("Reply message %s", e.Record.Id))
}

func onRemoveMessageReply(dao *daos.Dao, e *core.RecordDeleteEvent) error {
	expandMessageReply(dao, e.Record)

	msg, msgOk := e.Record.Expand()["message"].(*models.Record)
	if msgOk {
		msg.Set("replies_count", msg.GetInt("replies_count")-1)
		passivePrintError(dao.SaveRecord(msg))
	}

	return nil
}
