package dbdata

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bjdgyc/anylink/base"
)

const (
	LayoutTimeFormat    = "2006-01-02 15:04:05"
	LayoutTimeFormatMin = "2006-01-02 15:04"
	RealTimeMaxSize     = 120 // Maximum number of real-time data saved
)

type StatsInfo struct {
	RealtimeData map[string]*list.List
	Actions      []string
	Scopes       []string
	mux          sync.Mutex
}

type ScopeDetail struct {
	sTime   time.Time
	eTime   time.Time
	minutes int
	fsTime  string
	feTime  string
}

var StatsInfoIns *StatsInfo

func init() {
	StatsInfoIns = &StatsInfo{
		Actions:      []string{"online", "network", "cpu", "mem"},
		Scopes:       []string{"rt", "1h", "24h", "3d", "7d", "30d"},
		RealtimeData: make(map[string]*list.List),
	}
	for _, v := range StatsInfoIns.Actions {
		StatsInfoIns.RealtimeData[v] = list.New()
	}
}

// Check statistics type value
func (s *StatsInfo) ValidAction(action string) bool {
	for _, item := range s.Actions {
		if item == action {
			return true
		}
	}
	return false
}

// Check date range value
func (s *StatsInfo) ValidScope(scope string) bool {
	for _, item := range s.Scopes {
		if item == scope {
			return true
		}
	}
	return false
}

// Set up real-time statistics
func (s *StatsInfo) SetRealTime(action string, val interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.RealtimeData[action].Len() >= RealTimeMaxSize {
		ele := s.RealtimeData[action].Front()
		s.RealtimeData[action].Remove(ele)
	}
	s.RealtimeData[action].PushBack(val)
}

// Get real-time statistics
func (s *StatsInfo) GetRealTime(action string) (res []interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for e := s.RealtimeData[action].Front(); e != nil; e = e.Next() {
		res = append(res, e.Value)
	}
	return
}

// Save data to database
func (s *StatsInfo) SaveStatsInfo(so StatsOnline, sn StatsNetwork, sc StatsCpu, sm StatsMem) {
	if so.Num != 0 {
		_ = Add(so)
	}
	if sn.Up != 0 || sn.Down != 0 {
		_ = Add(sn)
	}
	if sc.Percent != 0 {
		_ = Add(sc)
	}
	if sm.Percent != 0 {
		_ = Add(sm)
	}
}

// Get statistics
func (s *StatsInfo) GetData(action string, scope string) (res []interface{}, err error) {
	if scope == "rt" {
		return s.GetRealTime(action), nil
	}
	statsMaps := make(map[string]interface{})
	currSec := fmt.Sprintf("%02d", time.Now().Second())

	// Get time period data
	sd := s.getScopeDetail(scope)
	timeList := s.getTimeList(sd)
	res = make([]interface{}, len(timeList))

	// Get database query conditions
	where := s.getStatsWhere(sd)
	if where == "" {
		return nil, errors.New("Unsupported database type: " + base.Cfg.DbType)
	}
	// Query data table
	switch action {
	case "online":
		statsRes := []StatsOnline{}
		FindWhere(&statsRes, 0, 0, where, sd.fsTime, sd.feTime)
		for _, v := range statsRes {
			t := v.CreatedAt.Format(LayoutTimeFormatMin)
			statsMaps[t] = v
		}
	case "network":
		statsRes := []StatsNetwork{}
		FindWhere(&statsRes, 0, 0, where, sd.fsTime, sd.feTime)
		for _, v := range statsRes {
			t := v.CreatedAt.Format(LayoutTimeFormatMin)
			statsMaps[t] = v
		}
	case "cpu":
		statsRes := []StatsCpu{}
		FindWhere(&statsRes, 0, 0, where, sd.fsTime, sd.feTime)
		for _, v := range statsRes {
			t := v.CreatedAt.Format(LayoutTimeFormatMin)
			statsMaps[t] = v
		}
	case "mem":
		statsRes := []StatsMem{}
		FindWhere(&statsRes, 0, 0, where, sd.fsTime, sd.feTime)
		for _, v := range statsRes {
			t := v.CreatedAt.Format(LayoutTimeFormatMin)
			statsMaps[t] = v
		}
	}
	// Integrate data
	for i, v := range timeList {
		if mv, ok := statsMaps[v]; ok {
			res[i] = mv
			continue
		}
		t, _ := time.ParseInLocation(LayoutTimeFormat, v+":"+currSec, time.Local)
		switch action {
		case "online":
			res[i] = StatsOnline{CreatedAt: t}
		case "network":
			res[i] = StatsNetwork{CreatedAt: t}
		case "cpu":
			res[i] = StatsCpu{CreatedAt: t}
		case "mem":
			res[i] = StatsMem{CreatedAt: t}
		}
	}
	return
}

// Get detailed values ​​for a date range
func (s *StatsInfo) getScopeDetail(scope string) (sd *ScopeDetail) {
	sd = &ScopeDetail{}
	t := time.Now()
	sd.eTime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 59, t.Nanosecond(), time.Local)
	sd.minutes = 0
	switch scope {
	case "1h":
		sd.sTime = sd.eTime.Add(-time.Minute * 60)
		sd.minutes = 1
	case "24h":
		sd.sTime = sd.eTime.AddDate(0, 0, -1)
		sd.minutes = 5
	case "7d":
		sd.sTime = sd.eTime.AddDate(0, 0, -7)
		sd.minutes = 30
	case "30d":
		sd.sTime = sd.eTime.AddDate(0, 0, -30)
		sd.minutes = 150
	}
	if sd.minutes != 0 {
		sd.sTime = sd.sTime.Add(-time.Minute * time.Duration(sd.minutes))
	}
	sd.fsTime = sd.sTime.Format(LayoutTimeFormat)
	sd.feTime = sd.eTime.Format(LayoutTimeFormat)
	return
}

// Split for date range
func (s *StatsInfo) getTimeList(sd *ScopeDetail) []string {
	subSec := int64(60 * sd.minutes)
	count := (sd.eTime.Unix()-sd.sTime.Unix())/subSec - 1
	eTime := sd.eTime.Unix() - subSec
	timeLists := make([]string, count)
	for i := count - 1; i >= 0; i-- {
		timeLists[i] = time.Unix(eTime, 0).Format(LayoutTimeFormatMin)
		eTime = eTime - subSec
	}
	return timeLists
}

// Get where condition
func (s *StatsInfo) getStatsWhere(sd *ScopeDetail) (where string) {
	where = "created_at BETWEEN ? AND ?"
	min := strconv.Itoa(sd.minutes)
	switch base.Cfg.DbType {
	case "mysql":
		where += " AND floor(TIMESTAMPDIFF(SECOND, created_at, '" + sd.feTime + "') / 60) % " + min + " = 0"
	case "sqlite3":
		where += " AND CAST(ROUND((JULIANDAY('" + sd.feTime + "') - JULIANDAY(created_at)) * 86400) / 60 as integer) % " + min + " = 0"
	case "postgres":
		where += " AND floor((EXTRACT(EPOCH FROM " + sd.feTime + ") - EXTRACT(EPOCH FROM created_at)) / 60) % " + min + " = 0"
	default:
		where = ""
	}
	return
}

func (s *StatsInfo) ClearStatsInfo(action string, ts string) (affected int64, err error) {
	switch action {
	case "online":
		affected, err = xdb.Where("created_at < '" + ts + "'").Delete(&StatsOnline{})
	case "network":
		affected, err = xdb.Where("created_at < '" + ts + "'").Delete(&StatsNetwork{})
	case "cpu":
		affected, err = xdb.Where("created_at < '" + ts + "'").Delete(&StatsCpu{})
	case "mem":
		affected, err = xdb.Where("created_at < '" + ts + "'").Delete(&StatsMem{})
	}
	return affected, err
}
