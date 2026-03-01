package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/alfattd/category-service/internal/domain"
	"github.com/alfattd/category-service/internal/repository"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var sharedDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fmt.Printf("failed to start postgres container: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate container: %v\n", err)
		}
	}()

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Printf("failed to get connection string: %v\n", err)
		os.Exit(1)
	}

	sharedDB, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer sharedDB.Close()

	if err := sharedDB.PingContext(ctx); err != nil {
		fmt.Printf("failed to ping db: %v\n", err)
		os.Exit(1)
	}

	if err := runMigrations(sharedDB); err != nil {
		fmt.Printf("failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func runMigrations(db *sql.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..", "..")
	migrationFile := filepath.Join(projectRoot, "postgres", "migrations", "000001_create_categories_table.up.sql")

	migration, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
	}

	if _, err := db.Exec(string(migration)); err != nil {
		return fmt.Errorf("failed to exec migration: %w", err)
	}

	return nil
}

func cleanupTable(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		_, err := sharedDB.Exec("DELETE FROM categories")
		if err != nil {
			t.Logf("failed to cleanup table: %v", err)
		}
	})
}

func newCategory(name string) *domain.Category {
	now := time.Now().UTC()
	return &domain.Category{
		ID:        fmt.Sprintf("test-id-%d", now.UnixNano()),
		Name:      name,
		CreatedAt: now.Truncate(time.Microsecond),
		UpdatedAt: now.Truncate(time.Microsecond),
	}
}

// ─── Create ───────────────────────────────────────────────────────────────────

func TestRepoCreate_Success(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat := newCategory("Electronics")

	err := repo.Create(ctx, cat)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, cat.ID, got.ID)
	assert.Equal(t, cat.Name, got.Name)
}

func TestRepoCreate_DuplicateName_ReturnsErrDuplicate(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat1 := newCategory("Electronics")
	require.NoError(t, repo.Create(ctx, cat1))

	cat2 := newCategory("Electronics")
	cat2.ID = "different-id"

	err := repo.Create(ctx, cat2)
	assert.ErrorIs(t, err, domain.ErrDuplicate)
}

func TestRepoCreate_DuplicateID_ReturnsErrDuplicate(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat := newCategory("Electronics")
	require.NoError(t, repo.Create(ctx, cat))

	cat.Name = "Books"
	err := repo.Create(ctx, cat)
	assert.ErrorIs(t, err, domain.ErrDuplicate)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestRepoUpdate_Success(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat := newCategory("Electronics")
	require.NoError(t, repo.Create(ctx, cat))

	cat.Name = "Updated Electronics"
	cat.UpdatedAt = time.Now().UTC().Truncate(time.Microsecond)

	err := repo.Update(ctx, cat)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Electronics", got.Name)
}

func TestRepoUpdate_NotFound_ReturnsErrNotFound(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat := newCategory("Electronics")
	cat.ID = "id-yang-tidak-ada"

	err := repo.Update(ctx, cat)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRepoUpdate_DuplicateName_ReturnsErrDuplicate(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat1 := newCategory("Electronics")
	require.NoError(t, repo.Create(ctx, cat1))

	cat2 := newCategory("Books")
	require.NoError(t, repo.Create(ctx, cat2))

	cat2.Name = "Electronics"
	err := repo.Update(ctx, cat2)
	assert.ErrorIs(t, err, domain.ErrDuplicate)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

func TestRepoDelete_Success(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat := newCategory("Electronics")
	require.NoError(t, repo.Create(ctx, cat))

	err := repo.Delete(ctx, cat.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, cat.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRepoDelete_NotFound_ReturnsErrNotFound(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	err := repo.Delete(ctx, "id-yang-tidak-ada")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// ─── GetByID ──────────────────────────────────────────────────────────────────

func TestRepoGetByID_Success(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat := newCategory("Electronics")
	require.NoError(t, repo.Create(ctx, cat))

	got, err := repo.GetByID(ctx, cat.ID)
	require.NoError(t, err)
	assert.Equal(t, cat.ID, got.ID)
	assert.Equal(t, cat.Name, got.Name)
	assert.Equal(t, cat.CreatedAt, got.CreatedAt.UTC())
	assert.Equal(t, cat.UpdatedAt, got.UpdatedAt.UTC())
}

func TestRepoGetByID_NotFound_ReturnsErrNotFound(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "id-yang-tidak-ada")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// ─── List ─────────────────────────────────────────────────────────────────────

func TestRepoList_Success(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	require.NoError(t, repo.Create(ctx, newCategory("Electronics")))
	time.Sleep(2 * time.Millisecond)
	require.NoError(t, repo.Create(ctx, newCategory("Books")))

	result, err := repo.List(ctx, domain.PaginationParams{Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRepoList_OrderByCreatedAtDesc(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	cat1 := newCategory("First")
	require.NoError(t, repo.Create(ctx, cat1))
	time.Sleep(2 * time.Millisecond)

	cat2 := newCategory("Second")
	require.NoError(t, repo.Create(ctx, cat2))

	result, err := repo.List(ctx, domain.PaginationParams{Page: 1, Limit: 10})
	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, "Second", result[0].Name)
	assert.Equal(t, "First", result[1].Name)
}

func TestRepoList_Pagination_ReturnsCorrectPage(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	names := []string{"A", "B", "C", "D", "E"}
	for _, name := range names {
		require.NoError(t, repo.Create(ctx, newCategory(name)))
		time.Sleep(2 * time.Millisecond)
	}

	// page 1, limit 2 → should return 2 items (most recent: E, D)
	page1, err := repo.List(ctx, domain.PaginationParams{Page: 1, Limit: 2})
	require.NoError(t, err)
	assert.Len(t, page1, 2)
	assert.Equal(t, "E", page1[0].Name)
	assert.Equal(t, "D", page1[1].Name)

	// page 2, limit 2 → should return 2 items (C, B)
	page2, err := repo.List(ctx, domain.PaginationParams{Page: 2, Limit: 2})
	require.NoError(t, err)
	assert.Len(t, page2, 2)
	assert.Equal(t, "C", page2[0].Name)
	assert.Equal(t, "B", page2[1].Name)

	// page 3, limit 2 → should return 1 item (A)
	page3, err := repo.List(ctx, domain.PaginationParams{Page: 3, Limit: 2})
	require.NoError(t, err)
	assert.Len(t, page3, 1)
	assert.Equal(t, "A", page3[0].Name)
}

func TestRepoList_Empty_ReturnsEmptySlice(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	result, err := repo.List(ctx, domain.PaginationParams{Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result)
}

// ─── Count ────────────────────────────────────────────────────────────────────

func TestRepoCount_ReturnsCorrectTotal(t *testing.T) {
	cleanupTable(t)
	repo := repository.NewPostgresCategoryRepo(sharedDB)
	ctx := context.Background()

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	require.NoError(t, repo.Create(ctx, newCategory("Electronics")))
	require.NoError(t, repo.Create(ctx, newCategory("Books")))
	require.NoError(t, repo.Create(ctx, newCategory("Fashion")))

	count, err = repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}
