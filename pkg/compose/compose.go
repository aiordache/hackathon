package compose

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
	"os/exec"
	"strings"
	"time"

	"github.com/aiordache/hackathon/pkg/blob"
)

var (
	ErrLayerNotFound       = errors.New("layer not found")
	ErrBlobNotFound        = errors.New("blob not found")
	ErrComposeFileNotFound = errors.New("compose file not found")
)

func EmbeddedCompose(ctx context.Context, image string) ([]byte, error) {
	tag := "latest"
	split := strings.Split(image, ":")
	if len(split) == 2 {
		tag = split[1]
	}
	layer := ComposeLayer(ctx, split[0], tag)
	if layer == nil {
		return nil, fmt.Errorf("%w: %s", ErrLayerNotFound, image)
	}
	b, err := blob.GetBlob(context.Background(), layer.Parent)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBlobNotFound, layer.Parent)
	}
	composeBytes, err := blob.Untar("docker-compose.yaml", b)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrComposeFileNotFound, err.Error())
	}
	return composeBytes, nil
}

func ComposeLayer(ctx context.Context, image, tag string) *AnnotationData {
	v, ok := GetAnnotations(ctx, image, tag, "")["composehackv1"]
	if !ok {
		return nil
	}
	return &v
}

func GetAnnotations(ctx context.Context, image, tag, sha string) map[string]AnnotationData {
	if ctx.Err() != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	out := map[string]AnnotationData{}
	m := Inspect(ctx, getId(image, tag, sha))
	parent := getId(image, tag, sha)
	if len(m.Layers) > 0 {
		parent = getId(image, tag, m.Layers[0].Digest)
	}
	appendAnnotations(out, parent, m.Annotations)
	if m.MediaType != "application/vnd.oci.image.index.v1+json" {
		return out
	}
	for _, manifest := range m.Manifests {
		if manifest.Digest != "" {
			maps.Copy(out, GetAnnotations(ctx, image, tag, manifest.Digest))
		}
	}
	return out
}

func getId(image, tag, sha string) string {
	if sha != "" {
		return image + "@" + sha
	}
	if tag != "" {
		return image + ":" + tag
	}
	return image
}

func appendAnnotations(dst map[string]AnnotationData, parent string, annotations map[string]string) {
	for k, v := range annotations {
		dst[k] = AnnotationData{
			Value:  v,
			Parent: parent,
		}
	}
}

type AnnotationData struct {
	Value  any
	Parent string
}

type OCIManifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Manifests     []struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
		Size      int    `json:"size"`
		Platform  struct {
			Architecture string `json:"architecture"`
			Os           string `json:"os"`
		} `json:"platform"`
		Annotations map[string]string `json:"annotations,omitempty"`
	} `json:"manifests"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
	} `json:"layers,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

func Inspect(ctx context.Context, image string) *OCIManifest {
	rawJson, err := exec.CommandContext(ctx, "docker", "buildx", "imagetools", "inspect", image, "--raw").Output()
	if err != nil {
		log.Fatalf("Error inspecting image: %v", err)
		return nil
	}
	var m OCIManifest
	err = json.Unmarshal(rawJson, &m)
	if err != nil {
		log.Fatalf("Error unmarshalling image: %v", err)
		return nil
	}
	return &m
}
