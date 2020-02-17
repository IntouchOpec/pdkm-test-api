package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	AfterTable()
	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	r.GET("/", WelcomHandler)
	user := r.Group("/user")
	user.GET("/", GetUserHanlder)
	user.GET("/:id", GetUserDetailHanlder)
	user.POST("/", CreateUserHanlder)
	user.PUT("/:id", UpdateUserHanlder)
	user.PATCH("/:id", UpdateUserHanlder)
	user.DELETE("/:id", DeleteUserHanlder)
	r.Run(":8000")
}

type Respone struct {
	Data  interface{}
	Page  int
	Total int
}

type Context struct {
	*gin.Context
	DB *sql.DB
}

func WelcomHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "hello")
}

func GetUserHanlder(c *gin.Context) {
	db := DB()
	var wh string
	u, err := url.Parse(fmt.Sprintf("?%s", c.Request.URL.RawQuery))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fild := []string{"id", "first_name", "last_name", "email", "gender", "age", "min_age", "max_age"}
	for key, value := range q {
		var chekFild bool
		for i := 0; i < len(fild); i++ {
			if fild[i] == key {
				chekFild = true
				break
			}
		}
		if !chekFild {
			c.JSON(http.StatusBadRequest, gin.H{"message": "fild not mach"})
			return
		}
		if key == "id" {
			wh += " " + key + " LIKE" + "'" + value[0] + "' and"
			continue
		}
		if key == "min_age" {
			wh += " " + "age" + " >=" + "'" + value[0] + "' and"
			continue
		}
		if key == "max_age" {
			wh += " " + "age" + " <=" + "'" + value[0] + "' and"
			continue
		}
		wh += " " + key + " LIKE" + "'%" + value[0] + "%' and"
	}
	if wh != "" {
		wh = wh[:len(wh)-3]
	}
	us, err := ListUser(db, wh)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, us)
	return
}

func GetUserDetailHanlder(c *gin.Context) {
	id := c.Param("id")
	db := DB()
	var u User
	if err := u.Detail(db, id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func CreateUserHanlder(c *gin.Context) {
	var u User
	if err := c.Bind(&u); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	db := DB()
	if err := u.Create(db); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, u)
}

func UpdateUserHanlder(c *gin.Context) {
	var u User
	id := c.Param("id")
	if err := c.Bind(&u); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	db := DB()
	if err := u.Update(db, id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, u)
	return
}

func DeleteUserHanlder(c *gin.Context) {
	var u User
	id := c.Param("id")
	db := DB()
	if err := u.Delete(db, id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, u)
	return
}
