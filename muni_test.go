package muni

import (
	"encoding/xml"
	"fmt"
	"sort"
	"testing"
)

func TestParsing(t *testing.T) {
	xmlString := `
	<RTT>
	    <AgencyList>
	        <Agency Name="SF-MUNI" HasDirection="True" Mode="Bus">
	            <RouteList>
	                <Route Name="5-Fulton" Code="5">
	                    <RouteDirectionList>
	                        <RouteDirection Code="Inbound" Name="Inbound to Downtown">
	                            <StopList>
	                                <Stop name="McAllister St and Fillmore St" StopCode="15392">
	                                    <DepartureTimeList>
	                                        <DepartureTime>11</DepartureTime>
	                                        <DepartureTime>28</DepartureTime>
	                                        <DepartureTime>45</DepartureTime>
	                                    </DepartureTimeList>
	                                </Stop>

	                            </StopList>
	                        </RouteDirection>
	                    </RouteDirectionList>
	                </Route>

	                <Route Name="Fulton Rapid" Code="5R">
	                    <RouteDirectionList>
	                        <RouteDirection Code="Inbound" Name="Inbound to Downtown">
	                            <StopList>
	                                <Stop name="McAllister St and Fillmore St" StopCode="15392">
	                                    <DepartureTimeList />
	                                </Stop>
	                            </StopList>
	                        </RouteDirection>
	                    </RouteDirectionList>
	                </Route>
	            </RouteList>
	        </Agency>
	    </AgencyList>
	</RTT>`

	rtt := &RTT{}

	err := xml.Unmarshal([]byte(xmlString), rtt)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	wanted := "SF-MUNI"
	got := rtt.AgencyList[0].Name

	if wanted != got {
		t.Fail()
	}

	if len(rtt.AgencyList[0].RouteList) != 2 {
		t.Fail()
	}
}

func TestSorting(t *testing.T) {
	routes := []HumanReadableSchedule{
		{
			Name: "Y",
			Next: []int{1, 2},
		},
		{
			Name: "X",
			Next: []int{1, 2},
		},
	}

	sort.Sort(HumanReadableSchedules(routes))
	fmt.Println(routes)
	if routes[0].Name != "X" {
		t.FailNow()
	}
}

func TestStupid(t *testing.T) {
	if "a" > "b" {
		t.FailNow()
	}
}

func TestStupidAgain(t *testing.T) {
	s := []string{"c", "b", "a"}
	sort.Strings(s)
	if s[0] != "a" {
		t.FailNow()
	}
}

func TestStupidAgainAgain(t *testing.T) {
	s := []string{"c", "b", "a"}
	sort.Reverse(sort.StringSlice(s))
	fmt.Println(s)
	if s[0] != "c" {
		t.FailNow()
	}
}

func TestSortingReverse(t *testing.T) {
	routes := []HumanReadableSchedule{
		{
			Name: "Y",
			Next: []int{1, 2},
		},
		{
			Name: "X",
			Next: []int{1, 2},
		},
		{
			Name: "Z",
			Next: []int{1, 99},
		},
	}

	sort.Sort(sort.Reverse(HumanReadableSchedules(routes)))
	fmt.Println(routes)
	if routes[0].Name != "Z" {
		t.FailNow()
	}
}
