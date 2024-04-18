# GaussDB采集器开发规范

### 1.命名规范

1. ***指标源文件命名***

​	`gaussdb_xxx.go`

2. ***指标命名***

​	`gaussDB_xxx_xxx_xxx{}`

### 2.文件目录规范

1. ***源文件总目录***

   `$HOME/openGuass_exporter/`

2. ***指标源文件相关目录***

   `$HOME/openGauss_exporter/collector/`

### 3. 代码提交规范【<font color=red>重要</font>】

1. 从远程将代码库克隆到本地，并且新建分支，修改完毕后将代码提交到远程仓库，然后统一合并。相关命令如下：

   1. 克隆代码：

      `git clone appops@8.134.95.226:/home/ming/gitRepos/openGauss_exporter.git`

   2. 创建新分支

      `git branch dev_xxx`

      ***其中，新分支名称自定义***

   3. 切换到新分支

      `git checkout dev_xxx`

   4. 提交远程仓库

      `git add remote origin appops@8.134.95.226:/home/appops2/gitRepos/openGuass_exporter.git`

2. 添加、修改、删除文件后在群中同步，并且登记到相关表格中

​	[【基础与数据运维组】采集器开发记录-(代码提交记录)](https://doc.weixin.qq.com/sheet/e3_Ae0AlwakAHMnpwZlYIvQX6Idi2HEN?scode=AFYArwcHAAsRJ4F78yAe0AlwakAHM&tab=BB08J2)

### 4. 指标文件编写指引

1. ***必需实现的函数、变量及结构体***

   1. 子系统名称（变量）

      `const xxxSubsystem=xxx`

   2. 接收器

      `type xxxCollector struct {}`

   3. 指标名称及SQL语句

      `var（）`

   4. Name方法

      原型：`func (xxxCollector) Name() string {}`

      作用：指明指标的名称

      是否可选：否

   5. Help方法

      原型：`func (xxxCollector) Help() string {}`

      作用：说明指标帮助
   
      是否可选：否
   
   6. Version方法
   
      原型：`func (xxxCollector) Version() string {}`
   
      作用：说明指标的版本号码
   
      是否可选：否
   
   7. Scrape方法
   
      原型：`func (xxxCollector) Scrape(ctx context,db *sql.DB,ch chan <- prometheus.Metric,logger log.Logger) error {}`
   
      作用：获取指标的函数
   
      是否可选：否
   
      