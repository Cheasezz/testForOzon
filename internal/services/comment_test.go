package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Cheasezz/testForOzon/internal/core"
	"github.com/Cheasezz/testForOzon/internal/errs"
	"github.com/Cheasezz/testForOzon/internal/repositories"
	"github.com/Cheasezz/testForOzon/pkg/logger"
	"github.com/Cheasezz/testForOzon/pkg/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type cmntSrvsDeps struct {
	r  *repositories.MockCommentRepo
	ps *pubsub.MockIPubSub
	l  *logger.MockLogger
}

func initDeps(c *gomock.Controller) cmntSrvsDeps {
	return cmntSrvsDeps{
		r:  repositories.NewMockCommentRepo(c),
		ps: pubsub.NewMockIPubSub(c),
		l:  logger.NewMockLogger(c),
	}
}

func newCommentCreateInput() core.CommentCreateInput {
	postId := uuid.New()
	parentId := uuid.New()
	return core.CommentCreateInput{
		UserId:   "UserNickName",
		PostId:   postId,
		ParentId: &parentId,
		Content:  "TestContentForComment",
	}
}

func newComment(input core.CommentCreateInput) *core.Comment {
	return &core.Comment{
		PostId:    input.PostId,
		Id:        uuid.New(),
		ParentId:  input.ParentId,
		UserId:    input.UserId,
		CreatedAt: time.Now(),
		Content:   input.Content,
	}
}

func TestCommentService_CreateComment(t *testing.T) {
	type mockBehavior func(d cmntSrvsDeps, input core.CommentCreateInput)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := initDeps(ctrl)

	cmntSrvs := NewCommentService(d.r, d.ps, d.l)
	input := newCommentCreateInput()
	cmnt := newComment(input)

	tests := []struct {
		name         string
		input        core.CommentCreateInput
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name:  "ok",
			input: input,
			mockBehavior: func(d cmntSrvsDeps, input core.CommentCreateInput) {
				d.r.EXPECT().CommentForPostAllowed(gomock.Any(), input.PostId).Return(true, nil)
				d.r.EXPECT().CreateComment(gomock.Any(), gomock.Any()).Return(cmnt, nil)
				d.ps.EXPECT().Publish(gomock.Any())
			},
		},
		{
			name:    "comments are prohibited",
			input:   input,
			wantErr: errCmntAreProh,
			mockBehavior: func(d cmntSrvsDeps, input core.CommentCreateInput) {
				d.r.EXPECT().CommentForPostAllowed(gomock.Any(), input.PostId).Return(false, nil)
			},
		},
		{
			name: "comment is too long",
			input: core.CommentCreateInput{
				UserId:   input.UserId,
				PostId:   input.PostId,
				ParentId: input.ParentId,
				Content:  strings.Repeat("a", 2001),
			},
			wantErr: errToLongtext,
			mockBehavior: func(d cmntSrvsDeps, input core.CommentCreateInput) {
				d.r.EXPECT().CommentForPostAllowed(gomock.Any(), input.PostId).Return(true, nil)
			},
		},
		// Не знаю, хорошая ли идея писать такие тесты??
		// 
		// 
		{
			name:    "err in repo CommentForPostAllowed",
			input:   input,
			wantErr: errs.ErrPostDsntExist,
			mockBehavior: func(d cmntSrvsDeps, input core.CommentCreateInput) {
				d.r.EXPECT().CommentForPostAllowed(gomock.Any(), input.PostId).Return(false, errs.ErrPostDsntExist)
				d.l.EXPECT().Error(gomock.Any(), gomock.Any())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(d, tt.input)

			cmnt, err := cmntSrvs.CreateComment(context.Background(), tt.input)
			if err != nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, cmnt)
			}
		})
	}
}
