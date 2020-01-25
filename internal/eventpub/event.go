package eventpub

const (
	TypeRepositoryAdded = "REPOSITORY_ADDED"
)

type RepositoryAddedEvent struct {
	Branch     string `json:"branch"`     // master
	Repository string `json:"repository"` // full path, github.com/foo/bar
	GoVersion  string `json:"goVersion"`  // 1.13.6
}
