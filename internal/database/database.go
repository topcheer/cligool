package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Session 终端会话数据库模型
type Session struct {
	ID        string    `db:"id" json:"id"`
	Owner     string    `db:"owner" json:"owner"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Active    bool      `db:"active" json:"active"`
	Metadata  JSONMap   `db:"metadata" json:"metadata,omitempty"`
}

// User 用户模型
type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// JSONMap 自定义JSON类型
type JSONMap map[string]interface{}

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONMap value: %v", value)
	}

	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (interface{}, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// DB 数据库连接
type DB struct {
	PG    *sqlx.DB
	Redis *redis.Client
	ctx   context.Context
}

// NewPostgresDB 创建新的数据库连接
func NewPostgresDB(databaseURL, redisURL string) (*DB, error) {
	// 连接PostgreSQL
	pgDB, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// 配置连接池
	pgDB.SetMaxOpenConns(25)
	pgDB.SetMaxIdleConns(25)
	pgDB.SetConnMaxLifetime(5 * time.Minute)

	// 解析Redis URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	redisClient := redis.NewClient(opt)

	// 测试连接
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &DB{
		PG:    pgDB,
		Redis: redisClient,
		ctx:   ctx,
	}, nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	var errs []error

	if err := db.PG.Close(); err != nil {
		errs = append(errs, fmt.Errorf("PostgreSQL close error: %w", err))
	}

	if err := db.Redis.Close(); err != nil {
		errs = append(errs, fmt.Errorf("Redis close error: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("database close errors: %v", errs)
	}

	return nil
}

// Migrate 运行数据库迁移
func (db *DB) Migrate() error {
	schema := `
	-- 用户表
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- 会话表
	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(36) PRIMARY KEY,
		owner VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
		active BOOLEAN DEFAULT false,
		metadata JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- 会话参与者表
	CREATE TABLE IF NOT EXISTS session_participants (
		id VARCHAR(36) PRIMARY KEY,
		session_id VARCHAR(36) REFERENCES sessions(id) ON DELETE CASCADE,
		user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
		role VARCHAR(50) DEFAULT 'viewer', -- 'owner', 'editor', 'viewer'
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(session_id, user_id)
	);

	-- 创建索引
	CREATE INDEX IF NOT EXISTS idx_sessions_owner ON sessions(owner);
	CREATE INDEX IF NOT EXISTS idx_sessions_active ON sessions(active);
	CREATE INDEX IF NOT EXISTS idx_session_participants_session ON session_participants(session_id);
	CREATE INDEX IF NOT EXISTS idx_session_participants_user ON session_participants(user_id);
	`

	_, err := db.PG.Exec(schema)
	return err
}

// CreateSession 创建新会话
func (db *DB) CreateSession(session *Session) error {
	query := `
		INSERT INTO sessions (id, owner, active, metadata)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.PG.Exec(query, session.ID, session.Owner, session.Active, session.Metadata)
	return err
}

// GetSession 获取会话
func (db *DB) GetSession(id string) (*Session, error) {
	var session Session
	query := `SELECT id, owner, active, metadata, created_at, updated_at FROM sessions WHERE id = $1`
	err := db.PG.Get(&session, query, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// ListSessions 列出所有会话
func (db *DB) ListSessions() ([]Session, error) {
	var sessions []Session
	query := `SELECT id, owner, active, metadata, created_at, updated_at FROM sessions ORDER BY created_at DESC`
	err := db.PG.Select(&sessions, query)
	return sessions, err
}

// DeleteSession 删除会话
func (db *DB) DeleteSession(id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := db.PG.Exec(query, id)
	return err
}

// UpdateSession 更新会话状态
func (db *DB) UpdateSession(session *Session) error {
	query := `
		UPDATE sessions
		SET active = $2, metadata = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := db.PG.Exec(query, session.ID, session.Active, session.Metadata)
	return err
}

// CreateUser 创建用户
func (db *DB) CreateUser(user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.PG.Exec(query, user.ID, user.Email, user.PasswordHash, user.Name)
	return err
}

// GetUser 获取用户
func (db *DB) GetUser(id string) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`
	err := db.PG.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (db *DB) GetUserByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`
	err := db.PG.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CacheSet 设置缓存
func (db *DB) CacheSet(key string, value interface{}, expiration time.Duration) error {
	return db.Redis.Set(db.ctx, key, value, expiration).Err()
}

// CacheGet 获取缓存
func (db *DB) CacheGet(key string) (string, error) {
	return db.Redis.Get(db.ctx, key).Result()
}

// CacheDelete 删除缓存
func (db *DB) CacheDelete(key string) error {
	return db.Redis.Del(db.ctx, key).Err()
}