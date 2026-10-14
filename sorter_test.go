package finder

import "testing"

func TestSortBuilderAppendsStableIDSort(t *testing.T) {
	tests := []struct {
		name     string
		sorts    string
		columns  []string
		idColumn *string
		want     string
	}{
		{
			name:    "default created at sort",
			columns: []string{"id", "created_at"},
			want:    "recipes.created_at DESC,recipes.id DESC",
		},
		{
			name:    "requested non unique sort",
			sorts:   "name",
			columns: []string{"id", "name"},
			want:    "recipes.name ASC,recipes.id DESC",
		},
		{
			name:    "requested descending non unique sort",
			sorts:   "-name",
			columns: []string{"id", "name"},
			want:    "recipes.name DESC,recipes.id DESC",
		},
		{
			name:    "requested id is not duplicated",
			sorts:   "id",
			columns: []string{"id", "name"},
			want:    "recipes.id ASC",
		},
		{
			name:    "invalid sort falls back to id",
			sorts:   "unknown",
			columns: []string{"id", "name"},
			want:    "recipes.id DESC",
		},
		{
			name:     "custom id column",
			sorts:    "name",
			columns:  []string{"name"},
			idColumn: stringPointer("recipe_id"),
			want:     "recipes.name ASC,recipes.recipe_id DESC",
		},
		{
			name:    "model without an id keeps existing sort",
			sorts:   "name",
			columns: []string{"name"},
			want:    "recipes.name ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortBuilder("recipes", tt.sorts, tt.idColumn, &tt.columns)
			if got != tt.want {
				t.Fatalf("sortBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAppendStableSortForOverride(t *testing.T) {
	columns := []string{"id", "created_at"}

	tests := []struct {
		name     string
		override string
		want     string
	}{
		{
			name:     "appends id to expression",
			override: "COALESCE(o.total_spent,0) DESC",
			want:     "COALESCE(o.total_spent,0) DESC,customers.id DESC",
		},
		{
			name:     "does not duplicate descending id",
			override: "customers.created_at DESC, customers.id DESC",
			want:     "customers.created_at DESC, customers.id DESC",
		},
		{
			name:     "does not override requested id direction",
			override: "customers.id ASC",
			want:     "customers.id ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendStableSort(tt.override, "customers", nil, &columns)
			if got != tt.want {
				t.Fatalf("appendStableSort() = %q, want %q", got, tt.want)
			}
		})
	}
}

func stringPointer(value string) *string {
	return &value
}
