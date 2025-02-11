package godemlia

type Node struct {
	ID   []byte `json:"id,omitempty"`
	IP   string `json:"ip,omitempty"`
	Port int    `json:"port,omitempty"`
}
