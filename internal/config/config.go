package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	OrchestratorPort      string `yaml:"orchestrator_port" env-default:"4040"`
	OrchestratorHost      string `yaml:"orchestrator_host" env-default:"localhost"`
	GrpcPort              string `yaml:"grpc_port" env-default:"4041"`
	GrpcHost              string `yaml:"grpc_host" env-default:"localhost"`
	ComputingPower        int    `yaml:"computing_power" env:"COMPUTING_POWER" env-default:"1"`
	TimeAdditionMs        int    `yaml:"time_addition_ms" env:"TIME_ADDITION_MS" env-default:"200"`
	TimeSubtractionMs     int    `yaml:"time_subtraction_ms" env:"TIME_SUBTRACTION_MS" env-default:"200"`
	TimeMultiplicationsMs int    `yaml:"time_multiplications_ms" env:"TIME_MULTIPLICATIONS_MS" env-default:"300"`
	TimeDivisionsMs       int    `yaml:"time_divisions_ms" env:"TIME_DIVISIONS_MS" env-default:"400"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load("local.env")

	var cfg Config
	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
