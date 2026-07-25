package migrations

import (
	"fmt"
	"sort"
)

var all []Migration

func register(m Migration) {
	all = append(all, m)
}

// All returns registered migrations sorted by version ascending.
func All() ([]Migration, error) {
	out := append([]Migration(nil), all...)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Version < out[j].Version
	})
	seen := map[int64]string{}
	for _, m := range out {
		if m.Version <= 0 {
			return nil, fmt.Errorf("migrations: invalid version %d for %q", m.Version, m.Name)
		}
		if prev, ok := seen[m.Version]; ok {
			return nil, fmt.Errorf("migrations: duplicate version %d (%q and %q)", m.Version, prev, m.Name)
		}
		if m.Up == nil {
			return nil, fmt.Errorf("migrations: %d_%s has nil Up", m.Version, m.Name)
		}
		seen[m.Version] = m.Name
	}
	return out, nil
}
