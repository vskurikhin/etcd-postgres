/*
 * This file was last modified at 2024-07-17 11:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * config.go
 * $Id$
 */

package memory

import "time"

// Config defines the config for storage.
type Config struct {
	// Time before deleting expired keys
	//
	// Default is 10 * time.Second
	GCInterval time.Duration
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	GCInterval: 10 * time.Second,
}

// configDefault is a helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if int(cfg.GCInterval.Seconds()) < int(time.Second) {
		cfg.GCInterval = ConfigDefault.GCInterval
	}
	return cfg
}
