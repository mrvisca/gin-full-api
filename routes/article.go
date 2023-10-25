package routes

import (
	"gin-full-api/config"
	"gin-full-api/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func GetHome(c *gin.Context) {
	items := []models.Artichel{}

	// Get all records
	config.DB.Find(&items)
	//// SELECT * FROM users;

	c.JSON(200, gin.H{
		"status": "Berhasil ke halaman home",
		"data":   items,
	})
}

func GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item models.Artichel

	// SELECT * FROM TABLE WHERE SLUG = "SLUG"
	if config.DB.First(&item, "slug = ?", slug).RecordNotFound() {
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

func PostArticle(c *gin.Context) {
	var olditem models.Artichel
	slug := slug.Make(c.PostForm("title"))
	// SELECT * FROM TABLE WHERE SLUG = "SLUG"
	if !config.DB.First(&olditem, "slug = ?", slug).RecordNotFound() {
		// Penambahan Generate string baru pada slug
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	item := models.Artichel{
		Title:  c.PostForm("title"),
		Desc:   c.PostForm("desc"),
		Tag:    c.PostForm("tag"),
		Slug:   slug,
		UserID: uint(c.MustGet("jwt_user_id").(float64)), // Mengambil ke yang di oper oleh autentikasi
	}

	// Kalau slugnya sama maka random generate slug
	// Cek database apakah sudah ada slug yang sama
	// judul-pertama
	// judul-pertama-stringrandom

	config.DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "Berhasil post data",
		"data":   item,
	})
}

func GetArticleByTag(c *gin.Context) {
	tag := c.Param("tag")
	items := []models.Artichel{}

	// config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&items)
	// Lakukan pencarian data
	if result := config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&items); result.Error != nil {
		// Terjadi kesalahan saat menjalankan query
		c.JSON(404, gin.H{
			"status":  "Elor",
			"message": "Terjadi kesalahan saat menjalankan query",
		})
	} else if result.RowsAffected == 0 {
		// Tidak ada data yang ditemukan
		c.JSON(404, gin.H{
			"status":  "Elor",
			"message": "Aku tidak menemukan datanya bosque",
		})
	} else {
		c.JSON(200, gin.H{"data": items})
	}
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	var item models.Artichel

	if config.DB.First(&item, "id = ?", id).RecordNotFound() {
		c.JSON(404, gin.H{"status": "Elor", "message": "Elor, data tidak ditemukan"})
		c.Abort()
		return
	}

	// Filter update sesuai dengan akses user_id
	if uint(c.MustGet("jwt_user_id").(float64)) != item.UserID {
		c.JSON(403, gin.H{"status": "Elor", "message": "This data is forbidden"})
		c.Abort()
		return
	}

	// Jika data ditemukan maka akan dilakukan pengupdatean data
	config.DB.Model(&item).Where("id = ?", id).Updates(models.Artichel{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Tag:   c.PostForm("tag"),
	})

	c.JSON(200, gin.H{
		"status": "Berhasil update data",
		"data":   item,
	})
}
