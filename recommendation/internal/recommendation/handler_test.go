package recommendation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler_GetRecommendation(t *testing.T) {
	ag := NewAvailabilityGetterMock(nil, nil)
	svc, _ := NewService(ag)

	sut, _ := NewHandler(*svc)

	timeFormat := "2006-01-02"
	now := time.Now()
	tripStart := now.Format(timeFormat)
	tripEnd := now.Add(time.Hour * 24).Format(timeFormat)
	location := "UK"
	budget := 50
	want := GetRecommendationResponse{
		HotelName: "test",
		TotalCost: struct {
			Cost     int64  `json:"cost"`
			Currency string `json:"currency"`
		}{
			Cost:     50,
			Currency: "USD",
		},
	}
	uri := fmt.Sprintf("location=%s&from=%s&to=%s&budget=%d", location, tripStart, tripEnd, budget)

	req := httptest.NewRequest(http.MethodGet, "/recommendation?"+uri, nil)
	w := httptest.NewRecorder()

	sut.GetRecommendation(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	rec := GetRecommendationResponse{}

	err = json.Unmarshal(data, &rec)
	if err != nil {
		t.Errorf("expected error to be nil unmarshalling response got %v", err)
	}
	assert.Equal(t, want, rec)
}

func TestHandler_GetRecommendation_Errors(t *testing.T) {
	ag := NewAvailabilityGetterMock(nil, nil)
	svc, _ := NewService(ag)

	sut, _ := NewHandler(*svc)

	timeFormat := "2006-01-02"
	now := time.Now()
	tripStart := now.Format(timeFormat)
	tripEnd := now.Add(time.Hour * 24).Format(timeFormat)
	location := "UK"
	budget := "50"

	testCases := []struct {
		desc       string
		start      string
		end        string
		location   string
		budget     string
		statusCode int
	}{
		{
			desc:       "start cannot be empty",
			start:      "",
			end:        tripEnd,
			location:   location,
			budget:     budget,
			statusCode: http.StatusBadRequest,
		},
		{
			desc:       "end cannot be empty",
			start:      tripStart,
			end:        "",
			location:   location,
			budget:     budget,
			statusCode: http.StatusBadRequest,
		},
		{
			desc:       "location cannot be empty",
			start:      tripStart,
			end:        tripEnd,
			location:   "",
			budget:     budget,
			statusCode: http.StatusBadRequest,
		},
		{
			desc:       "budget cannot be empty",
			start:      tripStart,
			end:        tripEnd,
			location:   location,
			budget:     "",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			uri := fmt.Sprintf("location=%s&from=%s&to=%s&budget=%s", tC.location, tC.start, tC.end, tC.budget)

			req := httptest.NewRequest(http.MethodGet, "/recommendation?"+uri, nil)
			w := httptest.NewRecorder()

			sut.GetRecommendation(w, req)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tC.statusCode, res.StatusCode)
		})
	}
}
