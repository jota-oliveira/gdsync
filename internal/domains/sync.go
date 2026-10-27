// Package domains - Sync
package domains

type SyncPlan struct {
	New      []*File
	Modified []*File
	Deleted  []*File
}
