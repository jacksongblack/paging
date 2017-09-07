# paging
golang the mvc framework beego paging tool

## demo
``` 
your models

type User struct {
  Id int `orm:""pk`
  Name string `orm:"size(200)"`
  Password string `orm:"size(200)"`
}

your controllers

import (
"github/jacksongblack/paging"
"github.com/astaxie/beego"
"github.com/astaxie/beego/orm"
)

type QueryCtr struct {
  beego.Controller
}

func (this *QueryCtr)Index(){
    users := []*models.User{}
    ormer:= orm.NewOrm()
    msg,err:=paging.Filter(this.Ctx,&users,this.ORM)
    if err!=nil{
    	this.Ctx.Output.Header("Content-Type", "application/json")
	    this.Ctx.ResponseWriter.Write(err)
      this.StopRun()
    }
    this.Ctx.Output.Header("Content-Type", "application/json")
	  this.Ctx.ResponseWriter.Write(err)
    
}

func (this *QueryCtr)PasswordIsNotNull(){
    users := []*models.User{}
    ormer:= orm.NewOrm()
    query:=paging.Query(new(modes.User))
    query:=query.And("password__isnull",true)
    data,err:=paging.QuerySetFilter(this.ctx,&users,query,ormer)
    if err!=nil{
    	this.Ctx.Output.Header("Content-Type", "application/json")
	    this.Ctx.ResponseWriter.Write(err)
      this.StopRun()
    }
    this.Ctx.Output.Header("Content-Type", "application/json")
	  this.Ctx.ResponseWriter.Write(data)
}
```
### Prams
url_demo : http://localhost:8080/user?page=1&field=id
```
@Param	page		    query 	string	  false		第几页
@Param	field		    query 	string	  false		排序字段
@Param	order		    query 	string	  false		排序：DESC降、ASC升
@Param	limit		    query 	string	  false		查询数量
```
