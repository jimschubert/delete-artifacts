package app

import (
	"testing"
	"time"

	"github.com/google/go-github/v32/github"
)

// Helper function to create a pointer to an int64
func int64Ptr(i int64) *int64 {
	return &i
}

// Helper to create a test artifact
func createArtifact(name string, sizeInBytes int64, createdAt time.Time) *github.Artifact {
	ts := github.Timestamp{Time: createdAt}
	return &github.Artifact{
		Name:        &name,
		SizeInBytes: &sizeInBytes,
		CreatedAt:   &ts,
	}
}

func TestFilterArtifacts_MinBytes(t *testing.T) {
	tests := []struct {
		name         string
		minBytes     int64
		artifactSize int64
		wantMatch    bool
	}{
		{
			name:         "artifact size equals MinBytes - should match",
			minBytes:     100,
			artifactSize: 100,
			wantMatch:    true,
		},
		{
			name:         "artifact size greater than MinBytes - should match",
			minBytes:     100,
			artifactSize: 200,
			wantMatch:    true,
		},
		{
			name:         "artifact size less than MinBytes - should not match",
			minBytes:     100,
			artifactSize: 50,
			wantMatch:    false,
		},
		{
			name:         "MinBytes is zero - should match any size",
			minBytes:     0,
			artifactSize: 0,
			wantMatch:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				MinBytes: tt.minBytes,
			}
			artifact := createArtifact("test-artifact", tt.artifactSize, time.Now().Add(-1*time.Hour))
			result := app.filterArtifacts([]*github.Artifact{artifact})

			if tt.wantMatch && len(result) != 1 {
				t.Errorf("expected artifact to match, but it didn't")
			}
			if !tt.wantMatch && len(result) != 0 {
				t.Errorf("expected artifact to not match, but it did")
			}
		})
	}
}

func TestFilterArtifacts_MaxBytes(t *testing.T) {
	tests := []struct {
		name         string
		minBytes     int64
		maxBytes     *int64
		artifactSize int64
		wantMatch    bool
	}{
		{
			name:         "artifact size equals MaxBytes - should not match (size > MaxBytes check)",
			minBytes:     0,
			maxBytes:     int64Ptr(100),
			artifactSize: 100,
			wantMatch:    true, // size is NOT > maxBytes, so it passes
		},
		{
			name:         "artifact size greater than MaxBytes - should not match",
			minBytes:     0,
			maxBytes:     int64Ptr(100),
			artifactSize: 200,
			wantMatch:    false,
		},
		{
			name:         "artifact size less than MaxBytes - should match",
			minBytes:     0,
			maxBytes:     int64Ptr(100),
			artifactSize: 50,
			wantMatch:    true,
		},
		{
			name:         "MaxBytes is nil - should match any size (if MinBytes passes)",
			minBytes:     0,
			maxBytes:     nil,
			artifactSize: 1000000,
			wantMatch:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				MinBytes: tt.minBytes,
				MaxBytes: tt.maxBytes,
			}
			artifact := createArtifact("test-artifact", tt.artifactSize, time.Now().Add(-1*time.Hour))
			result := app.filterArtifacts([]*github.Artifact{artifact})

			if tt.wantMatch && len(result) != 1 {
				t.Errorf("expected artifact to match, but it didn't")
			}
			if !tt.wantMatch && len(result) != 0 {
				t.Errorf("expected artifact to not match, but it did")
			}
		})
	}
}

func TestFilterArtifacts_Name(t *testing.T) {
	tests := []struct {
		name         string
		filterName   string
		artifactName string
		wantMatch    bool
	}{
		{
			name:         "exact name match - should match",
			filterName:   "my-artifact",
			artifactName: "my-artifact",
			wantMatch:    true,
		},
		{
			name:         "name does not match - should not match",
			filterName:   "my-artifact",
			artifactName: "other-artifact",
			wantMatch:    false,
		},
		{
			name:         "partial name match - should not match (requires exact)",
			filterName:   "my-artifact",
			artifactName: "my-artifact-extra",
			wantMatch:    false,
		},
		{
			name:         "empty filter name - should match any artifact",
			filterName:   "",
			artifactName: "any-artifact",
			wantMatch:    true,
		},
		{
			name:         "case sensitive name - should not match",
			filterName:   "My-Artifact",
			artifactName: "my-artifact",
			wantMatch:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				MinBytes: 0,
				Name:     tt.filterName,
			}
			artifact := createArtifact(tt.artifactName, 100, time.Now().Add(-1*time.Hour))
			result := app.filterArtifacts([]*github.Artifact{artifact})

			if tt.wantMatch && len(result) != 1 {
				t.Errorf("expected artifact to match, but it didn't")
			}
			if !tt.wantMatch && len(result) != 0 {
				t.Errorf("expected artifact to not match, but it did")
			}
		})
	}
}

func TestFilterArtifacts_Pattern(t *testing.T) {
	tests := []struct {
		name         string
		pattern      string
		artifactName string
		wantMatch    bool
	}{
		{
			name:         "pattern matches suffix - should match",
			pattern:      "\\.bin$",
			artifactName: "artifact.bin",
			wantMatch:    true,
		},
		{
			name:         "pattern does not match - should not match",
			pattern:      "\\.bin$",
			artifactName: "artifact.txt",
			wantMatch:    false,
		},
		{
			name:         "pattern matches prefix - should match",
			pattern:      "^test-",
			artifactName: "test-artifact",
			wantMatch:    true,
		},
		{
			name:         "pattern matches substring - should match",
			pattern:      "artifact",
			artifactName: "my-artifact-name",
			wantMatch:    true,
		},
		{
			name:         "empty pattern - should match any artifact",
			pattern:      "",
			artifactName: "any-artifact",
			wantMatch:    true,
		},
		{
			name:         "invalid pattern - should not match",
			pattern:      "[invalid",
			artifactName: "artifact",
			wantMatch:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				MinBytes: 0,
				Pattern:  tt.pattern,
			}
			artifact := createArtifact(tt.artifactName, 100, time.Now().Add(-1*time.Hour))
			result := app.filterArtifacts([]*github.Artifact{artifact})

			if tt.wantMatch && len(result) != 1 {
				t.Errorf("expected artifact to match, but it didn't")
			}
			if !tt.wantMatch && len(result) != 0 {
				t.Errorf("expected artifact to not match, but it did")
			}
		})
	}
}

func TestFilterArtifacts_ActiveDuration(t *testing.T) {
	tests := []struct {
		name           string
		activeDuration string
		artifactAge    time.Duration
		wantMatch      bool
	}{
		{
			name:           "artifact older than active duration - should match",
			activeDuration: "1h",
			artifactAge:    2 * time.Hour,
			wantMatch:      true,
		},
		{
			name:           "artifact newer than active duration - should not match (still active)",
			activeDuration: "1h",
			artifactAge:    30 * time.Minute,
			wantMatch:      false,
		},
		{
			name:           "artifact exactly at active duration boundary - should match (time passes during execution)",
			activeDuration: "1h",
			artifactAge:    1*time.Hour + 1*time.Millisecond, // slightly older to guarantee match
			wantMatch:      true,
		},
		{
			name:           "empty active duration - should match any artifact",
			activeDuration: "",
			artifactAge:    1 * time.Minute,
			wantMatch:      true,
		},
		{
			name:           "invalid duration string - should not match",
			activeDuration: "invalid",
			artifactAge:    2 * time.Hour,
			wantMatch:      false,
		},
		{
			name:           "negative duration - should not match",
			activeDuration: "-1h",
			artifactAge:    2 * time.Hour,
			wantMatch:      false,
		},
		{
			name:           "zero duration - should not match",
			activeDuration: "0s",
			artifactAge:    2 * time.Hour,
			wantMatch:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				MinBytes:       0,
				ActiveDuration: tt.activeDuration,
			}
			createdAt := time.Now().Add(-tt.artifactAge)
			artifact := createArtifact("test-artifact", 100, createdAt)
			result := app.filterArtifacts([]*github.Artifact{artifact})

			if tt.wantMatch && len(result) != 1 {
				t.Errorf("expected artifact to match, but it didn't")
			}
			if !tt.wantMatch && len(result) != 0 {
				t.Errorf("expected artifact to not match, but it did")
			}
		})
	}
}

func TestFilterArtifacts_CombinedFilters(t *testing.T) {
	tests := []struct {
		name           string
		minBytes       int64
		maxBytes       *int64
		filterName     string
		pattern        string
		activeDuration string
		artifact       *github.Artifact
		wantMatch      bool
	}{
		{
			name:           "all filters pass - should match",
			minBytes:       50,
			maxBytes:       int64Ptr(200),
			filterName:     "test-artifact.bin",
			pattern:        "\\.bin$",
			activeDuration: "30m",
			artifact:       createArtifact("test-artifact.bin", 100, time.Now().Add(-1*time.Hour)),
			wantMatch:      true,
		},
		{
			name:           "MinBytes fails - should not match",
			minBytes:       200,
			maxBytes:       int64Ptr(500),
			filterName:     "test-artifact.bin",
			pattern:        "\\.bin$",
			activeDuration: "30m",
			artifact:       createArtifact("test-artifact.bin", 100, time.Now().Add(-1*time.Hour)),
			wantMatch:      false,
		},
		{
			name:           "MaxBytes fails - should not match",
			minBytes:       50,
			maxBytes:       int64Ptr(80),
			filterName:     "test-artifact.bin",
			pattern:        "\\.bin$",
			activeDuration: "30m",
			artifact:       createArtifact("test-artifact.bin", 100, time.Now().Add(-1*time.Hour)),
			wantMatch:      false,
		},
		{
			name:           "Name fails - should not match",
			minBytes:       50,
			maxBytes:       int64Ptr(200),
			filterName:     "other-artifact.bin",
			pattern:        "\\.bin$",
			activeDuration: "30m",
			artifact:       createArtifact("test-artifact.bin", 100, time.Now().Add(-1*time.Hour)),
			wantMatch:      false,
		},
		{
			name:           "Pattern fails - should not match",
			minBytes:       50,
			maxBytes:       int64Ptr(200),
			filterName:     "test-artifact.bin",
			pattern:        "\\.txt$",
			activeDuration: "30m",
			artifact:       createArtifact("test-artifact.bin", 100, time.Now().Add(-1*time.Hour)),
			wantMatch:      false,
		},
		{
			name:           "ActiveDuration fails (artifact too new) - should not match",
			minBytes:       50,
			maxBytes:       int64Ptr(200),
			filterName:     "test-artifact.bin",
			pattern:        "\\.bin$",
			activeDuration: "2h",
			artifact:       createArtifact("test-artifact.bin", 100, time.Now().Add(-1*time.Hour)),
			wantMatch:      false,
		},
		{
			name:           "only MinBytes filter - should match",
			minBytes:       50,
			maxBytes:       nil,
			filterName:     "",
			pattern:        "",
			activeDuration: "",
			artifact:       createArtifact("any-artifact", 100, time.Now().Add(-1*time.Minute)),
			wantMatch:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				MinBytes:       tt.minBytes,
				MaxBytes:       tt.maxBytes,
				Name:           tt.filterName,
				Pattern:        tt.pattern,
				ActiveDuration: tt.activeDuration,
			}
			result := app.filterArtifacts([]*github.Artifact{tt.artifact})

			if tt.wantMatch && len(result) != 1 {
				t.Errorf("expected artifact to match, but it didn't")
			}
			if !tt.wantMatch && len(result) != 0 {
				t.Errorf("expected artifact to not match, but it did")
			}
		})
	}
}

func TestFilterArtifacts_MultipleArtifacts(t *testing.T) {
	app := &App{
		MinBytes: 50,
		MaxBytes: int64Ptr(200),
		Pattern:  "\\.bin$",
	}

	artifacts := []*github.Artifact{
		createArtifact("artifact1.bin", 100, time.Now().Add(-1*time.Hour)), // matches
		createArtifact("artifact2.txt", 100, time.Now().Add(-1*time.Hour)), // fails pattern
		createArtifact("artifact3.bin", 10, time.Now().Add(-1*time.Hour)),  // fails MinBytes
		createArtifact("artifact4.bin", 300, time.Now().Add(-1*time.Hour)), // fails MaxBytes
		createArtifact("artifact5.bin", 150, time.Now().Add(-1*time.Hour)), // matches
	}

	result := app.filterArtifacts(artifacts)

	if len(result) != 2 {
		t.Errorf("expected 2 matching artifacts, got %d", len(result))
	}

	expectedNames := map[string]bool{"artifact1.bin": true, "artifact5.bin": true}
	for _, a := range result {
		if !expectedNames[a.GetName()] {
			t.Errorf("unexpected artifact in result: %s", a.GetName())
		}
	}
}

func TestFilterArtifacts_EmptyInput(t *testing.T) {
	app := &App{
		MinBytes: 0,
	}

	result := app.filterArtifacts([]*github.Artifact{})

	if len(result) != 0 {
		t.Errorf("expected 0 artifacts for empty input, got %d", len(result))
	}
}

func TestFilterArtifacts_NilInput(t *testing.T) {
	app := &App{
		MinBytes: 0,
	}

	result := app.filterArtifacts(nil)

	if len(result) != 0 {
		t.Errorf("expected 0 artifacts for nil input, got %d", len(result))
	}
}
