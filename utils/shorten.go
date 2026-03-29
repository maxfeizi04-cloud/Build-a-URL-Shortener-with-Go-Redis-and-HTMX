package utils

import (
	"encoding/base64"
	"fmt"
	"time"
)

func GetShortCode() string {
	fmt.Println("Shortening URL")
	ts := time.Now().UnixNano()
	fmt.Println("Timestamp:", ts)

	// we convert the timestamp to byte slice and then encode it to base64 string
	// 我们将时间戳转换为字节切片，然后将其编码为 base64 字符串
	ts_bytes := []byte(fmt.Sprintf("%d", ts))
	key := base64.StdEncoding.EncodeToString(ts_bytes)
	fmt.Println("Key:", key)

	// we remove the last two chars since they are usuall always equal sings (==)
	// 我们移除最后两个字符，因为它们通常总是相等的符号（==)
	key = key[:len(key)-2]
	fmt.Println("Key:", key)
	
	// we return the last chars after 16 chars, these are almost always different
	// 我们返回最后的字符，超过16个字符后，这些字符几乎总是不同的
	return key[16:]
}
