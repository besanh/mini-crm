package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/besanh/mini-crm/common/constant"
	"github.com/besanh/mini-crm/models"
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
		field.JSON("user_profile", models.UserProfile{}),
		field.Enum("status").Values(constant.USER_STATUS_ACTIVE, constant.USER_STATUS_INACTIVE, constant.USER_STATUS_DELETED),
		field.Strings("scope").
			SchemaType(map[string]string{
				dialect.Postgres: "text[]",
			}).
			Optional(),
	}
}

// Edges of the Users.
func (Users) Edges() []ent.Edge {
	return nil
}

func (Users) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "status"),
	}
}
