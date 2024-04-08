package paginator

type Pagination struct {
	Limit  int
	Offset int
}

type PaginationResult struct {
	Limit  int
	Offset int
	Total  int
}
