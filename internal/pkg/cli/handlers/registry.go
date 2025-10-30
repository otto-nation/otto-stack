package handlers

import (
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

type HandlerFactory func(string) cliTypes.CommandHandler

var registry = make(map[string]HandlerFactory)

func Register(packageName string, factory HandlerFactory) {
	registry[packageName] = factory
}

func Get(packageName, commandName string) cliTypes.CommandHandler {
	if factory, exists := registry[packageName]; exists {
		return factory(commandName)
	}
	return nil
}
