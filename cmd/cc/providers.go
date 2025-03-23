package main

import (
	"github.com/fr0g-66723067/cc/internal/ai/claude"
	"github.com/fr0g-66723067/cc/internal/container/docker"
	"github.com/fr0g-66723067/cc/internal/vcs/git"
)

// This file ensures that all providers are registered by importing them.
// The init() function in each provider package registers the provider with
// its respective registry.

// Ensure all providers are registered
var (
	_ = claude.NewProvider
	_ = docker.NewProvider
	_ = git.NewProvider
)