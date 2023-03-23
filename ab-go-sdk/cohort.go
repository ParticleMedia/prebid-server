package ab

import (
	"github.com/RoaringBitmap/roaring"
)

type Cohort struct {
	Location   string          `json:"location"`
	Users      *roaring.Bitmap `json:"-"`
	UpdateTime string          `json:"update_time"`
}

func (this *Cohort) Contains(x uint32) bool {
	return this.Users != nil && this.Users.Contains(x)
}
