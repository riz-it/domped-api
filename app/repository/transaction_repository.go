package repository

import (
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/domain"
)

type TransactionRepository struct {
	Repository[domain.TransactionEntity]
	Log *logrus.Logger
}

func NewTransaction(log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		Log: log,
	}
}
