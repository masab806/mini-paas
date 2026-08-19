package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Deployments holds the schema definition for the Deployments entity.
type Deployments struct {
	ent.Schema
}

// Fields of the Deployments.
func (Deployments) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("branch").
			Default("main"),

		field.String("image_tag").
			NotEmpty(),

		field.String("repo_url").
			NotEmpty(),

		field.String("status").
			Default("QUEUED"),

		field.Time("created_at").
			Default(time.Now),

		field.Time("finished_at").
			Optional().
			Nillable(),
	}
}

// Edges of the Deployments.
func (Deployments) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).
			Ref("Deployments").
			Unique().
			Required(),
	}
}
