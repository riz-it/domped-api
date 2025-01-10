package repository

import (
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/domain"
)

type TopUpRepository struct {
	Repository[domain.TopUpEntity]
	Log *logrus.Logger
}

func NewTopUp(log *logrus.Logger) *TopUpRepository {
	return &TopUpRepository{
		Log: log,
	}
}
