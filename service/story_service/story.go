package story_service

import (
	"encoding/json"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
)

func StoryRankService() ([]models.Hot, error) {
	var hots []models.Hot
	if gredis.Exists(e.STORY_TAG_HOT_RANK) {
		data, err := gredis.Get(e.STORY_TAG_HOT_RANK)
		if err != nil {
			return nil, err
		} else {
			json.Unmarshal(data, &hots)
			return hots, nil
		}

	}
	hots, err := models.StoryHot()
	if err != nil {
		return nil, err
	}
	go func() {
		gredis.Set(e.STORY_TAG_HOT_RANK, hots, 3600)
	}()
	return hots, nil
}

func CategoryService() ([]string, error) {
	var cates []string
	if gredis.Exists(e.CATEGORY_STORY) {
		data, err := gredis.Get(e.CATEGORY_STORY)
		if err != nil {
			return nil, err
		} else {
			json.Unmarshal(data, &cates)
			return cates, nil
		}

	}
	cates, err := models.Category()
	if err != nil {
		return nil, err
	}
	go func() {
		gredis.Set(e.CATEGORY_STORY, cates, 3600)
	}()
	return cates, nil
}

func ReloadChapterByIdService(id string) (models.BookContent, error) {
	var content models.BookContent
	key := e.CHAPTER_CONTENT + id
	if gredis.Exists(key) {
		gredis.Delete(key)
	}
	content, err := models.ReloadChapterById(id)
	if err != nil {
		return models.BookContent{}, err
	}
	if content.Content != "" {
		go func() {
			gredis.Set(e.CHAPTER_CONTENT+id, content, 3600)
		}()
	}
	return content, nil
}
func GetChapterByIdService(id string) (models.BookContent, error) {
	var content models.BookContent
	//if gredis.Exists(e.CHAPTER_CONTENT + id) {
	//	data, err := gredis.Get(e.CHAPTER_CONTENT + id)
	//	if err != nil {
	//		return models.BookContent{}, err
	//	} else {
	//		json.Unmarshal(data, &content)
	//		if content.Content != "" {
	//			return content, nil
	//		}
	//	}
	//
	//}
	content, err := models.GetChapterById(id)
	if err != nil {
		return models.BookContent{}, err
	}
	//if content.Content != "" {
	//	go func() {
	//		gredis.Set(e.CHAPTER_CONTENT+id, content, 30)
	//	}()
	//}
	return content, nil
}
func GetChapterByIdServiceAsync(id string) (models.BookContent, error) {
	var content models.BookContent
	//if gredis.Exists(e.CHAPTER_CONTENT + id) {
	//	data, err := gredis.Get(e.CHAPTER_CONTENT + id)
	//	if err != nil {
	//		return models.BookContent{}, err
	//	} else {
	//		json.Unmarshal(data, &content)
	//		if content.Content != "" {
	//			return content, nil
	//		}
	//	}
	//
	//}
	content, err := models.GetChapterByIdAsync(id)
	if err != nil {
		return models.BookContent{}, err
	}
	//if content.Content != "" {
	//	go func() {
	//		gredis.Set(e.CHAPTER_CONTENT+id, content, 30)
	//	}()
	//}
	return content, nil
}
func StoryInfoService(id string) (models.BookDetail, error) {
	var detail models.BookDetail
	if gredis.Exists(e.BOOK_INFO + id) {
		data, err := gredis.Get(e.BOOK_INFO + id)
		if err != nil {
			return models.BookDetail{}, err
		} else {
			json.Unmarshal(data, &detail)
			return detail, nil
		}

	}
	detail, err := models.StoryInfo(id)
	if err != nil {
		return models.BookDetail{}, err
	}
	go func() {
		gredis.Set(e.BOOK_INFO+id, detail, 3600)
	}()
	return detail, nil
}
