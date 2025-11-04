package handlers

import (
	types "github.com/otto-nation/otto-stack/internal/pkg/types"
)

type HandlerFactory func(string) types.CommandHandler

var registry = make(map[string]HandlerFactory)

func Register(packageName string, factory HandlerFactory) {
	registry[packageName] = factory
}

func Get(packageName, commandName string) types.CommandHandler {
	if factory, exists := registry[packageName]; exists {
		return factory(commandName)
	}
	return nil
}
