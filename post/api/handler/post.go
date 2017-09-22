package handler

import (
	"encoding/json"

	"github.com/micro/go-log"

	"github.com/hb-go/micro/post/api/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/go-api/proto"
	post "github.com/hb-go/micro/post/api/proto/post"
	postSrv "github.com/hb-go/micro/post/srv/proto/post"
	commentSrv "github.com/hb-go/micro/post/srv/proto/comment"

	"golang.org/x/net/context"
)

type Post struct{}

// Example.Call is called by the API as /post/example/call with post body {"name": "foo"}
func (e *Post) Post(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Log("Received Post.Post request")

	// extract the client from the context

	postClient, ok := client.PostFromContext(ctx)
	if !ok {
		return errors.InternalServerError("go.micro.api.post.post", "post client not found")
	}

	id, ok := req.Get["id"]
	log.Log("Req post id:%v", id.String())

	response := post.Response{}

	// make request
	rspPost, err := postClient.GetPost(ctx, &postSrv.ReqPost{
		Id: 0,
	})
	if err != nil {
		return errors.InternalServerError("go.micro.api.post.post.GetPost", err.Error())
	}
	response.Post = rspPost

	commentClient, ok := client.CommentFromContext(ctx)
	if !ok {
		return errors.InternalServerError("go.micro.api.post.comment", "post client not found")
	}

	// make request
	rspComments, err := commentClient.GetComments(ctx, &commentSrv.ReqComment{
		Id: 0,
	})
	if err != nil {
		return errors.InternalServerError("go.micro.api.post.comment.GetComments", err.Error())
	}
	response.Comments = rspComments.Comments

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}