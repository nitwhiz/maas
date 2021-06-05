package minecraft

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
