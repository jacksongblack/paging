package paging
///分页用库
import (
	"math"
	"reflect"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/context"
	"strings"
	"strconv"
)
type FilterMsg struct {
	TotalData   int
	TotalPage   int
	CurrentPage int
	Total       int
	Data        interface{}
}
//Get element's type
func GetType(typ reflect.Type) reflect.Type {
	switch typ.Kind().String() {
	case "ptr":
		typ = GetType(typ.Elem())
		break

	case "slice":
		typ = GetType(typ.Elem())
		break

	default:
		return typ
	}

	return typ
}
//Query data
func Query(data interface{}) *orm.Condition {
	qs := orm.NewCondition()
	return qs
}
func GetDataType(data interface{}) reflect.Type {
	typeVal := reflect.TypeOf(data)
	type_name := typeVal.Kind().String()
	if type_name == "ptr" || type_name == "slice" {
		typeVal = GetType(typeVal)
	}
	return typeVal
}
func Filter(ctx *context.Context,data interface{},ormer orm.Ormer,cols...string ) (*FilterMsg, error) {
	q:=Query(data)
	return QuerySetFilter(ctx,data,q,ormer,cols...)

}

func QuerySetFilter(ctx *context.Context,data interface{},condition *orm.Condition,ormer orm.Ormer,cols...string )(*FilterMsg, error){
	if err := ctx.Request.ParseForm(); err != nil {
		beego.Error(err)
		return nil, err
	}

	params := make(map[string]string)
	for k, v := range ctx.Request.Form {
		params[k] = v[0]
	}

	// 对排序、分页、条数限制、排序字段进行过滤

	page := 1
	limit := 20
	var err error
	for k, v := range params {

		if k == "page" || k == "limit" || k == "order" || k == "field" {
			continue
		}
		if v != "" {
			if strings.Contains(k,"__isnull"){
				b,err:=strconv.ParseBool(v)
				if err!=nil{
					return nil,err
				}
				condition=condition.And(k,b)
				continue
			}
			if strings.Contains(k,"__exclude"){
				keys:=strings.Split(k,"__exclude")
				condition=condition.AndNot(keys[0],v)
				continue
			}
			if strings.Contains(k,"_or_"){
				cond:= orm.NewCondition()
				k_list:=strings.Split(k,"_or_")
				for n,key:=range k_list{
					if n==0{
						cond = cond.And(key,v)
					}else {
						cond=cond.Or(key,v)
					}

				}
				condition=condition.AndCond(cond)
				continue
			}
			if strings.Contains(k,"__in"){
				split := strings.Split(v, ",")
				condition=condition.And(k,split)
				continue
			}
			condition=condition.And(k,v)

		}
	}
	pageStr, ok := params["page"]
	if ok {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
	}

	page = page - 1
	if page < 0 {
		page = 0
	}
	q:=ormer.QueryTable(reflect.New(GetDataType(data)).Interface()).SetCond(condition)
	limitStr, ok := params["limit"]
	if ok {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}
		q = q.Limit(limit, page*limit)
	}

	// 排序判断
	field, ok := params["field"]
	if ok {
		order, ok := params["order"]
		if !ok {
			order = "DESC"
		}
		if order == "ASC" {
			q = q.OrderBy(field)
		}
		if order == "DESC" {
			q = q.OrderBy("-" + field)
		}
	}

	typeVal := reflect.TypeOf(data)
	if typeVal.Kind().String() == "ptr" || typeVal.Kind().String() == "slice" {
		typeVal = GetType(typeVal)
	}


	totalData, _ := q.Count() //总数
	var totalPage float64

	//如果url参数不存在limit , 则不进行总页数计算, 则总页数为0
	_, issetOk := params["limit"]
	if issetOk {
		if totalData > 0 {
			totalPage = math.Ceil(float64(totalData) / float64(limit)) //总页数
		}
	}

	n, err := q.All(data,cols...)
	if err != nil {
		return nil, err
	}

	filterMsg := &FilterMsg{}
	filterMsg.TotalData = int(totalData)
	filterMsg.TotalPage = int(totalPage)
	filterMsg.CurrentPage = page + 1
	filterMsg.Total = int(n)
	filterMsg.Data=data
	return filterMsg, nil
}