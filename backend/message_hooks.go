package main

import (
	"fmt"

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

	// fetch recipient except everyone
	if record.GetString("recipient") != "everyone" {
		if recipient, err := dao.FindFirstRecordByData("user_details", "student_id", record.GetString("recipient")); err == nil {
			dao.ExpandRecord(recipient, []string{"college_departments"}, func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
				return dao.FindRecordsByIds(relCollection.Name, relIds)
			})

			record.MergeExpand(map[string]any{"recipient": recipient})
		}
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
	if recipientId == "everyone" {
		return nil
	}

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

	if err := expandMessage(dao, e.Record); err != nil {
		return err
	}

	totalAmount, _ := computeGiftCost(e.Record)
	user := e.Record.Expand()["user"].(*models.Record)
	return checkSufficientFunds(dao, user.GetString("user"), sendPrice+totalAmount)
}

func onAddMessage(app core.App, e *core.RecordCreateEvent) error {
	dao := app.Dao()
	expandMessage(dao, e.Record)

	user := e.Record.Expand()["user"].(*models.Record)
	recipient, isRecipientAccessible := e.Record.Expand()["recipient"].(*models.Record)
	totalAmount, remittableAmount := computeGiftCost(e.Record)

	wallet, err := getWalletByUserId(dao, user.GetString("user"))
	if err != nil {
		return err
	}

	studentId := e.Record.GetString("recipient")
	if err := updateRanking(dao, studentId, totalAmount+sendPrice); err != nil {
		passivePrintError(err)
		return nil
	}

	if err := createTransaction(dao, wallet.Id, -sendPrice, fmt.Sprintf("Send message to %s", studentId)); err != nil {
		return err
	}

	if totalAmount != 0 {
		if err := createTransaction(dao,
			wallet.Id, -totalAmount,
			fmt.Sprintf("Sent virtual gifts for %s", studentId)); err != nil {
			passivePrintError(err)
		}
	}

	if isRecipientAccessible && remittableAmount != 0 {
		if err := createTransactionFromUser(dao, recipient.GetString("user"),
			remittableAmount, fmt.Sprintf("Gift message from message %s", e.Record.Id)); err != nil {
			passivePrintError(err)
		}
	}

	// send email
	if isRecipientAccessible {
		if msg, err := emailTemplates.message.With(map[string]any{
			"Email":      recipient.GetString("email"),
			"MessageURL": fmt.Sprintf("%s/wall/%s/%s", frontendUrl, studentId, e.Record.Id),
		}).Message(app.Settings().Meta, recipient.GetString("email")); err == nil {
			passivePrintError(app.NewMailClient().Send(msg))
		}
	}

	// update last active at
	user.Set("last_active", types.DateTime{})
	passivePrintError(dao.SaveRecord(user))

	return nil
}

func onRemoveMessage(dao *daos.Dao, e *core.RecordDeleteEvent) error {
	expandMessage(dao, e.Record)

	// deduct total_cost
	ranking, err := dao.FindFirstRecordByData("rankings", "recipient", e.Record.GetString("recipient"))
	if err != nil {
		passivePrintError(err)
		return nil
	}

	totalAmount, _ := computeGiftCost(e.Record)
	ranking.Set("total_coins", ranking.GetFloat("total_coins")-totalAmount)
	passivePrintError(dao.SaveRecord(ranking))

	return nil
}

func onAddMessageReply(app core.App, e *core.RecordCreateEvent) error {
	dao := app.Dao()
	expandMessageReply(dao, e.Record)
	user := e.Record.Expand()["sender"].(*models.Record)

	// update last active at
	if user != nil {
		user.Set("last_active", types.DateTime{})
		passivePrintError(dao.SaveRecord(user))
	}

	msg, msgOk := e.Record.Expand()["message"].(*models.Record)
	if msgOk {
		msg.Set("replies_count", msg.GetInt("replies_count")+1)
		passivePrintError(dao.SaveRecord(msg))
		expandMessage(dao, msg)
		recipient, isRecipientAccessible := msg.Expand()["recipient"].(*models.Record)
		if isRecipientAccessible {
			if mailMsg, err := emailTemplates.reply.With(map[string]any{
				"Email":       recipient.GetString("email"),
				"RecipientID": msg.GetString("recipient"),
				"MessageURL":  fmt.Sprintf("%s/wall/%s/%s", frontendUrl, msg.GetString("recipient"), e.Record.Id),
			}).Message(app.Settings().Meta, recipient.GetString("email")); err == nil {
				passivePrintError(app.NewMailClient().Send(mailMsg))
			}
		}
	}

	return createTransactionFromUser(dao, user.GetString("user"), -sendPrice, fmt.Sprintf("Reply message %s", e.Record.Id))
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

func onBeforeAddMessageReply(dao *daos.Dao, e *core.RecordCreateEvent) error {
	// check profanity content
	if err := checkProfanity(e.Record.GetString("content")); err != nil {
		return err.ToApiError()
	}

	if err := expandMessageReply(dao, e.Record); err != nil {
		return err
	}

	sender := e.Record.Expand()["sender"].(*models.Record)
	return checkSufficientFunds(dao, sender.GetString("user"), sendPrice)
}
