package itemdata

//go:generate mockgen -source=itemData.go -destination=mocks/mock.go

type ItemData interface {
	CreatePost(post Post) (Post, error)
	GetPosts() ([]Post, error)
	GetCategory(category string) ([]Post, error)
	GetName(login string) ([]Post, error)
	GetPostID(id string) (Post, error)
	SetPost(post Post) error
	DeletePost(postID string) error
}
