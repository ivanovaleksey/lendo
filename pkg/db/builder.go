package db

import "github.com/Masterminds/squirrel"

var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
