package main

import (
	"./common"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type forumInterface interface {
	GetListTopics() []topic
	GetListTopicsDate() []topic
	GetLatestTopicId() string
	CreateTopic(topic) forum
	UpdateTopic(topic) forum
}

type topicInterface interface {
	PostTo(post) post
	GetPostByID(string) post
	GetListPosts() []post
	GetLatestPostId() string
	GetUpdateTime() string
}

type forum struct {
	Name        string  `json:"name"`
	LatestTopic string  `json:"latestTopic"`
	Template    string  `json:"template"`
	Topics      []topic `json:"topics"`
}

type postPost struct {
	Captcha string `json:"Captcha"`
	text    string `json:"text"`
	author  string `json:"author"`
	Ok      string `json:"Ok"`
}

type topic struct {
	Id         string    `json:"id"`
	Title      string    `json:"title"`
	Name       string    `json:"name"`
	Author     string    `json:"author"`
	Comments   int       `json:"comments"`
	Created    time.Time `json:"created"`
	LastUpdate time.Time `json:"updated"`
	LatestPost string    `json:"latestPost"`
	Closed     bool      `json:"closed"`
	Template   string    `json:"template"`
	Posts      []post    `json:"posts"`
}

type post struct {
	Id      string    `json:"id"`
	Text    string    `json:"text"`
	Author  string    `json:"author"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}

type topicType []topic

func (a topicType) Len() int           { return len(a) }
func (a topicType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a topicType) Less(i, j int) bool { return a[i].LastUpdate.After(a[j].LastUpdate) }

func (f forum) GetListTopicsDate() []topic {
	tmp := f.Topics
	sort.Sort(topicType(tmp))
	return tmp
}

func (f forum) GetListTopics() []topic {
	return f.Topics
}

func (f forum) GetLatestTopicId() string {
	return f.LatestTopic
}

func (f forum) CreateTopic(t topic) forum {
	log.Debugf("Creating new Topic " + t.Title)
	max := 0
	for _, t := range f.GetListTopics() {
		i, _ := strconv.Atoi(t.Id)
		if i > max {
			max = i
		}
	}
	t.Id = strconv.Itoa(max + 1)
	t.Created = time.Now()
	t.LastUpdate = t.Created
	t.Closed = false
	if len(t.Posts) >= 0 {
		t.Comments = len(t.Posts) - 1
	} else {
		t.Comments = -1
	}
	for i, _ := range t.Posts {
		t.Posts[i].Id = strconv.Itoa(i)
		t.Posts[i].Created = t.Created
	}
	t.Name = f.Name
	f.LatestTopic = t.Id
	f.Topics = append(f.Topics, t)
	f.writeHTML("static/index.html")
	t.writeHTML("static/topic/" + t.Id + "/index.html")
	return f
}

// replace topic within forum with the updated topic
func (f forum) UpdateTopic(t topic) forum {
	log.Debugf("Updating Topic " + t.Title)
	topicI := -1
	for i, to := range f.GetListTopics() {
		if to.Id == t.Id {
			topicI = i
		}
	}
	if topicI >= 0 {
		f.Topics[topicI] = t
	}
	f.writeHTML("static/index.html")
	return f
}

func (f forum) GetTopicByID(id string) topic {
	var ret topic
	for _, t := range f.GetListTopics() {
		if t.Id == id {
			ret = t
			break
		}
	}
	return ret
}

func (t topic) GetUpdateTime() string {
	return t.LastUpdate.Format("Jan 02, 2006 15:04")
}

func (t topic) PostTo(p post) topic {
	log.Debugf("Adding new posting from " + p.Author + " to " + t.Title)
	max := 0
	for _, po := range t.GetListPosts() {
		i, _ := strconv.Atoi(po.Id)
		if i > max {
			max = i
		}
	}
	p.Id = strconv.Itoa(max + 1)
	p.Created = time.Now()
	t.LatestPost = p.Id
	t.LastUpdate = p.Created
	t.Posts = append(t.Posts, p)
	t.Comments = t.Comments + 1
	t.writeHTML("static/topic/" + t.Id + "/index.html")
	return t
}

func (t topic) GetListPosts() []post {
	return t.Posts
}

func (t topic) GetLatestPostId() string {
	return t.LatestPost
}

func (t topic) GetPostByID(id string) post {
	var ret post
	for _, p := range t.GetListPosts() {
		if p.Id == id {
			ret = p
			break
		}
	}
	return ret
}

func writeNewPostHtml(t topic) {
	log.Debugf("Creating Post page for Topic " + t.Id)
	var app map[string]interface{}
	output := "static/topic/" + t.Id + "/newpost.html"
	input_data, err := ioutil.ReadFile("templates/post.html") // To Do change default template
	check(err)

	b, _ := json.Marshal(t)
	err = json.Unmarshal(b, &app)
	check(err)

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"date": func(s string) string {
			t, _ := time.Parse(time.RFC3339, s)
			return t.Format("Jan 02, 2006 15:04")
		},
	}).Parse(string(input_data))
	check(err)
	fName := filepath.Base(output)
	path := output[:len(output)-len(fName)]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	fi, err := os.Create(output)
	check(err)

	err = tmpl.Execute(fi, app)
	check(err)
}

func writeNewTopicHtml(f forum) {
	log.Debugf("Creating post topic page")
	var app map[string]interface{}
	output := "static/topic/new/index.html"
	input_data, err := ioutil.ReadFile("templates/topic.html") // To Do change default template
	check(err)

	b, _ := json.Marshal(f)
	err = json.Unmarshal(b, &app)
	check(err)

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"date": func(s string) string {
			t, _ := time.Parse(time.RFC3339, s)
			return t.Format("Jan 02, 2006 15:04")
		},
	}).Parse(string(input_data))
	check(err)
	fName := filepath.Base(output)
	path := output[:len(output)-len(fName)]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	fi, err := os.Create(output)
	check(err)

	err = tmpl.Execute(fi, app)
	check(err)
}

func (t topic) writeHTML(output string) {
	log.Debugf("Creating discussion page for  " + t.Id)
	var app map[string]interface{}
	input_data, err := ioutil.ReadFile(t.Template)
	if err != nil {
		input_data, err = ioutil.ReadFile("templates/discussion.html") // To Do change default template
	}
	check(err)

	b, _ := json.Marshal(t)
	err = json.Unmarshal(b, &app)
	check(err)

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"date": func(s string) string {
			t, _ := time.Parse(time.RFC3339, s)
			return t.Format("Jan 02, 2006 15:04")
		},
	}).Parse(string(input_data))
	check(err)
	fName := filepath.Base(output)
	path := output[:len(output)-len(fName)]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	fi, err := os.Create(output)
	check(err)

	err = tmpl.Execute(fi, app)
	check(err)
	writeNewPostHtml(t)
}

func (f forum) writeHTML(output string) {
	log.Debugf("Writing Index page for forum")
	var app map[string]interface{}
	input_data, err := ioutil.ReadFile(f.Template)
	check(err)

	tmp := f
	tmp.Topics = f.GetListTopicsDate()
	b, _ := json.Marshal(tmp)
	err = json.Unmarshal(b, &app)
	check(err)

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"date": func(s string) string {
			t, _ := time.Parse(time.RFC3339, s)
			return t.Format("Jan 02, 2006 15:04")
		},
	}).Parse(string(input_data))
	check(err)

	fi, err := os.Create(output)
	check(err)

	err = tmpl.Execute(fi, app)
	check(err)
	writeNewTopicHtml(f)
}

var myforum forum

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	myforum = forum{
		Name:     "My Forum",
		Template: "templates/forum.html",
		Topics:   []topic{},
	}
	myforum.writeHTML("static/index.html")
	writeNewTopicHtml(myforum)
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("to be implemented"))
}

func Topics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, _ := json.Marshal(myforum.GetListTopics())
	w.Write(b)
}

func NewTopic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := topic{}
	json.NewDecoder(r.Body).Decode(&t)
	log.Debugf("Received API Request to add new Topic  " + t.Title)
	myforum = myforum.CreateTopic(t)
	b, _ := json.Marshal(myforum.GetLatestTopicId)
	log.Debugf("Created Topic with ID  " + myforum.GetLatestTopicId())
	w.Write(b)
}

func Topic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("topic")
	b, _ := json.Marshal(myforum.GetTopicByID(id))
	w.Write(b)
}

func Post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := ps.ByName("topic")
	p := ps.ByName("post")
	b, _ := json.Marshal(myforum.GetTopicByID(t).GetPostByID(p))
	w.Write(b)
}

func NewPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := ps.ByName("topic")
	p := post{}
	json.NewDecoder(r.Body).Decode(&p)
	log.Debugf("Received API Request to add new Post  " + p.Text)
	updatedTopic := myforum.GetTopicByID(t).PostTo(p)
	myforum = myforum.UpdateTopic(updatedTopic)
	log.Debugf("Updated Topic with ID  " + myforum.GetTopicByID(t).Id)
	b, _ := json.Marshal(myforum.GetTopicByID(t).GetLatestPostId)
	w.Write(b)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setupRoutes(router *httprouter.Router) {
	// login page
	// re-direct to login page
	router.GET("/topic", Topics)
	router.POST("/topic", NewTopic)
	router.GET("/topic/:topic", Topic)
	router.GET("/topic/:topic/post/:post", Post)
	router.POST("/topic/:topic", NewPost)
	router.ServeFiles("/forum/*filepath", http.Dir("static"))
}

func main() {
	var embeddedTLSserver common.EmbeddedServer
	embeddedTLSserver.New(common.WebserverCertificate, common.WebserverPrivateKey)

	router := httprouter.New()
	setupRoutes(router)

	log.Printf("Listening on HTTPS port: %s\n", common.ListenPort)
	log.Fatal(embeddedTLSserver.ListenAndServeTLS(":"+common.ListenPort, router))

}
