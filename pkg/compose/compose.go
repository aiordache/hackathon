package compose

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"
	"time"
)

func IsDockerComposeImage(image, tag string) bool {
	_, ok := GetAnnotations(image, tag, "")["composehackv1"]
	return ok
}

func GetAnnotations(image, tag, sha string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out := map[string]string{}
	m := Inspect(ctx, getId(image, tag, sha))
	mapAppend(out, m.Annotations)
	if m.MediaType != "application/vnd.oci.image.index.v1+json" {
		return out
	}
	for _, manifest := range m.Manifests {
		if manifest.Digest != "" {
			mapAppend(out, GetAnnotations(image, tag, manifest.Digest))
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

func mapAppend(m1, m2 map[string]string) {
	for k, v := range m2 {
		m1[k] = v
	}
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
