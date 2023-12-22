package port

type CloudProvider interface {
	ListRegions()

	// List servers
	ListServers(domain.Region) (domain.Servers, error)

	// GetImage
	GetImage(domain.Server) (domain.Image, error)
}
