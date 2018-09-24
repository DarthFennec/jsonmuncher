package benchmark

import (
	"time"
)

/*
   Small payload, http log like structure. Size: 190 bytes
*/

type SmallPayload struct {
	St   int
	Sid  int
	Tt   string
	Gr   int
	Uuid string
	Ip   string
	Ua   string
	Tz   int
	V    int
}

/*
   Medium payload, based on Clearbit API response. Size: 2.3kb
*/

type CBAvatar struct {
	Url string
}

type CBGravatar struct {
	Avatars []*CBAvatar
}

type CBGithub struct {
	Followers int
}

type CBName struct {
	FullName string
}

type CBPerson struct {
	Name     *CBName
	Github   *CBGithub
	Gravatar *CBGravatar
}

type MediumPayload struct {
	Person  *CBPerson
	Company map[string]interface{}
}

/*
   Large payload, based on Discourse API. Size: 41kb
*/

type DSUser struct {
	Username string
}

type DSTopic struct {
	Id   int
	Slug string
}

type DSTopicsList struct {
	Topics        []*DSTopic
	MoreTopicsUrl string
}

type LargePayload struct {
	Users  []*DSUser
	Topics *DSTopicsList
}

/*
   Huge payload, based on a large Helm index. Size: 333mb
*/

type IndexFile struct {
	APIVersion string                   `json:"apiVersion"`
	Generated  time.Time                `json:"generated"`
	Entries    map[string]ChartVersions `json:"entries"`
}

type ChartVersions []ChartVersion

type ChartVersion struct {
	Metadata
	URLs     []string  `json:"urls"`
	Created  time.Time `json:"created,omitempty"`
	Removed  bool      `json:"removed,omitempty"`
	Digest   string    `json:"digest,omitempty"`
}

type Metadata struct {
	Name          string            `json:"name,omitempty"`
	Home          string            `json:"home,omitempty"`
	Sources       []string          `json:"sources,omitempty"`
	Version       string            `json:"version,omitempty"`
	Description   string            `json:"description,omitempty"`
	Keywords      []string          `json:"keywords,omitempty"`
	Maintainers   []*Maintainer     `json:"maintainers,omitempty"`
	Engine        string            `json:"engine,omitempty"`
	Icon          string            `json:"icon,omitempty"`
	ApiVersion    string            `json:"apiVersion,omitempty"`
	Condition     string            `json:"condition,omitempty"`
	Tags          string            `json:"tags,omitempty"`
	AppVersion    string            `json:"appVersion,omitempty"`
	Deprecated    bool              `json:"deprecated,omitempty"`
	TillerVersion string            `json:"tillerVersion,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	KubeVersion   string            `json:"kubeVersion,omitempty"`
}

type Maintainer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Url   string `json:"url,omitempty"`
}
