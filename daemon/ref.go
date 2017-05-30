package daemon

import (
	"sync/atomic"

	"github.com/weaveworks/flux/remote"
)

type Ref struct {
	value atomic.Value // holds remote.Platform
}

func NewRef(platform remote.Platform) (pr *Ref) {
	defer func() { pr.Store(platform) }()
	return &Ref{}
}

func (pr *Ref) Load() remote.Platform {
	return pr.value.Load().(remote.Platform)
}

func (pr *Ref) Store(platform remote.Platform) {
	pr.value.Store(platform)
}

// remote.Platform implementation so clients don't need to be refactored around
// Load() API

func (pr *Ref) Ping() error {
	return pr.Load().Ping()
}

func (pr *Ref) Version() (string, error) {
	return pr.Load().Version()
}

func (pr *Ref) Export() ([]byte, error) {
	return pr.Load().Export()
}

func (pr *Ref) ListServices(namespace string) ([]flux.ServiceStatus, error) {
	return pr.Load().ListServices(namespace)
}

func (pr *Ref) ListImages(spec update.ServiceSpec) ([]flux.ImageStatus, error) {
	return pr.Load().ListImages(spec)
}

func (pr *Ref) UpdateManifests(spec update.Spec) (job.ID, error) {
	return pr.Load().UpdateManifests(spec)
}

func (pr *Ref) SyncNotify() error {
	return pr.Load().SyncNotify()
}

func (pr *Ref) JobStatus(id job.ID) (job.Status, error) {
	return pr.Load().JobStatus(id)
}

func (pr *Ref) SyncStatus(ref string) ([]string, error) {
	return pr.Load().SyncStatus(ref)
}

func (pr *Ref) PublicSSHKey(regenerate bool) (ssh.PublicKey, error) {
	return pr.Load().PublicSSHKey(regenerate)
}
