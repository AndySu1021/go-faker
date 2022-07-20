# 假資料測試工具

提供 CLI 操作介面，協助團隊快速建立一套假資料塞入工具，能夠自動解析 table schema 搭配 codegen 自動生成 model 文件模板，可以對其進行二次開發。  

## 項目說明
- 目前自動解析部分僅支援 MySQL
- ./faker/faker.go 中提供一些隨機生成方法可以使用

## 操作指令

### 創建 model
生成的檔案將會存在 ./model 目錄下
```shell
faker make $(table_name)
```

### 塞入資料
塞入資料前可以在 ./model/$(table_name).go 中的 Definition 方法自定義需要塞入的值
```shell
faker create $(table_name) -n=$(num)
```
也可以在最後的部分加上 -- --$(column_name)=$(column_value) 動態更改 Definition 塞入的值  
例如：
```shell
faker create $(table_name) -n=$(num) -- --status=2
```

### 清空資料
```shell
faker clear $(table_name)
```

### 清空資料
```shell
faker clear $(table_name)
```

### 更新 ModelMap
如果是手動新增 ./model 目錄下的檔案，需要更新 ./model/model.go 中的 ModelMap
```shell
faker refresh
```