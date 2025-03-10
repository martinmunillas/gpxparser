package gpx

import (
	"github.com/golang/geo/s2"
)

const earthMeterRadius = 6371000

func (gpx *GPX) WithDistances() *GPX {
	if gpx.computedDistances {
		return gpx
	}
	totalDistance := 0.0
	totalClimb := 0.0
	totalDescent := 0.0

	for i, track := range gpx.Tracks {
		for j, segment := range track.Segments {
			for k, b := range segment.TrackPoints {
				if k == 0 {
					gpx.Tracks[i].Segments[j].TrackPoints[i].RunningDistance = totalDistance
					continue
				}
				a := segment.TrackPoints[k-1]
				latLng1 := s2.LatLngFromDegrees(a.Latitude, a.Longitude)
				latLng2 := s2.LatLngFromDegrees(b.Latitude, b.Longitude)
				dist := s2.ChordAngleBetweenPoints(s2.PointFromLatLng(latLng1), s2.PointFromLatLng(latLng2)).Angle().Radians() * earthMeterRadius

				gpx.Tracks[i].Segments[j].TrackPoints[k].DistanceWithPrevious = dist
				totalDistance += dist
				gpx.Tracks[i].Segments[j].TrackPoints[k].RunningDistance = totalDistance

				elevationDiff := b.Elevation - a.Elevation
				if elevationDiff > 0 {
					totalClimb += elevationDiff
				} else {
					totalDescent -= elevationDiff
				}
			}
		}
	}

	gpx.TotalDistanceMeters = uint64(totalDistance)
	gpx.NetElevationMeters = int64(totalClimb - totalDescent)
	gpx.TotalClimbMeters = uint64(totalClimb)
	gpx.TotalDescentMeters = uint64(totalDescent)

	return gpx
}
