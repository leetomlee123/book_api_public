package models

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"gopkg.in/mgo.v2/bson"
)

type ArticleEs struct {
	ID         string `json:"id"`
	BookName   string `json:"title"`
	BookAuthor string `json:"author"`
}
type Book struct {
	Id          bson.ObjectId `bson:"_id"`
	BookName    string        `bson:"book_name" json:"Name"`
	Category    string        `json:"CName"`
	Rate        int           `bson:"rate"`
	Author      string        `bson:"author"`
	UTime       string        `bson:"u_time" json:"UTime"`
	BookDesc    string        `bson:"book_desc" json:"Desc,omitempty"`
	Status      string        `bson:"status" json:"BookStatus"`
	Cover       string        `bson:"cover" json:"Img"`
	LastChapter string        `bson:"last_chapter" `
	//FirstChapterId string `bson:"first_chapter_id"`
	//LastChapterId  string `bson:"last_chapter_id"`
}
type Rank struct {
	Id    string `bson:"id"`
	Cover string `bson:"cover" `
	Name  string `bson:"name"`
}

type RankBook struct {
	Id       string `bson:"_id"`
	BookName string `bson:"book_name" json:"Name"`
	Cover    string `bson:"cover" json:"Img"`
	Author   string `bson:"author"`
	Category string `json:"CName"`
}
type Hot struct {
	Id       bson.ObjectId `bson:"_id"`
	BookName string        `bson:"book_name" json:"Name"`
	Hot      int64
}
type Info struct {
	Date    string `bson:"date" json:"Date"`
	Title   string `bson:"title" json:"Title"`
	Content string `bson:"content" json:"Content"`
}
type User struct {
	Id       bson.ObjectId `bson:"_id" json:"id,omitempty"`
	Name     string        `form:"name" bson:"name" json:"name" binding:"required"`
	PassWord string        `form:"password" bson:"password" json:"password,omitempty" binding:"required"`
	EMail    string        `form:"email" bson:"email" json:"email"`
	Vip      int8          `form:"vip" bson:"vip" json:"vip"`
	State    int8          `bson:"state"`
}
type LoginUser struct {
	Name     string `form:"name" bson:"name" json:"name" binding:"required"`
	PassWord string `form:"password" bson:"password" json:"password,omitempty" binding:"required"`
}
type RegUser struct {
	Name     string `form:"name" bson:"name" json:"name" binding:"required"`
	PassWord string `form:"password" bson:"password" json:"password" binding:"required"`
	EMail    string `form:"email" bson:"email" json:"email" binding:"required"`
	State    int8   `bson:"state"`
	Vip      int8   `form:"vip" bson:"vip" json:"vip"`
}
type BookDetail struct {
	Id          bson.ObjectId `bson:"_id"`
	BookName    string        `bson:"book_name" json:"Name"`
	Category    string        `json:"CName"`
	Author      string        `bson:"author"`
	Hot         int           `bson:"hot"`
	Rate        float64       `bson:"rate"`
	UTime       string        `bson:"u_time" json:"LastTime"`
	BookDesc    string        `bson:"book_desc" json:"Desc,omitempty"`
	Status      string        `bson:"status" json:"BookStatus"`
	Cover       string        `bson:"cover" json:"Img"`
	LastChapter string        `bson:"last_chapter" `
	Count       int
	//LastChapterId   string `bson:"last_chapter_id"`
	//FirstChapterId  string `bson:"first_chapter_id"`
	SameAuthorBooks []Book
}

/**
根据category 分页books
*/
type CateBook struct {
	BookName string        `bson:"book_name" json:"bookName"`
	Author   string        `json:"author"`
	Cover    string        `json:"cover"`
	Id       bson.ObjectId `bson:"_id" json:"id"`
	BookDesc string        `bson:"book_desc" json:"BookDesc"`
}
type Chapter struct {
	ChapterName string        `bson:"chapter_name" json:"name"`
	ChapterId   bson.ObjectId `bson:"_id" json:"id"`
	HasContent  int           ` json:"hasContent"`
}
type BookContent struct {
	Id      string `bson:"_id" json:"id"`
	Content string `bson:"content" json:"content"`
	Link    string `bson:"link" json:"link"`
}
type Result struct {
	Id string `bson:"_id"' `
}

type Account struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string
	IdS  []string `bson:"ids"`
}
type Process struct {
	Account     string
	BookId      string
	BookProcess string
}

//book detail
func StoryInfo(id string) (BookDetail, error) {
	var bookDetail BookDetail
	bid := bson.ObjectIdHex(id)
	if e := bookDB.FindId(bid).One(&bookDetail); e != nil {
		log.Print(e)
		//log.Print(e)
		return BookDetail{}, e
	}
	var bks []Book
	m := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$ne": bookDetail.Id}, "author": bookDetail.Author}},
		{"$project": bson.M{"book_desc": 0}},
	}

	if e := bookDB.Pipe(m).All(&bks); e != nil {
		//log.Print(e)
		return BookDetail{}, e
	}
	bookDetail.SameAuthorBooks = bks
	count, _ := chapterDB.Find(bson.M{"book_id": id}).Count()
	bookDetail.Count = count
	return bookDetail, nil
}

//qidian hot rank
func StoryRank(tye int) (interface{}, error) {
	var _row interface{}

	m := []bson.M{
		{"$match": bson.M{"type": tye}},
		{"$project": bson.M{"_id": 0, "type": 0}},
		//{"$sort": bson.M{"chapter_id": 1}},
	}
	if err := rankDB.Pipe(m).One(&_row); err != nil {
		return nil, err
	}

	return _row, nil
}

//book cloud shelf need account
func Shelf(account Account) ([]Book, error) {
	var books []Book
	var objids []bson.ObjectId
	for _, id := range account.IdS {

		objids = append(objids, bson.ObjectIdHex(id))

	}
	if account.IdS != nil {
		if e := bookDB.Find(bson.M{"_id": bson.M{"$in": objids}}).All(&books); e != nil {
			return nil, e
		}
	}
	go func() {
		accountDB.UpsertId(account.Id, bson.M{"$set": bson.M{"last_alive_date": time.Now().Format(time.RFC3339)}})
	}()
	return books, nil

}

//use click hot rank
func StoryHot() ([]Hot, error) {
	var data []Hot
	m := []bson.M{
		{"$sort": bson.M{"hot": -1}},
		{"$limit": 30},
		{"$project": bson.M{"hot": 1, "_id": 1, "book_name": 1}},
	}
	if err := bookDB.Pipe(m).All(&data); err != nil {
		return nil, err
	}
	return data, nil
}

//hot gt 0 book
func HotGtZeroBooks() ([]Book, error) {
	var bks []Book
	if err := bookDB.Find(bson.M{"hot": bson.M{"$gt": 0}}).All(&bks); err != nil {
		log.Print(err)
		return nil, err
	}
	return bks, nil
}

//modify shelf
func ModifyShelf(bookId string, action string, account Account) error {

	if account.IdS != nil {
		i := 0
		f := true
		for ; i < len(account.IdS); i++ {
			if account.IdS[i] == bookId {
				f = false
				break
			}
		}
		if f {
			if action == "add" {
				account.IdS = append(account.IdS, bookId)
				if err := accountDB.UpdateId(account.Id, bson.M{"$set": bson.M{"ids": account.IdS}}); err != nil {
					return err
				}
			}
		} else {
			if action == "del" {
				if err := accountDB.UpdateId(account.Id, bson.M{"$set": bson.M{"ids": append(account.IdS[:i], account.IdS[i+1:]...)}}); err != nil {
					return err
				}
			}
		}
	} else {
		if action == "add" {
			if _, err := accountDB.UpsertId(account.Id, bson.M{"$set": bson.M{"ids": []string{bookId}}}); err != nil {
				return err
			}
		}
	}
	return nil
}

//get all book category
func Category() ([]string, error) {
	var data []Result

	m := []bson.M{
		{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}},
		//{"$project": bson.M{"count": 0, "_id": 1}},
	}
	if err := bookDB.Pipe(m).All(&data); err != nil {
		return nil, err
	}
	var tem []string
	for _, v := range data {
		tem = append(tem, v.Id)
	}
	return tem, nil
}

//get pages book by category
func GetStoryByCategoryWithPage(page int, cate string, size int) ([]CateBook, error) {

	var datas []CateBook
	//db.Table("book").Where("category=?", cate).Select("author,book_name, cover,id").Offset((page - 1) * size).Limit(size).Scan(&datas)
	m := []bson.M{
		{"$match": bson.M{"category": cate}},
		{"$project": bson.M{"author": 1, "book_name": 1, "cover": 1, "_id": 1, "book_desc": 1}},
		{"$skip": (page - 1) * size},
		{"$limit": size},
	}
	if err := bookDB.Pipe(m).All(&datas); err != nil {
		return nil, err
	}
	return datas, nil
}

//get book chapters by book id
func GetStoryChapters(id string, count int) ([]Chapter, error) {

	var chapters []Chapter
	//db.Table("chapter").Where("book_id=?", id).Select("chapter_id,chapter_name").Order("chapter_id asc").Scan(&chapters)
	m := []bson.M{
		{"$match": bson.M{"book_id": id}},
		{"$project": bson.M{"chapter_name": 1, "_id": 1}},
		//{"$sort": bson.M{"chapter_id": 1}},
		{"$skip": count},
	}
	if err := chapterDB.Pipe(m).All(&chapters); err != nil {
		return nil, err
	}
	var temp []Chapter
	go func() {
		if bson.IsObjectIdHex(id) {
			bookDB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"hot": 1}})
		}

	}()

	for _, chapter := range chapters {
		//go func() {
		//	GetChapterById(chapter.ChapterId.Hex())
		//}()
		temp = append(temp,
			Chapter{ChapterId: chapter.ChapterId, ChapterName: chapter.ChapterName, HasContent: 1})
	}

	return temp, nil
	//chaptersProto := ChaptersProto{}
	//for _, chapter := range chapters {
	//	//go func() {
	//	//	GetChapterById(chapter.ChapterId.Hex())
	//	//}()
	//	proto := ChapterProto{Id: chapter.ChapterId.String(), Name: chapter.ChapterName, HasContent: "1"}
	//
	//	chaptersProto.ChaptersProto = append(chaptersProto.ChaptersProto, &proto)
	//
	//}
	//
	//data, _ := proto.Marshal(&chaptersProto)
	//return data, nil
}

//remove chapter by id
func DeleteChapterById(id string) error {
	err2 := chapterDB.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err2 != nil {
		return err2
	}
	return nil
}

//get chapter by id
func GetChapterById(id string) (BookContent, error) {

	var result BookContent
	err := chapterDB.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return BookContent{}, err
	}
	content := result.Content
	if strings.Contains(result.Link, "paoshuzw") {
		result.Link = strings.Replace(result.Link, "paoshuzw.com", "xbiquge.la", 1)
	}
	if content == "" {
		retry := 3
		for {
			result.Content = GetContent(result.Link, false)
			if (result.Content != "") || (retry <= 0) {
				break
			}
			retry -= 1
		}

		if result.Content == "" {
			return BookContent{}, nil
		}

		go func() {
			if bson.IsObjectIdHex(id) {
				chapterDB.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"content": result.Content}})
			}
		}()
	}
	result.Link = "https://github.com/leetomlee123/book"
	//if result.Content != "" {
	//	//中文全角转半角
	//	result.Content = DBCtoSBC(result.Content)
	//}
	return result, nil
}
func DBCtoSBC(s string) string {
	retStr := ""
	for _, i := range s {
		insideCode := i
		if insideCode == 12288 {
			insideCode = 32
		} else {
			insideCode -= 65248
		}
		if insideCode < 32 || insideCode > 126 {
			retStr += string(i)
		} else {
			retStr += string(insideCode)
		}
	}
	return retStr
}

//get chapter by id
func GetChapterByIdAsync(id string) (BookContent, error) {

	var result BookContent
	err := chapterDB.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return BookContent{}, err
	}
	content := result.Content
	if strings.Contains(result.Link, "paoshuzw") {
		result.Link = strings.Replace(result.Link, "paoshuzw.com", "xbiquge.la", 1)
	}
	if content == "" {
		retry := 3
		for {
			result.Content = GetContent(result.Link, true)
			if (result.Content != "") || (retry <= 0) {
				break
			}
			retry -= 1
		}

		if result.Content == "" {
			return BookContent{}, nil
		}

		go func() {
			if bson.IsObjectIdHex(id) {
				chapterDB.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"content": result.Content}})

			}
		}()
	}
	result.Link = "https://github.com/leetomlee123/book"
	//if result.Content != "" {
	//	//中文全角转半角
	//	result.Content = DBCtoSBC(result.Content)
	//}
	return result, nil
}

//reload chapter content
func ReloadChapterById(id string) (BookContent, error) {
	var result BookContent
	err := chapterDB.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		return BookContent{}, err
	}
	if strings.Contains(result.Link, "paoshuzw") {
		result.Link = strings.Replace(result.Link, "paoshuzw.com", "xbiquge.la", 1)
	}
	content := GetContent(result.Link, false)

	if content != "" {
		go func() {
			if bson.IsObjectIdHex(id) {
				chapterDB.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"content": result.Content}})
			}
		}()
		result.Content = content
	}
	return result, nil
}
func GetBookByAuthorAndName(author string, name string) (BookDetail, error) {
	var bookDetail BookDetail
	if strings.TrimSpace(author) == "" {
		return bookDetail, nil
	}
	if strings.TrimSpace(name) == "" {
		return bookDetail, nil
	}
	var query []bson.M
	q1 := bson.M{"book_name": name}
	query = append(query, q1)
	q2 := bson.M{"author": author}
	query = append(query, q2)
	if e := bookDB.Find(bson.M{"$and": query}).One(&bookDetail); e != nil {
		//log.Print(e)
		//log.Print(e)
		return BookDetail{}, e
	}
	info, err := StoryInfo(bookDetail.Id.Hex())
	if err != nil {
		return BookDetail{}, err
	}
	//var bks []Book
	//m := []bson.M{
	//	{"$match": bson.M{"_id": bson.M{"$ne": bookDetail.Id}, "author": bookDetail.Author}},
	//	{"$project": bson.M{"book_desc": 0}},
	//}
	//
	//if e := bookDB.Pipe(m).All(&bks); e != nil {
	//	//log.Print(e)
	//	return BookDetail{}, e
	//}
	//bookDetail.SameAuthorBooks = bks
	//count, _ := chapterDB.Find(bson.M{"book_id": bookDetail.Id}).Count()
	//bookDetail.Count = count
	return info, nil

}

type Qry struct {
	Qry map[string]interface{}
}

//search service
func Search(key string, page int, size int) ([]Book, error) {

	var bks []Book
	var r map[string]interface{}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"dis_max": map[string]interface{}{
				"queries": []map[string]interface{}{
					{"match": map[string]string{
						"book_name": key,
					}},
					{"match": map[string]string{
						"book_author": key,
					}},
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	search, err := EsDB.Search(
		EsDB.Search.WithContext(context.Background()),
		EsDB.Search.WithIndex("book"),
		EsDB.Search.WithBody(&buf),
		EsDB.Search.WithTrackTotalHits(true),
		EsDB.Search.WithPretty(),
		EsDB.Search.WithFrom((page-1)*size),
		EsDB.Search.WithSize(size),
	)
	if err != nil {
		return nil, err
	}

	defer search.Body.Close()
	if search.IsError() {
		return nil, nil
	}

	if err := json.NewDecoder(search.Body).Decode(&r); err != nil {
		return bks, nil
	}
	// Print the ID and document source for each hit.
	var ids []string
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		//ids = append(ids, hit.(map[string]string)["_id"])
		i := hit.(map[string]interface{})["_id"]
		ids = append(ids, i.(string))
	}

	//if strings.TrimSpace(key) == "" {
	//	return bks, nil
	//}
	//var query []bson.M
	////var all []bson.M
	//q1 := bson.M{"book_name": bson.M{"$regex": key, "$options": "$i$m"}}
	//query = append(query, q1)
	//q2 := bson.M{"author": bson.M{"$regex": key, "$options": "$i$m"}}
	////
	//query = append(query, q2)
	//
	//
	for _, id := range ids {
		var bk Book
		if bson.IsObjectIdHex(id) {
			hex := bson.ObjectIdHex(id)
			err := bookDB.Find(bson.M{"_id": hex}).One(&bk)
			if err != nil {
				continue
			}
			bks = append(bks, bk)
		}
	}

	return bks, nil
}

//site notice
func Notice() (error, []Info) {
	var info []Info
	if err := infoDB.Find(bson.M{}).All(&info); err != nil {
		return err, nil
	}
	return nil, info
}

//get account info by token
func GetAccountByName(username string) (Account, error) {

	var account Account
	if e := accountDB.Find(bson.M{"name": username}).One(&account); e != nil {
		return Account{}, e
	}
	return account, nil
}

//阅读进度
func PersonBookReadProcess(account string, bookId string, process string) error {
	var processes []Process
	var query []bson.M
	q1 := bson.M{"account": account}
	query = append(query, q1)
	q2 := bson.M{"bookid": bookId}
	query = append(query, q2)
	if err := processDB.Find(bson.M{"$and": query}).All(&processes); err != nil {
		return err
	}
	if len(processes) > 0 {
		_, err := processDB.Upsert(bson.M{"$and": query}, bson.M{"$set": bson.M{"bookprocess": process}})
		if err != nil {
			logging.Error(err)
		}
	} else {
		if err := processDB.Insert(&Process{
			Account:     account,
			BookId:      bookId,
			BookProcess: process,
		}); err != nil {
			print(err)
			logging.Error(err)
		}
	}
	return nil
}

//get read process
func GetReadProcess(account string, bookId string) string {
	var processes []Process
	var query []bson.M
	q1 := bson.M{"account": account}
	query = append(query, q1)
	q2 := bson.M{"bookid": bookId}
	query = append(query, q2)
	if err := processDB.Find(bson.M{"$and": query}).All(&processes); err != nil {
		return ""
	}
	if len(processes) > 0 {
		return processes[0].BookProcess

	} else {
		return ""
	}
}

//get chapter content by web link
func GetContent(u string, asnyc bool) string {
	// Request the HTML page.
	//u = "http://localhost:8085/book/chapter?url=" + u

	if asnyc {
		u = setting.AppSetting.AsyncContentPath + "/book/chapter?url=" + u
	} else {
		u = setting.AppSetting.ContentPath + "/book/chapter?url=" + u

	}

	// u = "http://23.91.100.230:8085/book/chapter?url=" + u
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(u)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		//log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	ss, _ := ioutil.ReadAll(res.Body) //把  body 内容读入字符串 s

	return string(ss)

}

//get web encode
func determineEncoding(r io.Reader) encoding.Encoding {
	bytes, _ := bufio.NewReader(r).Peek(1024)

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
