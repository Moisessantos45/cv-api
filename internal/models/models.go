package models

import (
	"time"

	"gorm.io/gorm"
)

type Auth struct {
	ID               uint64 `json:"id" gorm:"primaryKey"`
	Email            string `json:"email" gorm:"unique;not null"`
	PasswordHash     string `json:"-" gorm:"not null"`
	TwoFactorSecret  string `json:"-" gorm:"type:varchar(26)"`
	TwoFactorEnabled bool   `json:"two_factor_enabled" gorm:"default:false"`
	FirstSession     bool   `json:"first_session" gorm:"default:true"`
	FullProfile      bool   `json:"full_profile" gorm:"default:false"`
	EmailConfirmed   bool   `json:"email_confirmed" gorm:"default:false"`
	Token            string `json:"token" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SocialLink struct {
	Type string `json:"type" gorm:"type:varchar(50)"`
	URL  string `json:"url" gorm:"type:varchar(500)"`
}

type Profile struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	LastName     string         `gorm:"type:varchar(255);not null" json:"last_name"`
	CVLink       string         `gorm:"type:varchar(500)" json:"cv_link"`
	LinkedInLink string         `gorm:"type:varchar(500)" json:"linkedin_link"`
	GitHubLink   string         `gorm:"type:varchar(500)" json:"github_link"`
	SocialLinks  JSONLinkSocial `gorm:"type:jsonb;default:'[]'" json:"social_links"`
	Description  string         `gorm:"type:text" json:"description"`

	AuthID    uint64         `gorm:"unique;not null" json:"auth_id"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Auth Auth `json:"-" gorm:"foreignKey:AuthID"`
}

type Experience struct {
	ID            uint64          `gorm:"primaryKey" json:"id"`
	Title         string          `gorm:"type:varchar(255);not null" json:"title"`
	Company       string          `gorm:"type:varchar(255);not null" json:"company"`
	Location      string          `gorm:"type:varchar(255)" json:"location"`
	StartDate     string          `gorm:"type:date;not null" json:"start_date"`
	EndDate       string          `gorm:"type:date" json:"end_date,omitempty"`
	Description   string          `gorm:"type:text" json:"description"`
	SkillsLearned JSONStringArray `gorm:"type:jsonb;default:'[]'" json:"skills_learned"`
	CreatedAt     time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

type Link struct {
	Frontend string `json:"frontend"`
	Backend  string `json:"backend"`
}

type StatePost struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null;type:varchar(50)"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Project struct {
	ID              uint64          `json:"id" gorm:"primaryKey"`
	Slug            string          `json:"slug" gorm:"unique;not null;type:varchar(150)"`
	Title           string          `json:"title" gorm:"type:varchar(125);not null"`
	TypeProject     string          `json:"type_project" gorm:"type:varchar(24);not null"`
	Description     string          `json:"description" gorm:"type:text;not null"`
	Technologies    JSONStringArray `json:"technologies" gorm:"type:jsonb;default:'[]'"`
	Characteristics JSONStringArray `json:"characteristics" gorm:"type:jsonb;default:'[]'"`
	Learning        JSONStringArray `json:"learning" gorm:"type:jsonb;default:'[]'"`
	Banner          string          `json:"banner" gorm:"type:varchar(255);not null"`
	Images          JSONStringArray `json:"images" gorm:"type:jsonb;default:'[]'"`
	Link            string          `json:"link" gorm:"type:varchar(150)"`
	CreatedAt       string          `json:"created_at" gorm:"type:date;not null"`
	LinkFrontend    JSONStringArray `json:"link_frontend" gorm:"type:jsonb;default:'[]'"`
	LinkBackend     JSONStringArray `json:"link_backend" gorm:"type:jsonb;default:'[]'"`
	StateID         uint64          `json:"state_id" gorm:"not null;default:1"`
	CounterLikes    int64           `json:"counter_likes" gorm:"type:int;default:0"`

	RecordedAt time.Time `json:"recorded_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	State StatePost `json:"-" gorm:"foreignKey:StateID"`
}

type Video struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"type:varchar(125);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	URL         string `json:"url"  gorm:"type:varchar(125)"`
	Status      bool   `json:"status" gorm:"default:true"`
	CreatedAt   string `json:"created_at" gorm:"type:date;not null"`

	RecordedAt time.Time `json:"recorded_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Post struct {
	ID        uint64          `json:"id" gorm:"primaryKey"`
	Slug      string          `json:"slug" gorm:"unique;not null;type:varchar(150)"`
	Title     string          `json:"title" gorm:"not null;type:varchar(200)"`
	Content   string          `json:"content" gorm:"type:text"`
	Banner    string          `json:"banner" gorm:"type:text"`
	AuthorID  uint64          `json:"author_id" gorm:"not null"`
	Tags      JSONStringArray `json:"tags" gorm:"type:jsonb;default:'[]'"`
	Category  string          `json:"category" gorm:"type:varchar(100)"`
	StateID   uint64          `json:"state_id" gorm:"not null"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`

	Profile Profile   `json:"-" gorm:"foreignKey:AuthorID"`
	State   StatePost `json:"-" gorm:"foreignKey:StateID"`
}

var Models = []any{
	&Auth{},
	&Project{},
	&Video{},
	&Profile{},
	&Experience{},
	&Post{},
}

var StatePosts = []StatePost{
	{Name: "draft"},
	{Name: "published"},
	{Name: "archived"},
	{Name: "private"},
	{Name: "unlisted"},
}
