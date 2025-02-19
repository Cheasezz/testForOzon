package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/pkg/postgres"
	"github.com/google/uuid"
)

var (
	errStartTransaction   = errors.New("failed to start transaction")
	errCheckCmtAllowed    = errors.New("failed to check comments_allowed")
	errCmntDisabled       = errors.New("comments are disabled for this post")
	errFailedCommit       = errors.New("failed to commit transaction")
	errFailedParentCmnt   = errors.New("failed to check parent comment")
	errParentCmntNotExist = errors.New("parent comment does not exist")
)

type CommentRepo struct {
	db *postgres.Postgres
}

func NewCommentRepo(db *postgres.Postgres) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) CreateComment(ctx context.Context, comment core.Comment) (*core.Comment, error) {
	fmt.Println("CreateComment pg repo func call")

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errStartTransaction, err)
	}

	// Проверяем, если comment.ParentId != nil, то родительский комментарий должен существовать
	if comment.ParentId != nil {
		var parentExists bool

		query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)`, commentsTable)
		err = tx.QueryRow(ctx, query, comment.ParentId).Scan(&parentExists)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("%w: %w", errFailedParentCmnt, err)
		}
		if !parentExists {
			tx.Rollback(ctx)
			return nil, errParentCmntNotExist
		}
	}

	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, post_id, parent_id, created_at, content) 
	          VALUES ($1, $2, $3, $4, $5, $6) 
	          RETURNING *`, commentsTable)

	var createdComment core.Comment

	err = r.db.Scany.Get(ctx, tx, &createdComment, query,
		comment.Id, comment.UserId, comment.PostId, comment.ParentId, comment.CreatedAt, comment.Content)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%w: %w", errFailedCommit, err)
	}
	return &createdComment, nil
}

func (r *CommentRepo) CommentForPostAllowed(ctx context.Context, postId uuid.UUID) (bool, error) {
	fmt.Println("CommentForPostAllowed pg repo func call")

	query := fmt.Sprintf(`SELECT comments_allowed FROM %s WHERE id = $1`, postsTable)
	var commentAllowed bool
	err := r.db.Scany.Get(ctx, r.db.Pool, &commentAllowed, query, postId)
	if err != nil {
		return false, errPostDsntExist
	}

	return commentAllowed, nil
}

func (r *CommentRepo) GetRootComments(ctx context.Context, postId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("GetRootComments pg repo func call")

	query := fmt.Sprintf(`SELECT * FROM %s WHERE post_id = $1 AND parent_id IS NULL ORDER BY created_at ASC LIMIT $2 OFFSET $3`,
		commentsTable)
	var comments []*core.Comment
	err := r.db.Scany.Select(ctx, r.db.Pool, &comments, query, postId, limit, offset)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepo) GetRepliesById(ctx context.Context, parentCommentId uuid.UUID, limit, offset int) ([]*core.Comment, error) {
	fmt.Println("GetReplies pg repo func call")

	query := fmt.Sprintf(`SELECT * FROM %s WHERE  parent_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`,
		commentsTable)

	var comments []*core.Comment
	err := r.db.Scany.Select(ctx, r.db.Pool, &comments, query, parentCommentId, limit, offset)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepo) RepliesCount(ctx context.Context, commentId uuid.UUID) (int, error) {
	fmt.Println("RepliesCount pg repo func call")

	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE parent_id = $1`, commentsTable)
	var count int
	err := r.db.Pool.QueryRow(ctx, query, commentId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CommentRepo) GetRepliesCounts(ctx context.Context, ids []uuid.UUID) (map[string]int, error) {
	query := fmt.Sprintf(`
	SELECT parent_id, COUNT(*) as cnt
	FROM %s
	WHERE parent_id = ANY($1)
	GROUP BY parent_id;
`, commentsTable)

	rows, err := r.db.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countMap := make(map[string]int)
	for rows.Next() {
		var parentID uuid.UUID
		var cnt int
		if err := rows.Scan(&parentID, &cnt); err != nil {
			return nil, err
		}
		countMap[parentID.String()] = cnt
	}
	return countMap, nil
}
