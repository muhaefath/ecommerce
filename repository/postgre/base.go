package postgre

import (
	"ecommerce/repository"
	sdkSql "ecommerce/utils/sql"
)

// baseRepo is base repo to store common func for repo.
type baseRepo struct {
	db sdkSql.DBer
}

func (b *baseRepo) DB() repository.QueryProvider {
	return b.db
}
