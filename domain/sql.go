package domain

import sq "github.com/Masterminds/squirrel"

var Psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
