package usecase

import (
	"SimpleForum/pkg/logger"
)

// TO DO
/*
1) Check whether there is reaction of reacted  post/comment in db already
	True:
		Is that implemented reaction same as the in db's reaction
			True:
				Delete that reaction from the DB
			False:
				Delete that opposite reaction from the DB
				Insert that current reaction to the DB
	False:
		Insert that current reaction to the DB
2) Impelement corresponding actions after executing the previous condition to Post/Comment table in DB


*/
func (app *Application) ExecutionOfReactionLD(userId int, identifier int, postOrcomment, reactionLD string) error {

	dbreaction, err := app.ServiceDB.RetrieveExistenceOfReactionLD(userId, identifier, postOrcomment)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem getting the reaction of the post from the db", err)
	}

	if dbreaction == nil {
		err = app.ServiceDB.InsertReaction(userId, identifier, postOrcomment, reactionLD)
		if err != nil {
			return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem creating the reaction for postOrComment in the db", err)
		}
		if postOrcomment == "post" {
			err = app.ServiceDB.UpdateReactionOfPost(identifier, reactionLD, "increment")
			if err != nil {
				return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem incrementing the number to the reaction counter of the post in the db", err)
			}
		} else if postOrcomment == "comment" {
			err = app.ServiceDB.UpdateReactionOfCommentary(identifier, reactionLD, "increment")
			if err != nil {
				return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem incrementing the number to the reaction counter of the post in the db", err)
			}
		}
	} else {
		var dbaction string

		if dbreaction.Action == "L" {
			dbaction = "like"
		} else if dbreaction.Action == "D" {
			dbaction = "dislike"
		}

		err = app.ServiceDB.DeleteReaction(userId, identifier, postOrcomment, dbaction)
		if err != nil {
			return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem deleting the reaction for postOrComment in the db", err)
		}

		if reactionLD != dbaction {
			err = app.ServiceDB.InsertReaction(userId, identifier, postOrcomment, reactionLD)
			if err != nil {
				return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem inserting the reaction for postOrComment in the db", err)
			}
		}

		if postOrcomment == "post" {
			err = app.ServiceDB.UpdateReactionOfPost(identifier, dbaction, "decrement")
			if err != nil {
				return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem incrementing the number to the reaction counter of the post in the db", err)
			}

			if reactionLD != dbaction {
				err = app.ServiceDB.UpdateReactionOfPost(identifier, reactionLD, "increment")
				if err != nil {
					return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem incrementing the number to the reaction counter of the post in the db", err)
				}
			}

		} else if postOrcomment == "comment" {

			err = app.ServiceDB.UpdateReactionOfCommentary(identifier, dbaction, "decrement")
			if err != nil {
				return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem incrementing the number to the reaction counter of the post in the db", err)
			}

			if reactionLD != dbaction {
				err = app.ServiceDB.UpdateReactionOfCommentary(identifier, reactionLD, "increment")
				if err != nil {
					return logger.ErrorWrapper("UseCase", "ExecutionOfReaction", "There is a problem incrementing the number to the reaction counter of the post in the db", err)
				}
			}

		}
	}
	return nil
}
