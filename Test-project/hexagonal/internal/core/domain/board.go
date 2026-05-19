package domain

type Board [][]string

type BoardSettings struct {
	Size  uint `json:"size"`
	Bombs uint `json:"bombs"`
}
