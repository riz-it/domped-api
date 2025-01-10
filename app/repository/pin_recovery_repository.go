package repository

import (
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/domain"
)

type PinRecoveryRepository struct {
	Repository[domain.PinRecoveryEntity]
	Log *logrus.Logger
}

func NewPinRecovery(log *logrus.Logger) *PinRecoveryRepository {
	return &PinRecoveryRepository{
		Log: log,
	}
}
