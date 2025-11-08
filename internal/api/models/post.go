package models

// ----------------- Structures -----------------

type Post struct {
	PostId      string `json:"post_id"`
	UserId      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
	Tags        []Tag  `json:"tags"`
}

type Board struct {
	BoardId     string `json:"board_id"`
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `json:"posts"`
}

type Tag struct {
	TagId string `json:"tag_id"`
	Name  string `json:"name"`
}

// --------------- Request Models ---------------

type CreatePostRequest struct {
	PostId      string `json:"post_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdatePostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetPostsRequest struct {
	PostsIds []string `json:"posts_ids" binding:"required"`
}

type SearchPostsRequest struct {
	Query   string   `json:"query"`
	UserIds []string `json:"user_ids"`
	TagIds  []string `json:"tag_ids"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

type CreateBoardRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetBoardsRequest struct {
	BoardsIds []string `json:"boards_ids" binding:"required"`
}

type SearchBoardsRequest struct {
	Query   string   `json:"query"`
	UserIds []string `json:"user_ids"`
	PostIds []string `json:"post_ids"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type GetTagRequest struct {
	TagId string `json:"tag_id" binding:"required"`
}

type SearchTagsRequest struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type PinPostToBoardRequest struct {
	PostId  string `json:"post_id" binding:"required"`
	BoardId string `json:"board_id" binding:"required"`
}

type UnpinPostFromBoardRequest struct {
	PostId  string `json:"post_id" binding:"required"`
	BoardId string `json:"board_id" binding:"required"`
}

type AddTagToPostRequest struct {
	PostId string `json:"post_id" binding:"required"`
	TagId  string `json:"tag_id" binding:"required"`
}

type RemoveTagFromPostRequest struct {
	PostId string `json:"post_id" binding:"required"`
	TagId  string `json:"tag_id" binding:"required"`
}

// --------------- Response Models ---------------

type CreatePostResponse struct {
	PostId  string `json:"post_id"`
	Message string `json:"message"`
}

type UpdatePostResponse struct {
	Message string `json:"message"`
}

type GetPostResponse struct {
	Post Post `json:"post"`
}

type GetPostsResponse struct {
	Posts []Post `json:"posts"`
}

type SearchPostsResponse struct {
	Posts []Post `json:"posts"`
}

type DeletePostResponse struct {
	Message string `json:"message"`
}

type CreateBoardResponse struct {
	BoardId string `json:"board_id"`
}

type UpdateBoardResponse struct {
	Message string `json:"message"`
}

type GetBoardResponse struct {
	Board Board `json:"board"`
}

type GetBoardsResponse struct {
	Boards []Board `json:"boards"`
}

type SearchBoardsResponse struct {
	Boards []Board `json:"boards"`
}

type DeleteBoardResponse struct {
	Message string `json:"message"`
}

type CreateTagResponse struct {
	TagId string `json:"tag_id"`
}

type GetTagResponse struct {
	Tag Tag `json:"tag"`
}

type SearchTagsResponse struct {
	Tags []Tag `json:"tags"`
}

type PinPostToBoardResponse struct {
	Message string `json:"message"`
}

type UnpinPostFromBoardResponse struct {
	Message string `json:"message"`
}

type AddTagToPostResponse struct {
	Message string `json:"message"`
}

type RemoveTagFromPostResponse struct {
	Message string `json:"message"`
}
