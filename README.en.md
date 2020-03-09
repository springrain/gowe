# readygo

#### 介绍
golang开发脚手架

#### 软件架构
基于[gin](https://github.com/gin-gonic/gin)和[zorm](ttps://gitee.com/chunanyong/zorm)  
[自带代码生成器](https://gitee.com/chunanyong/readygo/tree/master/codegenerator)    
zorm基于原生sql语句编写,是springrain的精简和优化  
zorm代码精简,总计2000行左右,注释详细,方便定制修改.  
zorm支持事务传播,这是zorm诞生的主要原因  


#### 例子
具体可以参照 [UserStructService.go](https://gitee.com/chunanyong/readygo/tree/master/permission/permservice)



 1.  生成实体类或手动编写,建议使用代码生成器 https://gitee.com/chunanyong/readygo/tree/master/codegenerator
  ```  

//UserOrgStructTableName 表名常量,方便直接调用
const UserOrgStructTableName = "t_user_org"

// UserOrgStruct 用户部门中间表
type UserOrgStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id 编号
	Id string `column:"id"`

	//UserId 用户编号
	UserId string `column:"userId"`

	//OrgId 机构编号
	OrgId string `column:"orgId"`

	//ManagerType 0会员,1员工,2主管
	ManagerType int `column:"managerType"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//GetTableName 获取表名称
func (entity *UserOrgStruct) GetTableName() string {
	return UserOrgStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *UserOrgStruct) GetPKColumnName() string {
	return "id"
}

  ```  
2.  初始化zorm

    ```  
    dataSourceConfig := zorm.DataSourceConfig{
	Host:     "127.0.0.1",
	Port:     3306,
	DBName:   "readygo",
	UserName: "root",
	PassWord: "root",
	DBType:   "mysql",
     }
     zorm.NewBaseDao(&dataSourceConfig)
    ```  
3.  增
    ```
    var user permstruct.UserStruct
    err := zorm.SaveStruct(context.Background(), &user)
    ```
4.  删
    ```
    err := zorm.DeleteStruct(context.Background(),&user)
    ```
  
5.  改
    ```
    err := zorm.UpdateStruct(context.Background(),&user)
    //finder更新
    err := zorm.UpdateFinder(context.Background(),finder)
    ```
6.  查
    ```
	finder := zorm.NewSelectFinder(permstruct.UserStructTableName)
	page := zorm.NewPage()
	var users = make([]permstruct.UserStruct, 0)
	err := zorm.QueryStructList(context.Background(), finder, &users, &page)
    ```
7.  事务传播
    ```
    //匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始
		errSaveUserStruct := zorm.SaveStruct(ctx, userStruct)

		if errSaveUserStruct != nil {
			return nil, errSaveUserStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})
    ```
8.  生产示例
    ```  
    //FindUserOrgByUserId 根据userId查找部门UserOrg中间表对象
    func FindUserOrgByUserId(ctx context.Context, userId string, page *zorm.Page) ([]permstruct.UserOrgStruct, error) {
	if len(userId) < 1 {
		return nil, errors.New("userId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.* FROM  ").Append(permstruct.UserOrgStructTableName).Append(" re ")
	finder.Append("   WHERE re.userId=?    order by re.managerType desc   ", userId)

	userOrgs := make([]permstruct.UserOrgStruct, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &userOrgs, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return userOrgs, nil
    }
    ```  

9.  [测试](https://www.jianshu.com/p/1adc69468b6f)
    ```  
    //函数测试
    go test -run TestAdd2
    //性能测试
    go test -bench=.
    go test -v -bench=. -cpu=8 -benchtime="3s" -timeout="5s" -benchmem
    ```  

