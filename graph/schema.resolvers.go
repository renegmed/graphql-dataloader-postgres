package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/tinrab/curly-waddle/graph/generated"
	"github.com/tinrab/curly-waddle/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	user := &model.User{
		ID:   ksuid.New().String(),
		Name: input.Name,
	}
	log.Printf("++++ id: %v  name: %v\n", user.ID, user.Name)

	conn, err := r.Db.Conn(ctx)
	defer conn.Close() // Return the connection to the pool.

	_, err = conn.ExecContext(ctx, `INSERT INTO users (id, name) VALUES ($1, $2);`, user.ID, user.Name)
	if err != nil {
		log.Println("+++++ CreateUser error: ", err)
		return nil, err
	}
	return user, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	user := &model.User{}
	err := r.Db.
		QueryRowContext(ctx, "SELECT id, name FROM users WHERE id=$1", input.UserID).
		Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		Body:      input.Body,
		User:      user,
	}
	_, err = r.Db.ExecContext(
		ctx,
		"INSERT INTO posts(id, user_id, created_at, body) VALUES($1, $2, $3, $4)",
		post.ID,
		post.User.ID,
		post.CreatedAt,
		post.Body,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *queryResolver) Users(ctx context.Context, pagination *model.Pagination) ([]*model.User, error) {
	if pagination == nil {
		pagination = &model.Pagination{
			Skip: 0,
			Take: 100,
		}
	}

	rows, err := r.Db.QueryContext(
		ctx,
		"SELECT id, name FROM users OFFSET $1 LIMIT $2",
		pagination.Skip,
		pagination.Take,
	)
	if err != nil {
		return nil, err
	}

	var users []*model.User
	var userIDs []string

	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
		userIDs = append(userIDs, user.ID)
	}

	return users, nil
}

func (r *queryResolver) Posts(ctx context.Context, pagination *model.Pagination) ([]*model.Post, error) {
	if pagination == nil {
		pagination = &model.Pagination{
			Skip: 0,
			Take: 100,
		}
	}

	rows, err := r.Db.QueryContext(
		ctx,
		"SELECT id, user_id, created_at, body FROM posts OFFSET $1 LIMIT $2",
		pagination.Skip,
		pagination.Take,
	)
	if err != nil {
		return nil, err
	}

	var posts []*model.Post
	var userIDs []string
	for rows.Next() {

		post := &model.Post{
			User: &model.User{},
		}
		err = rows.Scan(&post.ID, &post.User.ID, &post.CreatedAt, &post.Body)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
		userIDs = append(userIDs, post.User.ID)
	}

	return posts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
