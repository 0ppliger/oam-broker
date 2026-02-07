package main

import (
	"context"
	
	"github.com/owasp-amass/asset-db/repository"
	"github.com/sirupsen/logrus"
)

type ApiV1 struct {
	ctx context.Context
	store repository.Repository
	bus *EventBus
	logger *logrus.Logger
}

type Serializable interface {
	JSON() []byte
}
