package parser

import (
	"testing"

	"github.com/dharmab/skyeye/pkg/brevity"
	"github.com/martinlindhe/unit"
)

func TestParserSpiked(t *testing.T) {
	testCases := []parserTestCase{
		{
			text: "ANYFACE, EAGLE 1 SPIKED 2-7-0",
			expectedRequest: &brevity.SpikedRequest{
				Callsign: "eagle 1",
				Bearing:  unit.Angle(270) * unit.Degree,
			},
			expectedOk: true,
		},
		{
			text: "Anyface Raven 1-4, Spike 0-2-0",
			expectedRequest: &brevity.SpikedRequest{
				Callsign: "raven 1 4",
				Bearing:  unit.Angle(20) * unit.Degree,
			},
			expectedOk: true,
		},
	}
	runParserTestCases(t, New(TestCallsign), testCases)
}
