package utils

import (
	"fsd-backend/internal/models"
	db "fsd-backend/prisma/db"
	"sort"
	"time"
)

func PeriodToStartTime(period db.StatisticPeriod) time.Time {

	now := time.Now()
	var startTime time.Time

	switch period {
	case db.StatisticPeriodWeekly:
		startTime = now.AddDate(0, 0, -7)
	case db.StatisticPeriodMonthly:
		startTime = now.AddDate(0, -1, 0)
	case db.StatisticPeriodSemiannual:
		startTime = now.AddDate(0, -6, 0)
	case db.StatisticPeriodAnnual:
		startTime = now.AddDate(-1, 0, 0)
	default:
		startTime = now.AddDate(0, 0, -7)

	}

	return startTime
}

func TopNKeys(counter map[string]int, n int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range counter {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	var top []string
	for i := 0; i < len(sorted) && i < n; i++ {
		top = append(top, sorted[i].Key)
	}
	return top
}

func ComputeUserStatistics(rp *models.RecentlyPlayed) models.ComputedStats {
	trackMap := map[string]int{}
	artistMap := map[string]int{}
	albumMap := map[string]int{}
	durationSum := 0
	for _, item := range rp.Items {
		track := item.Track
		trackMap[track.ID]++
		durationSum += track.DurationMs
		for _, artist := range track.Artists {
			artistMap[artist.ID]++
		}
		albumMap[track.Album.ID]++
	}
	return models.ComputedStats{
		TotalTracks:   len(rp.Items),
		TotalDuration: durationSum / 1000,
		UniqueArtists: len(artistMap),
		Vibe:          "happy go lucky", // TODO: infer vibe from rp
		TopArtistIDs:  TopNKeys(artistMap, 5),
		TopTrackIDs:   TopNKeys(trackMap, 5),
		TopAlbumIDs:   TopNKeys(albumMap, 5),
	}
}
