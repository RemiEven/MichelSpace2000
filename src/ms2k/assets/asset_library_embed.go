//go:build !wasm

package assets

import (
	"context"
	"fmt"
)

func (al *Library) getFileData(_ context.Context, path string) ([]byte, error) {
	fileData, err := assetFS.ReadFile("files/" + path)
	if err != nil {
		return nil, fmt.Errorf("failed to read path [%s]: %w", path, err)
	}
	return fileData, nil
}
