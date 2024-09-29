# Pagination

## Overview

The pagination system in PixivFE is implemented through a combination of Go functions and [Jet HTML templates](https://github.com/CloudyKit/jet).

## Pagination data structure

The core of the pagination system is the `PaginationData` struct defined in `templateFunctions.go`. This structure encapsulates all the necessary information for rendering pagination controls:

```go
type PaginationData struct {
    CurrentPage int
    MaxPage     int
    Pages       []PageInfo
    HasPrevious bool
    HasNext     bool
    PreviousURL string
    NextURL     string
    FirstURL    string
    LastURL     string
    HasMaxPage  bool
    LastPage    int
}
```

The `PageInfo` struct is used to represent individual page links:

```go
type PageInfo struct {
    Number int
    URL    string
}
```

## Pagination logic

The `CreatePaginator` function in `templateFunctions.go` is responsible for generating the `PaginationData` structure. It takes four parameters:

1. `base`: The base URL for pagination links
2. `ending`: A string to append to the end of each pagination URL
3. `current_page`: The current page number
4. `max_page`: The maximum number of pages (`-1` if unknown)

The function calculates the range of pages to display, typically showing two pages on either side of the current page. It generates URLs for each page and populates the `PaginationData` structure with all necessary information for rendering the pagination controls.

## Template rendering

The pagination controls are rendered using the `pagination` block defined in `pagination.jet.html`. This template takes a `PaginationData` object and generates the HTML for the pagination controls.

The template includes logic for:

1. Rendering "Previous" and "Next" links
2. Displaying the first page link if not in the current range
3. Rendering ellipsis (...) when there are gaps in the page range
4. Highlighting the current page
5. Displaying the last page link if not in the current range

## Usage in views

To implement pagination in a view:

1. First, generate the pagination data using the `createPaginator` template function:
  ```html
  {{- url := unfinishedQuery(.QueriesC, "page") }}
  {{- paginationData := createPaginator(url, "#checkpoint", .Page, -1) }}
  ```

2. Then, render the pagination controls by yielding to the pagination block:
  ```html
  {{- yield pagination(data=paginationData) }}
  ```

## Customisation

The pagination system can be customised by modifying the `pagination.jet.html` template; changes to the HTML structure, CSS classes, and overall appearance of the pagination controls can be made without altering the underlying logic.

The `CreatePaginator` function can also be modified to change the number of pages displayed on either side of the current page by adjusting the `peek` constant.
