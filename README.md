# SOSS (Secure Object Storage Service)

一个可以在把文件上传到s3存储, 上传之前加密文件内容，下载时自动解密的小工具。  
目前适配了阿里云OSS, 后续再扩展.    

## 参考
**[参考高天大佬的 soss 工具](https://github.com/gaogaotiantian/soss)**
此项目为soss的go版本, 逻辑作出了如下几点的修改:
1. 将oss-sdk与controller逻辑分离, 后续可以方便的扩展更多的s3对象存储支持;
2. 将原来的**加密**的逻辑, 改为**加密+压缩**, 在文件较大的情况下可以节省空间, *但可能会导致工具运行时候消耗更多内存;*
3. 上传、下载改用了并行处理, 处理多个文件时候能提高性能


## 安装
### [Download the latest binary](https://github.com/linlanniao/soss/releases/latest)
**wget**  
使用 wget 下载预编译的二进制压缩文件     
例子: `VERSION=v4.2.0` `BINARY=soss_Linux_x86_64`

```
wget https://github.com/linlanniao/soss/releases/download//${VERSION}/${BINARY}.tar.gz -O - |\
  tar xz && mv ${BINARY} /usr/local/bin/soss
```


## 准备工作
### AccessKey 和 Access Key Secret
以阿里云为例, 在你的阿里云管理系统内，找到下面的内容：
* OSS Bucket的endpoint（例如`oss-cn-hangzhou.aliyuncs.com`）
* OSS Bucket的名字
* 你的用户的access key（推荐使用RAM用户）
    * `export S3_ACCESS_KEY_ID=<KEY ID>`
    * `export S3_ACCESS_KEY_SECRET=<KEY SECRET>`


### 配置文件
在`config.yaml`中，配置好`client_type` `endpoint`和`bucket`

### config.yaml example
```yaml
# oss client 的类型, 目前只支持 “阿里云oss” 后续可能会扩展更多
client_type: oss

# bucket的名字
bucket: ppops-bucket

# endpoint 地址, 注意带上http/https
endpoint: https://oss-cn-guangzhou.aliyuncs.com 
```
* 将配置文件保存在 `$HOME/.soss/config.yaml` 或者当前目录 `./config.yaml`  

如果不想使用`config.yaml`，也可以在命令行作为参数输入。


## 使用说明

### 文件列表

```
# 如果配置好了config.json
soss list
或
soss ls

soss ls --prefix data/

# 如果想在命令行输入bucket和endpoint
soss ls -b bucket_name -e endpoint
```

### 上传文件

```
soss upload -k my_password text.txt image.png
或
soss up -k my_password  text.txt image.png

# 支持上传整个文件夹的内容，文件夹所有内容会保持结构上传到bucket根目录
soss upload -k my_password data/

# 设置bucket保存路径的prefix，文件夹所有内容会保持结构上传到data/目录
soss upload -k my_password --prefix data/ data/

# 如果encrypt key是一个32或者64位的hex，则直接作为AES的key使用，否则进行SHA256，转换成32 byte的key
soss upload -k deadbeef12345678deadbeef87654321 text.txt

# 同样也可以传入bucket和endpoint
soss upload -b bucket -e endpoint -k my_password text.txt
```

### 下载文件

```
soss download -k my_password text.txt image.png

# 指定保存文件夹
soss download -k my_password --output_dir ./data text.txt image.png

# 剩下的参数和upload一样, 具体可以通过-h参数查看
```

### LICENSE

Copyright 2024 linlanniao.

Distributed under the terms of the [MIT License](LICENSE)