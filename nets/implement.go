package nets

import (
	"context"

	"github.com/EreminDm/golang_basic_crud/entity"
)

// Provider describes provider methods.
type Provider interface {
	One(ctx context.Context, value string) (entity.PersonalData, error)
	All(ctx context.Context) ([]entity.PersonalData, error)
	Remove(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, p entity.PersonalData) (int64, error)
	Insert(ctx context.Context, document entity.PersonalData) (entity.PersonalData, error)
}
