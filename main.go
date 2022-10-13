package main

import (
	"github.com/gin-gonic/gin"
	//"gorm.io/driver/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password string `gorm:"size:255;not null"`
}


func main() {

	db := InitDB()
	defer db.DB()

	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		password := ctx.PostForm("password")
		telephone := ctx.PostForm("telephone")
		//数据验证
		if len(telephone) != 11{
			ctx.JSON(http.StatusUnprocessableEntity,map[string]interface{}{
				"code":422,
				"msg":"手机号必须为11位",
			})
			return
		}
		if len(password) < 6{
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"code":422,
				"msg":"密码不能少于6位",
			})
			return
		}
		//如果名称没有传，给一个10位的字符串
		if len(name) == 0 {
			name = RandomString(10)
		}

		log.Println(name,password,telephone)
		//判断手机号是否存在
		if isTelephoneExist(db,telephone){
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"code":422,
				"msg":"手机号已经存在",
			})
			return
		}

		//创建用户
		newUser := User{
			Name:name,
			Telephone: telephone,
			Password: password,
		}
		db.Create(&newUser)

		//返回结果
		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	r.Run() // 监听并在 0.0.0.0	:8080 上启动服务
}

func RandomString(n int) string {
	var letters = []byte("djisdjsdjkcmiwjiowpqpq")
	result := make([]byte,n)
	rand.Seed(time.Now().Unix()) //给定一个初始的随机种子

	for i := range result{
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func isTelephoneExist(db *gorm.DB,telephone string) bool {
	var user User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID != 0{
		return true
	}
	return false

}

func InitDB() *gorm.DB {
	//driverName := "mysql"
	//host := "localhost"
	//port :="3306"
	//database := "ginessential"
	//username := "root"
	//password := "root"
	//charset := "utf8"
	dsn := "root:root@tcp(localhost:3306)/ginessential?charset=utf8mb4&parseTime=True&loc=Local"
	//args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
	//	username,
	//	password,
	//	host,
	//	port,
	//	database,
	//	charset,
	//	)
	db, err := gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	//自动创建数据表
	db.AutoMigrate(&User{})
	//db.Migrator().CreateTable(&User{})
	return db
}