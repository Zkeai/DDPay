package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Zkeai/DDPay/common/logger"
)

func BigFloat2BigInt(f *big.Float) *big.Int {
	scale := new(big.Float).SetFloat64(1e9)
	scale.Mul(f, scale)

	i := new(big.Int)
	scale.Int(i)
	return i

}
func FloatToString(value float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, value)
}
func BigFloat2Int(f *big.Float) int64 {

	scale := new(big.Float).SetFloat64(1e9)
	f.Mul(f, scale)

	// 转换为 int64 类型
	i, _ := f.Int64()
	return i
}

//跨域

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin) // ✅ 动态设置源，避免 "null"
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, credentials")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// FileExist 判断文件是否存在及是否有权限访问
func FileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if os.IsPermission(err) {
		return false
	}

	return true
}

// CreateDirIfNotExists 检测目录是否存在
func CreateDirIfNotExists(path ...string) {
	for _, value := range path {
		if FileExist(value) {
			continue
		}
		err := os.Mkdir(value, 0755)
		if err != nil {
			logger.Error(fmt.Sprintf("创建目录失败:%s", err.Error()))
		}
	}
}
func TimeToUTC8(timeStamp int64) string {
	// 示例时间戳（秒级）
	timestamp := int64(timeStamp)

	// 将时间戳转换为 time.Time 类型
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1000000
	t := time.Unix(seconds, nanoseconds)
	// 定义北京时间的时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Error(err)
		return ""
	}

	// 将时间转换为北京时间
	beijingTime := t.In(location)
	return beijingTime.Format("2006-01-02 15:04:05")
}

// formatValue 函数用于格式化资产值
func formatValue(value string) string {
	if value == "0.00000000" {
		return "0"
	}
	return value
}

func GenerateFormattedBinanceText(assets []struct {
	Asset  string `json:"a"`
	Free   string `json:"f"`
	Locked string `json:"l"`
}) string {
	var sb strings.Builder
	for _, asset := range assets {
		free := formatValue(asset.Free)
		locked := formatValue(asset.Locked)
		sb.WriteString(fmt.Sprintf("<font color=\"blue\">%s｜可用：%s｜锁定：%s</font><br>", asset.Asset, free, locked))
	}
	return sb.String()
}

// DDPaySign 用于对参数进行 MD5 签名
func DDPaySign(params map[string]string, signKey string) string {
	// 获取所有的 key 并排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signParts []string

	for _, k := range keys {
		val := params[k]
		if val == "" || k == "signature" {
			continue
		}
		signParts = append(signParts, k+"="+val)
	}

	// 拼接字符串后加上 signKey
	signStr := strings.Join(signParts, "&") + signKey

	// 计算 MD5 签名
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}

func RestoreBody(data []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewBuffer(data))
}
