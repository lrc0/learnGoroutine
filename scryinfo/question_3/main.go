package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := ReadAndWriteData("./test1.txt", []byte("aaaaaaaa"))
	fmt.Println("err: ", err)
	fmt.Println("data: ", string(data))
}

//ReadAndWriteData read and write about the file
func ReadAndWriteData(filename string, data []byte) ([]byte, error) {
	if len(data) < 8 {
		log.Println("数据长度小于8")
		return nil, fmt.Errorf("数据长度小于8")
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	//文件不存在
	if os.IsNotExist(err) {
		_, err = f.Write(data)
		if err != nil {
			log.Println("err: ", err)
			return nil, err
		}
		return nil, nil
	}
	if err != nil {
		log.Println("err: ", err.Error())
		return nil, err
	}

	//文件存在
	fInfo, err := f.Stat()
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}

	buf := make([]byte, fInfo.Size())
	_, err = f.Read(buf)
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}

	err = f.Truncate(0)
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}

	_, err = f.WriteAt(data, 0)
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}

	return buf, nil
}
