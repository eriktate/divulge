package pg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/eriktate/divulge"
	"github.com/eriktate/divulge/pg"
	"github.com/google/uuid"
)

func Test_Users(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	hostname := "localhost"
	userName := "postgres"
	password := "password"
	db, err := pg.New(hostname, userName, password)
	if err != nil {
		t.Fatal(err)
	}

	user1 := divulge.User{
		Name:  "User One",
		Email: fmt.Sprintf("%s@test.com", uuid.New().String()),
	}

	user2 := divulge.User{
		Name:  "User Two",
		Email: fmt.Sprintf("%s@test.com", uuid.New().String()),
	}

	// RUN
	startingUsers, err := db.ListUsers(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing baseline: %s", err)
	}

	id1, err := db.SaveUser(ctx, user1)
	if err != nil {
		t.Fatalf("unexpected error creating user 1: %s", err)
	}

	id2, err := db.SaveUser(ctx, user2)
	if err != nil {
		t.Fatalf("unexpected error creating user 2: %s", err)
	}

	fetchedUser1, err := db.FetchUser(ctx, id1)
	if err != nil {
		t.Fatalf("unexpected error fetching user 1: %s", err)
	}

	allUsers, err := db.ListUsers(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing all users: %s", err)
	}

	if err := db.RemoveUser(ctx, id1); err != nil {
		t.Fatalf("unexpected error removing user1: %s", err)
	}

	if err := db.RemoveUser(ctx, id2); err != nil {
		t.Fatalf("unexpected error removing user2: %s", err)
	}

	// ASSERT
	if divulge.IsEmpty(id1) {
		t.Fatal("unexpected empty id1")
	}

	if divulge.IsEmpty(id2) {
		t.Fatal("unexpected empty id2")
	}

	newUserCount := len(allUsers) - len(startingUsers)
	if newUserCount != 2 {
		t.Fatalf("unexpected number of new users: %d", newUserCount)
	}

	if fetchedUser1.Name != user1.Name {
		t.Fatal("expected fetchedUser1.ID to equal user1.ID")
	}

	if fetchedUser1.Email != user1.Email {
		t.Fatal("expected fetchedUser1.ID to equal user1.ID")
	}
}
