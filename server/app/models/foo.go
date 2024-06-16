package models

import (
	"context"

	"github.com/abibby/remote-input/server/app/providers"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/database/model/modeldi"
)

//go:generate spice generate:migration
type Foo struct {
	model.BaseModel

	ID int `json:"id" db:"id,primary,autoincrement"`
}

func init() {
	providers.Add(modeldi.Register[*Foo])
}

func FooQuery(ctx context.Context) *builder.ModelBuilder[*Foo] {
	return builder.From[*Foo]().WithContext(ctx)
}