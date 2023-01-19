package data

import (
	"errors"
	"incrediblefour/features/content"
	"log"

	"gorm.io/gorm"
)

type contentData struct {
	db *gorm.DB
}

func New(db *gorm.DB) content.ContentData {
	return &contentData{
		db: db,
	}
}

func (cd *contentData) Add(userID uint, newContent content.Core) (content.Core, error) {
	cnv := CoreToData(newContent)
	cnv.UserID = uint(userID)
	err := cd.db.Create(&cnv).Error
	if err != nil {
		return content.Core{}, err
	}

	newContent.ID = cnv.ID

	return newContent, nil
}

func (cd *contentData) ContentDetail(contentID uint) (content.Core, error) {
	res := Contents{}
	if err := cd.db.Table("contents").Joins("JOIN users ON users.id = contents.user_id").Select("contents.id, users.avatar as avatar, users.username as username, contents.image, contents.caption, contents.created_at as CreatedAt").Where("contents.id = ?", contentID).Find(&res).Error; err != nil {
		log.Println("Get User Content by User ID query error : ", err.Error())
		return content.Core{}, err
	}

	if res.Image == "" {
		return content.Core{}, errors.New("Data not found")
	}

	return ToCore(res), nil
}

func (cd *contentData) ContentList() ([]content.Core, error) {
	res := []AllContents{}
	if err := cd.db.Table("contents").Joins("JOIN users ON users.id = contents.user_id").Select("contents.id, users.avatar as avatar, users.username as username, contents.image, contents.caption, contents.created_at as CreatedAt").Find(&res).Error; err != nil {
		log.Println("get all content query error : ", err.Error())
		return []content.Core{}, err
	}
	return AllListToCore(res), nil
}

func (cd *contentData) Update(userID uint, contentID uint, updatedContent content.Core) (content.Core, error) {
	getID := Contents{}
	err := cd.db.Where("id = ?", contentID).First(&getID).Error

	if err != nil {
		log.Println("get content error : ", err.Error())
		return content.Core{}, err
	}

	if getID.UserID != userID {
		log.Println("Unauthorized request")
		return content.Core{}, errors.New("Unauthorized request")
	}

	cnv := CoreToData(updatedContent)
	qry := cd.db.Where("id = ?", contentID).Updates(&cnv)
	if qry.RowsAffected <= 0 {
		log.Println("update content query error : data not found")
		return content.Core{}, errors.New("not found")
	}

	if err := qry.Error; err != nil {
		log.Println("update content query error : ", err.Error())
	}
	return updatedContent, nil
}

func (cd *contentData) Delete(userID uint, contentID uint) error {
	getID := Contents{}
	err := cd.db.Where("id = ?", contentID).First(&getID).Error

	if err != nil {
		log.Println("get content error : ", err.Error())
		return errors.New("failed to get content data")
	}

	if getID.UserID != userID {
		log.Println("unauthorized request")
		return errors.New("inauthorized request")
	}

	qryDelete := cd.db.Delete(&Contents{}, contentID)

	affRow := qryDelete.RowsAffected

	if affRow <= 0 {
		log.Println("No rows affected")
		return errors.New("failed to delete user content, data not found")
	}

	return nil
}

func (cd *contentData) GetProfile(username string) ([]content.Core, error) {
	res := []Contents{}
	if err := cd.db.Table("contents").Joins("JOIN users ON users.id = contents.user_id").Select("users.id, users.avatar as avatar, users.username as username, contents.image, contents.caption, contents.user_id, contents.created_at").Where("users.username = ?", username).Find(&res).Error; err != nil {
		log.Println("Get User Content by username query error : ", err.Error())
		return []content.Core{}, err
	}

	return ListToCore(res), nil

}

// 	// res := map[string]interface{}{}

// 	// if err := cd.db.Raw("SELECT users.id, users.avatar as pp_path, users.banner, users.name, users.username, users.bio from users where username = ?", username).Find(&res).Error; err != nil {
// 	// 	log.Println("get user content by username query error : ", err.Error())
// 	// 	return nil, err
// 	// }

// 	// resultContent := []map[string]interface{}{}
// 	// // resContent := Contents{}

// 	// if err := cd.db.Raw("SELECT cn.id, cn.image, cn.caption, u.username, u.avatar, c.text FROM contents cn JOIN users u1 ON u1.id = cn.user_id JOIN comments c ON c.content_id = cn.id JOIN users u ON u.id = c.user_id ").Where("u1.username = ?", username).Find(&resultContent).Error; err != nil {
// 	// 	log.Println("get user content by username query error : ", err.Error())
// 	// 	return nil, err
// 	// }

// 	// resComment := map[string]interface{}{}

// 	// for _, v := range resultContent{
// 	// 	if v["id"] = resultContent.ID{
// 	// 		resComment = v
// 	// 	}
// 	// }
// 	// // resultContent["id"] = resContent.ID
// 	// // resultContent["avatar"] = resContent.Avatar
// 	// // resultContent["username"] = resContent.Username
// 	// // resultContent["image"] = resContent.Image
// 	// // resultContent["caption"] = resContent.Caption
// 	// // resultContent["user_id"] = resContent.UserID

// 	// return res, nil

// }
