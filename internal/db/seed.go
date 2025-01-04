package db

import (
	"context"
	"log"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/indiecodermm/go-social/internal/store"
)

const (
	usersCount    = 10
	postsCount    = 100
	commentsCount = 500
)

func Seed(store store.Store) {
	log.Println("Seeding database...")
	ctx := context.Background()

	users := generateUsers(usersCount)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("Error creating user: %v", err)
			return
		}
	}

	posts := generatePosts(postsCount, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("Error creating post: %v", err)
			return
		}
	}

	comments := generateComments(postsCount, posts, users)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("Error creating comment: %v", err)
			return
		}
	}

	log.Println("Database seeded")
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	for i := 0; i < count; i++ {
		users[i] = &store.User{
			Username: gofakeit.Gamertag(),
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, false, true, true, 6),
		}
	}
	return users
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   gofakeit.BookTitle(),
			Content: gofakeit.Quote(),
			Tags:    []string{gofakeit.Hobby(), gofakeit.HackerNoun()},
		}
	}
	return posts
}

func generateComments(count int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := 0; i < count; i++ {
		post := posts[rand.Intn(len(posts))]
		user := users[rand.Intn(len(users))]
		comments[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: gofakeit.HackerPhrase(),
		}
	}
	return comments
}
