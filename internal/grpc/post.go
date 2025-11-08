package grpc

import (
	"context"
	"fmt"
	postProto "gateway/internal/proto/post"
	"gateway/pkg/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostClient struct {
	api postProto.PostServiceClient
	cc  *grpc.ClientConn
}

func NewPostClient(addr string, metrics *middleware.Metrics) (*PostClient, error) {
	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(metrics.GRPCClientMetricsInterceptor),
	)

	if err != nil {
		log.Error().Err(err).Msg("failed to connect to Post gRPC service")
		return nil, fmt.Errorf("post grpc client: dial: %w", err)
	}

	log.Debug().Msg("connected to Post gRPC service with address: " + addr)

	return &PostClient{
		api: postProto.NewPostServiceClient(cc),
		cc:  cc,
	}, nil
}

func (c *PostClient) CreatePost(ctx context.Context, postId, userId, title, description string) (string, error) {
	req := &postProto.CreatePostRequest{
		PostId:      postId,
		UserId:      userId,
		Title:       title,
		Description: description,
	}

	res, err := c.api.CreatePost(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create post via gRPC")
		return "", fmt.Errorf("create post grpc: %w", err)
	}

	return res.PostId, nil
}

func (c *PostClient) UpdatePost(ctx context.Context, postId, userId, title, description string) error {
	req := &postProto.UpdatePostRequest{
		PostId:      postId,
		UserId:      userId,
		Title:       title,
		Description: description,
	}

	_, err := c.api.UpdatePost(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to update post via gRPC")
		return fmt.Errorf("update post grpc: %w", err)
	}

	return nil
}

func (c *PostClient) GetPost(ctx context.Context, postId string) (string, string, string, string, string, [][]string, error) {
	req := &postProto.GetPostRequest{
		PostId: postId,
	}

	res, err := c.api.GetPost(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get post via gRPC")
		return "", "", "", "", "", nil, fmt.Errorf("get post grpc: %w", err)
	}

	tags := make([][]string, len(res.Post.Tags))
	for i, tag := range res.Post.Tags {
		tags[i] = []string{tag.TagId, tag.Name}
	}

	return res.Post.PostId, res.Post.UserId, res.Post.Title, res.Post.Description, res.Post.Extension, tags, nil
}

func (c *PostClient) GetPosts(ctx context.Context, postIds []string) ([]string, []string, []string, []string, []string, [][][]string, error) {
	req := &postProto.GetPostsRequest{
		PostIds: postIds,
	}

	res, err := c.api.GetPosts(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get posts via gRPC")
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("get posts grpc: %w", err)
	}

	var postIdsRes, userIds, titles, descriptions, extensions []string
	var tagIds [][][]string

	for _, post := range res.Posts {
		tags := make([][]string, len(post.Tags))
		for i, tag := range post.Tags {
			tags[i] = []string{tag.TagId, tag.Name}
		}

		postIdsRes = append(postIdsRes, post.PostId)
		userIds = append(userIds, post.UserId)
		titles = append(titles, post.Title)
		descriptions = append(descriptions, post.Description)
		extensions = append(extensions, post.Extension)
		tagIds = append(tagIds, tags)
	}

	return postIdsRes, userIds, titles, descriptions, extensions, tagIds, nil
}

func (c *PostClient) SearchPosts(ctx context.Context, query string, userIds, tagIds []string, limit, offset int) ([]string, []string, []string, []string, []string, [][][]string, error) {
	req := &postProto.SearchPostsRequest{
		Query:  query,
		UserId: userIds,
		TagId:  tagIds,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	res, err := c.api.SearchPosts(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to search posts via gRPC")
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("search posts grpc: %w", err)
	}

	var postIdsRes, userIdsRes, titles, descriptions, extensions []string
	var tagIdsRes [][][]string

	for _, post := range res.Posts {
		tags := make([][]string, len(post.Tags))
		for i, tag := range post.Tags {
			tags[i] = []string{tag.TagId, tag.Name}
		}

		postIdsRes = append(postIdsRes, post.PostId)
		userIdsRes = append(userIdsRes, post.UserId)
		titles = append(titles, post.Title)
		descriptions = append(descriptions, post.Description)
		extensions = append(extensions, post.Extension)
		tagIdsRes = append(tagIdsRes, tags)
	}

	return postIdsRes, userIdsRes, titles, descriptions, extensions, tagIdsRes, nil
}

func (c *PostClient) DeletePost(ctx context.Context, postId, userId string) error {
	req := &postProto.DeletePostRequest{
		PostId: postId,
		UserId: userId,
	}

	_, err := c.api.DeletePost(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete post via gRPC")
		return fmt.Errorf("delete post grpc: %w", err)
	}

	return nil
}

func (c *PostClient) CreateBoard(ctx context.Context, userId, name, description string) (string, error) {
	req := &postProto.CreateBoardRequest{
		UserId:      userId,
		Name:        name,
		Description: description,
	}

	res, err := c.api.CreateBoard(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create board via gRPC")
		return "", fmt.Errorf("create board grpc: %w", err)
	}

	return res.BoardId, nil
}

func (c *PostClient) UpdateBoard(ctx context.Context, boardId, userId, name, description string) (string, error) {
	req := &postProto.UpdateBoardRequest{
		BoardId:     boardId,
		UserId:      userId,
		Name:        name,
		Description: description,
	}

	res, err := c.api.UpdateBoard(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to update board via gRPC")
		return "", fmt.Errorf("update board grpc: %w", err)
	}

	return res.Message, nil
}

func (c *PostClient) GetBoard(ctx context.Context, boardId string) (string, string, string, string, [][]string, error) {
	req := &postProto.GetBoardRequest{
		BoardId: boardId,
	}

	res, err := c.api.GetBoard(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get board via gRPC")
		return "", "", "", "", nil, fmt.Errorf("get board grpc: %w", err)
	}

	posts := make([][]string, len(res.Board.Posts))
	for i, p := range res.Board.Posts {
		posts[i] = []string{p.PostId, p.UserId, p.Title, p.Description, p.Extension}
	}

	return res.Board.BoardId, res.Board.UserId, res.Board.Name, res.Board.Description, posts, nil
}

func (c *PostClient) GetBoards(ctx context.Context, boardIds []string) ([]string, []string, []string, []string, [][][]string, error) {
	req := &postProto.GetBoardsRequest{
		BoardIds: boardIds,
	}

	res, err := c.api.GetBoards(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get boards via gRPC")
		return nil, nil, nil, nil, nil, fmt.Errorf("get boards grpc: %w", err)
	}

	var boardIdsRes, userIds, names, descriptions []string
	posts := make([][][]string, len(res.Boards))

	for _, board := range res.Boards {
		post := make([][]string, len(board.Posts))
		for i, p := range board.Posts {
			post[i] = []string{p.PostId, p.UserId, p.Title, p.Description, p.Extension}
		}

		boardIdsRes = append(boardIdsRes, board.BoardId)
		userIds = append(userIds, board.UserId)
		names = append(names, board.Name)
		descriptions = append(descriptions, board.Description)
		posts = append(posts, post)
	}

	return boardIdsRes, userIds, names, descriptions, posts, nil
}

func (c *PostClient) SearchBoards(ctx context.Context, query string, userIds, postIds []string, limit, offset int) ([]string, []string, []string, []string, [][][]string, error) {
	req := &postProto.SearchBoardsRequest{
		Query:  query,
		UserId: userIds,
		PostId: postIds,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	res, err := c.api.SearchBoards(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to search boards via gRPC")
		return nil, nil, nil, nil, nil, fmt.Errorf("search boards grpc: %w", err)
	}

	var boardIdsRes, userIdsRes, names, descriptions []string
	posts := make([][][]string, len(res.Boards))

	for _, board := range res.Boards {
		post := make([][]string, len(board.Posts))
		for i, p := range board.Posts {
			post[i] = []string{p.PostId, p.UserId, p.Title, p.Description, p.Extension}
		}

		boardIdsRes = append(boardIdsRes, board.BoardId)
		userIdsRes = append(userIdsRes, board.UserId)
		names = append(names, board.Name)
		descriptions = append(descriptions, board.Description)
		posts = append(posts, post)
	}

	return boardIdsRes, userIdsRes, names, descriptions, posts, nil
}

func (c *PostClient) DeleteBoard(ctx context.Context, boardId, userId string) (string, error) {
	req := &postProto.DeleteBoardRequest{
		BoardId: boardId,
		UserId:  userId,
	}

	resp, err := c.api.DeleteBoard(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete board via gRPC")
		return "", fmt.Errorf("delete board grpc: %w", err)
	}

	return resp.Message, nil
}

func (c *PostClient) CreateTag(ctx context.Context, name string) (string, error) {
	req := &postProto.CreateTagRequest{
		Name: name,
	}

	res, err := c.api.CreateTag(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create tag via gRPC")
		return "", fmt.Errorf("create tag grpc: %w", err)
	}

	return res.TagId, nil
}

func (c *PostClient) GetTag(ctx context.Context, tagId string) (string, string, error) {
	req := &postProto.GetTagRequest{
		TagId: tagId,
	}

	res, err := c.api.GetTag(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get tag via gRPC")
		return "", "", fmt.Errorf("get tag grpc: %w", err)
	}

	return res.TagId, res.Name, nil
}

func (c *PostClient) SearchTags(ctx context.Context, query string, limit, offset int) ([]string, []string, error) {
	req := &postProto.SearchTagsRequest{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	res, err := c.api.SearchTags(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to search tags via gRPC")
		return nil, nil, fmt.Errorf("search tags grpc: %w", err)
	}

	var tagIds, names []string

	for _, tag := range res.Tags {
		tagIds = append(tagIds, tag.TagId)
		names = append(names, tag.Name)
	}

	return tagIds, names, nil
}

func (c *PostClient) PinPostToBoard(ctx context.Context, postId, boardId, userId string) error {
	req := &postProto.PinPostToBoardRequest{
		PostId:  postId,
		BoardId: boardId,
		UserId:  userId,
	}

	_, err := c.api.PinPostToBoard(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to pin post to board via gRPC")
		return fmt.Errorf("pin post to board grpc: %w", err)
	}

	return nil
}

func (c *PostClient) UnpinPostFromBoard(ctx context.Context, postId, boardId, userId string) error {
	req := &postProto.UnpinPostFromBoardRequest{
		PostId:  postId,
		BoardId: boardId,
		UserId:  userId,
	}

	_, err := c.api.UnpinPostFromBoard(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to unpin post from board via gRPC")
		return fmt.Errorf("unpin post from board grpc: %w", err)
	}

	return nil
}

func (c *PostClient) AddTagToPost(ctx context.Context, postId, tagId, userId string) error {
	req := &postProto.AddTagToPostRequest{
		PostId: postId,
		TagId:  tagId,
		UserId: userId,
	}

	_, err := c.api.AddTagToPost(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to add tag to post via gRPC")
		return fmt.Errorf("add tag to post grpc: %w", err)
	}

	return nil
}

func (c *PostClient) RemoveTagFromPost(ctx context.Context, postId, tagId, userId string) error {
	req := &postProto.RemoveTagFromPostRequest{
		PostId: postId,
		TagId:  tagId,
		UserId: userId,
	}

	_, err := c.api.RemoveTagFromPost(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove tag from post via gRPC")
		return fmt.Errorf("remove tag from post grpc: %w", err)
	}

	return nil
}

func (c *PostClient) Close() {
	log.Debug().Msg("closing Post gRPC client connection with address: " + c.cc.Target())
	err := c.cc.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close Post gRPC client connection")
	}
}
