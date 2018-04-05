package httpshell

type Winsize struct {
	Rows    uint16 `json:"rows"`
	Columns uint16 `json:"columns"`
	// unused
	x uint16
	y uint16
}