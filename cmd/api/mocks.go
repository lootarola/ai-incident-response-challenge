//go:build mocks

package main

import (
	_ "github.com/lootarola/ai-incident-response-challenge/pkg/database/mock"
	_ "github.com/lootarola/ai-incident-response-challenge/pkg/server/mock"
)
