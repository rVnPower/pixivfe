package template

import (
	"reflect"
	"testing"
)

func TestCreatePaginator(t *testing.T) {
	tests := []struct {
		name           string
		base           string
		ending         string
		currentPage    int
		maxPage        int
		pageMargin     int
		dropdownOffset int
		expectError    bool
		expectedData   PaginationData
	}{
		{
			name:           "Normal case",
			base:           "/ranking?content=all&date=20240101&mode=daily&page=",
			ending:         "#checkpoint",
			currentPage:    3,
			maxPage:        10,
			pageMargin:     1,
			dropdownOffset: 2,
			expectError:    false,
			expectedData: PaginationData{
				CurrentPage: 3,
				MaxPage:     10,
				Pages: []PageInfo{
					{Number: 2, URL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint"},
					{Number: 3, URL: "/ranking?content=all&date=20240101&mode=daily&page=3#checkpoint"},
					{Number: 4, URL: "/ranking?content=all&date=20240101&mode=daily&page=4#checkpoint"},
				},
				HasPrevious: true,
				HasNext:     true,
				PreviousURL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint",
				NextURL:     "/ranking?content=all&date=20240101&mode=daily&page=4#checkpoint",
				FirstURL:    "/ranking?content=all&date=20240101&mode=daily&page=1#checkpoint",
				LastURL:     "/ranking?content=all&date=20240101&mode=daily&page=10#checkpoint",
				HasMaxPage:  true,
				LastPage:    4,
				DropdownPages: []PageInfo{
					{Number: 1, URL: "/ranking?content=all&date=20240101&mode=daily&page=1#checkpoint"},
					{Number: 2, URL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint"},
					{Number: 3, URL: "/ranking?content=all&date=20240101&mode=daily&page=3#checkpoint"},
					{Number: 4, URL: "/ranking?content=all&date=20240101&mode=daily&page=4#checkpoint"},
					{Number: 5, URL: "/ranking?content=all&date=20240101&mode=daily&page=5#checkpoint"},
				},
			},
		},
		{
			name:           "Unknown max page",
			base:           "/ranking?content=all&date=20240101&mode=daily&page=",
			ending:         "#checkpoint",
			currentPage:    3,
			maxPage:        -1,
			pageMargin:     1,
			dropdownOffset: 2,
			expectError:    false,
			expectedData: PaginationData{
				CurrentPage: 3,
				MaxPage:     -1,
				Pages: []PageInfo{
					{Number: 2, URL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint"},
					{Number: 3, URL: "/ranking?content=all&date=20240101&mode=daily&page=3#checkpoint"},
					{Number: 4, URL: "/ranking?content=all&date=20240101&mode=daily&page=4#checkpoint"},
				},
				HasPrevious: true,
				HasNext:     true,
				PreviousURL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint",
				NextURL:     "/ranking?content=all&date=20240101&mode=daily&page=4#checkpoint",
				FirstURL:    "/ranking?content=all&date=20240101&mode=daily&page=1#checkpoint",
				LastURL:     "/ranking?content=all&date=20240101&mode=daily&page=-1#checkpoint",
				HasMaxPage:  false,
				LastPage:    4,
				DropdownPages: []PageInfo{
					{Number: 1, URL: "/ranking?content=all&date=20240101&mode=daily&page=1#checkpoint"},
					{Number: 2, URL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint"},
					{Number: 3, URL: "/ranking?content=all&date=20240101&mode=daily&page=3#checkpoint"},
					{Number: 4, URL: "/ranking?content=all&date=20240101&mode=daily&page=4#checkpoint"},
					{Number: 5, URL: "/ranking?content=all&date=20240101&mode=daily&page=5#checkpoint"},
				},
			},
		},
		{
			name:           "First page",
			base:           "/ranking?content=all&date=20240101&mode=daily&page=",
			ending:         "#checkpoint",
			currentPage:    1,
			maxPage:        10,
			pageMargin:     1,
			dropdownOffset: 2,
			expectError:    false,
			expectedData: PaginationData{
				CurrentPage: 1,
				MaxPage:     10,
				Pages: []PageInfo{
					{Number: 1, URL: "/ranking?content=all&date=20240101&mode=daily&page=1#checkpoint"},
					{Number: 2, URL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint"},
				},
				HasPrevious: false,
				HasNext:     true,
				PreviousURL: "",
				NextURL:     "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint",
				FirstURL:    "/ranking?content=all&date=20240101&mode=daily&page=1#checkpoint",
				LastURL:     "/ranking?content=all&date=20240101&mode=daily&page=10#checkpoint",
				HasMaxPage:  true,
				LastPage:    2,
				DropdownPages: []PageInfo{
					{Number: 1, URL: "/ranking?content=all&date=20240101&mode=daily&page=1#checkpoint"},
					{Number: 2, URL: "/ranking?content=all&date=20240101&mode=daily&page=2#checkpoint"},
					{Number: 3, URL: "/ranking?content=all&date=20240101&mode=daily&page=3#checkpoint"},
				},
			},
		},
		{
			name:           "Invalid current page",
			base:           "/ranking?content=all&date=20240101&mode=daily&page=",
			ending:         "#checkpoint",
			currentPage:    0,
			maxPage:        10,
			pageMargin:     1,
			dropdownOffset: 2,
			expectError:    true,
		},
		{
			name:           "Invalid page margin",
			base:           "/ranking?content=all&date=20240101&mode=daily&page=",
			ending:         "#checkpoint",
			currentPage:    1,
			maxPage:        10,
			pageMargin:     -1,
			dropdownOffset: 2,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := CreatePaginator(tt.base, tt.ending, tt.currentPage, tt.maxPage, tt.pageMargin, tt.dropdownOffset)

			if tt.expectError {
				if err == nil {
					t.Errorf("CreatePaginator() error = nil, expected an error")
				}
			} else {
				if err != nil {
					t.Errorf("CreatePaginator() unexpected error = %v", err)
				}

				if !reflect.DeepEqual(gotData, tt.expectedData) {
					t.Errorf("CreatePaginator() gotData = %v, want %v", gotData, tt.expectedData)
				}
			}
		})
	}
}
