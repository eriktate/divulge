package divulge

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var zeroUUID uuid.UUID

// An Account is the owner of a blog/publication.
type Account struct {
	ID        uuid.UUID  `json:"id,omitempty" db:"id"`
	OwnerID   uuid.UUID  `json:"ownerId" db:"owner_id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

// A User is a member of an account. Responsible for creating blogs.
type User struct {
	ID        uuid.UUID
	Accounts  []uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

// A Post is just a blog post.
type Post struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	AccountID   uuid.UUID  `json:"accountId" db:"account_id"`
	AuthorID    uuid.UUID  `json:"authorId" db:"author_id"`
	Title       string     `json:"title" db:"title"`
	Summary     string     `json:"summary" db:"summary"`
	ContentPath string     `json:"contentPath" db:"content_path"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" db:"published_at"`
}

// An AccountService knows how to work with Accounts.
type AccountService interface {
	SaveAccount(ctx context.Context, account Account) (uuid.UUID, error)
	FetchAccount(ctx context.Context, id uuid.UUID) (Account, error)
	ListAccounts(ctx context.Context) ([]Account, error)
	RemoveAccount(ctx context.Context, id uuid.UUID) error
}

// A UserService knows how to work with Users.
type UserService interface {
	SaveUser(ctx context.Context, user User) (uuid.UUID, error)
	FetchUser(ctx context.Context, id uuid.UUID) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	RemoveUser(ctx context.Context, id uuid.UUID) error
}

// IsEmpty returns true if the id provided is empty.
func IsEmpty(id uuid.UUID) bool {
	return id == zeroUUID
}