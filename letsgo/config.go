package letsgo

type Configeration struct {
	Round     int     `json:"round"`
	StartX    int     `json:"startX"`
	StartY    int     `json:"startY"`
	Rmin      int     `json:"rmin"`
	Rmax      int     `json:"rmax"`
	Gmin      int     `json:"gmin"`
	Gmax      int     `json:"gmax"`
	Bmin      int     `json:"bmin"`
	Bmax      int     `json:"bmax"`
	RoundTime float64 `json:"roundTime"`
	SpellTime int64   `json:"spellTime"`
	Interval  int64   `json:"interval"`
	Key       string  `json:"key"`
}
