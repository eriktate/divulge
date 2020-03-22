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
	ContentPath string     `json:"contentPath,omitempty" db:"content_path"`
	Content     string     `json:"content,omitempty" db:"-"`
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

// A PostService knows how to work with Posts.
type PostService interface {
	SavePost(ctx context.Context, post Post) (uuid.UUID, error)
	PublishPost(ctx context.Context, id uuid.UUID) error
	RedactPost(ctx context.Context, id uuid.UUID) error

	FetchPost(ctx context.Context, id uuid.UUID) (Post, error)
	ListPostsByAccount(ctx context.Context, accountID uuid.UUID) ([]Post, error)

	RemovePost(ctx context.Context, id uuid.UUID) error
}

// FileStore knows how to work with post content.
type FileStore interface {
	Write(ctx context.Context, key string, data []byte) error
	Read(ctx context.Context, key string) ([]byte, error)
}

// IsEmpty returns true if the id provided is empty.
func IsEmpty(id uuid.UUID) bool {
	return id == zeroUUID
}
