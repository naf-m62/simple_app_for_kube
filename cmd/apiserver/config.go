package apiserver

import "simple_app_for_kube/cmd/database"

type Config struct {
	Host string           `yaml:"host"`
	Port string           `yaml:"port"`
	Db   *database.Config `yaml:"db"`
}
