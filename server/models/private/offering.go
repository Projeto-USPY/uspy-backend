package private

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	db_utils "github.com/Projeto-USPY/uspy-backend/db/utils"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/private"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetComment retrieves the comment associated with an offering made by a given user
func GetComment(ctx *gin.Context, DB db.Env, userID string, off *controllers.Offering) {
	mask := "subjects/%s/offerings/%s/comments/%s"
	userHash := models.User{ID: userID}.Hash()
	subHash := models.Subject{
		Code:           off.Subject.Code,
		CourseCode:     off.Subject.CourseCode,
		Specialization: off.Subject.Specialization,
	}.Hash()

	var comment models.Comment
	snap, err := DB.Restore(fmt.Sprintf(mask, subHash, off.Hash, userHash))

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error getting comment: (sub:%s/%s, user:%s): %s", subHash, off.Hash, userHash, err.Error()),
		)

		return
	}

	if err := snap.DataTo(&comment); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error binding comment: %s", err.Error()))
		return
	}

	private.GetComment(ctx, &comment)
}

// GetCommentRating retrieves the rating made for a comment by a given user
func GetCommentRating(
	ctx *gin.Context,
	DB db.Env,
	userID string,
	comment *controllers.CommentRating,
) {
	userHash := models.User{
		ID: userID,
	}.Hash()

	var model models.CommentRating
	collectionMask := "users/%s/comment_ratings/%s"

	snap, err := DB.Restore(fmt.Sprintf(collectionMask, userHash, comment.ID))

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error looking up comment rating: %s", err.Error()))
		return
	}

	if snap.DataTo(&model); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error binding comment: %s", err.Error()))
		return
	}

	private.GetCommentRating(ctx, &model)
}

// RateComment takes a user's rating and applies it to a given comment
func RateComment(
	ctx *gin.Context,
	DB db.Env,
	userID string,
	comment *controllers.CommentRating,
	body *controllers.CommentRateBody,
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
		ProfessorHash:  comment.Offering.Hash,
		Subject:        comment.Offering.Subject.Code,
		Course:         comment.Offering.Subject.CourseCode,
		Specialization: comment.Offering.Subject.Specialization,
	}.Hash()

	err := DB.Client.RunTransaction(ctx, func(txCtx context.Context, tx *firestore.Transaction) error {
		commentsCol := "subjects/%s/offerings/%s/comments"
		target := DB.Client.Collection(
			fmt.Sprintf(commentsCol, subHash, comment.Offering.Hash),
		).Where("id", "==", uuid.MustParse(comment.ID)).Limit(1)

		snaps, err := tx.Documents(target).GetAll()
		if err != nil {
			return err
		} else if len(snaps) == 0 {
			return db_utils.ErrCommentNotFound
		}

		var modelComment models.Comment
		var targetRef *firestore.DocumentRef

		if err := snaps[0].DataTo(&modelComment); err != nil {
			return err
		}

		targetRef = snaps[0].Ref

		commentRatingMask := "users/%s/comment_ratings/%s"
		ratingRef := DB.Client.Doc(fmt.Sprintf(commentRatingMask, userHash, comment.ID))

		commentRating := models.CommentRating{
			ID:     uuid.MustParse(comment.ID),
			Upvote: body.Type == "upvote",

			ProfessorHash:  comment.Offering.Hash,
			Subject:        comment.Offering.Subject.Code,
			Course:         comment.Offering.Subject.CourseCode,
			Specialization: comment.Offering.Subject.Specialization,
		}

		type update struct {
			ref     *firestore.DocumentRef
			changes []firestore.Update
		}

		updates := make([]update, 0, 10)

		replicaMask := "users/%s/user_comments/%s"
		replicaRef := DB.Client.Doc(fmt.Sprintf(replicaMask, targetRef.ID, userCommentHash))

		// if rating already exists and it's different, we add the decrement updates for the comment and replica's count
		if ratingDoc, err := tx.Get(ratingRef); err == nil {
			storedUpvote, err := ratingDoc.DataAt("upvote")

			if err == nil && body.Type != "none" && storedUpvote.(bool) == commentRating.Upvote {
				// rating did not change
				return nil
			} else if err != nil {
				return err
			}

			// if rating changed and it's none (user removed their rating)
			if body.Type == "none" {
				if storedUpvote.(bool) {
					updates = append(updates,
						update{
							ref:     targetRef,
							changes: []firestore.Update{{Path: "upvotes", Value: firestore.Increment(-1)}},
						},
						update{
							ref:     replicaRef,
							changes: []firestore.Update{{Path: "comment.upvotes", Value: firestore.Increment(-1)}},
						},
					)
				} else {
					updates = append(updates,
						update{
							ref:     targetRef,
							changes: []firestore.Update{{Path: "downvotes", Value: firestore.Increment(-1)}},
						},
						update{
							ref:     replicaRef,
							changes: []firestore.Update{{Path: "comment.downvotes", Value: firestore.Increment(-1)}},
						},
					)
				}
			} else { // rating changed to the opposite type
				if commentRating.Upvote {
					updates = append(updates,
						update{
							ref:     targetRef,
							changes: []firestore.Update{{Path: "downvotes", Value: firestore.Increment(-1)}},
						},
						update{
							ref:     replicaRef,
							changes: []firestore.Update{{Path: "comment.downvotes", Value: firestore.Increment(-1)}},
						},
					)
				} else {
					updates = append(updates,
						update{
							ref:     targetRef,
							changes: []firestore.Update{{Path: "upvotes", Value: firestore.Increment(-1)}},
						},
						update{
							ref:     replicaRef,
							changes: []firestore.Update{{Path: "comment.upvotes", Value: firestore.Increment(-1)}},
						},
					)
				}
			}

		} else if status.Code(err) != codes.NotFound {
			return err
		}

		// now we must add the updates to increment the comment and replica's count (only if the type isnt none)
		if body.Type != "none" {
			if commentRating.Upvote {
				updates = append(updates,
					update{
						ref:     targetRef,
						changes: []firestore.Update{{Path: "upvotes", Value: firestore.Increment(1)}},
					},
					update{
						ref:     replicaRef,
						changes: []firestore.Update{{Path: "comment.upvotes", Value: firestore.Increment(1)}},
					},
				)
			} else {
				updates = append(updates,
					update{
						ref:     targetRef,
						changes: []firestore.Update{{Path: "downvotes", Value: firestore.Increment(1)}},
					},
					update{
						ref:     replicaRef,
						changes: []firestore.Update{{Path: "comment.downvotes", Value: firestore.Increment(1)}},
					},
				)
			}
		}

		var wg sync.WaitGroup
		updateErrors := make(chan error, len(updates))
		wg.Add(len(updates))

		// perform updates in parallel
		for _, u := range updates {
			go func(upd update, wg *sync.WaitGroup) {
				defer wg.Done()
				updateErrors <- tx.Update(upd.ref, upd.changes)
			}(u, &wg)
		}

		wg.Wait()
		close(updateErrors)

		for err := range updateErrors {
			if err != nil {
				return err
			}
		}

		// upsert comment rating if type isnt none
		if body.Type != "none" {
			return tx.Set(ratingRef, commentRating)
		}

		// delete comment rating
		return tx.Delete(ratingRef)
	})

	if err != nil {
		if status.Code(err) == codes.NotFound || err == db_utils.ErrCommentNotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(
			http.StatusInternalServerError,
			fmt.Errorf("error rating comment: (sub:%s, prof:%s): %s", subHash, comment.Offering.Hash, err.Error()),
		)
		return
	}

	private.RateComment(ctx)
}

// ReportComment takes a user's report and applies it to a given comment
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
		ProfessorHash:  comment.Offering.Hash,
		Subject:        comment.Offering.Subject.Code,
		Course:         comment.Offering.Subject.CourseCode,
		Specialization: comment.Offering.Subject.Specialization,
	}.Hash()

	err := DB.Client.RunTransaction(ctx, func(txCtx context.Context, tx *firestore.Transaction) error {
		commentsCol := "subjects/%s/offerings/%s/comments"
		target := DB.Client.Collection(
			fmt.Sprintf(commentsCol, subHash, comment.Offering.Hash),
		).Where("id", "==", uuid.MustParse(comment.ID)).Limit(1)

		snaps, err := tx.Documents(target).GetAll()
		if err != nil {
			return err
		} else if len(snaps) == 0 {
			return db_utils.ErrCommentNotFound
		}

		var modelComment models.Comment
		var targetUserID string

		if err := snaps[0].DataTo(&modelComment); err != nil {
			return err
		}

		targetUserID = snaps[0].Ref.ID

		commentReportMask := "users/%s/comment_reports/%s"
		reportRef := DB.Client.Doc(fmt.Sprintf(commentReportMask, userHash, comment.ID))

		modelCommentReport := models.CommentReport{
			ID:     uuid.MustParse(comment.ID),
			Report: body.Body,

			ProfessorHash:  comment.Offering.Hash,
			Subject:        comment.Offering.Subject.Code,
			Course:         comment.Offering.Subject.CourseCode,
			Specialization: comment.Offering.Subject.Specialization,
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
			} else {
				return err
			}
		}

		return tx.Set(reportRef, modelCommentReport)
	})

	if err != nil {
		if status.Code(err) == codes.NotFound || err == db_utils.ErrCommentNotFound {
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

// PublishComment publishes a comment and its associated reaction made by a given user
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
	if err := db_utils.CheckSubjectPermission(DB, userHash, modelSub.Hash()); err != nil {
		if err == db_utils.ErrSubjectNotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find subject %v: %s", modelSub, err.Error()))
			return
		}

		if err == db_utils.ErrNoPermission {
			ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("user %v has no permission to comment: %s", userID, err.Error()))
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error checking subject permission: %s", err.Error()))
		return
	} else if _, err := DB.Restore("subjects/" + modelSub.Hash() + "/offerings/" + off.Hash); err != nil {
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
			ProfessorHash:  off.Hash,
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

	private.PublishComment(ctx, &newComment)
}

// DeleteComment deletes a comment made by the given user
//
// It deletes not only the comment associated with the offering, but also the replica in user comments
func DeleteComment(
	ctx *gin.Context,
	DB db.Env,
	userID string,
	off *controllers.Offering,
) {
	subModel := models.NewSubjectFromController(&off.Subject)
	userHash := models.User{ID: userID}.Hash()

	// get comment ref from offerings
	commentRef := DB.Client.Doc(fmt.Sprintf(
		"subjects/%s/offerings/%s/comments/%s",
		subModel.Hash(),
		off.Hash,
		userHash,
	))

	// get user comment ref from user
	userCommentModel := models.UserComment{
		ProfessorHash:  off.Hash,
		Subject:        off.Code,
		Course:         off.CourseCode,
		Specialization: off.Specialization,
	}

	userCommentRef := DB.Client.Doc(fmt.Sprintf(
		"users/%s/user_comments/%s",
		userHash,
		userCommentModel.Hash(),
	))

	// Use batch write to delete documents atomically
	batch := DB.Client.Batch()
	batch.Delete(commentRef)
	batch.Delete(userCommentRef)

	if _, err := batch.Commit(ctx); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error deleting comment: %s", err.Error()))
		return
	}

	private.DeleteComment(ctx)
}
