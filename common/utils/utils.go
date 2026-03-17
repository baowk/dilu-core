package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// cryptoRandIntn 使用 crypto/rand 生成 [0, n) 范围的安全随机数
// 对常用的小范围值缓存 big.Int 避免每次分配
var smallBigInts [256]*big.Int

func init() {
	for i := range smallBigInts {
		smallBigInts[i] = big.NewInt(int64(i))
	}
}

func cryptoRandIntn(n int) int {
	var bound *big.Int
	if n > 0 && n < len(smallBigInts) {
		bound = smallBigInts[n]
	} else {
		bound = big.NewInt(int64(n))
	}
	val, err := rand.Int(rand.Reader, bound)
	if err != nil {
		return mathrand.Intn(n)
	}
	return int(val.Int64())
}

var characters string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 生成随机字符串（使用 crypto/rand）
func RandStringByLen(size int) string {
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		result[i] = characters[cryptoRandIntn(len(characters))]
	}
	return string(result)
}

// 生成随机数字字符串（使用 crypto/rand）
func RandNumberByLen(size int) string {
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		result[i] = characters[cryptoRandIntn(10)]
	}
	return string(result)
}

// 在数组中随机取一个
func RandFromArray(array *[]any) any {
	return (*array)[mathrand.Intn(len(*array))]
}

func RandFromArrayString(array []string) string {
	return array[mathrand.Intn(len(array))]
}

// 指定区间随机生成
func RandNumber(min, max float64, precision int) (float64, error) {
	r := mathrand.Float64()*(max-min) + min
	return RoundFloat(r, precision)
}

// Float 精度格式化
func RoundFloat(num float64, precision int) (float64, error) {
	roundedStr := strconv.FormatFloat(num, 'f', precision, 64)
	rounded, err := strconv.ParseFloat(roundedStr, 64)
	if err != nil {
		return 0, err
	}
	return rounded, nil
}

// 使用空接口实现数组去重
func Deduplicate(array []any) []any {
	encountered := make(map[any]bool)
	result := []any{}

	for _, item := range array {
		if !encountered[item] {
			encountered[item] = true
			result = append(result, item)
		}
	}
	return result
}

/**
 *	字符串素组去重
 */
func DeduplicateString(array []string) []string {
	encountered := make(map[string]bool)
	result := []string{}

	for _, item := range array {
		if !encountered[item] {
			encountered[item] = true
			result = append(result, item)
		}
	}
	return result
}

// string数组转接口数组
func StrToInterfaceArray(stringSlice []string) []any {
	interfaceSlice := make([]any, len(stringSlice))
	for i, v := range stringSlice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

// 结构体转string数组
func InterfaceToStrArray(interfaceArr []any) []string {
	stringArr := make([]string, len(interfaceArr))
	for i, v := range interfaceArr {
		stringArr[i] = fmt.Sprintf("%v", v)
	}
	return stringArr
}

// 结构体转map
func Struct2map(obj any) (data map[string]any, err error) {
	// 通过反射将结构体转换成map
	data = make(map[string]any)
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	for i := 0; i < objT.NumField(); i++ {
		fileName, ok := objT.Field(i).Tag.Lookup("json")
		if ok {
			data[fileName] = objV.Field(i).Interface()
		} else {
			data[objT.Field(i).Name] = objV.Field(i).Interface()
		}
	}
	return data, nil
}

// 生成无-的uuid
func GenUUid() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}

// 结构体转int
func GetInterfaceToInt(t1 any) int {
	switch t1.(type) {
	case uint:
		return int(t1.(uint))
	case int8:
		return int(t1.(int8))
	case uint8:
		return int(t1.(uint8))
	case int16:
		return int(t1.(int16))
	case uint16:
		return int(t1.(uint16))
	case int32:
		return int(t1.(int32))
	case uint32:
		return int(t1.(uint32))
	case int64:
		return int(t1.(int64))
	case uint64:
		return int(t1.(uint64))
	case float32:
		return int(t1.(float32))
	case float64:
		return int(t1.(float64))
	case string:
		t2, _ := strconv.Atoi(t1.(string))
		if t2 == 0 && len(t1.(string)) > 0 {
			f, _ := strconv.ParseFloat(t1.(string), 64)
			t2 = int(f)
		}
		return t2
	case nil:
		return 0
	case json.Number:
		t3, _ := t1.(json.Number).Int64()
		return int(t3)
	default:
		return t1.(int)
	}
}

// interface 转结构体
func InterfaceToStruct(i1 any, i2 *any) error {
	d, err := json.Marshal(i1)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, i2)
}

// string 数组转int数组
func String2Int(strArr []string) []int {
	res := make([]int, len(strArr))

	for index, val := range strArr {
		res[index], _ = strconv.Atoi(val)
	}
	return res
}

var baseChars = "pgabcdefhiuvwxyztjkmnqrs53926487VFGHLJKMCDABEUNPQRSTWXYZ"

// 把10进制转N进制 baseN最大56,24存小写字母,32为数字和小写大于32大小字符
func BaseDecimalToN(num int, baseN int) string {
	result := ""

	for num > 0 {
		remainder := num % baseN
		result = string(baseChars[remainder]) + result
		num = num / baseN
	}

	return result
}

// 把N进制转10进制 baseN最大56,24存小写字母,32为数字和小写大于32大小字符
func BaseNToDecimal(baseNum string, baseN int) int {
	result := 0
	power := 1
	for i := len(baseNum) - 1; i >= 0; i-- {
		char := string(baseNum[i])
		index := strings.Index(baseChars, char)
		result += index * power
		power *= baseN
	}
	return result
}

// MaskSensitiveInfo 对于字符串脱敏
// s 需要脱敏的字符串
// start 从第几位开始脱敏
// maskNumber 需要脱敏长度
// maskChars 掩饰字符串，替代需要脱敏处理的字符串
func MaskSensitiveInfo(s string, start int, maskNumber int, maskChars ...string) string {
	// 将字符串s的[start, end)区间用maskChar替换，并返回替换后的结果。
	maskChar := "*"
	if maskChars != nil {
		maskChar = maskChars[0]
	}
	// 处理起始位置超出边界的情况
	if start < 0 {
		start = 0
	}
	// 处理结束位置超出边界的情况
	end := start + maskNumber
	if end > len(s) {
		end = len(s)
	}
	return s[:start] + strings.Repeat(maskChar, end-start) + s[end:]
}

/**
 * int转比特数组
 */
func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

/**
 * 比特数组转int
 */
func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

/**
 * 驼峰转蛇形 snake string
 * @description XxYy to xx_yy , XxYY to xx_y_y
 * @date 2023/10/13
 * @param s 需要转换的字符串
 * @param allMode true XxYY to xx_y_y false XxYY to xx_yy
 * @return string
 **/
func SnakeCase(s string, allMode bool) string {
	num := len(s)
	data := make([]byte, 0, num*2)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if d >= 'A' && d <= 'Z' {
			if i > 0 {
				if allMode {
					data = append(data, '_', d+32)
				} else {
					if s[i-1] >= 'A' && s[i-1] <= 'Z' {
						data = append(data, d+32)
					} else {
						data = append(data, '_', d+32)
					}
				}
			} else {
				data = append(data, d+32)
			}
		} else {
			data = append(data, d)
		}
	}
	//ToLower把大写字母统一转小写
	return string(data[:])
}

/**
 * 蛇形转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY
 * @date 2023/10/13
 * @param s要转换的字符串
 * @return string
 **/
func CamelCase(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
