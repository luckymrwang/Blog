package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"strconv"
	"time"
)

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Content         string
	Attachment      string
	Created         time.Time
	Updated         time.Time
	Views           int64
	Author          string
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
	ClassID         int64
}

// 分类
type Classify struct {
	Id    int64
	Title string
}

// 统计文章分类
type ClassifyCount struct {
	Id      int64
	ClassID int64
	Num     int64
	Title   string
}

func init() {
	// orm.Debug = true
	orm.RegisterModel(new(Topic), new(Classify), new(ClassifyCount))
}

// 添加文章
func AddTopic(title, content string, classID int64) (int64, error) {
	o := orm.NewOrm()
	topic := &Topic{
		Title:     title,
		Content:   content,
		Created:   time.Now(),
		Updated:   time.Now(),
		ReplyTime: time.Now(),
		ClassID:   classID,
	}
	id, err := o.Insert(topic)
	return id, err
}

// 获取所有文章
func GetAllTopic(p int, cid int64) (topic []*Topic, err error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	pageNum, _ := beego.AppConfig.Int("pageNum")
	if cid > 0 {
		_, err = qs.OrderBy("-Created").Filter("class_i_d", cid).Limit(pageNum, (p-1)*pageNum).All(&topics)
	} else {
		_, err = qs.OrderBy("-Created").Limit(pageNum, (p-1)*pageNum).All(&topics)
	}
	fmt.Println(os.Stderr)
	for _, value := range topics {
		fmt.Println("title:", value.Title, "id:", value.Id)
	}
	return topics, err
}

// 获取单个文章信息
func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)
	return topic, err
}

// 获取文章总数
func GetTopicCount(title string, cid int64) (count int64, err error) {
	fmt.Println("cid: ", cid)
	o := orm.NewOrm()
	if len(title) > 0 {
		count, err = o.QueryTable("topic").Filter("title__icontains", title).Count()
	} else if cid > 0 {
		count, err = o.QueryTable("topic").Filter("class_i_d", cid).Count()
	} else {
		count, err = o.QueryTable("topic").Count()
	}

	if err != nil {
		return 0, err
	}
	fmt.Println("count: ", count)
	return count, err
}

// 文章搜索
func SearchTopic(title string) (topic []*Topic, err error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	_, err = o.QueryTable("topic").Filter("title__icontains", title).All(&topics)
	return topics, err
}

// 添加分类
func AddClassify(title string) error {
	o := orm.NewOrm()
	classify := &Classify{
		Title: title,
	}
	_, err := o.Insert(classify)
	return err
}

// 查询分类
func FindClassify(id int64) (classify []*Classify, err error) {
	o := orm.NewOrm()
	classify = make([]*Classify, 0)
	if id == 0 {
		_, err = o.QueryTable("classify").All(&classify)
	} else {
		_, err = o.QueryTable("classify").Filter("id", id).All(&classify)
	}
	return classify, err
}

// 统计分类总数
func GetClassifyCount() (classifyCount []*ClassifyCount, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t.class_i_d, count(1) num", "c.title").From("topic t").LeftJoin("classify c").On("t.class_i_d = c.id").Where("t.class_i_d in (c.id)").GroupBy("t.class_i_d").OrderBy("num desc")

	sql := qb.String()
	o := orm.NewOrm()

	_, err = o.Raw(sql).QueryRows(&classifyCount)
	return classifyCount, err
}

// 获取分类最大ID
func GetClassifyMaxID() (int, error) {
	var maps []orm.Params
	o := orm.NewOrm()
	num, err := o.Raw("select max(id) maxid from classify").Values(&maps)
	if err == nil && num > 0 {
		return strconv.Atoi(maps[0]["maxid"].(string))
	}
	return 0, err
}

// 删除分类
func DeleteClassify(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Classify{Id: id})
	return err
}

// 修改分类
func UpdateClassify(id int64, title string) error {
	o := orm.NewOrm()
	classify := Classify{Id: id}
	if o.Read(&classify) == nil {
		classify.Title = title
		_, err := o.Update(&classify)
		return err
	}
	return nil
}
