package myutils

import (
	"fmt"
	"os"
	_ "testing"
)

func TestDes() {
	var key string
	var deserr error
	// text := "root"
	// key := "12345678"
	// iv := key
	// fmt.Println(DesEncrypt(text, key, iv))

	// var s1 string
	// // var s2 string
	// s1 = "e.X5T@h" + string([]byte{27})
	// fmt.Println(DesEncrypt(text, "e.X5T@h" + string([]byte{27}), s1))
	// // s2, _ = DesEncrypt(text, s1, s1)

	if _, err := os.Stat("HISTORY"); err == nil || os.IsExist(err) {
		data_encrypt, err := os.ReadFile("HISTORY")
		if err != nil {
			fmt.Println("读取HISTORY文件失败，请检查文件权限")
		}
		key, deserr = DesDecrypt(string(data_encrypt), "e.X5T@h"+string([]byte{27}), "")
		if deserr != nil {
			fmt.Println("解密HISTORY文件失败，检查部署步骤，重置此文件")
		}
		fmt.Println(key)
	} else {
		if err != nil {
			fmt.Println(err)
			fmt.Println("HISTORY文件找不到，请检查部署步骤")
		}
	}

	fmt.Print("||||||")
	//fmt.Println(DesDecrypt("2afb686b6e2a2e03b85612fd7d41ac94", key, ""))
	fmt.Println(DesEncrypt("eo!9Fs@=5Av", key, ""))
	fmt.Print("||||||")
	fmt.Println(DesEncrypt("root", key, ""))

}
