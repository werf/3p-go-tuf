package tuf

import (
	. "gopkg.in/check.v1"
)

func genKeyErr(_ []string, err error) error {
	return err
}

func initRepo(c *C, r *Repo, paths []string) {
	c.Assert(r.Init(false), IsNil)
	c.Assert(genKeyErr(r.GenKey("root")), IsNil)
	c.Assert(genKeyErr(r.GenKey("targets")), IsNil)
	c.Assert(genKeyErr(r.GenKey("snapshot")), IsNil)
	c.Assert(genKeyErr(r.GenKey("timestamp")), IsNil)
	c.Assert(r.AddTargets(paths, nil), IsNil)
	c.Assert(r.Snapshot(), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)
}

func (RepoSuite) TestRolesExpirationsRotation(c *C) {
	local := MemoryStore(nil, map[string][]byte{
		"myfile": []byte("DATA\n"),
	})

	r, err := NewRepo(local)
	c.Assert(err, IsNil)

	initRepo(c, r, []string{"myfile"})

	prevRootVersion, err := r.RootVersion()
	c.Assert(err, IsNil)
	prevRootExpires, err := r.RootExpires()
	c.Assert(err, IsNil)

	prevTargetsVersion, err := r.TargetsVersion()
	c.Assert(err, IsNil)
	prevTargetsExpires, err := r.TargetsExpires()
	c.Assert(err, IsNil)

	prevSnapshotVersion, err := r.SnapshotVersion()
	c.Assert(err, IsNil)
	prevSnapshotExpires, err := r.SnapshotExpires()
	c.Assert(err, IsNil)

	prevTimestampVersion, err := r.TimestampVersion()
	c.Assert(err, IsNil)
	prevTimestampExpires, err := r.TimestampExpires()
	c.Assert(err, IsNil)

	// Update timestamp
	for i := 0; i < 10; i++ {
		r, err := NewRepo(local)
		c.Assert(err, IsNil)

		newTimestampExpires := prevTimestampExpires.AddDate(0, 0, 1)

		c.Assert(r.IncrementTimestampVersionWithExpires(newTimestampExpires), IsNil)
		c.Assert(r.Commit(), IsNil)

		root, err := r.root()
		c.Assert(err, IsNil)
		c.Assert(root.Version, Equals, prevRootVersion)
		c.Assert(root.Expires, Equals, prevRootExpires)

		targets, err := r.topLevelTargets()
		c.Assert(err, IsNil)
		c.Assert(targets.Version, Equals, prevTargetsVersion)
		c.Assert(targets.Expires, Equals, prevTargetsExpires)

		snapshot, err := r.snapshot()
		c.Assert(err, IsNil)
		c.Assert(snapshot.Version, Equals, prevSnapshotVersion)
		c.Assert(snapshot.Expires, Equals, prevSnapshotExpires)

		timestamp, err := r.timestamp()
		c.Assert(err, IsNil)
		c.Assert(timestamp.Version, Equals, prevTimestampVersion+1)
		c.Assert(timestamp.Expires, Equals, newTimestampExpires)
		prevTimestampVersion = timestamp.Version
		prevTimestampExpires = timestamp.Expires
	}

	// Update snapshot, depends on timestamp
	for i := 0; i < 10; i++ {
		r, err := NewRepo(local)
		c.Assert(err, IsNil)

		newSnapshotExpires := prevSnapshotExpires.AddDate(0, 0, 7)
		newTimestampExpires := prevTimestampExpires.AddDate(0, 0, 1)

		c.Assert(r.IncrementSnapshotVersionWithExpires(newSnapshotExpires), IsNil)
		c.Assert(r.IncrementTimestampVersionWithExpires(newTimestampExpires), IsNil)
		c.Assert(r.Commit(), IsNil)

		root, err := r.root()
		c.Assert(err, IsNil)
		c.Assert(root.Version, Equals, prevRootVersion)
		c.Assert(root.Expires, Equals, prevRootExpires)

		targets, err := r.topLevelTargets()
		c.Assert(err, IsNil)
		c.Assert(targets.Version, Equals, prevTargetsVersion)
		c.Assert(targets.Expires, Equals, prevTargetsExpires)

		snapshot, err := r.snapshot()
		c.Assert(err, IsNil)
		c.Assert(snapshot.Version, Equals, prevSnapshotVersion+1)
		c.Assert(snapshot.Expires, Equals, newSnapshotExpires)
		prevSnapshotVersion = snapshot.Version
		prevSnapshotExpires = snapshot.Expires

		timestamp, err := r.timestamp()
		c.Assert(err, IsNil)
		c.Assert(timestamp.Version, Equals, prevTimestampVersion+1)
		c.Assert(timestamp.Expires, Equals, newTimestampExpires)
		prevTimestampVersion = timestamp.Version
		prevTimestampExpires = timestamp.Expires
	}

	// Update targets, depends on snapshot and timestamp
	for i := 0; i < 10; i++ {
		r, err := NewRepo(local)
		c.Assert(err, IsNil)

		newTargetsExpires := prevTargetsExpires.AddDate(0, 3, 0)
		newSnapshotExpires := prevSnapshotExpires.AddDate(0, 0, 7)
		newTimestampExpires := prevTimestampExpires.AddDate(0, 0, 1)

		c.Assert(r.IncrementTargetsVersionWithExpires(newTargetsExpires), IsNil)
		c.Assert(r.IncrementSnapshotVersionWithExpires(newSnapshotExpires), IsNil)
		c.Assert(r.IncrementTimestampVersionWithExpires(newTimestampExpires), IsNil)
		c.Assert(r.Commit(), IsNil)

		root, err := r.root()
		c.Assert(err, IsNil)
		c.Assert(root.Version, Equals, prevRootVersion)
		c.Assert(root.Expires, Equals, prevRootExpires)

		targets, err := r.topLevelTargets()
		c.Assert(err, IsNil)
		c.Assert(targets.Version, Equals, prevTargetsVersion+1)
		c.Assert(targets.Expires, Equals, newTargetsExpires)
		prevTargetsVersion = targets.Version
		prevTargetsExpires = targets.Expires

		snapshot, err := r.snapshot()
		c.Assert(err, IsNil)
		c.Assert(snapshot.Version, Equals, prevSnapshotVersion+1)
		c.Assert(snapshot.Expires, Equals, newSnapshotExpires)
		prevSnapshotVersion = snapshot.Version
		prevSnapshotExpires = snapshot.Expires

		timestamp, err := r.timestamp()
		c.Assert(err, IsNil)
		c.Assert(timestamp.Version, Equals, prevTimestampVersion+1)
		c.Assert(timestamp.Expires, Equals, newTimestampExpires)
		prevTimestampVersion = timestamp.Version
		prevTimestampExpires = timestamp.Expires
	}

	// Update root, depends on snapshot and timestamp
	for i := 0; i < 10; i++ {
		r, err := NewRepo(local)
		c.Assert(err, IsNil)

		newRootExpires := prevRootExpires.AddDate(1, 0, 0)
		newSnapshotExpires := prevSnapshotExpires.AddDate(0, 0, 7)
		newTimestampExpires := prevTimestampExpires.AddDate(0, 0, 1)

		c.Assert(r.IncrementRootVersionWithExpires(newRootExpires), IsNil)
		c.Assert(r.IncrementSnapshotVersionWithExpires(newSnapshotExpires), IsNil)
		c.Assert(r.IncrementTimestampVersionWithExpires(newTimestampExpires), IsNil)
		c.Assert(r.Commit(), IsNil)

		root, err := r.root()
		c.Assert(err, IsNil)
		c.Assert(root.Version, Equals, prevRootVersion+1)
		c.Assert(root.Expires, Equals, newRootExpires)
		prevRootVersion = root.Version
		prevRootExpires = root.Expires

		targets, err := r.topLevelTargets()
		c.Assert(err, IsNil)
		c.Assert(targets.Version, Equals, prevTargetsVersion)
		c.Assert(targets.Expires, Equals, prevTargetsExpires)

		snapshot, err := r.snapshot()
		c.Assert(err, IsNil)
		c.Assert(snapshot.Version, Equals, prevSnapshotVersion+1)
		c.Assert(snapshot.Expires, Equals, newSnapshotExpires)
		prevSnapshotVersion = snapshot.Version
		prevSnapshotExpires = snapshot.Expires

		timestamp, err := r.timestamp()
		c.Assert(err, IsNil)
		c.Assert(timestamp.Version, Equals, prevTimestampVersion+1)
		c.Assert(timestamp.Expires, Equals, newTimestampExpires)
		prevTimestampVersion = timestamp.Version
		prevTimestampExpires = timestamp.Expires
	}

	// Update root, targets, snapshot and timestamp at the same time
	for i := 0; i < 10; i++ {
		r, err := NewRepo(local)
		c.Assert(err, IsNil)

		newRootExpires := prevRootExpires.AddDate(1, 0, 0)
		newTargetsExpires := prevTargetsExpires.AddDate(0, 3, 0)
		newSnapshotExpires := prevSnapshotExpires.AddDate(0, 0, 7)
		newTimestampExpires := prevTimestampExpires.AddDate(0, 0, 1)

		c.Assert(r.IncrementRootVersionWithExpires(newRootExpires), IsNil)
		c.Assert(r.IncrementTargetsVersionWithExpires(newTargetsExpires), IsNil)
		c.Assert(r.IncrementSnapshotVersionWithExpires(newSnapshotExpires), IsNil)
		c.Assert(r.IncrementTimestampVersionWithExpires(newTimestampExpires), IsNil)
		c.Assert(r.Commit(), IsNil)

		root, err := r.root()
		c.Assert(err, IsNil)
		c.Assert(root.Version, Equals, prevRootVersion+1)
		c.Assert(root.Expires, Equals, newRootExpires)
		prevRootVersion = root.Version
		prevRootExpires = root.Expires

		targets, err := r.topLevelTargets()
		c.Assert(err, IsNil)
		c.Assert(targets.Version, Equals, prevTargetsVersion+1)
		c.Assert(targets.Expires, Equals, newTargetsExpires)
		prevTargetsVersion = targets.Version
		prevTargetsExpires = targets.Expires

		snapshot, err := r.snapshot()
		c.Assert(err, IsNil)
		c.Assert(snapshot.Version, Equals, prevSnapshotVersion+1)
		c.Assert(snapshot.Expires, Equals, newSnapshotExpires)
		prevSnapshotVersion = snapshot.Version
		prevSnapshotExpires = snapshot.Expires

		timestamp, err := r.timestamp()
		c.Assert(err, IsNil)
		c.Assert(timestamp.Version, Equals, prevTimestampVersion+1)
		c.Assert(timestamp.Expires, Equals, newTimestampExpires)
		prevTimestampVersion = timestamp.Version
		prevTimestampExpires = timestamp.Expires
	}
}
