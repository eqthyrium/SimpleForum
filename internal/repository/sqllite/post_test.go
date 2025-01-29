package sqllite_test

import (
	"SimpleForum/internal/repository/sqllite"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// **Helper function** to set up an in-memory SQLite database for testing.
func setupTestDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "../../../mydb.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// **Test for CreatePost**
func TestCreatePost(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}

	postID, err := repo.CreatePost(1, "Test Title", "Test Content", "")
	if err != nil {
		t.Errorf("Failed to create post: %v", err)
	}

	if postID <= 0 {
		t.Errorf("Expected valid post ID, got %d", postID)
	}
}

// **Test for GetLatestAllPosts**

func TestGetLatestAllPosts(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	_, _ = repo.CreatePost(1, "First Post", "Content", "")

	_, err = repo.GetLatestAllPosts([]string{})
	if err != nil {
		t.Errorf("Failed to fetch latest posts: %v", err)
	}

}

// **Test for GetPostsByCertainUser**

func TestGetPostsByCertainUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	_, _ = repo.CreatePost(1, "User Post", "User Content", "")

	_, err = repo.GetPostsByCertainUser(1)
	if err != nil {
		t.Errorf("Failed to fetch posts by user: %v", err)
	}

}

// **Test for GetCertainPostInfo**

func TestGetCertainPostInfo(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	postID, _ := repo.CreatePost(1, "Test Title", "Test Content", "")

	post, err := repo.GetCertainPostInfo(postID)
	if err != nil {
		t.Errorf("Failed to get post info: %v", err)
	}

	if post.PostId != postID {
		t.Errorf("Expected post ID %d, got %d", postID, post.PostId)
	}
}

// **Test for UpdateReactionOfPost**

func TestUpdateReactionOfPost(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	postID, _ := repo.CreatePost(1, "Title", "Content", "")

	err = repo.UpdateReactionOfPost(postID, "like", "increment")
	if err != nil {
		t.Errorf("Failed to update reaction: %v", err)
	}
}

// **Test for UpdateEditedPost**

func TestUpdateEditedPost(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	postID, _ := repo.CreatePost(1, "Title", "Content", "")

	err = repo.UpdateEditedPost(1, postID, "Updated Content")
	if err != nil {
		t.Errorf("Failed to update post: %v", err)
	}
}

// **Test for DeleteCertainPost**

func TestDeleteCertainPost(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	postID, _ := repo.CreatePost(1, "Title", "Content", "")

	err = repo.DeleteCertainPost(postID)
	if err != nil {
		t.Errorf("Failed to delete post: %v", err)
	}
}

// **Test for ValidateOfExistenceCertainPost**

func TestValidateOfExistenceCertainPost(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	postID, _ := repo.CreatePost(1, "Title", "Content", "")

	exists, err := repo.ValidateOfExistenceCertainPost(1, postID)
	if err != nil {
		t.Errorf("Failed to check post existence: %v", err)
	}

	if !exists {
		t.Errorf("Expected post to exist, but it doesn't")
	}
}

// **Test for GetMyCommentedPosts**

func TestGetMyCommentedPosts(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	repo := &sqllite.Repository{DB: db}
	postID, _ := repo.CreatePost(1, "Title", "Content", "")
	_, _ = db.Exec("INSERT INTO Commentaries (UserId, PostId, Content) VALUES (1, ?, 'Comment')", postID)

	_, err = repo.GetMyCommentedPosts(1)
	if err != nil {
		t.Errorf("Failed to fetch commented posts: %v", err)
	}

}
