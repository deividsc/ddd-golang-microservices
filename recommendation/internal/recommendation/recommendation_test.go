package recommendation

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

type AvailabilityGetterMock struct {
	opts []Option
	err  error
}

func (m *AvailabilityGetterMock) GetAvailability(ctx context.Context, tripStart, tripEnd time.Time, location string) ([]Option, error) {
	return m.opts, m.err
}

func NewAvailabilityGetterMock(opts []Option, err error) *AvailabilityGetterMock {
	if opts == nil {
		opts = []Option{
			{
				HotelName:     "test",
				Location:      "UK",
				PricePerNight: *money.New(50, "USD"),
			},
			{
				HotelName:     "test2",
				Location:      "UK",
				PricePerNight: *money.New(500, "USD"),
			},
			{
				HotelName:     "test3",
				Location:      "UK",
				PricePerNight: *money.New(100, "USD"),
			},
		}
	}
	return &AvailabilityGetterMock{
		opts: opts,
		err:  err,
	}
}

func TestServiceGet(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour * 24)
	budget := money.New(100, "USD")
	testCases := []struct {
		desc            string
		tripStart       time.Time
		tripEnd         time.Time
		location        string
		budget          *money.Money
		errAvailability error
		want            struct {
			err            error
			recommendation Recommendation
		}
	}{
		{
			desc:      "Recomendation found",
			tripStart: start,
			tripEnd:   end,
			location:  "UK",
			budget:    budget,
			want: struct {
				err            error
				recommendation Recommendation
			}{
				recommendation: Recommendation{
					TripStart: start,
					TripEnd:   end,
					HotelName: "test",
					Location:  "UK",
					TripPrice: *money.New(50, "USD"),
				},
			},
		},
		{
			desc:     "TripStart cannot be empty",
			tripEnd:  end,
			location: "UK",
			budget:   budget,
			want: struct {
				err            error
				recommendation Recommendation
			}{
				err: errors.New("trip start cannot be empty"),
			},
		},
		{
			desc:      "TripEnd cannot be empty",
			tripStart: start,
			location:  "UK",
			budget:    budget,
			want: struct {
				err            error
				recommendation Recommendation
			}{
				err: errors.New("trip end cannot be empty"),
			},
		},
		{
			desc:      "location cannot be empty",
			tripStart: start,
			tripEnd:   end,
			location:  "",
			budget:    budget,
			want: struct {
				err            error
				recommendation Recommendation
			}{
				err: errors.New("location cannot be empty"),
			},
		},
		{
			desc:            "error getting availability",
			tripStart:       start,
			tripEnd:         end,
			location:        "UK",
			budget:          budget,
			errAvailability: errors.New("test error"),
			want: struct {
				err            error
				recommendation Recommendation
			}{
				err: errors.New("error getting availability: test error"),
			},
		},
		{
			desc:      "no trips within budget",
			tripStart: start,
			tripEnd:   end,
			location:  "UK",
			budget:    money.New(10, "USD"),
			want: struct {
				err            error
				recommendation Recommendation
			}{
				err: errors.New("no trips within budget"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ag := NewAvailabilityGetterMock(nil, tC.errAvailability)
			sut, _ := NewService(ag)

			res, err := sut.Get(context.TODO(), tC.tripStart, tC.tripEnd, tC.location, tC.budget)

			assert.Equal(t, fmt.Sprint(tC.want.err), fmt.Sprint(err))
			if res != nil {
				assert.Equal(t, tC.want.recommendation, *res)
			}

		})
	}
}
