package navitia

import (
	"context"
	"fmt"
	"github.com/aabizri/navitia/types"
	"net/url"
	"strconv"
	"time"
)

// JourneyResults countains the results of a Journey request
//
// Warning: types.Journey.From / types.Journey.To aren't garanteed to be filled.
// Based on very basic inspection, it seems they aren't filled when there are sections...
type JourneyResults struct {
	Journeys []types.Journey

	Logging

	session *Session
}

// String satisfies stringer, pretty-prints JourneyResults
func (jr JourneyResults) String() string {
	var msg string
	for i, journey := range jr.Journeys {
		if i != 0 {
			msg += "\n"
		}
		msg += fmt.Sprintf("Journey #%d: %s\n", i, journey.String())
	}
	return msg
}

// JourneyRequest countain the parameters needed to make a Journey request
type JourneyRequest struct {
	// There must be at least one From or To parameter defined
	// When used with just one of them, the resulting Journey won't have a populated Sections field.
	From types.QueryEscaper
	To   types.QueryEscaper

	// When do you want to depart ? Or is DateIsArrival when do you want to arrive at your destination.
	Date          time.Time
	DateIsArrival bool

	// The traveler's type
	Traveler types.TravelerType

	// Define the freshness of data to use to compute journeys
	Freshness types.DataFreshness

	// Forbidden public transport objects
	Forbidden []types.ID

	// Allowed public transport objects
	// Note: This counstraint intersects with Forbidden
	Allowed []types.ID

	// Force the first section mode if it isn't a public transport mode
	// Note: The parameter is inclusive, not exclusive. As such if you want to forbid a mode you have to include all modes except that one.
	FirstSectionModes []types.Mode

	// Same, but for the last section
	LastSectionModes []types.Mode

	// MaxDurationToPT is the maximum allowed duration to reach the public transport.
	// Use this to limit the walking/biking part.
	MaxDurationToPT time.Duration

	// These four following parameters set the speed of each mode (Walking, Bike, BSS & car)
	// In meters per second
	WalkingSpeed   float64
	BikeSpeed      float64
	BikeShareSpeed float64
	CarSpeed       float64

	// Minimum and maximum amounts of journeys suggested
	MinJourneys uint
	MaxJourneys uint

	// Count fixes the amount of journeys to be returned, overriding minimum & maximum amount
	// Note: if Count=0 then it isn't taken into account
	Count uint

	// Maximum number of transfers in each journey
	MaxTransfers uint

	// Maximum duration of a trip
	MaxDuration time.Duration // To seconds

	// Wheelchair restricts the answer to accessible public transports
	Wheelchair bool
}

// toURL formats a journey request to url
// Should be refactored using a switch statement
func (req JourneyRequest) toURL() (url.Values, error) {
	params := url.Values{}

	// Define a few useful functions
	addUint := func(key string, amount uint64) {
		if amount != 0 {
			str := strconv.FormatUint(amount, 10)
			params.Add(key, str)
		}
	}
	addInt := func(key string, amount int64) {
		if amount != 0 {
			str := strconv.FormatInt(amount, 10)
			params.Add(key, str)
		}
	}
	addString := func(key string, str string) {
		if str != "" {
			params.Add(key, str)
		}
	}
	addIDSlice := func(key string, ids []types.ID) {
		if len(ids) != 0 {
			for _, id := range ids {
				params.Add(key, id.QueryEscape())
			}
		}
	}
	addModes := func(key string, modes []types.Mode) {
		if len(modes) != 0 {
			for _, mode := range modes {
				params.Add(key, string(mode))
			}
		}
	}
	addFloat := func(key string, amount float64) {
		if amount != 0 {
			speedStr := strconv.FormatFloat(amount, 'f', 3, 64)
			params.Add(key, speedStr)
		}
	}

	// Encode the from and to
	if from := req.From; from != nil {
		params.Add("from", from.QueryEscape())
	}
	if to := req.To; to != nil {
		params.Add("to", to.QueryEscape())
	}

	if datetime := req.Date; !datetime.IsZero() {
		str := datetime.Format(types.DateTimeFormat)
		params.Add("datetime", str)
		if req.DateIsArrival {
			params.Add("datetime_represents", "arrival")
		}
	}

	addString("traveler_type", string(req.Traveler))

	addString("data_freshness", string(req.Freshness))

	addIDSlice("forbidden_uris[]", req.Forbidden)

	addIDSlice("allowed_id[]", req.Allowed)

	addModes("first_section_mode[]", req.FirstSectionModes)

	addModes("last_section_mode[]", req.LastSectionModes)

	// max_duration_to_pt
	addInt("max_duration_to_pt", int64(req.MaxDurationToPT/time.Second))

	// walking_speed, bike_speed, bss_speed & car_speed
	addFloat("walking_speed", req.WalkingSpeed)
	addFloat("bike_speed", req.BikeSpeed)
	addFloat("bss_speed", req.BikeShareSpeed)
	addFloat("car_speed", req.CarSpeed)

	// If count is defined don't bother with the minimimal and maximum amount of items to return
	if count := req.Count; count != 0 {
		addUint("count", uint64(count))
	} else {
		addUint("min_nb_journeys", uint64(req.MinJourneys))
		addUint("max_nb_journeys", uint64(req.MaxJourneys))
	}

	// max_nb_transfers
	addUint("max_nb_transfers", uint64(req.MaxTransfers))

	// max_duration
	addInt("max_duration", int64(req.MaxDuration/time.Second))

	// wheelchair
	if req.Wheelchair {
		params.Add("wheelchair", "true")
	}

	return params, nil
}

// journeys is the internal function used by Journeys functions
func (s *Session) journeys(ctx context.Context, url string, params JourneyRequest) (*JourneyResults, error) {
	var results = &JourneyResults{session: s}
	err := s.request(ctx, url, params, results)
	return results, err
}

const journeysEndpoint string = "journeys"

// Journeys computes a list of journeys according to the parameters given
func (s *Session) Journeys(ctx context.Context, params JourneyRequest) (*JourneyResults, error) {
	// Create the URL
	url := s.APIURL + "/" + journeysEndpoint

	// Call
	return s.journeys(ctx, url, params)
}

// JourneysR computes a list of journeys according to the parameters given and in a specific region
func (s *Session) JourneysR(ctx context.Context, params JourneyRequest, regionID string) (*JourneyResults, error) {
	// Create the URL
	url := s.APIURL + "/coverage/" + regionID + "/" + journeysEndpoint

	// Call
	return s.journeys(ctx, url, params)
}
