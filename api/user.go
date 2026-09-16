package api

/*

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	// 导入自己的包
	"gin-study/config"
	"gin-study/model"
	"gin-study/utils"
)

// RegisterUserRoutes 注册用户相关路由
func registerUserRoutes(r *gin.RouterGroup) {

	// 增
	r.POST("/create", func(c *gin.Context) {
		var user model.User            //定义user结构体，方便传
		err := c.ShouldBindJSON(&user) //这是接收前端发的信息的
		if err != nil {
			utils.Log.Error("数据错误：", zap.Error(err))
			utils.Fail(c, "数据错误："+err.Error())
			return
		}
		encryptPwd, err := utils.EncryptPassword(user.Password) //先把前端发的密码加密为一个副本
		if err != nil {
			utils.Log.Error("加密失败：", zap.Error(err))
			utils.Fail(c, "加密失败")
			return
		}
		user.Password = encryptPwd                               //把副本赋值给原本的密码，这时候user。password已经变为加密状态
		txErr := config.DB.Transaction(func(tx *gorm.DB) error { //先声明事务在对数据库的操作中全部替换为tx
			if err := tx.Create(&user).Error; err != nil { //创建
				return err
			}
			DefaultArticle := model.Article{
				Title:   "这是欢迎文章",
				Context: "这是欢迎页",
				UserID:  user.ID,
			}
			if err := tx.Create(&DefaultArticle).Error; err != nil {
				return err
			}
			return nil
		})
		//=====================================事务完成====================
		if txErr != nil {
			utils.Log.Error("创建用户+默认文章事务执行失败", zap.Error(txErr), zap.String("username", user.Username))
			utils.Fail(c, "用户创建失败")
			return
		}

		utils.Success(c, user, "已完成创建")
		utils.Log.Info("创建成功", zap.String("username", user.Username), zap.Uint("userID", user.ID))
	})

	// 查单个
	r.GET("/search/:id", func(c *gin.Context) {
		var user model.User
		is := c.Param("id")
		id, err := strconv.Atoi(is)
		if err != nil {
			utils.Log.Warn("ID必须是纯数字", zap.String("id", is), zap.Error(err))
			utils.Fail(c, "ID必须是纯数字")
			return
		}
		result := config.DB.First(&user, id)
		if result.Error != nil {
			utils.Log.Error("查询错误", zap.Error(result.Error))
			utils.Fail(c, "查询错误"+result.Error.Error())
			return
		}
		utils.Log.Info("成功查询", zap.Int("id", id))
		utils.Success(c, user, "成功查询")
	})

	// 分页查所有
	r.GET("/list", func(c *gin.Context) {
		var list []model.User
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		offset := (page - 1) * size
		result := config.DB.Offset(offset).Limit(size).Find(&list)
		if result.Error != nil {
			utils.Log.Error("查询列表失败：", zap.Error(result.Error))
			utils.Fail(c, "查询列表失败")
			return
		}
		utils.Log.Info("成功查询列表", zap.Int("page", page), zap.Int("size", size))
		utils.Success(c, list, "查询列表成功")
	})

	// 删
	//好像缺了个查软删，懒得写了，窗帘updated_at
	r.DELETE("/delete/:id", func(c *gin.Context) {
		is := c.Param("id")
		id, err := strconv.Atoi(is)
		if err != nil {
			utils.Log.Warn("ID必须是纯数字", zap.String("id", is), zap.Error(err))
			utils.Fail(c, "ID必须是纯数字")
			return
		}
		result := config.DB.Delete(&model.User{}, id)
		if result.Error != nil {
			utils.Log.Error("删除错误", zap.Error(result.Error))
			utils.Fail(c, "删除错误"+result.Error.Error())
			return
		}
		if result.RowsAffected == 0 {
			utils.Log.Warn("删除用户:目标ID不存在", zap.Int("id", id))
			utils.Fail(c, "用户不存在")
			return
		}
		utils.Log.Info("成功删除", zap.Int("id", id))
		utils.Success(c, nil, "成功删除")
	})

	// 改
	r.PUT("/update/:id", func(c *gin.Context) {
		var user model.User
		is := c.Param("id")
		id, err := strconv.Atoi(is)
		if err != nil {
			utils.Log.Warn("ID必须是纯数字", zap.String("id", is), zap.Error(err))
			utils.Fail(c, "ID必须是纯数字")
			return
		}
		err1 := c.ShouldBindJSON(&user)
		if err1 != nil {
			utils.Log.Error("参数错误：", zap.Error(err1))
			utils.Fail(c, "参数错误："+err1.Error())
			return
		}
		Result := config.DB.Model(&model.User{}).Where("id = ?", id).Select("username", "nickname").Updates(user)
		if Result.Error != nil {
			utils.Log.Error("修改错误：", zap.Error(Result.Error))
			utils.Fail(c, "修改错误："+Result.Error.Error())
			return
		}
		utils.Log.Info("修改成功", zap.Int("id", id))
		utils.Success(c, user, "修改完成")
	})
	//校验密码
	r.POST("/login", func(c *gin.Context) {
		//1先定义一个短暂的结构体接数据
		var loginUser struct { //这里注意是var而非type
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		//2绑定shouldbindwith
		err := c.ShouldBindJSON(&loginUser)
		if err != nil {
			utils.Log.Error("参数错误：", zap.Error(err))
			utils.Fail(c, "参数错误："+err.Error())
			return
		}
		//3验证密码对不对，对应的是查
		var user model.User
		result := config.DB.Where("username = ? ", loginUser.Username).First(&user)
		if result.Error != nil {
			utils.Log.Warn("用户错误，可能不存在：", zap.String("loginUser_Username", loginUser.Username), zap.Error(result.Error))
			utils.Fail(c, "用户错误，可能不存在")
			return
		}
		isOk := utils.ComparePassword(loginUser.Password, user.Password) //这次判断正误
		if !isOk {
			utils.Fail(c, "账号或密码错误")
			return
		}
		//4开始引用generatetoken函数了，这里注意只有user这单个用户信息的变量，但之后可能会有多个注意，而且只是判断这个合法不合法无法判断正误
		token, refreashtoken, err := utils.GenerateToken(user.ID)
		if err != nil {
			utils.Log.Error("生成Token失败:", zap.Error(err))
			utils.Fail(c, "生成Token失败")
			return
		}

		//5返回Token给前端

		utils.Success(c, gin.H{"token": token, "refreashtoken": refreashtoken, "user": user}, "登录成功")
		utils.Log.Info("登录成功", zap.String("username", user.Username), zap.Uint("userID", user.ID))
	})
}

func Int64(s string) int64 {
	utils.Log.Panic("unimplemented")
	return 0
}
func RegisterArticleRoutes(r *gin.RouterGroup) {

	// 1. 发布文章（必须登录）
	r.POST("/articles/create", func(c *gin.Context) {
		// 从JWT中间件拿用户ID
		userID, _ := c.Get("userID")
		uid := userID.(uint)

		var article model.Article
		err := c.ShouldBindJSON(&article)
		if err != nil {
			utils.Log.Error("参数错误:", zap.Error(err))
			utils.Fail(c, "参数错误："+err.Error())
			return
		}

		// 绑定当前登录用户ID
		article.UserID = uid
		config.DB.Create(&article)
		utils.Success(c, article, "发布文章成功")
		utils.Log.Info("发布文章成功", zap.String("title", article.Title), zap.Uint("userID", article.UserID), zap.Uint("articleID", article.ID))
	})

	// 2. 查询单个用户 + 他的所有文章（一对多核心！）
	r.GET("/user/:id/articles", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Log.Warn("ID格式错误", zap.String("id", idStr), zap.Error(err))
			utils.Fail(c, "ID格式错误")
			return
		}
		var user model.User
		// 🔥 核心：Preload 加载关联文章，仅此一行实现一对多
		result := config.DB.Preload("Articles").First(&user, id)

		if result.Error != nil {
			utils.Fail(c, "查询失败："+result.Error.Error())
			utils.Log.Error("查询错误:", zap.Error(result.Error))
			return
		}

		utils.Success(c, user, "查询用户+文章成功")
		utils.Log.Info("查询用户+文章成功", zap.Int("id", id))
	})

	// 3. 查询所有文章
	r.GET("/list/articles", func(c *gin.Context) {
		page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, err2 := strconv.Atoi(c.DefaultQuery("size", "10"))
		searchTitle := c.Query("title")
		if err1 != nil || err2 != nil {
			utils.Log.Error("页码参数匹配错误:", zap.Error(err1), zap.Error(err2))
			utils.Fail(c, "页码参数匹配错误")
			return
		}
		if page < 1 {
			page = 1
		}
		if size > 50 || size < 1 {
			size = 10
		}

		var total int64 //重要
		offset := (page - 1) * size
		db := config.DB
		if searchTitle != "" {
			db = db.Where("title LIKE ?", "%"+searchTitle+"%")
		}
		db.Model(&model.Article{}).Count(&total)
		var list []model.Article
		result := db.Offset(offset).Limit(size).Find(&list)
		if result.Error != nil {
			utils.Log.Error("查询失败:", zap.Error(result.Error))
			utils.Fail(c, "查询失败")
			return
		}
		totalPage := (total + int64(size) - 1) / int64(size)
		resData := gin.H{
			"list": list, // 单个用户
			"page_info": gin.H{
				"current":  page,
				"size":     size,
				"total":    total,     // 文章总数
				"page_num": totalPage, // 文章总页数
			},
		}

		utils.Success(c, resData, "查询文章列表成功")
		utils.Log.Info("查询文章列表成功", zap.Int("page", page), zap.Int("size", size), zap.String("searchTitle", searchTitle))
	})

}
*/
