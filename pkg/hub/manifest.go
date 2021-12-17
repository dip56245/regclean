package hub

import (
	"context"
	"encoding/json"
	"net/http"
)

type Manifest struct {
	Digest       string
	ConfigDigest string
	Size         uint64
}

type hubResponse struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

func (hr *hubResponse) GetSizeInBytes() uint64 {
	var out uint64
	for i := 0; i < len(hr.Layers); i++ {
		out += uint64(hr.Layers[i].Size)
	}
	return out
}

func (h *Hub) Manifest(repo string, tag string) (*Manifest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.RequestTimeOut)
	defer cancel()
	resp, err := h.methodGetEx(ctx,
		func(req *http.Request) {
			req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
		},
		"/v2/%s/manifests/%s", repo, tag)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var item hubResponse
	if err = decoder.Decode(&item); err != nil {
		return nil, err
	}
	return &Manifest{
		ConfigDigest: item.Config.Digest,
		Digest:       resp.Header.Get("Docker-Content-Digest"),
		Size:         item.GetSizeInBytes(),
	}, nil
}

func (h *Hub) DeleteManifest(repo string, digest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.RequestTimeOut)
	defer cancel()
	resp, err := h.methodDelete(ctx, "/v2/%s/manifests/%s", repo, digest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
