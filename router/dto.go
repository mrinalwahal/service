package router

// ListOptions represents the options for listing records.
type ListOptions struct {

	//	Number of records to skip.
	Skip int `query:"skip" validate:"gte=0"`

	//	Number of records to return.
	Limit int `query:"limit" validate:"gte=0,lte=100"`

	//	Order by field.
	OrderBy string `query:"orderBy" validate:"oneof=created_at updated_at title"`

	//	Order by direction.
	OrderDirection string `query:"orderDirection" validate:"oneof=asc desc"`

	//	Title of the record.
	Title string `query:"name"`
}

// UpdateOptions represents the options for updating a record.
type UpdateOptions struct {

	//	Title of the record.
	Title string `json:"title"`
}
