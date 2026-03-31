package main

import (
	"github.com/zerx-lab/wordZero/apps/wordzero/conf"
	"github.com/zerx-lab/wordZero/pkg/s3storage"
)

func initS3Uploader() error {
	cfg, err := conf.LoadConfig(configFile)
	if err != nil {
		return err
	}

	return s3storage.InitGlobalUploader(cfg.GetS3Config())
}
