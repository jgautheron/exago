package repository

import "github.com/exago/svc/repository/model"

type Worker interface {
	GetImports() (model.Imports, error)
}
