package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Artichel struct {
	gorm.Model
	Title string
	Slug  string `gorm:"unique_index"`
	Desc  string `sql:"type:text;"`
}

var DB *gorm.DB

func Bellow() {
	var err error
	DB, err = gorm.Open("mysql", "root:@(localhost)/learngin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer DB.Close()

	DB.AutoMigrate(&Artichel{})

	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		artikel := v1.Group("/artikel")
		{
			artikel.GET("/", getHome)
			artikel.GET("/article/:slug", getArticle)
			artikel.POST("/articles", postArticle)
		}
	}

	router.Run()
}

func getHome(c *gin.Context) {
	items := []Artichel{}

	// Get all records
	DB.Find(&items)
	//// SELECT * FROM users;

	c.JSON(200, gin.H{
		"status": "Berhasil ke halaman home",
		"data":   items,
	})
}

func getArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item Artichel

	// SELECT * FROM TABLE WHERE SLUG = "SLUG"
	if DB.First(&item, "slug = ?", slug).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  "Elor",
			"message": "Data tidak ditemukan",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "Berhasil ges",
		"message": "Selamat ya kamu berhasil sampai tahap ini",
		"data":    item,
	})
}

func postArticle(c *gin.Context) {
	item := Artichel{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	// Kalau slugnya sama maka random generate slug
	// Cek database apakah sudah ada slug yang sama
	// judul-pertama
	// judul-pertama-stringrandom

	DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "Berhasil post data",
		"data":   item,
	})
}
