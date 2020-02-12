package backend

import (
	"context"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/google/go-cmp/cmp"
	"github.com/sourcegraph/sourcegraph/internal/db/dbtesting"
)

func TestUpdateServiceVersion(t *testing.T) {
	dbtesting.SetupGlobalTestDB(t)

	ctx := context.Background()
	for _, tc := range []struct {
		version string
		err     error
	}{
		{"0.0.0", nil},
		{"0.0.1", nil},
		{"0.1.0", nil},
		{"0.2.0", nil},
		{"1.0.0", nil},
		{"1.2.0", &UpgradeError{
			Service:  "service",
			Previous: semver.MustParse("1.0.0"),
			Latest:   semver.MustParse("1.2.0"),
		}},
		{"2.1.0", &UpgradeError{
			Service:  "service",
			Previous: semver.MustParse("1.0.0"),
			Latest:   semver.MustParse("2.1.0"),
		}},
	} {
		have := UpdateServiceVersion(ctx, "service", tc.version)
		want := tc.err

		if diff := cmp.Diff(have, want); diff != "" {
			t.Fatal(diff)
		}

		t.Logf("version = %q", tc.version)
	}
}
func TestIsValidUpgrade(t *testing.T) {
	for _, tc := range []struct {
		name     string
		previous string
		latest   string
		want     bool
	}{{
		name:     "no versions",
		previous: "",
		latest:   "",
		want:     true,
	}, {
		name:     "no previous version",
		previous: "",
		latest:   "v3.13.0",
		want:     true,
	}, {
		name:     "same version",
		previous: "v3.13.0",
		latest:   "v3.13.0",
		want:     true,
	}, {
		name:     "one minor version up",
		previous: "v3.12.4",
		latest:   "v3.13.1",
		want:     true,
	}, {
		name:     "one major version up",
		previous: "v3.13.1",
		latest:   "v4.0.0",
		want:     true,
	}, {
		name:     "more than one minor version up",
		previous: "v3.9.4",
		latest:   "v3.11.0",
		want:     false,
	}, {
		name:     "major jump",
		previous: "v3.9.4",
		latest:   "v4.1.0",
		want:     false,
	},
	} {
		t.Run(tc.name, func(t *testing.T) {
			previous, _ := semver.NewVersion(tc.previous)
			latest, _ := semver.NewVersion(tc.latest)

			if got := IsValidUpgrade(previous, latest); got != tc.want {
				t.Errorf(
					"IsValidUpgrade(previous: %s, latest: %s) = %t, want %t",
					tc.latest,
					tc.latest,
					got,
					tc.want,
				)
			}
		})
	}
}
