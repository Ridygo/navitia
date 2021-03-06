package navitia

import (
	"context"
	"reflect"
	"testing"
)

func Test_Regions(t *testing.T) {
	if *apiKey == "" {
		t.Skip(skipNoKey)
	}

	ctx := context.Background()
	req := RegionRequest{}

	// Run the query with GeoJSON
	t.Run("with_geojson", func(t *testing.T) {
		req := req
		req.Geo = true
		res, err := testSession.Regions(ctx, req)
		if err != nil {
			t.Fatalf("error in Regions: %v\n\tParameters: %#v\n\tReceived: %#v", err, req, res)
		}
	})

	// Run the query without GeoJSON
	t.Run("without_geojson", func(t *testing.T) {
		req := req
		res, err := testSession.Regions(ctx, req)
		if err != nil {
			t.Fatalf("error in Regions: %v\n\tParameters: %#v\n\tReceived: %#v", err, req, res)
		}
	})
}

// Test_RegionResults_Unmarshal tests unmarshalling for RegionResults.
// As the unmarshalling is done by encoding/json, this allows us to check that the input can be reliably unmarshalled into the structure we have for it.
//
// This launches both a "correct" and "incorrect" subtest, allowing us to test both cases.
// 	If we expect no errors but we get one, the test fails
//	If we expect an error but we don't get one, the test fails
func Test_RegionResults_Unmarshal(t *testing.T) {
	testUnmarshal(t, testData["coverage"], reflect.TypeOf(RegionResults{}))
}
