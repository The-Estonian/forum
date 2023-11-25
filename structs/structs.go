package structs

type Post struct {
	Id      string
	Title   string
	User    string
	Post    string
	Created string
}

type Category struct {
	Id       string
	Category string
}

type Comment struct {
	Id      string
	PostId  string
	UserId  string
	Comment string
	Created string
}

type User struct {
	Id       string
	Username string
	Email    string
}

type AccessRights struct {
	AccessRight string
}

type MegaData struct {
	User        User
	Post        Post
	AllPosts    []Post
	AllComments []Comment
	Access      AccessRights
	Errors		[]string
}
