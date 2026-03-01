package request

type PathID struct {
	ID string `uri:"id" binding:"required,len=36"`
}
