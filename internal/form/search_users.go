package form

// SearchUsers represents a user search form.
type SearchUsers struct {
	Query   string `form:"q"`
	User    string `form:"user"`
	Name    string `form:"name"`
	Email   string `form:"email"`
	All     bool   `form:"all"`
	Deleted bool   `form:"deleted"`
	Count   int    `form:"count" binding:"required" serialize:"-"`
	Offset  int    `form:"offset" serialize:"-"`
	Order   string `form:"order" serialize:"-"`
	Reverse bool   `form:"reverse" serialize:"-"`
}

func (f *SearchUsers) GetQuery() string {
	return f.Query
}

func (f *SearchUsers) SetQuery(q string) {
	f.Query = q
}

func (f *SearchUsers) ParseQueryString() error {
	return ParseQueryString(f)
}
