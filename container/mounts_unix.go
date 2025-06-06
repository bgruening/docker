//go:build !windows

package container

// Mount contains information for a mount operation.
type Mount struct {
	Source                 string `json:"source"`
	Destination            string `json:"destination"`
	Writable               bool   `json:"writable"`
	Data                   string `json:"data"`
	Propagation            string `json:"mountpropagation"`
	NonRecursive           bool   `json:"nonrecursive"`
	ReadOnlyNonRecursive   bool   `json:"readonlynonrecursive"`
	ReadOnlyForceRecursive bool   `json:"readonlyforcerecursive"`
}
