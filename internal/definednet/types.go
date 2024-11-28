package definednet

// ConfigOverride is a data model for Defined.net host configuration override.
type ConfigOverride struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
