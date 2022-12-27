//go:build wasm

package assets

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"syscall/js"
)

func (al *Library) getFileData(ctx context.Context, path string) ([]byte, error) {
	href := js.Global().Get("location").Get("href").String()
	href = strings.TrimSuffix(href, "/index.html") // TODO: avoid doing this for each file

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, href+"/assets/"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for path [%s]: %w", path, err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file for path [%s]: %w", path, err)
	}
	defer response.Body.Close()

	fileData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read fetched data for path [%s]: %w", path, err)
	}

	return fileData, nil
}
