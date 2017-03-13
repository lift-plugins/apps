package apps

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

const (
	mdPath = ".lift.json"
)

// DeployMetadata defines the information that is stored in the current working directory related to the created app and cache used to synchronize
// its source code.
type DeployMetadata struct {
	AppID   string `json:"app_id"`
	CacheID string `json:"cache_id"`
}

// Read loads deploy metadata.
func (m *DeployMetadata) Read() error {
	data, err := ioutil.ReadFile(mdPath)
	if err != nil {
		return errors.Wrapf(err, "failed deploy metadata tokens file at %q", mdPath)
	}

	if err := json.Unmarshal(data, &m); err != nil {
		return errors.Wrapf(err, "failed unmarshaling tokens file at %q", mdPath)
	}
	return nil
}

// Writes deploy metadata.
func (m *DeployMetadata) Write() error {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return errors.Wrap(err, "failed marshaling deploy metadata")
	}

	if err := ioutil.WriteFile(mdPath, data, os.FileMode(0600)); err != nil {
		return errors.Wrapf(err, "failed writing deploy metadata to %q", mdPath)
	}
	return nil
}
