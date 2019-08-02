package comments

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/comments/internal"
)

func CommentActorFromContext(ctx context.Context) (authorUserID int32, err error) {
	actor, err := graphqlbackend.CurrentUser(ctx)
	if err != nil {
		return 0, err
	}
	if actor == nil {
		return 0, errors.New("authenticated required to create comment")
	}
	return actor.DatabaseID(), nil
}

func (GraphQLResolver) CreateComment(ctx context.Context, arg *graphqlbackend.CreateCommentArgs) (graphqlbackend.Comment, error) {
	// TODO!(sqs): add auth checks
	authorUserID, err := CommentActorFromContext(ctx)
	if err != nil {
		return nil, err
	}

	v := &internal.DBComment{
		AuthorUserID: authorUserID,
		Body:         arg.Input.Body,
	}
	v.Object, err = commentObjectFromGQLID(arg.Input.Node)
	if err != nil {
		return nil, err
	}

	comment, err := internal.DBComments{}.Create(ctx, nil, v)
	if err != nil {
		return nil, err
	}
	return newGQLToComment(ctx, comment)
}

func (GraphQLResolver) EditComment(ctx context.Context, arg *graphqlbackend.EditCommentArgs) (graphqlbackend.Comment, error) {
	v, err := commentByGQLID(ctx, arg.Input.ID)
	if err != nil {
		return nil, err
	}
	comment, err := internal.DBComments{}.Update(ctx, v.ID, internal.DBCommentUpdate{
		Body: &arg.Input.Body,
	})
	if err != nil {
		return nil, err
	}
	return newGQLToComment(ctx, comment)
}

func (GraphQLResolver) DeleteComment(ctx context.Context, arg *graphqlbackend.DeleteCommentArgs) (*graphqlbackend.EmptyResponse, error) {
	v, err := commentByGQLID(ctx, arg.Comment)
	if err != nil {
		return nil, err
	}
	return nil, internal.DBComments{}.DeleteByID(ctx, v.ID)
}