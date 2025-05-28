package domain

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// タグの最大数
const MaxTagsPerLeaf = 10

// ドメイン固有エラー
var (
	ErrTagLimitExceeded = errors.New("タグは" + strconv.Itoa(MaxTagsPerLeaf) + "個までです。")
	ErrAlreadyRead      = errors.New("既読状態です。")
)

// LeafID Value Object
// 不変性を担保し、ID生成・バリデーションに凝集

type LeafID struct {
	value string
}

func NewLeafID(value string) (LeafID, error) {
	if value == "" {
		return LeafID{}, errors.New("LeafIDは空にできません")
	}
	return LeafID{value: value}, nil
}

func NewLeafIDFromUUID() LeafID {
	return LeafID{value: uuid.NewString()}
}

func (id LeafID) String() string {
	return id.value
}

func (id LeafID) Equals(other LeafID) bool {
	return id.value == other.value
}

// URL Value Object

type LeafURL struct {
	value string
}

func NewLeafURL(value string) (LeafURL, error) {
	if value == "" {
		return LeafURL{}, errors.New("URLは空にできません")
	}
	// URLバリデーション
	if len(value) < 3 || len(value) > 2048 {
		return LeafURL{}, errors.New("URLは3文字以上2048文字以下である必要があります")
	}
	// URL形式の検証
	if _, err := url.ParseRequestURI(value); err != nil {
		return LeafURL{}, fmt.Errorf("URLの形式が無効です: %w", err)
	}
	return LeafURL{value: value}, nil
}

func (u LeafURL) String() string {
	return u.value
}

func (u LeafURL) Equals(other LeafURL) bool {
	return u.value == other.value
}

// Tag Value Object

type Tag struct {
	value string
}

func NewTag(value string) (Tag, error) {
	if value == "" {
		return Tag{}, errors.New("Tagは空にできません")
	}
	return Tag{value: value}, nil
}

func (t Tag) String() string {
	return t.value
}

func (t Tag) Equals(other Tag) bool {
	return t.value == other.value
}

// Leaf集約ルート
// コメントで集約ルートであることを明示
// タグ重複禁止・長さ制限も追加

type Leaf struct {
	id       LeafID
	note     string
	url      LeafURL
	platform string
	tags     []Tag
	read     bool
	syncedAt time.Time
}

// Getter
func (l *Leaf) ID() LeafID          { return l.id }
func (l *Leaf) Note() string        { return l.note }
func (l *Leaf) URL() LeafURL        { return l.url }
func (l *Leaf) Platform() string    { return l.platform }
func (l *Leaf) Tags() []Tag         { return l.tags }
func (l *Leaf) Read() bool          { return l.read }
func (l *Leaf) SyncedAt() time.Time { return l.syncedAt }

// ファクトリ
// ID生成
// バリデーション一括
func NewLeaf(id string, note string, url string, platform string, tagValues []string) (*Leaf, error) {
	if note == "" {
		return nil, errors.New("Noteは空にできません")
	}
	if platform == "" {
		return nil, errors.New("Platformは空にできません")
	}
	// IDが空の場合は新規生成
	if id == "" {
		id = NewLeafIDFromUUID().String()
	}
	leafID, err := NewLeafID(id)
	if err != nil {
		return nil, err
	}
	leafURL, err := NewLeafURL(url)
	if err != nil {
		return nil, err
	}
	tags := make([]Tag, 0, len(tagValues))
	tagSet := make(map[string]struct{})
	for _, v := range tagValues {
		t, err := NewTag(v)
		if err != nil {
			return nil, err
		}
		if _, exists := tagSet[t.value]; exists {
			return nil, errors.New("タグが重複しています")
		}
		tagSet[t.value] = struct{}{}
		tags = append(tags, t)
	}
	if len(tags) > MaxTagsPerLeaf {
		return nil, ErrTagLimitExceeded
	}
	return &Leaf{
		id:       leafID,
		note:     note,
		url:      leafURL,
		platform: platform,
		tags:     tags,
		read:     false,
		syncedAt: time.Now().UTC(),
	}, nil
}

// 既存のLeafを再構築するためのファクトリ
func ReconstructLeaf(id string, note string, url string, platform string, tagValues []string, read bool, syncedAt time.Time) (*Leaf, error) {
	leafID, err := NewLeafID(id)
	if err != nil {
		return nil, err
	}
	if note == "" {
		return nil, errors.New("Noteは空にできません")
	}
	if platform == "" {
		return nil, errors.New("Platformは空にできません")
	}
	leafURL, err := NewLeafURL(url)
	if err != nil {
		return nil, err
	}
	tags := make([]Tag, 0, len(tagValues))
	tagSet := make(map[string]struct{})
	for _, v := range tagValues {
		t, err := NewTag(v)
		if err != nil {
			return nil, err
		}
		if _, exists := tagSet[t.value]; exists {
			return nil, errors.New("タグが重複しています")
		}
		tagSet[t.value] = struct{}{}
		tags = append(tags, t)
	}
	if len(tags) > MaxTagsPerLeaf {
		return nil, ErrTagLimitExceeded
	}
	return &Leaf{
		id:       leafID,
		note:     note,
		url:      leafURL,
		platform: platform,
		tags:     tags,
		read:     read,
		syncedAt: syncedAt,
	}, nil
}

// ノート内容の変更
func (l *Leaf) UpdateNote(note string) error {
	if note == "" {
		return errors.New("Noteは空にできません")
	}
	l.note = note
	return nil
}

// プラットフォームの変更
func (l *Leaf) UpdatePlatform(platform string) error {
	if platform == "" {
		return errors.New("Platformは空にできません")
	}
	l.platform = platform
	return nil
}

// 既読マーク
func (l *Leaf) MarkAsRead() error {
	if l.read {
		return ErrAlreadyRead
	}
	l.read = true
	return nil
}

// タグのバリデーション付き更新（重複・上限チェック）
func (l *Leaf) UpdateTags(tags []Tag) error {
	if len(tags) > MaxTagsPerLeaf {
		return ErrTagLimitExceeded
	}
	tagSet := make(map[string]struct{})
	for _, t := range tags {
		if _, exists := tagSet[t.value]; exists {
			return errors.New("タグが重複しています")
		}
		tagSet[t.value] = struct{}{}
	}
	l.tags = make([]Tag, len(tags))
	copy(l.tags, tags)
	return nil
}
