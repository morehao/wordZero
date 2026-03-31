package main

import (
	"github.com/ygpkg/yg-go/logs"
	"github.com/zerx-lab/wordZero/apps/wordzero/internal/config"
	"github.com/zerx-lab/wordZero/pkg/s3"
)

func initS3Uploader(cfg *config.WordZeroConfig) error {
	s3Cfg := cfg.GetS3Config()
	if s3Cfg == nil || s3Cfg.Bucket == "" {
		logs.Warnf("[main] S3 config not found, skip uploader init")
		return nil
	}

	uploader, err := s3.NewUploader(s3Cfg)
	if err != nil {
		logs.Errorf("[main] init S3 uploader failed: %s", err)
		return err
	}
	s3.SetGlobalUploader(uploader)
	logs.Infof("[main] S3 uploader initialized, bucket: %s", s3Cfg.Bucket)
	return nil
}
