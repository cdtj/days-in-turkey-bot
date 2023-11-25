package model

// Country is struct to store User Country settings,
// fields are exported in case to store them in DB
type Country struct {
	Code          string
	Name          string
	Flag          string
	DaysContinual int
	DaysLimit     int
	ResetInterval int
	VisaFree      bool
}

func NewCountry(code, flag, name string, daysContinual, daysLimit, resetInterval int, visaFree bool) *Country {
	return &Country{
		Name:          name,
		Code:          code,
		Flag:          flag,
		DaysContinual: daysContinual,
		DaysLimit:     daysLimit,
		ResetInterval: resetInterval,
		VisaFree:      visaFree,
	}
}

func (c *Country) GetResetInterval() int {
	return c.ResetInterval
}

func (c *Country) GetDaysCont() int {
	return c.DaysContinual
}

func (c *Country) GetDaysLimit() int {
	return c.DaysLimit
}

func (c *Country) GetCode() string {
	return c.Code
}

func (c *Country) GetFlag() string {
	return c.Flag
}

func (c *Country) GetName() string {
	return c.Name
}

func (c *Country) GetVisaFree() bool {
	return c.VisaFree
}
