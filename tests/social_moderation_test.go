package tests

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/repository"
)

func TestCommentModerationPersistence(t *testing.T) {
	ctx:=context.Background()
	db,err:=repository.NewDB("sqlite",filepath.Join(t.TempDir(),"moderation.db"))
	if err!=nil{t.Fatal(err)}
	defer db.Close()
	users:=repository.NewUserRepo(db)
	for _,u:=range []domain.User{
		{ID:"author",DisplayName:"Tác giả",Email:"moderation-author@example.com",PasswordHash:"hash",CreatedAt:time.Now(),UpdatedAt:time.Now()},
		{ID:"reader",DisplayName:"Người đọc",Email:"moderation-reader@example.com",PasswordHash:"hash",CreatedAt:time.Now(),UpdatedAt:time.Now()},
	}{
		if err:=users.Create(ctx,&u);err!=nil{t.Fatal(err)}
	}
	stories:=repository.NewStoryRepo(db)
	story:=&domain.Story{ID:"story-"+uuid.NewString(),UserID:"author",Title:"Đang xuất bản",Status:"published",CreatedAt:time.Now(),UpdatedAt:time.Now()}
	if err:=stories.CreateStory(ctx,story);err!=nil{t.Fatal(err)}
	comment,err:=stories.CreateComment(ctx,story.ID,"reader","Truyện rất hay!")
	if err!=nil{t.Fatal(err)}
	queue,err:=stories.ListPendingComments(ctx,20)
	if err!=nil||len(queue)!=1||queue[0].ID!=comment.ID{t.Fatalf("pending queue: %+v %v",queue,err)}
	visible,err:=stories.ListComments(ctx,story.ID,"author")
	if err!=nil||len(visible)!=0{t.Fatal("unapproved comment leaked")}
	if err:=stories.ModerateComment(ctx,comment.ID,"approved");err!=nil{t.Fatal(err)}
	visible,err=stories.ListComments(ctx,story.ID,"author")
	if err!=nil||len(visible)!=1||visible[0].Body!="Truyện rất hay!"{t.Fatal("approved comment not visible")}
	if err:=stories.ModerateComment(ctx,comment.ID,"rejected");err!=repository.ErrCommentNotFound{t.Fatalf("moderation must be one-shot: %v",err)}
	next,err:=stories.CreateComment(ctx,story.ID,"reader","Bình luận bị từ chối")
	if err!=nil{t.Fatal(err)}
	if err:=stories.ModerateComment(ctx,next.ID,"rejected");err!=nil{t.Fatal(err)}
	visible,err=stories.ListComments(ctx,story.ID,"author")
	if err!=nil||len(visible)!=1{t.Fatal("rejected comment became visible")}
}
