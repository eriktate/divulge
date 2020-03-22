// +build integration

package pg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/eriktate/divulge"
	"github.com/eriktate/divulge/pg"
	"github.com/google/uuid"
)

func Test_Accounts(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	hostname := "localhost"
	username := "postgres"
	password := "password"
	db, err := pg.New(hostname, username, password)
	if err != nil {
		t.Fatal(err)
	}

	accountUser := divulge.User{
		Name:  "Account User",
		Email: fmt.Sprintf("%s@test.com", uuid.New().String()),
	}

	accountUserID, err := db.SaveUser(ctx, accountUser)
	if err != nil {
		t.Fatal(err)
	}

	account1 := divulge.Account{
		Name:    "Account One",
		OwnerID: accountUserID,
	}

	account2 := divulge.Account{
		Name:    "Account Two",
		OwnerID: accountUserID,
	}

	// RUN
	startingAccounts, err := db.ListAccounts(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing baseline: %s", err)
	}

	id1, err := db.SaveAccount(ctx, account1)
	if err != nil {
		t.Fatalf("unexpected error creating account 1: %s", err)
	}

	id2, err := db.SaveAccount(ctx, account2)
	if err != nil {
		t.Fatalf("unexpected error creating account 2: %s", err)
	}

	fetchedAccount1, err := db.FetchAccount(ctx, id1)
	if err != nil {
		t.Fatalf("unexpected error fetching account 1: %s", err)
	}

	allAccounts, err := db.ListAccounts(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing all accounts: %s", err)
	}

	if err := db.RemoveAccount(ctx, id1); err != nil {
		t.Fatalf("unexpected error removing account1: %s", err)
	}

	if err := db.RemoveAccount(ctx, id2); err != nil {
		t.Fatalf("unexpected error removing account2: %s", err)
	}

	// ASSERT
	if divulge.IsEmpty(id1) {
		t.Fatal("unexpected empty id1")
	}

	if divulge.IsEmpty(id2) {
		t.Fatal("unexpected empty id2")
	}

	newAccountCount := len(allAccounts) - len(startingAccounts)
	if newAccountCount != 2 {
		t.Fatalf("unexpected number of new accounts: %d", newAccountCount)
	}

	if fetchedAccount1.Name != account1.Name {
		t.Fatal("expected fetchedAccount1.ID to equal account1.ID")
	}

	if fetchedAccount1.OwnerID != accountUserID {
		t.Fatal("expected fetchedAccount1.ID to equal accountUserID")
	}
}
