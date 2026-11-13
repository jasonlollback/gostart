package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasonlollback/gostart/v4/generator"
)

const maxLockSize = 64 << 10

func ReadLock(root string) (generator.LockFile, error) {
	file, err := os.Open(filepath.Join(root, "gostart.lock"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return generator.LockFile{}, errors.New("gostart.lock was not found")
		}
		return generator.LockFile{}, fmt.Errorf("open gostart.lock: %w", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return generator.LockFile{}, fmt.Errorf("read gostart.lock size: %w", err)
	}
	if info.Size() > maxLockSize {
		return generator.LockFile{}, errors.New("gostart.lock is too large")
	}

	decoder := json.NewDecoder(io.LimitReader(file, maxLockSize))
	decoder.DisallowUnknownFields()
	// Projects made before v4 carry "edition" in their config. It is read and
	// dropped, so doctor and update still understand those projects.
	var lock struct {
		SchemaVersion   int    `json:"schema_version"`
		TemplateVersion string `json:"template_version"`
		Config          struct {
			generator.Config
			Edition string `json:"edition"`
		} `json:"config"`
	}
	if err := decoder.Decode(&lock); err != nil {
		return generator.LockFile{}, fmt.Errorf("read gostart.lock: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return generator.LockFile{}, errors.New("gostart.lock must contain one JSON object")
	}
	if lock.SchemaVersion != generator.LockSchemaVersion {
		return generator.LockFile{}, fmt.Errorf("gostart.lock uses unsupported schema version %d", lock.SchemaVersion)
	}
	if strings.TrimSpace(lock.TemplateVersion) == "" {
		return generator.LockFile{}, errors.New("gostart.lock contains no template version")
	}
	return generator.LockFile{
		SchemaVersion:   lock.SchemaVersion,
		TemplateVersion: lock.TemplateVersion,
		Config:          lock.Config.Config,
	}, nil
}
