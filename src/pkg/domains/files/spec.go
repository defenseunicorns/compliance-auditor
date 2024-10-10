package files

type Spec struct {
	Filepaths []FileInfo `json:"filepaths" yaml:"filepaths"`
}

type FileInfo struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}
