package suggest

// Dismissals are the sets of records that someone marked "Not duplicates".
type Dismissals struct {
	sets     []map[string]bool
	byRecord map[string][]int
}

// NewDismissals indexes the dismissed sets of record ids.
func NewDismissals(sets [][]string) *Dismissals {
	d := &Dismissals{byRecord: map[string][]int{}}
	for _, ids := range sets {
		set := map[string]bool{}
		for _, id := range ids {
			set[id] = true
		}
		index := len(d.sets)
		d.sets = append(d.sets, set)
		for id := range set {
			d.byRecord[id] = append(d.byRecord[id], index)
		}
	}
	return d
}

// Covers reports whether one dismissed set holds every id. A bucket that it
// covers is hidden. A new record in the bucket is in no set with the others,
// so the bucket shows again.
func (d *Dismissals) Covers(ids []string) bool {
	if len(ids) == 0 {
		return false
	}
	for _, index := range d.byRecord[ids[0]] {
		all := true
		for _, id := range ids[1:] {
			if !d.sets[index][id] {
				all = false
				break
			}
		}
		if all {
			return true
		}
	}
	return false
}
