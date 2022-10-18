package commit

type Commit struct {
	Hash         string              `json:"commitHash"`
	AuthorName   string              `json:"authorName"`
	AuthorEmail  string              `json:"authorEmail"`
	Date         string              `json:"createdAt"`
	ChangedFiles []*ChangedFile      `json:"changedFiles"`
	Libraries    map[string][]string `json:"libraries"`
}

type ChangedFile struct {
	Path       string `json:"fileName"`
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
	Language   string `json:"language"`
}
