package handlers

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
)

type HandlerFactory func(string) base.CommandHandler

var registry = make(map[string]HandlerFactory)

func Register(packageName string, factory HandlerFactory) {
	registry[packageName] = factory
}

func Get(packageName, commandName string) base.CommandHandler {
	if factory, exists := registry[packageName]; exists {
		return factory(commandName)
	}
	return nil
}
