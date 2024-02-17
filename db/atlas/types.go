package atlas

type AtlasConfig struct {
	URL             string `json:"url"`
	DevUrl          string `json:"dev"`
	RevisionsSchema string `json:"revisionsSchema"`
}

type AtlasOptions struct {
	config     AtlasConfig
	initialSql string
	dialect    string
}
