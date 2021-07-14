package controller

import (
	"github.com/cnmade/bsmi-kb/app/orm/model"
	"github.com/cnmade/bsmi-kb/app/vo"
	"github.com/cnmade/bsmi-kb/pkg/common"
	"github.com/gin-gonic/gin"

	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

type Api struct {
}

type apiBlogList struct {
	Aid   string `form:"aid" json:"aid"  binding:"required"`
	Title string `form:"title" json:"title"  binding:"required"`
}

func (a *Api) NavAll(c *gin.Context) {

	var articleList []model.Article
	common.NewDb.Where("p_aid = 0").Find(&articleList)

	var na []vo.Nav_item
	for _, s := range articleList {
		na = append(na, vo.Nav_item{
			Name:         s.Title,
			Id: uint64(s.Aid),
			LoadOnDemand: true,
		})
	}
	c.JSON(http.StatusOK, na)
}


func (a *Api) NavLoad(c *gin.Context) {
	rawAid := c.Query("node")
	if rawAid == "" {
		c.JSON(http.StatusOK, []string{})
		return
	}
	aid, err := strconv.Atoi(rawAid)
	fmt.Println(aid)
	if err != nil {
		common.Sugar.Fatal(err)
		c.JSON(http.StatusOK, []string{})
		return
	}

	var articleList []model.Article
	common.NewDb.Where("p_aid = ?", aid).Find(&articleList)

	var na []vo.Nav_item
	for _, s := range articleList {
		na = append(na, vo.Nav_item{
			Name:         s.Title,
			Id: uint64(s.Aid),
			LoadOnDemand: true,
		})
	}
	if (len(na) <= 0) {
		 c.JSON(http.StatusOK, []string{})
	} else {
		c.JSON(http.StatusOK, na)
	}
}


func (a *Api) Resort(c *gin.Context) {

	c.JSON(http.StatusOK, []string{})
}



func (a *Api) Index(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		common.Sugar.Fatal(err)
	}
	page -= 1
	if page < 0 {
		page = 0
	}

	prev_page := page
	if prev_page < 1 {
		prev_page = 1
	}

	rpp := 20
	offset := page * rpp
	var blogListSlice []apiBlogList

	rows, err := common.DB.Query("Select aid, title from bk_article where publish_status = 1 order by aid desc limit ? offset ? ", &rpp, &offset)
	if err != nil {
		common.Sugar.Fatal(err)
	}
	defer common.CloseRowsDefer(rows)
	if rows != nil {
		var (
			aid   sql.NullString
			title sql.NullString
		)
		blogListSlice = make([]apiBlogList, 0) //Must be zero slice
		var aBlog apiBlogList
		for rows.Next() {
			err := rows.Scan(&aid, &title)
			if err != nil {
				common.Sugar.Fatal(err)
			}
			aBlog.Aid = aid.String
			aBlog.Title = title.String
			blogListSlice = append(blogListSlice, aBlog)
		}
		err = rows.Err()
		if err != nil {
			common.Sugar.Fatal(err)
		}
	}
	c.JSON(http.StatusOK, blogListSlice)
}

type apiBlogItem struct {
	Aid     string `form:"aid" json:"aid"  binding:"required"`
	Title   string `form:"title" json:"title"  binding:"required"`
	Content string `form:"content" json:"content"  binding:"required"`
}

func (a *Api) View(c *gin.Context) {
	aid, err := strconv.Atoi(c.Param("id"))
	fmt.Println(aid)
	if err != nil {
		common.Sugar.Fatal(err)
	}
	var b apiBlogItem

	rows, err := common.DB.Query("Select aid, title, content from bk_article where aid =  ? limit 1 ", &aid)
	if err != nil {
		common.Sugar.Fatal(err)
	}
	defer common.CloseRowsDefer(rows)
	if rows != nil {
		var (
			aid     sql.NullString
			title   sql.NullString
			content sql.NullString
		)
		for rows.Next() {
			err := rows.Scan(&aid, &title, &content)
			if err != nil {
				fmt.Println(err)
			}
			b.Aid = aid.String
			b.Title = title.String
			b.Content = content.String
		}
		fmt.Println(b)
		err = rows.Err()
		if err != nil {
			common.Sugar.Fatal(err)
		}
	}
	fmt.Println(b)
	c.JSON(http.StatusOK, b)
}
