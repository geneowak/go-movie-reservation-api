package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geneowak/go-expense-tracker/internal/database"
	"github.com/geneowak/go-expense-tracker/internal/testdb"
	"github.com/google/go-cmp/cmp"
)

type mockQuerier struct {
	database.Querier
	createMovieFn func(ctx context.Context, arg database.CreateMovieParams) (database.Movie, error)
}

func (m *mockQuerier) CreateMovie(ctx context.Context, arg database.CreateMovieParams) (database.Movie, error) {
	return m.createMovieFn(ctx, arg)
}

func TestHandleCreateMovie(t *testing.T) {
	db := testdb.OpenAndMigrate(t)

	validBody := `{
		"name": "Inception",
		"description": "A mind-bending thriller about dreams within dreams",
		"duration_in_mins": 148,
		"trailer_url": "https://example.com/inception-trailer",
		"genre": ["Sci-Fi", "Action"],
		"pg_rating": "PG-13",
		"experience_types": ["IMAX", "3D"]
	}`

	tests := []struct {
		name       string
		querier    database.Querier // nil = real DB with fresh tx
		reqBody    string
		wantStatus int
		wantBody   map[string]any
		wantErrors map[string][]string
	}{
		{
			name:       "Valid request creates movie and returns 201",
			reqBody:    validBody,
			wantStatus: http.StatusCreated,
		},
		{
			name: "Malformed JSON returns 500",
			querier: &mockQuerier{createMovieFn: func(ctx context.Context, arg database.CreateMovieParams) (database.Movie, error) {
				return database.Movie{}, nil
			}},
			reqBody:    `{`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing required fields returns 422 with errors",
			querier: &mockQuerier{createMovieFn: func(ctx context.Context, arg database.CreateMovieParams) (database.Movie, error) {
				return database.Movie{}, nil
			}},
			reqBody: `{
				"name": "ab",
				"description": "short",
				"duration_in_mins": 0,
				"trailer_url": "not-a-url",
				"genre": [],
				"pg_rating": "",
				"experience_types": []
			}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErrors: map[string][]string{
				"name":             {"The name field must be at least 3 characters long."},
				"description":      {"The description field must be at least 10 characters long."},
				"duration_in_mins": {"The duration_in_mins field is required"},
				"experience_types": {"The experience_types field must contain atleast 1 items."},
				"genre":            {"The genre field must contain atleast 1 items."},
				"trailer_url":      {"The trailer_url must be a valid URL."},
				"pg_rating":        {"The pg_rating field is required"},
			},
		},
		{
			name: "DB error returns 500",
			querier: &mockQuerier{
				createMovieFn: func(ctx context.Context, arg database.CreateMovieParams) (database.Movie, error) {
					return database.Movie{}, errors.New("Database connection lost")
				},
			},
			reqBody:    validBody,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var querier database.Querier
			if tt.querier != nil {
				querier = tt.querier
			} else {
				if db == nil {
					t.Skip("TEST_DB_URL not set, skipping real DB test")
				}
				scope := testdb.Begin(t, db)
				querier = scope.Queries
			}

			cfg := newTestApiConfig(querier)
			req := httptest.NewRequest(http.MethodPost, "/api/movies", strings.NewReader(tt.reqBody))
			rec := httptest.NewRecorder()
			cfg.handleCreateMovie(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			// test if the data was saved to the DB
			if tt.wantStatus == http.StatusCreated && tt.querier == nil {
				var created database.Movie
				if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}
				// Verify the row exists in the database (within the same tx)
				movie, err := querier.GetMovieById(context.Background(), created.ID)
				if err != nil {
					t.Fatalf("failed to find movie: %v", err)
				}

				if movie.Name != "Inception" {
					t.Errorf("Movie title does not match: want %s, got %s", "Inception", movie.Name)
				}
			}

			var respBody map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&respBody); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if tt.wantErrors != nil {
				errorsMap, ok := respBody["errors"].(map[string]any)
				if !ok {
					t.Fatalf("response missing 'errors' field: %v", respBody)
				}
				// convert map[string]any -> map[string][]string for comparison
				got := make(map[string][]string)
				for k, v := range errorsMap {
					if arr, ok := v.([]any); ok {
						for _, item := range arr {
							got[k] = append(got[k], item.(string))
						}
					}
				}

				if diff := cmp.Diff(tt.wantErrors, got); diff != "" {
					t.Errorf("Errors mismatch (-want +got):\n%s", diff)
				}
			}

			if tt.wantBody != nil {
				if diff := cmp.Diff(tt.wantBody, respBody); diff != "" {
					t.Errorf("body mismatch (-want +got)\n%s", diff)
				}
			}
		})
	}
}
