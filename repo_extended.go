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
	if !r.local.FileIsStaged("root.json") {
		role.Version++
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
	return r.Snapshot()
}

func (r *Repo) IncrementSnapshotVersionWithExpires(expires time.Time) error {
	return r.SnapshotWithExpires(expires)
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
	role, err := r.topLevelTargets()
	if err != nil {
		return err
	}

	if !validExpires(expires) {
		return ErrInvalidExpires{expires}
	}

	role.Expires = expires.Round(time.Second)
	if !r.local.FileIsStaged("targets.json") {
		role.Version++
	}

	return r.setMeta("targets.json", role)
}

func (r *Repo) TargetsExpires() (time.Time, error) {
	role, err := r.topLevelTargets()
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
