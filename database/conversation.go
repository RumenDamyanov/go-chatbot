// Package database provides conversation persistence functionality for the go-chatbot package.
package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Conversation represents a chat conversation.
type Conversation struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Title     string                 `json:"title" db:"title"`
	Metadata  map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// Message represents a single message in a conversation.
type Message struct {
	ID             string                 `json:"id" db:"id"`
	ConversationID string                 `json:"conversation_id" db:"conversation_id"`
	Role           string                 `json:"role" db:"role"` // "user" or "assistant"
	Content        string                 `json:"content" db:"content"`
	Metadata       map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
}

// ConversationStore defines the interface for conversation persistence.
type ConversationStore interface {
	// CreateConversation creates a new conversation.
	CreateConversation(ctx context.Context, conv *Conversation) error

	// GetConversation retrieves a conversation by ID.
	GetConversation(ctx context.Context, id string) (*Conversation, error)

	// UpdateConversation updates an existing conversation.
	UpdateConversation(ctx context.Context, conv *Conversation) error

	// DeleteConversation deletes a conversation and all its messages.
	DeleteConversation(ctx context.Context, id string) error

	// ListConversations lists conversations for a user.
	ListConversations(ctx context.Context, userID string, limit, offset int) ([]*Conversation, error)

	// AddMessage adds a message to a conversation.
	AddMessage(ctx context.Context, msg *Message) error

	// GetMessages retrieves messages for a conversation.
	GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*Message, error)

	// DeleteMessage deletes a specific message.
	DeleteMessage(ctx context.Context, messageID string) error

	// GetConversationHistory retrieves the full conversation history.
	GetConversationHistory(ctx context.Context, conversationID string) ([]*Message, error)

	// SearchConversations searches conversations by content or title.
	SearchConversations(ctx context.Context, userID, query string, limit int) ([]*Conversation, error)
}

// SQLConversationStore implements ConversationStore using SQL database.
type SQLConversationStore struct {
	db     *sql.DB
	driver string // "postgres" or "sqlite3"
}

// NewSQLConversationStore creates a new SQL-based conversation store.
func NewSQLConversationStore(db *sql.DB, driver string) *SQLConversationStore {
	return &SQLConversationStore{
		db:     db,
		driver: driver,
	}
}

// Initialize creates the necessary database tables.
func (s *SQLConversationStore) Initialize(ctx context.Context) error {
	// Create conversations table
	conversationsSQL := `
		CREATE TABLE IF NOT EXISTS conversations (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			title TEXT NOT NULL,
			metadata TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	// Create messages table
	messagesSQL := `
		CREATE TABLE IF NOT EXISTS messages (
			id VARCHAR(255) PRIMARY KEY,
			conversation_id VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			metadata TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
		)`

	// Create indexes
	indexSQL := []string{
		"CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id)",
		"CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at)",
	}

	// Execute table creation
	if _, err := s.db.ExecContext(ctx, conversationsSQL); err != nil {
		return fmt.Errorf("failed to create conversations table: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, messagesSQL); err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}

	// Execute index creation
	for _, idx := range indexSQL {
		if _, err := s.db.ExecContext(ctx, idx); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// CreateConversation creates a new conversation.
func (s *SQLConversationStore) CreateConversation(ctx context.Context, conv *Conversation) error {
	metadataJSON, err := json.Marshal(conv.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	conv.CreatedAt = time.Now()
	conv.UpdatedAt = conv.CreatedAt

	query := `
		INSERT INTO conversations (id, user_id, title, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = s.db.ExecContext(ctx, query, conv.ID, conv.UserID, conv.Title, string(metadataJSON), conv.CreatedAt, conv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

// GetConversation retrieves a conversation by ID.
func (s *SQLConversationStore) GetConversation(ctx context.Context, id string) (*Conversation, error) {
	query := `
		SELECT id, user_id, title, metadata, created_at, updated_at
		FROM conversations WHERE id = $1`

	var conv Conversation
	var metadataJSON string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&conv.ID, &conv.UserID, &conv.Title, &metadataJSON, &conv.CreatedAt, &conv.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Parse metadata
	if metadataJSON != "" {
		if err := json.Unmarshal([]byte(metadataJSON), &conv.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &conv, nil
}

// UpdateConversation updates an existing conversation.
func (s *SQLConversationStore) UpdateConversation(ctx context.Context, conv *Conversation) error {
	metadataJSON, err := json.Marshal(conv.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	conv.UpdatedAt = time.Now()

	query := `
		UPDATE conversations
		SET user_id = $2, title = $3, metadata = $4, updated_at = $5
		WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, conv.ID, conv.UserID, conv.Title, string(metadataJSON), conv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// DeleteConversation deletes a conversation and all its messages.
func (s *SQLConversationStore) DeleteConversation(ctx context.Context, id string) error {
	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			// Ignore the specific error that occurs when transaction is already committed
			expectedErr := "sql: transaction has already been committed or rolled back"
			if err.Error() != expectedErr {
				// Log the rollback error, but don't override the main error
			}
		}
	}()

	// Delete messages first (due to foreign key constraint)
	_, err = tx.ExecContext(ctx, "DELETE FROM messages WHERE conversation_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	// Delete conversation
	result, err := tx.ExecContext(ctx, "DELETE FROM conversations WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return tx.Commit()
}

// ListConversations lists conversations for a user.
func (s *SQLConversationStore) ListConversations(ctx context.Context, userID string, limit, offset int) ([]*Conversation, error) {
	query := `
		SELECT id, user_id, title, metadata, created_at, updated_at
		FROM conversations
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		var conv Conversation
		var metadataJSON string

		err := rows.Scan(&conv.ID, &conv.UserID, &conv.Title, &metadataJSON, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		// Parse metadata
		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &conv.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		conversations = append(conversations, &conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate conversations: %w", err)
	}

	return conversations, nil
}

// AddMessage adds a message to a conversation.
func (s *SQLConversationStore) AddMessage(ctx context.Context, msg *Message) error {
	metadataJSON, err := json.Marshal(msg.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	msg.CreatedAt = time.Now()

	query := `
		INSERT INTO messages (id, conversation_id, role, content, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = s.db.ExecContext(ctx, query, msg.ID, msg.ConversationID, msg.Role, msg.Content, string(metadataJSON), msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}

	// Update conversation's updated_at timestamp
	_, err = s.db.ExecContext(ctx, "UPDATE conversations SET updated_at = $1 WHERE id = $2", msg.CreatedAt, msg.ConversationID)
	if err != nil {
		return fmt.Errorf("failed to update conversation timestamp: %w", err)
	}

	return nil
}

// GetMessages retrieves messages for a conversation.
func (s *SQLConversationStore) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*Message, error) {
	query := `
		SELECT id, conversation_id, role, content, metadata, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		var metadataJSON string

		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &metadataJSON, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse metadata
		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &msg.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate messages: %w", err)
	}

	return messages, nil
}

// DeleteMessage deletes a specific message.
func (s *SQLConversationStore) DeleteMessage(ctx context.Context, messageID string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM messages WHERE id = $1", messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// GetConversationHistory retrieves the full conversation history.
func (s *SQLConversationStore) GetConversationHistory(ctx context.Context, conversationID string) ([]*Message, error) {
	query := `
		SELECT id, conversation_id, role, content, metadata, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC`

	rows, err := s.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		var metadataJSON string

		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &metadataJSON, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse metadata
		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &msg.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate messages: %w", err)
	}

	return messages, nil
}

// SearchConversations searches conversations by content or title.
func (s *SQLConversationStore) SearchConversations(ctx context.Context, userID, query string, limit int) ([]*Conversation, error) {
	// Use database-agnostic case-insensitive search
	var searchQuery string
	if s.driver == "postgres" {
		searchQuery = `
			SELECT DISTINCT c.id, c.user_id, c.title, c.metadata, c.created_at, c.updated_at
			FROM conversations c
			LEFT JOIN messages m ON c.id = m.conversation_id
			WHERE c.user_id = $1 AND (
				c.title ILIKE $2 OR
				m.content ILIKE $2
			)
			ORDER BY c.updated_at DESC
			LIMIT $3`
	} else {
		// SQLite and MySQL compatible syntax
		searchQuery = `
			SELECT DISTINCT c.id, c.user_id, c.title, c.metadata, c.created_at, c.updated_at
			FROM conversations c
			LEFT JOIN messages m ON c.id = m.conversation_id
			WHERE c.user_id = ? AND (
				LOWER(c.title) LIKE LOWER(?) OR
				LOWER(m.content) LIKE LOWER(?)
			)
			ORDER BY c.updated_at DESC
			LIMIT ?`
	}

	searchPattern := "%" + query + "%"

	var rows *sql.Rows
	var err error

	if s.driver == "postgres" {
		rows, err = s.db.QueryContext(ctx, searchQuery, userID, searchPattern, limit)
	} else {
		rows, err = s.db.QueryContext(ctx, searchQuery, userID, searchPattern, searchPattern, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to search conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		var conv Conversation
		var metadataJSON string

		err := rows.Scan(&conv.ID, &conv.UserID, &conv.Title, &metadataJSON, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		// Parse metadata
		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &conv.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		conversations = append(conversations, &conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate conversations: %w", err)
	}

	return conversations, nil
}

// ConversationManager provides high-level conversation management.
type ConversationManager struct {
	store ConversationStore
}

// NewConversationManager creates a new conversation manager.
func NewConversationManager(store ConversationStore) *ConversationManager {
	return &ConversationManager{
		store: store,
	}
}

// CreateConversationWithMessage creates a new conversation with an initial message.
func (cm *ConversationManager) CreateConversationWithMessage(ctx context.Context, userID, title, initialMessage string) (*Conversation, *Message, error) {
	// Generate IDs
	convID := generateID()
	msgID := generateID()

	// Create conversation
	conv := &Conversation{
		ID:       convID,
		UserID:   userID,
		Title:    title,
		Metadata: make(map[string]interface{}),
	}

	if err := cm.store.CreateConversation(ctx, conv); err != nil {
		return nil, nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Add initial message if provided
	var msg *Message
	if initialMessage != "" {
		msg = &Message{
			ID:             msgID,
			ConversationID: convID,
			Role:           "user",
			Content:        initialMessage,
			Metadata:       make(map[string]interface{}),
		}

		if err := cm.store.AddMessage(ctx, msg); err != nil {
			return nil, nil, fmt.Errorf("failed to add initial message: %w", err)
		}
	}

	return conv, msg, nil
}

// AddUserMessage adds a user message to a conversation.
func (cm *ConversationManager) AddUserMessage(ctx context.Context, conversationID, content string) (*Message, error) {
	msg := &Message{
		ID:             generateID(),
		ConversationID: conversationID,
		Role:           "user",
		Content:        content,
		Metadata:       make(map[string]interface{}),
	}

	if err := cm.store.AddMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to add user message: %w", err)
	}

	return msg, nil
}

// AddAssistantMessage adds an assistant message to a conversation.
func (cm *ConversationManager) AddAssistantMessage(ctx context.Context, conversationID, content string) (*Message, error) {
	msg := &Message{
		ID:             generateID(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        content,
		Metadata:       make(map[string]interface{}),
	}

	if err := cm.store.AddMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to add assistant message: %w", err)
	}

	return msg, nil
}

// GetConversationContext retrieves recent messages for context.
func (cm *ConversationManager) GetConversationContext(ctx context.Context, conversationID string, maxMessages int) ([]*Message, error) {
	// Get the most recent messages
	messages, err := cm.store.GetConversationHistory(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Return the last N messages
	if len(messages) <= maxMessages {
		return messages, nil
	}

	return messages[len(messages)-maxMessages:], nil
}

// generateID generates a unique ID for conversations and messages.
func generateID() string {
	// Simple timestamp-based ID generation
	// In production, consider using UUID or more sophisticated ID generation
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
