package tuf

import (
	"time"

	"github.com/theupdateframework/go-tuf/data"
)

func (r *Repo) IncrementRootVersion() error {
	return r.IncrementRootVersionWithExpires(data.DefaultExpires("root"))
}

func (r *Repo) IncrementRootVersionWithExpires(expires time.Time) error {
	role, err := r.root()
	if err != nil {
		return err
	}

	if !validExpires(expires) {
		return ErrInvalidExpires{expires}
	}

	role.Expires = expires.Round(time.Second)
	if _, ok := r.versionUpdated["root.json"]; !ok {
		role.Version++
		r.versionUpdated["root.json"] = struct{}{}
	}

	return r.setMeta("root.json", role)
}

func (r *Repo) RootExpires() (time.Time, error) {
	role, err := r.root()
	if err != nil {
		return time.Time{}, err
	}
	return role.Expires, nil
}

func (r *Repo) IncrementSnapshotVersion() error {
	return r.Snapshot(CompressionTypeNone)
}

func (r *Repo) IncrementSnapshotVersionWithExpires(expires time.Time) error {
	return r.SnapshotWithExpires(CompressionTypeNone, expires)
}

func (r *Repo) SnapshotExpires() (time.Time, error) {
	role, err := r.snapshot()
	if err != nil {
		return time.Time{}, err
	}
	return role.Expires, nil
}

func (r *Repo) IncrementTargetsVersion() error {
	return r.IncrementTargetsVersionWithExpires(data.DefaultExpires("targets"))
}

func (r *Repo) IncrementTargetsVersionWithExpires(expires time.Time) error {
	role, err := r.targets()
	if err != nil {
		return err
	}

	if !validExpires(expires) {
		return ErrInvalidExpires{expires}
	}

	role.Expires = expires.Round(time.Second)
	if _, ok := r.versionUpdated["targets.json"]; !ok {
		role.Version++
		r.versionUpdated["targets.json"] = struct{}{}
	}

	return r.setMeta("targets.json", role)
}

func (r *Repo) TargetsExpires() (time.Time, error) {
	role, err := r.targets()
	if err != nil {
		return time.Time{}, err
	}
	return role.Expires, nil
}

func (r *Repo) IncrementTimestampVersion() error {
	return r.Timestamp()
}

func (r *Repo) IncrementTimestampVersionWithExpires(expires time.Time) error {
	return r.TimestampWithExpires(expires)
}

func (r *Repo) TimestampExpires() (time.Time, error) {
	role, err := r.timestamp()
	if err != nil {
		return time.Time{}, err
	}
	return role.Expires, nil
}
