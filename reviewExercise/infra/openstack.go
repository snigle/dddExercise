package infra

import (
	"github.com/snigle/dddExercise/reviewExercise/bootstrap"
	"github.com/snigle/dddExercise/reviewExercise/domain/port"
)

type OpenStack struct {
	// Conf qu'on récupère du bootstrap
	conf bootstrap.Config
}

func Constructor() port.CloudProvider {
	return OpenStack{
		conf: bootstrap.Config{},
	}
}

// constructeur --> retourner une instance qui créer le client

func ListRegions() {

}

// List servers
func ListServers(domain.Region) {

}

// GetImage
func GetImage() {}
