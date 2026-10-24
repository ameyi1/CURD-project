package utils

import (
	"context"
	"errors"
	"fmt"
	"gin-study/config"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Claims struct {
	UserID uint   `json:"userid"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (accesstoken, refreashtoken string, err error) {

	//短期token
	//结构体
	accessclaims := &Claims{
		UserID: userID,
		Type:   "Accesstoken",
		//忘了
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GlobalConfig.Jwt.Expire) * time.Hour)),
		},
	}
	//创建token
	accessjwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessclaims)
	//加密token,注意遗忘了，是对[]byte(config.GlobalConfig.Jwt.Secret))注意这个[]byte()
	accesstoken, err = accessjwt.SignedString([]byte(config.GlobalConfig.Jwt.Secret))
	if err != nil {
		return "", "", err
	}
	//生成长期token
	refreashclaims := &Claims{
		UserID: userID,
		Type:   "Refrashtoken",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GlobalConfig.Jwt.Expire) * 7 * time.Hour)),
		}}
	refreashjwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreashclaims)
	refreashtoken, err = refreashjwt.SignedString([]byte(config.GlobalConfig.Jwt.Secret))
	//加上redis,把生成的refreshtoken存到redis里面
	ctx := context.Background()
	key := "jwt:refreashtoken" + refreashtoken
	ttl := time.Duration(config.GlobalConfig.Jwt.Expire) * 7 * time.Hour
	if err := config.RedisClient.Set(ctx, key, userID, ttl).Err(); err != nil { //主要是这句，config.RedisClient.Set
		Log.Error("保存refreashtoken到redis失败", zap.Error(err))
		return "", "", err
	}

	if err != nil {
		return "", "", err
	}
	return accesstoken, refreashtoken, nil
} //解析短期token
func ParseToken(accesstoken string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(accesstoken, claims, func(t *jwt.Token) (interface{}, error) {
		//只接受 HMAC 系列算法，防止alg.none攻击
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", t.Header["alg"])
		}
		return []byte(config.GlobalConfig.Jwt.Secret), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("token无效")
	}
	if claims.Type != "Accesstoken" {
		return 0, errors.New("token类型错误")
	}
	return claims.UserID, nil
}

// 解析长期token
func ParseRefreashToken(refreashtoken string) (uint, error) {
	claims := &Claims{}
	refreshtoken, err := jwt.ParseWithClaims(refreashtoken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.Jwt.Secret), nil
	})
	if err != nil {
		return 0, err
	}
	if !refreshtoken.Valid {
		return 0, errors.New("token无效")
	}
	if claims.Type != "Refreashtoken" {
		return 0, errors.New("token类型错误")
	}
	return claims.UserID, nil

}

// 根据refreashtoken生成新token,新知识点
func Refrashtoken(c *gin.Context) {
	var req struct {
		Refreashtoken string `json:"refreashtoken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误", CodeInvalidParam)
		return
	}
	//接入redis来检查这个refreashtoken
	ctx := context.Background()
	key := "jwt:refreashtoken" + req.Refreashtoken
	_, err := config.RedisClient.Get(ctx, key).Result()
	//如果没查到东西
	if err == redis.Nil {
		Fail(c, "token过期或者失效了", CodeUnauthorized)
		return
	} else if err != nil { //或者服务出错了
		Log.Error("检查redis失败", zap.Error(err))
		Fail(c, "检查redis失败可能是服务出现问题", CodeServerError)
		return
	}

	userId, err2 := ParseRefreashToken(req.Refreashtoken)
	if err2 != nil {
		Fail(c, "refresh token无效", CodeUnauthorized)
		return
	}
	newaccesstoken, newrefreashtoken, err3 := GenerateToken(userId)
	if err3 != nil {
		Fail(c, "生成新token失败", CodeServerError)
		return
	}
	Success(c, gin.H{
		"access_token":  newaccesstoken,
		"refresh_token": newrefreashtoken,
		"user":          userId,
	}, "登录成功")
}

// 再添加一个推出主动删除refreashtoken的函数
func Logout(c *gin.Context) {
	var req struct {
		Refreashtoken string `json:"refreashtoken" binding:"required"`
	}
	key := "jwt:refreashtoken:" + req.Refreashtoken
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误", CodeInvalidParam)
		return
	}
	ctx := context.Background()
	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		Log.Error("退出登录删除 Redis Key 失败", zap.Error(err))
		Fail(c, "服务器内部错误", CodeServerError)
		return
	}
	Success(c, nil, "退出成功")
}

// 这是加自动改变全局变量的东西的，但懒得改了，而且那个方便复习
/*type Claims struct {
	UserID uint `json:"userid"`
	jwt.RegisteredClaims
}

// 生成Token
func GenerateToken(userID uint) (string, error) {
	expireHour := config.GlobalConfig.Jwt.Expire
	expirationTime := time.Now().Add(time.Duration(expireHour) * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.Jwt.Secret))
	return tokenString, err
}

// 解析Token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.Jwt.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

//=========================================================================================
// 1先创造密钥
/*var jwtKey = []byte("Ameyi")


// 2定义身份证的信息，便于后续验证
type Claims struct {
	UserID uint `json:"userid"`
	jwt.RegisteredClaims
}

// 3重要   生产token函数
func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) //这是生产过期时间
	//然后开始填这个身份证的信息了有点不懂为什么加&
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), //这一步是固定的，转译为jwt能识别的数字
		},
	}
	//inportant，创建token对象，带着身份证信息，有点不懂这个整体什么时候用
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //创建了

	//最看不懂的一集，加密生成最终token
	//加密密钥，是乱码的那种看不懂的东西
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// 4解析token函数，验证身份证是否有效，这个也不懂
func ParseToken(tokenString string) (*Claims, error) {
	// 空结构体：用来接收解析后的身份证信息（对标 ShouldBindJSON）
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { //这就开始看看你的密码是不是伪造的或者过期的，但无法识别正误
		return jwtKey, nil
	})
	// ===================== 验证Token是否有效 =====================
	// 如果报错 或 Token无效/过期
	if err != nil || !token.Valid {
		return nil, err // 返回错误
	}

	// 验证成功：返回身份证信息（里面有UserID）
	return claims, nil
}
//================================================================================================================
/*type Claims struct{
	UserID uint
	jwt.RegisteredClaims
}
var jwtKey =[]byte("Ameyi")
//func GenerateToken(/*这部分没写出来*/ /*userID uint)(string,error){
claims:=&Claims{//一定要&
UserID: userID,
RegisteredClaims: jwt.RegisteredClaims{
	ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 *time.Hour)),
},
}
//关键来了生产token（密钥）
token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)//这里写错过，只有一个返回值
tokenString,err:=token.SignedString(jwtKey)//这里写错了，是token.SignedString不是其他的
return tokenString ,err
}
func ParseToken(tokenString string)(*Claims/*这个写错了，无法理解逻辑*/ /*,error){
var claims Claims
token,err:=jwt.ParseWithClaims(tokenString,&claims,func(t *jwt.Token) (interface{}, error) {
	return jwtKey,nil
})//这一部分一点都没记清，纯粹的乱写，也是不理解逻辑
if err!=nil||!token.Valid{
	return nil,err//这也记得不是一啥
}
return &claims/*这也记错了，无法理解为什么写这个*/ /*,err
}*/
/*var jwtkey=[]byte("Ameyi")
var RedisClient *redis.Client
const AccessExpirate=time.Hour*1
const RefreashExpirate=time.Hour*24*7
type Claims struct{
	UserID uint`json:"userid"`
	jwt.RegisteredClaims
}
func GenerateToken(userID uint)(accessExpirate,refreashExpirate string,err error){
//生成短期token
accessclaims:=&Claims{
	UserID: userID,
	RegisteredClaims: jwt.RegisteredClaims{
ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessExpirate)),
IssuedAt: jwt.NewNumericDate(time.Now()),
	},

}
accesstoken:=jwt.NewWithClaims(jwt.SigningMethodES256,accessclaims)
accesstokenString,err:=accesstoken.SignedString(jwtkey)
if err!=nil{
return"", "", err
}
//生成长期token
refreashclaims:=&Claims{
	UserID: userID,
	RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreashExpirate)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	},}
refreashtoken:=jwt.NewWithClaims(jwt.SigningMethodES256,refreashclaims)
refreashtokenString,err:=refreashtoken.SignedString(jwtkey)
if err!=nil{
	return"","",err
}
}
func ParseToken(tokenStr string)(userId uint ,err error){
//为redis搭建平台，验证是否再黑名单上
ctx:=context.Background()
exists,err:=RedisClient.Exists(ctx,"token_black"+tokenStr).Result()
if err!=nil{
	return 0,err
}
if exists>0{
	return 0, errors.New("token已失效,重新登陆")
}
//解析
tokenobj,err:=jwt.ParseWithClaims(tokenStr,&Claims{},func(t *jwt.Token) (interface{}, error) {
	return jwtkey,nil
})
if err!=nil{
	return 0,err
}
claims,ok:=tokenobj.Claims.(*Claims)
if !ok||!tokenobj.Valid{
return 0,errors.New("错误凭证")
}
return claims.UserID,nil
}*/
