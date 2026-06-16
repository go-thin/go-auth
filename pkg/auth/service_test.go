package auth

import (
	"errors"
	"testing"

	goevents "github.com/go-thin/go-events"
	"github.com/go-thin/go-auth/pkg/models"
)

// recordingBus captures every event passed to Publish.
type recordingBus struct {
	published []goevents.Event
}

func (b *recordingBus) Publish(e goevents.Event) []error {
	b.published = append(b.published, e)
	return []error{}
}

func (b *recordingBus) Subscribe(topic string, h func(goevents.Envelope) error) goevents.Subscription {
	return noopSub{}
}

func (b *recordingBus) SubscribeAll(h func(goevents.Envelope) error) goevents.Subscription {
	return noopSub{}
}

type noopSub struct{}

func (noopSub) Unsubscribe() {}

// simpleStorage is a minimal in-memory storage for service tests.
type simpleStorage struct {
	users map[string]*models.User
}

func newSimpleStorage() *simpleStorage {
	return &simpleStorage{users: make(map[string]*models.User)}
}

func (s *simpleStorage) CreateUser(user models.User) error {
	s.users[user.Username] = &user
	return nil
}

func (s *simpleStorage) GetUserByUsername(username string) (*models.User, error) {
	u, ok := s.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func newTestService(t *testing.T, bus goevents.Bus) *Service {
	t.Helper()
	svc, err := New(Config{
		Storage: newSimpleStorage(),
		JWT: JWTConfig{
			AccessSecret:  []byte("access-secret"),
			RefreshSecret: []byte("refresh-secret"),
			SigningMethod:  HS256,
		},
		EventBus: bus,
	})
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}
	return svc
}

func TestUserRegisteredEvent_Topic(t *testing.T) {
	e := UserRegistered{}
	if e.Topic() != "auth.user.registered" {
		t.Errorf("Topic() = %q, want %q", e.Topic(), "auth.user.registered")
	}
}

func TestUserLoggedInEvent_Topic(t *testing.T) {
	e := UserLoggedIn{}
	if e.Topic() != "auth.user.logged_in" {
		t.Errorf("Topic() = %q, want %q", e.Topic(), "auth.user.logged_in")
	}
}

func TestRegister_PublishesUserRegisteredEvent(t *testing.T) {
	bus := &recordingBus{}
	svc := newTestService(t, bus)

	_, err := svc.Register(RegisterPayload{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	if len(bus.published) != 1 {
		t.Fatalf("expected 1 event published, got %d", len(bus.published))
	}

	e, ok := bus.published[0].(UserRegistered)
	if !ok {
		t.Fatalf("expected UserRegistered event, got %T", bus.published[0])
	}
	if e.Username != "alice" {
		t.Errorf("Username = %q, want %q", e.Username, "alice")
	}
	if e.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", e.Email, "alice@example.com")
	}
	if e.UserID == "" {
		t.Error("UserID is empty")
	}
}

func TestLogin_PublishesUserLoggedInEvent(t *testing.T) {
	bus := &recordingBus{}
	svc := newTestService(t, bus)

	_, err := svc.Register(RegisterPayload{
		Username: "bob",
		Email:    "bob@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	bus.published = nil // reset — we only care about the Login event

	_, err = svc.Login("bob", "password123", nil)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if len(bus.published) != 1 {
		t.Fatalf("expected 1 event published, got %d", len(bus.published))
	}

	e, ok := bus.published[0].(UserLoggedIn)
	if !ok {
		t.Fatalf("expected UserLoggedIn event, got %T", bus.published[0])
	}
	if e.Username != "bob" {
		t.Errorf("Username = %q, want %q", e.Username, "bob")
	}
	if e.UserID == "" {
		t.Error("UserID is empty")
	}
}

func TestRegister_NilBus_DoesNotPanic(t *testing.T) {
	svc := newTestService(t, nil)

	_, err := svc.Register(RegisterPayload{
		Username: "charlie",
		Email:    "charlie@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Errorf("Register with nil bus returned unexpected error: %v", err)
	}
}

func TestLogin_NilBus_DoesNotPanic(t *testing.T) {
	svc := newTestService(t, nil)

	if _, err := svc.Register(RegisterPayload{
		Username: "dave",
		Email:    "dave@example.com",
		Password: "password123",
	}); err != nil {
		t.Fatalf("Register: %v", err)
	}

	_, err := svc.Login("dave", "password123", nil)
	if err != nil {
		t.Errorf("Login with nil bus returned unexpected error: %v", err)
	}
}
