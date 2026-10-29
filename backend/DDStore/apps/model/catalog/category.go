package catalog

// Category is the copy of a category embedded in a catalog product.
type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
