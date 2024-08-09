/*
 * This file was last modified at 2024-08-06 18:25 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * yaml_load.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"os"

	"github.com/spf13/viper"
)

// LoadConfig
// Example yaml config file:
//
// etcdclient:
//
//	enabled: true
//	cache:
//	  enabled: true
//	  expire_ms: 1000
//	  gc_interval_sec: 10
//	grpc:
//	  address: localhost
//	  enabled: true
//	  port: 8442
//	  proto: tcp
//	  tls:
//	    enabled: true
//	    ca_file: cert/grpc-test_ca-cert.pem
//	    cert_file: cert/grpc-test_server-cert.pem
//	    key_file: cert/grpc-test_server-key.pem
//	http:
//	  address: localhost
//	  enabled: true
//	  port: 8443
//	  tls:
//	    enabled: true
//	    ca_file: cert/http-test_ca-cert.pem
//	    cert_file: cert/http-test_server-cert.pem
//	    key_file: cert/http-test_server-key.pem
func LoadConfig(path string) (cfg YamlConfig, err error) {

	if os.Getenv("GO_FAVORITES_SKIP_LOAD_CONFIG") != "" {
		return &yamlConfig{}, err
	}
	viper.SetConfigName("etcd-client")       // мя файла yamlConfig
	viper.SetConfigType("yaml")              // REQUIRED если файл yamlConfig не имеет расширения в имени
	viper.AddConfigPath("/etc/etcd-client/") // путь для поиска файла yamlConfig
	viper.AddConfigPath(path)                // несколько раз, чтобы добавить несколько путей поиска

	err = viper.ReadInConfig() // Find and read the yamlConfig file
	if err != nil {            // Handle errors reading the yamlConfig file
		return nil, err
	}
	var c yamlConfig
	err = viper.Unmarshal(&c)
	cfg = &c

	return cfg, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
