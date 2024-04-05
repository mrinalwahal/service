package service

type CreateOptions struct {

	//	Title of the todo.
	Title string
}

type ListOptions struct {

	//	Title of the todo.
	Title string
	//	Skip for pagination.
	Skip int
	//	Limit for pagination.
	Limit int
	//	Order by field.
	OrderBy string
	//	Order by direction.
	OrderDirection string
}

type UpdateOptions struct {

	//	Title of the todo.
	Title string
}
