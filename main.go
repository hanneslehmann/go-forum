package main

import (
	"log"
  "time"
  "fmt"
  "os"
  "io/ioutil"
  "net/http"
  "strconv"
  "path/filepath"
  "encoding/json"
  "text/template"
  "github.com/julienschmidt/httprouter"
  "./common"
)

type forumInterface interface {
    GetListTopics() []topic
    GetTopicByID(string) topic
    CreateTopic(topic) forum
    UpdateTopic(topic) forum
    GetLatestTopicId() string
}

type topicInterface interface {
    PostTo(post) post
    GetPostByID(string) post
    GetListPosts() []post
    GetLatestPostId() string
}

type forum struct {
  Name string         `json:"name"`
  LatestTopic string  `json:"latestTopic"`
  Template string     `json:"template"`
  Topics []topic      `json:"topics"`
}

type topic struct {
  Id string           `json:"id"`
  Title string        `json:"title"`
  Author string       `json:"author"`
  Comments int        `json:"comments"`
  Created time.Time   `json:"created"`
  LatestPost string   `json:"latestPost"`
  Closed bool         `json:"closed"`
  Template string     `json:"template"`
  Posts []post        `json:"posts"`
}

type post struct {
  Id string           `json:"id"`
  Text string         `json:"text"`
  Author string       `json:"author"`
  Created time.Time   `json:"created"`
}

func (f forum) GetListTopics() []topic {
  return f.Topics
}

func (f forum) GetLatestTopicId() string {
  return f.LatestTopic
}

func (f forum) CreateTopic(t topic) forum{
  max:=0
  for _,t := range f.GetListTopics() {
    i,_ :=strconv.Atoi(t.Id)
    if  i> max {
        max=i
    }
  }
  t.Id=strconv.Itoa(max+1)
  t.Created=time.Now()
  t.Closed=false
  f.LatestTopic=t.Id
  f.Topics = append(f.Topics, t)
  f.writeHTML("static/index.html")
  return f
}

func (f forum) UpdateTopic(t topic) forum{
  topicI:=-1
  for i,to := range f.GetListTopics() {
    if  to.Id == t.Id {
        topicI=i
    }
  }
  if topicI>=0{
    f.Topics[topicI]=t
  }
  f.writeHTML("static/index.html")
  return f
}

func (f forum) GetTopicByID(id string) topic {
  var ret topic
  for _,t := range f.GetListTopics() {
    if t.Id == id {
        ret=t
        break
    }
  }
  return ret
}

func (t topic) PostTo(p post) topic {
  max:=0
  for _,po := range t.GetListPosts() {
    i,_ :=strconv.Atoi(po.Id)
    if  i> max {
        max=i
    }
  }
  p.Id=strconv.Itoa(max+1)
  p.Created=time.Now()
  t.LatestPost=p.Id
  t.Posts = append(t.Posts, p)
  t.writeHTML("static/topic/"+t.Id+"/index.html")
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
  for _,p := range t.GetListPosts() {
    if p.Id == id {
        ret=p
        break
    }
  }
  return ret
}

var myforum forum

func init() {
  t:= topic {
    Id: "1",
    Title: "Topic 1",
    Author: "Hannes",
    Comments: 0,
    Created: time.Now(),
    Template: "templates/discussion.html",
    Posts: []post {
      post{ Id:"0",Text:"Topic1 - No text 1",Author:"Hannes1",Created:time.Now() },
      post{ Id:"1",Text:"Topic1 - No text 2",Author:"Hannes2",Created:time.Now(),},
    },
  }
  t2:= topic {
    Id: "2",
    Title: "Topic 2",
    Author: "Hannes",
    Comments: 0,
    Created: time.Now(),
    Template: "templates/discussion.html",
    Posts: []post {
      post{ Id:"0",Text:"Topic2 - No text 1",Author:"Hannes1",Created:time.Now() },
      post{ Id:"1",Text:"Topic2 - No text 2",Author:"Hannes2",Created:time.Now(),},
    },
  }
  myforum = forum{
    Name: "My Forum",
    Template: "templates/forum.html",
    Topics: []topic{t,t2},
  }
  myforum.writeHTML("static/index.html")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "to be implemented")
}

func Topics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  b, _ := json.Marshal(myforum.GetListTopics())
  w.Write(b)
}

func NewTopic (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  t:=topic{}
  json.NewDecoder(r.Body).Decode(&t)
  myforum=myforum.CreateTopic(t)
  b, _ := json.Marshal(myforum.GetLatestTopicId)
  w.Write(b)
}


func Topic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  id:=ps.ByName("topic")
  b, _ := json.Marshal(myforum.GetTopicByID(id))
  w.Write(b)
}

func Post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  t:=ps.ByName("topic")
  p:=ps.ByName("post")
  b, _ := json.Marshal(myforum.GetTopicByID(t).GetPostByID(p))
  w.Write(b)
}

func NewPost (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  t:=ps.ByName("topic")
  p:=post{}
  json.NewDecoder(r.Body).Decode(&p)
  updatedTopic:=myforum.GetTopicByID(t).PostTo(p)
  myforum=myforum.UpdateTopic(updatedTopic)
  b, _ := json.Marshal(myforum.GetTopicByID(t).GetLatestPostId)
  w.Write(b)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func writeForumToJsonFile(target string, f forum){

}

func (t topic) writeHTML(output string){
  var app map[string]interface{}
  input_data, err := ioutil.ReadFile(t.Template)
  check(err)

  b, _ := json.Marshal(t)
  err = json.Unmarshal(b, &app)
  check(err)

  tmpl, err := template.New("test").Parse(string(input_data))
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

func (f forum) writeHTML(output string){
  var app map[string]interface{}
  input_data, err := ioutil.ReadFile(f.Template)
  check(err)

  b, _ := json.Marshal(f)
  err = json.Unmarshal(b, &app)
  check(err)

  tmpl, err := template.New("test").Parse(string(input_data))
  check(err)

  fi, err := os.Create(output)
  check(err)

  err = tmpl.Execute(fi, app)
  check(err)
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
	embeddedTLSserver := &embeddedServer{
		webserverCertificate: common.WebserverCertificate,
		webserverKey: common.WebserverPrivateKey,
  }

  router := httprouter.New()
  setupRoutes(router)


	log.Printf("Listening on HTTPS port: %s\n", common.ListenPort)
	log.Fatal(embeddedTLSserver.ListenAndServeTLS(":"+common.ListenPort, router))
}
