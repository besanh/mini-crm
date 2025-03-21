package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	GModel interface {
		GetId() string
		SetId(id string)
		SetCreatedAt(t time.Time)
		SetUpdatedAt(t time.Time)
	}

	GBase struct {
		Id        uuid.UUID `json:"id" bun:"id,pk,type:uuid"`
		CreatedAt time.Time `json:"created_at" bun:"created_at,type:timestamp,notnull"`
		UpdatedAt time.Time `json:"updated_at" bun:"updated_at,type:timestamp,notnull"`
	}
)

func (b *GBase) GetId() string {
	return b.Id.String()
}

func (b *GBase) SetId(id uuid.UUID) {
	b.Id = id
}

func (b *GBase) SetCreatedAt(t time.Time) {
	b.CreatedAt = t
}

func (b *GBase) SetUpdatedAt(t time.Time) {
	b.UpdatedAt = t
}

func InitPgBase() *GBase {
	return &GBase{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
