package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	DateTime string `json:"date_time"`
	PostId   int    `json:"post_id"`
	Content  string `json:"content"`
	FileUrl  string `json:"file_url"`
}

type ThreadApiResponse struct {
	Posts []struct {
		No          int    `json:"no"`
		Now         string `json:"now"`
		Name        string `json:"name"`
		Sub         string `json:"sub,omitempty"`
		Com         string `json:"com,omitempty"`
		Filename    string `json:"filename,omitempty"`
		Ext         string `json:"ext,omitempty"`
		W           int    `json:"w,omitempty"`
		H           int    `json:"h,omitempty"`
		TnW         int    `json:"tn_w,omitempty"`
		TnH         int    `json:"tn_h,omitempty"`
		Tim         int64  `json:"tim,omitempty"`
		Time        int    `json:"time"`
		Md5         string `json:"md5,omitempty"`
		Fsize       int    `json:"fsize,omitempty"`
		Resto       int    `json:"resto"`
		Bumplimit   int    `json:"bumplimit,omitempty"`
		Imagelimit  int    `json:"imagelimit,omitempty"`
		SemanticURL string `json:"semantic_url,omitempty"`
		Replies     int    `json:"replies,omitempty"`
		Images      int    `json:"images,omitempty"`
		UniqueIps   int    `json:"unique_ips,omitempty"`
		Trip        string `json:"trip,omitempty"`
	} `json:"posts"`
}

type ThreadHandler struct {
	ThreadExists  bool
	CurrentThread string
	Posts         map[int]*Post
}

func CreateThreadHandler() *ThreadHandler {
	t := new(ThreadHandler)
	t.Posts = make(map[int]*Post)
	t.ThreadExists = false
	return t
}

func (t *ThreadHandler) UpdateThread() {
	// @todo consider not relying on chinkshit.xyz because xyz-anon seems dead
	resp, err := http.Get("http://chinkshit.xyz")

	// if a thread doesn't exist, flag the current status as thread-down and
	// stop bothering.
	if err != nil {
		t.ThreadExists = false
		t.CurrentThread = ""
		return
	}

	finalUrl := resp.Request.URL.String()

	match, _ := regexp.MatchString("^(http|https)://boards.4channel.org/g/thread/[0-9]+/", finalUrl)

	if match {
		t.ThreadExists = true
		t.CurrentThread = finalUrl
	}

	t.FetchPosts()
}

func (t *ThreadHandler) FetchPosts() {
	threadData := ThreadApiResponse{}

	if !t.ThreadExists {
		return
	}

	re := regexp.MustCompile("[0-9]+/")
	threadId := re.FindString(t.CurrentThread)
	threadId = strings.Trim(threadId, "/")

	threadUrl := fmt.Sprintf("http://a.4cdn.org/g/thread/%s.json", threadId)
	resp, err := http.Get(threadUrl)
	if err != nil {
		fmt.Printf("Can't fetch posts: %s \n", err.Error())
		return
	}

	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	err = json.Unmarshal(respBytes, &threadData)
	if err != nil {
		fmt.Printf("Can't parse json from thread response: %s \n", err.Error())
	}

	for _, post := range threadData.Posts {
		fileUrl := ""

		if post.Ext != "" && post.Tim != 0 && post.Fsize != 0 {
			fileName := strconv.FormatInt(post.Tim, 10) + "" + post.Ext
			fileUrl = fmt.Sprintf("https://i.4cdn.org/g/%s", fileName)
		}

		t.Posts[post.No] = &Post{DateTime: post.Now, PostId: post.No, Content: post.Com, FileUrl: fileUrl}
		postFile := fmt.Sprintf("posts/%s.json", strconv.Itoa(t.Posts[post.No].PostId))
		if _, err := os.Stat(postFile); os.IsNotExist(err) {
			post, _ := json.Marshal(t.Posts[post.No])
			ioutil.WriteFile(postFile, post, 0644)
		}
	}
}

func (t *ThreadHandler) GetPost(postId int) *Post {
	// try from cache
	if post, ok := t.Posts[postId]; ok {
		return post
	}

	// try from fs
	fileName := fmt.Sprintf("posts/%s.json", strconv.Itoa(postId))
	if _, err := os.Stat(fileName); err == nil {
		bytes, _ := ioutil.ReadFile(fileName)
		post := Post{}
		json.Unmarshal(bytes, &Post{})
		return &post
	}

	// not found
	return nil
}

func (t *ThreadHandler) StartWatcher() {
	go t.ThreadWatcher()
	go t.PostWatcher()
}

func (t *ThreadHandler) ThreadWatcher() {
	t.UpdateThread()
	time.Sleep(30 * time.Minute)
	t.ThreadWatcher()
}

func (t *ThreadHandler) PostWatcher() {
	t.FetchPosts()
	time.Sleep(5 * time.Minute)
	t.PostWatcher()
}
