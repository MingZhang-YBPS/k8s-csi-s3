package mounter

import (
	"fmt"
	"os"

	"github.com/ctrox/csi-s3/pkg/s3"
)

// Implements Mounter
type rcloneMounter struct {
	bucket          *s3.Bucket
	url             string
	region          string
	accessKeyID     string
	secretAccessKey string
}

const (
	rcloneCmd = "rclone"
)

func newRcloneMounter(b *s3.Bucket, cfg *s3.Config) (Mounter, error) {
	return &rcloneMounter{
		bucket:          b,
		url:             cfg.Endpoint,
		region:          cfg.Region,
		accessKeyID:     cfg.AccessKeyID,
		secretAccessKey: cfg.SecretAccessKey,
	}, nil
}

func (rclone *rcloneMounter) Stage(stageTarget string) error {
	return nil
}

func (rclone *rcloneMounter) Unstage(stageTarget string) error {
	return nil
}

func (rclone *rcloneMounter) Mount(source string, target string) error {
	args := []string{
		"mount",
		fmt.Sprintf(":s3:%s/%s", rclone.bucket.Name, rclone.bucket.FSPath),
		fmt.Sprintf("%s", target),
		"--daemon",
		"--s3-provider=AWS",
		"--s3-env-auth=true",
		fmt.Sprintf("--s3-region=%s", rclone.region),
		fmt.Sprintf("--s3-endpoint=%s", rclone.url),
		"--allow-other",
		// TODO: make this configurable
		"--vfs-cache-mode=writes",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", rclone.accessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", rclone.secretAccessKey)
	return fuseMount(target, rcloneCmd, args)
}