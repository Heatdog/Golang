package itemdata

type Author struct {
	ID       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
}

type Votes struct {
	User string `json:"user" bson:"user"`
	Vote int    `json:"vote" valid:"in(-1|1)" bson:"vote"`
}

type Post struct {
	ID               string    `json:"id" bson:"_id"`
	Ath              Author    `json:"author" bson:"author"`
	Comments         []Comment `json:"comments" bson:"comments"`
	Cat              string    `json:"category" valid:"in(music|funny|videos|programming|news|fashion)" bson:"category"`
	Score            int       `json:"score" bson:"score"`
	Type             string    `json:"type" valid:"in(text|link)" bson:"type"`
	Title            string    `json:"title" bson:"title"`
	Created          string    `json:"created" bson:"created"`
	UpvotePercentage int       `json:"upvotePercentage" bson:"upvotePercentage"`
	Views            int64     `json:"views" bson:"views"`
	Text             string    `json:"-" bson:"text"`
	Vote             []Votes   `json:"votes" bson:"vote"`
}

type Comment struct {
	Ath     Author `json:"author" bson:"author"`
	Body    string `json:"body" bson:"body"`
	Created string `json:"created" bson:"created"`
	ID      string `json:"id" bson:"id"`
}

type CreatePost struct {
	Cat   string `json:"category" valid:"in(music|funny|videos|programming|news|fashion),required"`
	Title string `json:"title,required"`
	Type  string `json:"type" valid:"in(text|link),required"`
	Text  string `json:"-"`
}
