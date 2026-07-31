// 包名：工具包，全局调用
package utils

// 导入加密库：官方安全库

// ====================== 函数1：密码加密 ======================
// 输入：明文密码  输出：加密后的乱码  错误
import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	//GenerateFromPassword
	//= 给密码上锁 → 生成乱码
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashBytes), err
}

// ====================== 函数2：密码校验 ======================
// 输入：用户输入的明文密码、数据库存的加密乱码
// 输出：true=密码正确  false=错误
func ComparePassword(password, hash string) bool {
	// 陌生函数2：bcrypt.CompareHashAndPassword
	// 作用：拿明文 和 加密乱码 比对
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
