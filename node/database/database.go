package database

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/zpatrick/go-config"
)

// PegnetNodeDatabase is the sql database for a full pegnet node
// The original primary purpose for this node is to provide time series data information
// TODO: Add more functionality to this db
type PegnetNodeDatabase struct {
	DB *gorm.DB
}

func NewPegnetNodeDatabase(config *config.Config) (*PegnetNodeDatabase, error) {
	n := new(PegnetNodeDatabase)
	var err error
	n.DB, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	return n, nil
}

// InsertTimeSeries will insert a time series if it is not already added
func InsertTimeSeries(db *gorm.DB, r interface{}) error {
	// If we have a conflict, the data should already be there
	err := db.Set("gorm:insert_option", "ON CONFLICT(block_height) DO NOTHING").Create(r)
	return err.Error
}

func FetchTimeSeries(db *gorm.DB, target interface{}, opts *FetchTimeSeriesParams) error {
	//var test []DifficultyTimeSeries
	res := opts.Apply(db).Find(target)
	return res.Error
}

func (p *PegnetNodeDatabase) Migrate() {
	AutoMigrateTimeSeries(p.DB)
}

type FetchTimeSeriesParams struct {
	Limit       int   `json:"limit"`
	Offset      int   `json:"offset"`
	StartHeight int   `json:"startheight"`
	StopHeight  int   `json:"stopheight"`
	UnixStart   int64 `json:"unixstart"`
	UnixStop    int64 `json:"unixstop"`
	Descend     bool  `json:"descend"`
}

func (p *FetchTimeSeriesParams) Apply(db *gorm.DB) *gorm.DB {
	r := db
	if p.Limit != 0 {
		r = r.Limit(p.Limit)
	}

	if p.Offset != 0 {
		r = r.Offset(p.Offset)
	}

	if p.StartHeight != 0 {
		r = r.Where("block_height >= ?", p.StartHeight)
	}

	if p.StopHeight != 0 {
		r = r.Where("block_height < ?", p.StopHeight)
	}

	if p.UnixStart != 0 {
		r = r.Where("timestamp >= ?", time.Unix(p.UnixStart, 0))
	}

	if p.UnixStop != 0 {
		r = r.Where("timestamp < ?", time.Unix(p.UnixStop, 0))
	}

	if p.Descend {
		r = r.Order("block_height DESC")
	} else {
		r = r.Order("block_height ASC")
	}

	return r
}
