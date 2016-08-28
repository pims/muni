package muni

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type RTT struct {
	AgencyList []*Agency `xml:"AgencyList>Agency"`
}

type Agency struct {
	Name         string   `xml:"Name,attr"`
	HasDirection bool     `xml:"HasDirection,attr"`
	Mode         string   `xml:"Mode,attr"`
	RouteList    []*Route `xml:"RouteList>Route"`
}

type Route struct {
	Name               string            `xml:"Name,attr"`
	Code               int               `xml:"Code,attr"`
	RouteDirectionList []*RouteDirection `xml:"RouteDirectionList>RouteDirection"`
}

type RouteDirection struct {
	Code     string  `xml:"Code,attr"`
	Name     string  `xml:"Name,attr"`
	StopList []*Stop `xml:"StopList>Stop"`
}

type Stop struct {
	XMLName           xml.Name `xml:"Stop"`
	Name              string   `xml:"name,attr"`
	StopCode          int      `xml:"StopCode,attr"`
	DepartureTimeList []int    `xml:"DepartureTimeList>DepartureTime"`
}

type HumanReadableSchedule struct {
	Name string
	Next []int
}

type HumanReadableSchedules []HumanReadableSchedule

func (s HumanReadableSchedules) Len() int {
	return len(s)
}
func (s HumanReadableSchedules) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s HumanReadableSchedules) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

type Muni struct {
	Token    string
	client   *http.Client
	endpoint string
}

// New creates a new Muni struct
func New(token string) *Muni {
	return &Muni{
		Token:    token,
		client:   http.DefaultClient,
		endpoint: "http://services.my511.org/Transit2.0/",
	}
}

func (m *Muni) GetNextDeparturesByStopCode(stopcode int) ([]*Route, error) {
	url := fmt.Sprintf(
		"%sGetNextDeparturesByStopCode.aspx?stopcode=%d&token=%s",
		m.endpoint,
		stopcode,
		m.Token)

	httpResp, err := m.client.Get(url)
	if err != nil {
		return []*Route{}, err
	}

	defer httpResp.Body.Close()
	buff, _ := ioutil.ReadAll(httpResp.Body)

	rtt := &RTT{}

	xmlErr := xml.Unmarshal(buff, rtt)
	if err != nil {
		return []*Route{}, xmlErr
	}

	if len(rtt.AgencyList) == 0 {
		msg := fmt.Sprintf("Failed to get route for stop %d", stopcode)
		return []*Route{}, errors.New(msg)
	}
	return rtt.AgencyList[0].RouteList, nil
}

// Next returns a list of human readable schedules
func (m *Muni) Next(stop int) ([]HumanReadableSchedule, error) {
	routes, err := m.GetNextDeparturesByStopCode(stop)
	schedules := make([]HumanReadableSchedule, 0)

	if err != nil {
		return schedules, err
	}

	for _, route := range routes {

		for _, routeDirection := range route.RouteDirectionList {
			for _, stop := range routeDirection.StopList {
				departures := make([]int, 0)

				for _, departureTime := range stop.DepartureTimeList {
					departures = append(departures, departureTime)
				}
				schedule := HumanReadableSchedule{
					Name: route.Name,
					Next: departures,
				}

				schedules = append(schedules, schedule)
			}
		}
	}
	sort.Sort(sort.Reverse(HumanReadableSchedules(schedules)))
	return schedules, nil
}
