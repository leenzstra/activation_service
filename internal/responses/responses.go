package responses

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	licenseDuratonPattern, _ = regexp.Compile(`(\d*)y(\d*)m(\d*)d`)
)

type ResponseBase struct {
	Result  bool        `json:"result"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type LicenseDuration struct {
	Years  int
	Months int
	Days   int
}

func (d *LicenseDuration) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	return d.FromString(string(data))
}

func (d *LicenseDuration) MarshalJSON() ([]byte, error) {
	data := fmt.Sprintf("%dy%dm%dd", d.Years, d.Months, d.Days)
	return []byte(data), nil
}

func (d *LicenseDuration) FromString(period string) error {
	groups := licenseDuratonPattern.FindStringSubmatch(period)
	fmt.Println(groups)
	if len(groups) < 4 {
		return errors.New("bad input " + period)
	}
	years, err := strconv.Atoi(groups[1])
	if err != nil {
		return err
	}
	months, err := strconv.Atoi(groups[2])
	if err != nil {
		return err
	}
	days, err := strconv.Atoi(groups[3])
	if err != nil {
		return err
	}
	*d = LicenseDuration{Years: years, Months: months, Days: days}
	return nil
}