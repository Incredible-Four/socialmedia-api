package content

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID       uint
	Avatar   string
	Username string
	Image    string
	Caption  string
	UserID   uint
	CreatedAt time.Time
}

type ContentHandler interface {
	Add() echo.HandlerFunc
	ContentDetail() echo.HandlerFunc
	ContentList() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type ContentService interface {
	Add(token interface{}, newContent Core) (Core, error)
	ContentDetail(contentID uint) (Core, error)
	ContentList() ([]Core, error)
	Update(token interface{}, contentID uint, updatedContent Core) (Core, error)
	Delete(token interface{}, contentID uint) error
}

type ContentData interface {
	Add(userID uint, newContent Core) (Core, error)
	ContentDetail(contentID uint) (Core, error)
	ContentList() ([]Core, error)
	Update(userID uint, contentID uint, updatedContent Core) (Core, error)
	Delete(userID uint, contentID uint) error
}
