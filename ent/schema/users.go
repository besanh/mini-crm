package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/besanh/mini-crm/common/constant"
	"github.com/google/uuid"
)

// Users holds the schema definition for the Users entity.
type Users struct {
	ent.Schema
}

// Fields of the Users.
func (Users) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()),
		field.Enum("status").Values(constant.USER_STATUS_ACTIVE, constant.USER_STATUS_INACTIVE, constant.USER_STATUS_DELETED),
		field.Strings("roles"),
	}
}

// Edges of the Users.
func (Users) Edges() []ent.Edge {
	return nil
}
