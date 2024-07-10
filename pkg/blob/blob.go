package blob

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetBlob(ctx context.Context, ref string) ([]byte, error) {
	split := strings.Split(ref, "@")
	token, err := token(ctx, split[0])
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://registry-1.docker.io/v2/%s/blobs/%s", split[0], split[1]), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.docker.image.rootfs.diff.tar.gzip")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func token(ctx context.Context, repo string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", repo), nil)
	if err != nil {
		return "", err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	authJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	payload := struct {
		Token string `json:"token"`
	}{}
	err = json.Unmarshal(authJson, &payload)
	if err != nil {
		return "", err
	}
	return payload.Token, nil
}

func Untar(file string, blob []byte) ([]byte, error) {
	uncompressedStream, err := gzip.NewReader(bytes.NewBuffer(blob))
	if err != nil {
		return nil, err
	}
	defer uncompressedStream.Close()
	tarReader := tar.NewReader(uncompressedStream)
	for {
		h, err := tarReader.Next()
		if err != nil {
			return nil, err
		}
		fmt.Println(h.Name)
		if h.Name == file {
			return io.ReadAll(tarReader)
		}
	}
}
