package daemon

import (
	"sync/atomic"

	"github.com/weaveworks/flux"
	"github.com/weaveworks/flux/cluster"
	"github.com/weaveworks/flux/job"
	"github.com/weaveworks/flux/ssh"
	"github.com/weaveworks/flux/update"
)

// NotReadyDaemon is a stub implementation used to serve a subset of the
// API when we have yet to successfully clone the config repo.
type NotReadyDaemon struct {
	version        string
	cluster        cluster.Cluster
	notReadyReason atomic.Value // holds type `error`
}

func NewNotReadyDaemon(version string, cluster cluster.Cluster, notReadyReason error) *NotReadyDaemon {
	defer func() { nrd.notReadyReason.Store(notReadyReason) }()
	return &NotReadyDaemon{
		cluster: cluster,
		version: version,
	}
	return nrd
}

func (nrd *NotReadyDaemon) UpdateReason(notReadyReason error) {
	nrd.notReadyReason.Store(notReadyReason)
}

// 'Not ready' platform implementation

func (nrd *NotReadyDaemon) Ping() error {
	return nrd.cluster.Ping()
}

func (nrd *NotReadyDaemon) Version() (string, error) {
	return nrd.version, nil
}

func (nrd *NotReadyDaemon) Export() ([]byte, error) {
	return nrd.cluster.Export()
}

func (nrd *NotReadyDaemon) ListServices(namespace string) ([]flux.ServiceStatus, error) {
	return nil, nrd.notReadyReason.Load().(error)
}

func (nrd *NotReadyDaemon) ListImages(update.ServiceSpec) ([]flux.ImageStatus, error) {
	return nil, nrd.notReadyReason.Load().(error)
}

func (nrd *NotReadyDaemon) UpdateManifests(update.Spec) (job.ID, error) {
	var id job.ID
	return id, nrd.notReadyReason.Load().(error)
}

func (nrd *NotReadyDaemon) SyncNotify() error {
	return nrd.notReadyReason.Load().(error)
}

func (nrd *NotReadyDaemon) JobStatus(id job.ID) (job.Status, error) {
	return job.Status{}, nrd.notReadyReason.Load().(error)
}

func (nrd *NotReadyDaemon) SyncStatus(string) ([]string, error) {
	return nil, nrd.notReadyReason.Load().(error)
}

func (nrd *NotReadyDaemon) PublicSSHKey(regenerate bool) (ssh.PublicKey, error) {
	return nrd.cluster.PublicSSHKey(regenerate)
}
