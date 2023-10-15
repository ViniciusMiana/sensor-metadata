package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapbox(t *testing.T) {
	service := NewMapBox("pk.eyJ1IjoidmlubnltaWFuYSIsImEiOiJjbG5udHIza3gwOGlvMndwMTQzM3prdTdnIn0.mlCNLgTF3Ctvyslbao5pRw")
	loc, err := service.FindLatLon("Los Angeles")
	require.NoError(t, err)
	require.Equal(t, Location{Lat: "34.053691", Lon: "-118.242766"}, *loc)
}
