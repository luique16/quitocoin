package user_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/domain/user"
	errorpkg "github.com/luique16/quitocoin/internal/error"
	"github.com/luique16/quitocoin/internal/provider"
)

func newService(repo user.Repository) user.Service {
	return user.NewService(repo, provider.NewPasswordHasher(), provider.NewIdGenerator())
}

// -- assertion helpers ---------------------------------------------------

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func assertErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %q, got %q", target, err)
	}
}

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func assertNotEqual[T comparable](t *testing.T, a, b T) {
	t.Helper()
	if a == b {
		t.Fatalf("expected different values, both are %v", a)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		t.Fatal("expected non-nil, got nil")
	}
}

func assertNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		return
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if rv.IsNil() {
			return
		}
	}
	t.Fatalf("expected nil, got %v", v)
}

func assertTrue(t *testing.T, v bool, msg string) {
	t.Helper()
	if !v {
		t.Fatalf("expected true: %s", msg)
	}
}

func assertDeepEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %#v, got %#v", expected, actual)
	}
}

func assertNotCalled(t *testing.T, name string, calls int) {
	t.Helper()
	if calls > 0 {
		t.Fatalf("expected %s not to be called, called %d time(s)", name, calls)
	}
}

// -- test data helpers ---------------------------------------------------

func now() time.Time {
	return time.Now().Truncate(time.Second)
}

func makeUser(overrides ...func(*ent.User)) *ent.User {
	u := &ent.User{
		ID:        "c0a80121-0000-4000-8000-000000000001",
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "$2a$10$hashedpassword",
		PublicID:  "pub_c0a80121",
		CreatedAt: now(),
	}
	for _, o := range overrides {
		o(u)
	}
	return u
}

func makeUsers(n int) []*ent.User {
	users := make([]*ent.User, n)
	for i := 0; i < n; i++ {
		idx := i
		users[idx] = makeUser(func(u *ent.User) {
			u.ID = uuidFromInt(idx)
			u.Email = emailFromInt(idx)
		})
	}
	return users
}

func uuidFromInt(i int) string {
	return "00000000-0000-0000-0000-" + padInt(i, 12)
}

func emailFromInt(i int) string {
	return "user" + itoa(i) + "@example.com"
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

func padInt(i int, n int) string {
	s := itoa(i)
	for len(s) < n {
		s = "0" + s
	}
	return s
}

func strPtr(s string) *string { return &s }

// -- Create tests --------------------------------------------------------

func TestCreate_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	input := user.CreateUserInput{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "securePass123!",
	}

	repo.CreateFn = func(_ context.Context, u *ent.User) (*ent.User, error) {
		assertNotEqual(t, u.Password, "securePass123!") // must be hashed
		assertEqual(t, u.Name, "Jane Doe")
		assertEqual(t, u.Email, "jane@example.com")
		assertTrue(t, u.ID != "", "id must be generated")
		assertTrue(t, u.PublicID != "", "public_id must be generated")
		assertTrue(t, !u.CreatedAt.IsZero(), "created_at must be set")

		return makeUser(func(u2 *ent.User) {
			u2.Name = "Jane Doe"
			u2.Email = "jane@example.com"
		}), nil
	}

	result, err := svc.Create(ctx, input)

	assertNoError(t, err)
	assertNotNil(t, result)
	assertEqual(t, "Jane Doe", result.Name)
	assertEqual(t, "jane@example.com", result.Email)
	assertTrue(t, result.ID != "", "id must be returned")
	assertTrue(t, result.PublicID != "", "public_id must be returned")
	assertTrue(t, !result.CreatedAt.IsZero(), "created_at must be returned")
	assertEqual(t, 1, repo.CreateCallCount())
}

func TestCreate_EmailAlreadyExists(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CreateFn = func(_ context.Context, _ *ent.User) (*ent.User, error) {
		return nil, errorpkg.ErrEmailExists
	}

	result, err := svc.Create(ctx, user.CreateUserInput{
		Name: "John Doe", Email: "existing@example.com", Password: "StrongP@ss1",
	})

	assertErrorIs(t, err, errorpkg.ErrEmailExists)
	assertNil(t, result)
	assertEqual(t, 1, repo.CreateCallCount())
}

func TestCreate_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.CreateFn = func(_ context.Context, _ *ent.User) (*ent.User, error) {
		return nil, errorpkg.ErrInternal
	}

	result, err := svc.Create(ctx, user.CreateUserInput{
		Name: "John Doe", Email: "john@example.com", Password: "StrongP@ss1",
	})

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertNil(t, result)
	assertEqual(t, 1, repo.CreateCallCount())
}

func TestCreate_EmptyName(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	_, err := svc.Create(context.Background(), user.CreateUserInput{
		Name: "", Email: "john@example.com", Password: "pass123",
	})

	assertError(t, err)
	assertEqual(t, 0, repo.CreateCallCount())
}

func TestCreate_EmptyEmail(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	_, err := svc.Create(context.Background(), user.CreateUserInput{
		Name: "John Doe", Email: "", Password: "pass123",
	})

	assertError(t, err)
	assertEqual(t, 0, repo.CreateCallCount())
}

func TestCreate_EmptyPassword(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	_, err := svc.Create(context.Background(), user.CreateUserInput{
		Name: "John Doe", Email: "john@example.com", Password: "",
	})

	assertError(t, err)
	assertEqual(t, 0, repo.CreateCallCount())
}

func TestCreate_InvalidEmail(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	invalidEmails := []struct {
		name  string
		email string
	}{
		{"missing @", "invalid-email"},
		{"missing domain", "user@"},
		{"missing local", "@domain.com"},
		{"double @", "user@@domain.com"},
		{"spaces", "user @domain.com"},
	}

	for _, tt := range invalidEmails {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), user.CreateUserInput{
				Name: "John Doe", Email: tt.email, Password: "pass123",
			})
			assertError(t, err)
		})
	}

	assertEqual(t, 0, repo.CreateCallCount())
}

func TestCreate_WeakPassword(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	weakPasswords := []struct {
		name string
		pwd  string
	}{
		{"too short", "Ab1!"},
		{"no uppercase", "abcdef1!"},
		{"no lowercase", "ABCDEF1!"},
		{"no digit", "Abcdefgh!"},
		{"no special char", "Abcdefgh1"},
		{"lowercase only", "abcdefgh"},
		{"spaces only", "        "},
	}

	for _, tt := range weakPasswords {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), user.CreateUserInput{
				Name: "John Doe", Email: "john@example.com", Password: tt.pwd,
			})
			assertErrorIs(t, err, errorpkg.ErrWeakPassword)
		})
	}

	assertEqual(t, 0, repo.CreateCallCount())
}

// -- Get tests -----------------------------------------------------------

func TestGet_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	expected := makeUser()
	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		assertEqual(t, expected.ID, id)
		return expected, nil
	}

	result, err := svc.Get(ctx, expected.ID)

	assertNoError(t, err)
	assertDeepEqual(t, expected, result)
	assertEqual(t, 1, repo.GetCallCount())
}

func TestGet_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetFn = func(_ context.Context, _ string) (*ent.User, error) {
		return nil, errorpkg.ErrUserNotFound
	}

	result, err := svc.Get(ctx, "nonexistent-id")

	assertErrorIs(t, err, errorpkg.ErrUserNotFound)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetCallCount())
}

func TestGet_EmptyID(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	_, err := svc.Get(context.Background(), "")

	assertError(t, err)
	assertEqual(t, 0, repo.GetCallCount())
}

func TestGet_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetFn = func(_ context.Context, _ string) (*ent.User, error) {
		return nil, errorpkg.ErrInternal
	}

	result, err := svc.Get(ctx, "some-id")

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetCallCount())
}

// -- List tests ----------------------------------------------------------

func TestList_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	expected := makeUsers(3)
	repo.ListFn = func(_ context.Context) ([]*ent.User, error) {
		return expected, nil
	}

	result, err := svc.List(ctx)

	assertNoError(t, err)
	if len(result) != 3 {
		t.Fatalf("expected 3 users, got %d", len(result))
	}
	assertDeepEqual(t, expected, result)
	assertEqual(t, 1, repo.ListCallCount())
}

func TestList_Empty(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.ListFn = func(_ context.Context) ([]*ent.User, error) {
		return []*ent.User{}, nil
	}

	result, err := svc.List(ctx)

	assertNoError(t, err)
	assertEqual(t, 0, len(result))
	assertEqual(t, 1, repo.ListCallCount())
}

func TestList_Nil(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.ListFn = func(_ context.Context) ([]*ent.User, error) {
		return nil, nil
	}

	result, err := svc.List(ctx)

	assertNoError(t, err)
	assertNil(t, result)
	assertEqual(t, 1, repo.ListCallCount())
}

func TestList_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.ListFn = func(_ context.Context) ([]*ent.User, error) {
		return nil, errorpkg.ErrInternal
	}

	result, err := svc.List(ctx)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertNil(t, result)
	assertEqual(t, 1, repo.ListCallCount())
}

// -- Update tests --------------------------------------------------------

func TestUpdate_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		assertEqual(t, existing.ID, id)
		return existing, nil
	}

	repo.UpdateFn = func(_ context.Context, u *ent.User) (*ent.User, error) {
		assertEqual(t, existing.ID, u.ID)
		assertEqual(t, "Jane Updated", u.Name)
		assertEqual(t, "jane.updated@example.com", u.Email)
		assertEqual(t, existing.Password, u.Password)
		return makeUser(func(u2 *ent.User) {
			u2.ID = existing.ID
			u2.Name = "Jane Updated"
			u2.Email = "jane.updated@example.com"
		}), nil
	}

	result, err := svc.Update(ctx, existing.ID, user.UpdateUserInput{
		Name:  strPtr("Jane Updated"),
		Email: strPtr("jane.updated@example.com"),
	})

	assertNoError(t, err)
	assertEqual(t, "Jane Updated", result.Name)
	assertEqual(t, "jane.updated@example.com", result.Email)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 1, repo.UpdateCallCount())
}

func TestUpdate_PartialName(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	repo.UpdateFn = func(_ context.Context, u *ent.User) (*ent.User, error) {
		assertEqual(t, existing.ID, u.ID)
		assertEqual(t, "Only Name Change", u.Name)
		assertEqual(t, existing.Email, u.Email)
		assertEqual(t, existing.Password, u.Password)
		return makeUser(func(u2 *ent.User) {
			u2.Name = "Only Name Change"
		}), nil
	}

	result, err := svc.Update(ctx, existing.ID, user.UpdateUserInput{
		Name: strPtr("Only Name Change"),
	})

	assertNoError(t, err)
	assertEqual(t, "Only Name Change", result.Name)
	assertEqual(t, existing.Email, result.Email)
}

func TestUpdate_PartialEmail(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	repo.UpdateFn = func(_ context.Context, u *ent.User) (*ent.User, error) {
		assertEqual(t, existing.ID, u.ID)
		assertEqual(t, existing.Name, u.Name)
		assertEqual(t, "new.email@example.com", u.Email)
		assertEqual(t, existing.Password, u.Password)
		return makeUser(func(u2 *ent.User) {
			u2.Email = "new.email@example.com"
		}), nil
	}

	result, err := svc.Update(ctx, existing.ID, user.UpdateUserInput{
		Email: strPtr("new.email@example.com"),
	})

	assertNoError(t, err)
		assertEqual(t, existing.Name, result.Name)
	assertEqual(t, "new.email@example.com", result.Email)
}

func TestUpdate_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetFn = func(_ context.Context, _ string) (*ent.User, error) {
		return nil, errorpkg.ErrUserNotFound
	}

	result, err := svc.Update(ctx, "nonexistent-id", user.UpdateUserInput{
		Name: strPtr("Ghost User"),
	})

	assertErrorIs(t, err, errorpkg.ErrUserNotFound)
	assertNil(t, result)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 0, repo.UpdateCallCount())
}

func TestUpdate_EmptyID(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	_, err := svc.Update(context.Background(), "", user.UpdateUserInput{
		Name: strPtr("No Matter"),
	})

	assertError(t, err)
	assertEqual(t, 0, repo.GetCallCount())
	assertEqual(t, 0, repo.UpdateCallCount())
}

func TestUpdate_NoFields(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	repo.UpdateFn = func(_ context.Context, u *ent.User) (*ent.User, error) {
		assertDeepEqual(t, existing, u)
		return existing, nil
	}

	result, err := svc.Update(ctx, existing.ID, user.UpdateUserInput{})

	assertNoError(t, err)
	assertNotNil(t, result)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 1, repo.UpdateCallCount())
}

func TestUpdate_EmailConflict(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	repo.UpdateFn = func(_ context.Context, _ *ent.User) (*ent.User, error) {
		return nil, errorpkg.ErrEmailExists
	}

	result, err := svc.Update(ctx, existing.ID, user.UpdateUserInput{
		Email: strPtr("taken@example.com"),
	})

	assertErrorIs(t, err, errorpkg.ErrEmailExists)
	assertNil(t, result)
}

func TestUpdate_InvalidEmail(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	invalidEmails := []struct {
		name  string
		email string
	}{
		{"missing @", "invalid-email"},
		{"missing domain", "user@"},
		{"missing local", "@domain.com"},
		{"double @", "user@@domain.com"},
		{"spaces", "user @domain.com"},
		{"empty string", ""},
	}

	for _, tt := range invalidEmails {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.Update(ctx, existing.ID, user.UpdateUserInput{
				Email: strPtr(tt.email),
			})
			assertError(t, err)
			assertNil(t, result)
		})
	}

	assertEqual(t, 0, repo.UpdateCallCount())
}

// -- UpdatePassword tests ------------------------------------------------

func TestUpdatePassword_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()
	hasher := provider.NewPasswordHasher()

	existing := makeUser()
	oldHash, _ := hasher.Hash("OldPass123!")
	existing.Password = oldHash

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		assertEqual(t, existing.ID, id)
		return existing, nil
	}

	repo.UpdateFn = func(_ context.Context, u *ent.User) (*ent.User, error) {
		assertEqual(t, existing.ID, u.ID)
		assertTrue(t, u.Password != oldHash, "password must change")
		assertTrue(t, hasher.Compare("NewStr0ng!", u.Password) == nil, "new password must be hashed correctly")
		existing.Password = u.Password
		return existing, nil
	}

	err := svc.UpdatePassword(ctx, existing.ID, user.UpdatePasswordInput{
		OldPassword: "OldPass123!",
		NewPassword: "NewStr0ng!",
	})

	assertNoError(t, err)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 1, repo.UpdateCallCount())
}

func TestUpdatePassword_InvalidOldPassword(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()
	hasher := provider.NewPasswordHasher()

	existing := makeUser()
	existing.Password, _ = hasher.Hash("OldPass123!")

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	err := svc.UpdatePassword(ctx, existing.ID, user.UpdatePasswordInput{
		OldPassword: "WrongOldPass!",
		NewPassword: "NewStr0ng!",
	})

	assertErrorIs(t, err, errorpkg.ErrIncorrectPassword)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 0, repo.UpdateCallCount())
}

func TestUpdatePassword_WeakNewPassword(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()
	hasher := provider.NewPasswordHasher()

	existing := makeUser()
	existing.Password, _ = hasher.Hash("OldPass123!")

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	weakPasswords := []struct {
		name string
		pwd  string
	}{
		{"too short", "Ab1!"},
		{"no uppercase", "abcdef1!"},
		{"no lowercase", "ABCDEF1!"},
		{"no digit", "Abcdefgh!"},
		{"no special char", "Abcdefgh1"},
	}

	for _, tt := range weakPasswords {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.UpdatePassword(ctx, existing.ID, user.UpdatePasswordInput{
				OldPassword: "OldPass123!",
				NewPassword: tt.pwd,
			})
			assertErrorIs(t, err, errorpkg.ErrWeakPassword)
		})
	}

	assertEqual(t, 0, repo.UpdateCallCount())
}

func TestUpdatePassword_EmptyID(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	err := svc.UpdatePassword(context.Background(), "", user.UpdatePasswordInput{
		OldPassword: "OldPass123!",
		NewPassword: "NewStr0ng!",
	})

	assertError(t, err)
	assertEqual(t, 0, repo.GetCallCount())
	assertEqual(t, 0, repo.UpdateCallCount())
}

func TestUpdatePassword_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetFn = func(_ context.Context, _ string) (*ent.User, error) {
		return nil, errorpkg.ErrUserNotFound
	}

	err := svc.UpdatePassword(ctx, "nonexistent-id", user.UpdatePasswordInput{
		OldPassword: "OldPass123!",
		NewPassword: "NewStr0ng!",
	})

	assertErrorIs(t, err, errorpkg.ErrUserNotFound)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 0, repo.UpdateCallCount())
}

func TestUpdatePassword_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()
	hasher := provider.NewPasswordHasher()

	existing := makeUser()
	existing.Password, _ = hasher.Hash("OldPass123!")

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}

	repo.UpdateFn = func(_ context.Context, _ *ent.User) (*ent.User, error) {
		return nil, errorpkg.ErrInternal
	}

	err := svc.UpdatePassword(ctx, existing.ID, user.UpdatePasswordInput{
		OldPassword: "OldPass123!",
		NewPassword: "NewStr0ng!",
	})

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 1, repo.UpdateCallCount())
}

// -- Delete tests --------------------------------------------------------

func TestDelete_Success(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()

	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		assertEqual(t, existing.ID, id)
		return existing, nil
	}

	repo.DeleteFn = func(_ context.Context, id string) error {
		assertEqual(t, existing.ID, id)
		return nil
	}

	err := svc.Delete(ctx, existing.ID)

	assertNoError(t, err)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 1, repo.DeleteCallCount())
}

func TestDelete_NotFound(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	repo.GetFn = func(_ context.Context, _ string) (*ent.User, error) {
		return nil, errorpkg.ErrUserNotFound
	}

	err := svc.Delete(ctx, "nonexistent-id")

	assertErrorIs(t, err, errorpkg.ErrUserNotFound)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 0, repo.DeleteCallCount())
}

func TestDelete_EmptyID(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)

	err := svc.Delete(context.Background(), "")

	assertError(t, err)
	assertEqual(t, 0, repo.GetCallCount())
	assertEqual(t, 0, repo.DeleteCallCount())
}

func TestDelete_RepositoryError(t *testing.T) {
	repo := NewMockRepository()
	svc := newService(repo)
	ctx := context.Background()

	existing := makeUser()
	repo.GetFn = func(_ context.Context, id string) (*ent.User, error) {
		return existing, nil
	}
	repo.DeleteFn = func(_ context.Context, _ string) error {
		return errorpkg.ErrInternal
	}

	err := svc.Delete(ctx, existing.ID)

	assertErrorIs(t, err, errorpkg.ErrInternal)
	assertEqual(t, 1, repo.GetCallCount())
	assertEqual(t, 1, repo.DeleteCallCount())
}
