package main

import (
	"github.com/zerx-lab/wordZero/apps/wordzero/conf"
	"github.com/zerx-lab/wordZero/pkg/s3"
)

func initS3Uploader() error {
	cfg, err := conf.LoadConfig(configFile)
	if err != nil {
		return err
	}

	if err := s3.InitGlobalUploader(cfg.GetS3Config()); err != nil {
		return err
	}
	return nil
}
