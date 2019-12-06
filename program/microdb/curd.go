package microdb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/micro-kit/micro-cli/program/common"
)

// 对数据文件的操作

// GetDbFilePath 获取存储配置db文件路径
func (db *MicroDB) GetDbFilePath() string {
	return strings.TrimRight(db.DbPath, "/") + "/db.json"
}

// SaveToFile 保存数据到db文件，覆盖
func (db *MicroDB) SaveToFile() (err error) {
	// 判断kit db信息是否正确
	if db.PackageName == "" {
		return errors.New("package 包名不能为空")
	}
	if db.Service.Name == "" {
		return errors.New("service 名为空，可能是未定义rpc service")
	}
	if db.Service.Rpcs == nil || len(db.Service.Rpcs) == 0 {
		return errors.New("rpc 未定义任何rpc服务")
	}
	// 查看db文件是否存在，如果存在则取出project_info信息，保存到新对象

	isExist, err := common.PathExists(db.GetDbFilePath())
	if err != nil {
		return err
	}
	if isExist == true {
		oldDb := new(MicroDB)
		oldBody, err := ioutil.ReadFile(db.GetDbFilePath())
		if err != nil {
			return err
		}
		err = json.Unmarshal(oldBody, oldDb)
		if err != nil {
			return err
		}
		db.ProjectInfo = oldDb.ProjectInfo
	}

	err = db.SaveToFileNotCheck()
	return nil
}

// SaveToFileNotCheck 保存数据文件到文件
func (db *MicroDB) SaveToFileNotCheck() error {
	log.Println("dbPath", db.DbPath)
	// 判断目录是否存在
	isExist, err := common.PathExists(db.DbPath)
	if err != nil {
		return err
	}
	if isExist == false {
		err = os.MkdirAll(db.DbPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	body, err := json.Marshal(db)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(db.GetDbFilePath(), body, 0666)
	if err != nil {
		return err
	}
	return nil
}

// LoadToFile 从db文件加载
func (db *MicroDB) LoadToFile() (*MicroDB, error) {
	// log.Println(db.GetDbFilePath())
	isExist, err := common.PathExists(db.GetDbFilePath())
	if err != nil {
		return nil, err
	}
	if isExist == false {
		return nil, errors.New("数据文件不存在，请先执行 init指令")
	}
	body, err := ioutil.ReadFile(db.GetDbFilePath())
	if err != nil {
		return nil, err
	}
	// 解析db
	err = json.Unmarshal(body, db)
	if err != nil {
		return nil, err
	}
	db.isInit = true
	return db, nil
}
