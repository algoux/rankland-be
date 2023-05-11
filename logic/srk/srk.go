package srk

import (
	"encoding/json"
	"errors"
	"rankland/util"
	"reflect"
	"time"
)

var TimeUnit = [...]string{"ms", "s", "min", "h", "d"}

// Duration SRK 格式的时长
// [25, 'ms'] 25ms
// [30, 's'] 30s
// [60, 's'] 1m
type Duration []interface{}

func NewDuration(d time.Duration) *Duration {
	dn := &Duration{}
	dn.SetDuration(d)
	return dn
}

func (d Duration) Duration() (time.Duration, error) {
	if len(d) == 0 {
		return 0, nil
	}
	if len(d) != 2 {
		return 0, errors.New("srk time duration format error")
	}
	if !util.ContainsInSlice([]reflect.Kind{reflect.Int, reflect.Int64, reflect.Float64}, reflect.TypeOf(d[0]).Kind()) {
		return 0, errors.New("parse failed, type of: " + reflect.TypeOf(d[0]).Kind().String())
	}
	if reflect.TypeOf(d[1]).Kind() != reflect.String {
		return 0, errors.New("parse failed, type of: " + reflect.TypeOf(d[1]).Key().String())
	}
	val := time.Duration(int64(reflect.ValueOf(d[0]).Float()))
	typ := reflect.ValueOf(d[1]).String()
	switch typ {
	case TimeUnit[0]:
		return val * time.Millisecond, nil
	case TimeUnit[1]:
		return val * time.Second, nil
	case TimeUnit[2]:
		return val * time.Minute, nil
	case TimeUnit[3]:
		return val * time.Hour, nil
	case TimeUnit[4]:
		return val * 24 * time.Hour, nil
	}

	return 0, errors.New("srk time duration format error, time unit: " + typ)
}

// SetDuration 设置 SRK 时长
// 取最大整数单位，如：60s 可以是 [60, "s"]、[1, "min"] 最大整数单位为 min，所以取后者
func (d *Duration) SetDuration(m time.Duration) {
	if m%time.Second != 0 {
		*d = Duration{m / time.Millisecond, TimeUnit[0]}
	} else if m%time.Minute != 0 {
		*d = Duration{m / time.Second, TimeUnit[1]}
	} else if m%time.Hour != 0 {
		*d = Duration{m / time.Minute, TimeUnit[2]}
	} else if m%(24*time.Hour) != 0 {
		*d = Duration{m / time.Hour, TimeUnit[3]}
	} else {
		*d = Duration{m / 24 / time.Hour, TimeUnit[4]}
	}
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	temp := make([]interface{}, 0)
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	*d = temp
	return nil
}

// type Contest struct {
// 	Title          map[string]string `json:"title"`
// 	StartAt        time.Time         `json:"startAt"`
// 	Duration       Duration          `json:"duration"`
// 	FrozenDuration Duration          `json:"frozenDuration"`
// }

// type Problem struct {
// 	Alias      string            `json:"alias"`
// 	Statistics map[string]int64  `json:"statistics"` // {"accepted":10, "submitted":100}
// 	Style      map[string]string `json:"style"`      // {"backgroundColor":"rgba(189, 14, 14, 0.7)"}
// }

// type Marker struct {
// 	ID    string `json:"id"`
// 	Label string `json:"label"`
// 	Style string `json:"style"`
// }

// type Serie struct {
// 	Title    string    `json:"title"`
// 	Rule     Rule      `json:"rule"`
// 	Segments []Segment `json:"segments"`
// }

// type Rule struct {
// 	Preset  string                 `json:"preset"`
// 	Options map[string]interface{} `json:"options"`
// }

// type Segment struct {
// 	Title string `json:"title"`
// 	Style string `json:"style"`
// }

// type Row struct {
// }

// type Sorter struct {
// 	Algorithm string `json:"algorithm"`
// 	Config    Config `json:"config"`
// }

// type Config struct {
// 	Penalty Duration `json:"penalty"`
// }

// type SRK struct {
// 	Type         string   `json:"type"`
// 	Version      string   `json:"version"`
// 	Contributors []string `json:"contributors"`

// 	Contest  Contest   `json:"contest"`
// 	Problems []Problem `json:"problems"`
// 	Markers  []Marker  `json:"markers"`
// 	Series   []Serie   `json:"series"`
// 	Rows     []Row     `json:"rows"`
// 	Sorter   Sorter    `json:"sorter"`
// }

// func New(options ...Option) *SRK {
// 	srk := &SRK{
// 		Type:         "general",
// 		Version:      "0.3.0",
// 		Contributors: []string{"algoUX (https://algoux.org)"},
// 	}
// 	for _, opt := range options {
// 		opt(srk)
// 	}
// 	return srk
// }

// type Option func(*SRK)

// func WithType(typ string) Option {
// 	return func(s *SRK) {
// 		s.Type = typ
// 	}
// }

// func WithVersion(version string) Option {
// 	return func(s *SRK) {
// 		s.Version = version
// 	}
// }
