package benchmark

import (
	"time"
)

// Small payload, http log like structure. Size: 190 bytes

// SmallPayload contains the small payload data.
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

// Medium payload, based on Clearbit API response. Size: 2.3kb

// CBAvatar represents one Gravatar avatar for a Clearbit person.
type CBAvatar struct {
	Url string
}

// CBGravatar represents a Clearbit person's Gravatar account data.
type CBGravatar struct {
	Avatars []*CBAvatar
}

// CBGithub represents a Clearbit person's Github account data.
type CBGithub struct {
	Followers int
}

// CBName represents a Clearbit person's name.
type CBName struct {
	FullName string
}

// CBPerson represents a Clearbit person.
type CBPerson struct {
	Name     *CBName
	Github   *CBGithub
	Gravatar *CBGravatar
}

// MediumPayload contains the medium payload data.
type MediumPayload struct {
	Person  *CBPerson
	Company map[string]interface{}
}

// Large payload, based on Discourse API. Size: 41kb

// DSUser represents a Discourse user.
type DSUser struct {
	Username string
}

// DSTopic represents one Discourse topic.
type DSTopic struct {
	Id   int
	Slug string
}

// DSTopicsList represents a paginated set of Discourse topics.
type DSTopicsList struct {
	Topics        []*DSTopic
	MoreTopicsUrl string
}

// LargePayload contains the large payload data.
type LargePayload struct {
	Users  []*DSUser
	Topics *DSTopicsList
}

// Huge payload, based on a large Helm index. Size: 333mb

// IndexFile contains the huge payload data.
type IndexFile struct {
	APIVersion string                   `json:"apiVersion"`
	Generated  time.Time                `json:"generated"`
	Entries    map[string]ChartVersions `json:"entries"`
}

// ChartVersions is a list of versions of a chart.
type ChartVersions []ChartVersion

// ChartVersion is one version of a chart.
type ChartVersion struct {
	Metadata
	URLs    []string  `json:"urls"`
	Created time.Time `json:"created,omitempty"`
	Removed bool      `json:"removed,omitempty"`
	Digest  string    `json:"digest,omitempty"`
}

// Metadata is the metadata for a chart.
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

// Maintainer is the information of a chart maintainer.
type Maintainer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Url   string `json:"url,omitempty"`
}
