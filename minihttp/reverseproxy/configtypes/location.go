package configtypes

// Location is used to map some host(Target) through a HTTP path
type Location struct {
	Path   string
	Target string
}
