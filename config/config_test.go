package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3ParseEnv(t *testing.T) {
	s3Config := S3Config{}
	err := s3Config.ParseEnv()
	assert.NoError(t, err)
	if s3Config.AwsAccessKeyID == "" {
		t.Skip("no aws creds, skipping tests")
	}
	assert.Equal(t, s3Config.Prefix, "/")
	assert.NotEqual(t, s3Config.AwsAccessKeyID, "")
}
