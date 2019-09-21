package commonconfig

import "time"

var defaultDateFormat = "2006-01-02"
var defaultDatetimeFormat = "2006-01-02 15:04:05"
var defaultTimeFormat = "15:04:05"

// TimeConfig app time config
type TimeConfig struct {
	//Timezone time zone.
	Timezone string
	//TimeFormat  format used when converting time in day.
	TimeFormat string
	//DateFormat  format used when converting date.
	DateFormat string
	//DatetimeFormata format used when converting date and time.
	DatetimeFormat string
	location       *time.Location
}

//Parse time string to local time.
//Panic if  Timezone error.
func (c *TimeConfig) Parse(s string) (time.Time, error) {
	format := c.DatetimeFormat
	if format == "" {
		format = defaultDatetimeFormat
	}
	return time.ParseInLocation(format, s, c.loadLocation())
}
func (c *TimeConfig) loadLocation() *time.Location {
	if c.location == nil {
		if c.Timezone == "" {
			c.location = time.Local
		} else {
			var err error
			c.location, err = time.LoadLocation(c.Timezone)
			if err != nil {
				panic(err)
			}
		}
	}
	return c.location
}

//TimeInLocation set time location to given time zone.
//Panic if  Timezone error.
func (c *TimeConfig) TimeInLocation(t time.Time) time.Time {
	return t.In(c.loadLocation())
}

//DateUnix format date from unix timestamp
func (c *TimeConfig) DateUnix(ts int64) string {
	return c.Date(time.Unix(ts, 0))
}

//Date format date.
//Panic if  Timezone error.
func (c *TimeConfig) Date(t time.Time) string {
	localTime := c.TimeInLocation(t)
	if c.DateFormat == "" {
		return localTime.Format(defaultDateFormat)
	}
	return localTime.Format(c.DateFormat)
}

//TimeUnix format time from unix timestamp
func (c *TimeConfig) TimeUnix(ts int64) string {
	return c.Time(time.Unix(ts, 0))
}

//Time format time.
func (c *TimeConfig) Time(t time.Time) string {
	localTime := c.TimeInLocation(t)

	if c.TimeFormat == "" {
		return localTime.Format(defaultTimeFormat)
	}
	return localTime.Format(c.TimeFormat)
}

//DatetimeUnix format date and time from unix timestamp
//Panic if  Timezone error.
func (c *TimeConfig) DatetimeUnix(ts int64) string {
	return c.Datetime(time.Unix(ts, 0))
}

//Datetime format date and time
func (c *TimeConfig) Datetime(t time.Time) string {
	localTime := c.TimeInLocation(t)

	if c.DatetimeFormat == "" {
		return localTime.Format(defaultDatetimeFormat)
	}
	return localTime.Format(c.DatetimeFormat)
}
