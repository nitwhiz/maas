package minecraft

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

const ManifestUrl = "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json"

type LatestVersion struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type Version struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Url             string `json:"url"`
	Time            string `json:"time"`
	ReleaseTime     string `json:"releaseTime"`
	SHA1            string `json:"sha1"`
	ComplianceLevel int    `json:"complianceLevel"`
}

type Manifest struct {
	Latest   LatestVersion `json:"latest"`
	Versions []Version     `json:"versions"`
}

func unmarshalManifest(bs []byte) (*Manifest, error) {
	var manifest Manifest

	err := json.Unmarshal(bs, &manifest)

	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

// DownloadManifest downloads and parses the manifest file
func DownloadManifest() (*Manifest, error) {
	r, err := http.Get(ManifestUrl)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r.Body)

	if err != nil {
		return nil, err
	}

	return unmarshalManifest(buf.Bytes())
}

// ReadManifestFromFile parses a manifest JSON file
func ReadManifestFromFile(path string) (*Manifest, error) {
	bs, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return unmarshalManifest(bs)
}

// FilterVersionsByType filters out all versions not containing the string `t`
func (m *Manifest) FilterVersionsByType(t string) *Manifest {
	if t == "all" {
		return m
	}

	var result []Version

	for _, v := range m.Versions {
		if v.Type == t {
			result = append(result, v)
		}
	}

	m.Versions = result

	return m
}

// FilterVersionsByIdSubstring filters out all versions not containing the string `sid`
func (m *Manifest) FilterVersionsByIdSubstring(id string) *Manifest {
	if id == "" {
		return m
	}

	var result []Version

	for _, v := range m.Versions {
		if strings.Contains(v.Id, id) {
			result = append(result, v)
		}
	}

	m.Versions = result

	return m
}

// OrderVersionsByFieldAndDirection orders the versions by `field` and `direction`
func (m *Manifest) OrderVersionsByFieldAndDirection(field string, direction string) *Manifest {
	asc := direction == "asc"

	if field == "releaseTime" {
		sort.SliceStable(m.Versions, func(i int, j int) bool {
			if asc {
				return m.Versions[i].ReleaseTime < m.Versions[j].ReleaseTime
			} else {
				return m.Versions[i].ReleaseTime > m.Versions[j].ReleaseTime
			}
		})
	} else {
		sort.SliceStable(m.Versions, func(i int, j int) bool {
			if asc {
				return m.Versions[i].Time < m.Versions[j].Time
			} else {
				return m.Versions[i].Time > m.Versions[j].Time
			}
		})
	}

	return m
}

// SliceVersions keeps maximum `size` versions in the manifest
func (m *Manifest) SliceVersions(size int) *Manifest {
	if len(m.Versions) >= size {
		m.Versions = m.Versions[:size]
	}

	return m
}
