package types

type Entry struct {
	Name         string `json:"name"`
	Time         string `json:"time"`
	Date         string `json:"Date"`
	ErrorWarning bool   `json:"errorWarning"`
}

type Record struct {
	ID       int32   `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Name     string  `json:"name"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Entries  []Entry `json:"entries"`
}

type File struct {
	Name string
	Path string
}

type DeploySuccess struct {
	Success bool
	Entry   Entry
}
