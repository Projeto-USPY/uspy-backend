package private

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/private"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetComment(ctx *gin.Context, DB db.Env, userID string, off *controllers.Offering) {
	mask := "subjects/%s/offerings/%s/comments"
	userHash := models.User{ID: userID}.Hash()
	subHash := models.Subject{
		Code:           off.Subject.Code,
		CourseCode:     off.Subject.CourseCode,
		Specialization: off.Subject.Specialization,
	}.Hash()

	var comment models.Comment
	if snap, err := DB.Restore(fmt.Sprintf(mask, subHash, off.Hash), userHash); err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error getting comment: (sub:%s/%s, user:%s): %s", subHash, off.Hash, userHash, err.Error()),
		)
		return
	} else {
		if err := snap.DataTo(&comment); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error binding comment: %s", err.Error()))
			return
		}
	}

	private.GetComment(ctx, &comment)
}

func ReportComment(
	ctx *gin.Context,
	DB db.Env,
	userID string,
	comment *controllers.CommentRating,
	body *controllers.CommentReportBody,
) {
	subHash := models.Subject{
		Code:           comment.Offering.Subject.Code,
		CourseCode:     comment.Offering.Subject.CourseCode,
		Specialization: comment.Offering.Subject.Specialization,
	}.Hash()
	userHash := models.User{
		ID: userID,
	}.Hash()

	userCommentHash := models.UserComment{
		ProfessorCode:  comment.Offering.Hash,
		Subject:        comment.Offering.Subject.Code,
		Course:         comment.Offering.Subject.CourseCode,
		Specialization: comment.Offering.Subject.Specialization,
	}.Hash()

	err := DB.Client.RunTransaction(ctx, func(txCtx context.Context, tx *firestore.Transaction) error {
		commentsCol := "subjects/%s/offerings/%s/comments"
		target := DB.Client.Collection(
			fmt.Sprintf(commentsCol, subHash, comment.Offering.Hash),
		).Where("id", "==", uuid.MustParse(comment.Comment)).Limit(1)

		snaps, err := tx.Documents(target).GetAll()
		if err != nil {
			return err
		} else if len(snaps) == 0 {
			return utils.ErrCommentNotFound
		}

		var modelComment models.Comment
		var targetUserID string

		if err := snaps[0].DataTo(&modelComment); err != nil {
			return err
		} else {
			targetUserID = snaps[0].Ref.ID
		}

		commentReportMask := "users/%s/comment_reports/%s"
		reportRef := DB.Client.Doc(fmt.Sprintf(commentReportMask, userHash, comment.Comment))

		modelCommentReport := models.CommentReport{
			ID:     uuid.MustParse(comment.Comment),
			Report: body.Body,
		}

		if _, err := tx.Get(reportRef); err != nil {
			if status.Code(err) == codes.NotFound { // comment has not been reported by this user yet
				// increment comment report count by 1
				if updateErr := tx.Update(snaps[0].Ref, []firestore.Update{{Path: "reports", Value: firestore.Increment(1)}}); updateErr != nil {
					return updateErr
				}

				// increment replica report count by 1
				replicaMask := "users/%s/user_comments/%s"
				replicaRef := DB.Client.Doc(fmt.Sprintf(replicaMask, targetUserID, userCommentHash))
				if updateErr := tx.Update(replicaRef, []firestore.Update{{Path: "comment.reports", Value: firestore.Increment(1)}}); updateErr != nil {
					return updateErr
				}
			}
		}

		return tx.Set(reportRef, modelCommentReport)
	})

	if err != nil {
		if status.Code(err) == codes.NotFound || err == utils.ErrCommentNotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error reporting comment: (sub:%s, prof:%s): %s", subHash, comment.Offering.Hash, err.Error()),
		)
		return
	}

	private.ReportComment(ctx)
}

func PublishComment(
	ctx *gin.Context,
	DB db.Env,
	userID string,
	off *controllers.Offering,
	comment *controllers.Comment,
) {
	modelSub := models.NewSubjectFromController(&off.Subject)
	userHash := models.User{ID: userID}.Hash()

	// check if subject exists and if user has permission to comment
	if err := utils.CheckSubjectPermission(DB, userHash, modelSub.Hash()); err != nil {
		if err == utils.ErrSubjectNotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", modelSub, err.Error()))
			return
		}

		if err == utils.ErrNoPermission {
			ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("user %v has no permission to comment: %s", userID, err.Error()))
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error checking subject permission: %s", err.Error()))
		return
	} else if _, err := DB.Restore("subjects/"+modelSub.Hash()+"/offerings", off.Hash); err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error getting offering: (sub:%s, prof:%s): %s", modelSub.Hash(), off.Hash, err.Error()),
		)
		return
	}

	// create new comment object
	newComment := models.Comment{
		ID:        uuid.New(),
		Rating:    comment.Rating,
		Body:      comment.Body,
		Edited:    false,
		Timestamp: time.Now(),
		Upvotes:   0,
		Downvotes: 0,
		Reports:   0,
	}

	err := DB.Client.RunTransaction(ctx, func(txCtx context.Context, tx *firestore.Transaction) error {
		collectionMask := "subjects/%s/offerings/%s/comments/%s"
		commentRef := DB.Client.Doc(
			fmt.Sprintf(
				collectionMask,
				modelSub.Hash(),
				off.Hash,
				userHash,
			),
		)

		var storedComment *models.Comment
		if snap, err := tx.Get(commentRef); err != nil {
			if status.Code(err) != codes.NotFound {
				return err
			}
		} else {
			if err := snap.DataTo(&storedComment); err != nil {
				return err
			}

			// overwrite new object with stored values
			newComment.Edited = true
			newComment.Upvotes = storedComment.Upvotes
			newComment.Downvotes = storedComment.Downvotes
			newComment.Reports = storedComment.Reports
			newComment.ID = storedComment.ID
		}

		// upsert comment in database
		tx.Set(commentRef, newComment)

		// upsert replica in user comments (will be used in the future)
		replica := models.UserComment{
			Comment:        newComment,
			ProfessorCode:  off.Hash,
			Subject:        off.Subject.Code,
			Course:         off.Subject.CourseCode,
			Specialization: off.Subject.Specialization,
		}

		replicaMask := "users/%s/user_comments/%s"
		replicaRef := DB.Client.Doc(
			fmt.Sprintf(
				replicaMask,
				userHash,
				replica.Hash(),
			),
		)

		return tx.Set(replicaRef, replica)
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error creating/updating comment: %s", err.Error()))
		return
	}

	private.PublishComment(ctx)
}
