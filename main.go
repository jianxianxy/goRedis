package main

import (
    "bufio"
	"fmt"
	"os"
    "strings"
    "strconv"
    "github.com/garyburd/redigo/redis"
)

type Redigo struct{
    conn redis.Conn
}

func main(){
    var cache = new(Redigo)
    tip()
    var seldb int = -1 
    input := bufio.NewScanner(os.Stdin) //接受运行是输入
	for input.Scan() {
		val := input.Text() //input.Text()终端输入的内容
		if val == "exit" {
			break
		}
        arr := strings.Split(val," ")
        if arr[0] == "db"{
            ind,_:=strconv.Atoi(arr[1])  
            seldb = ind
            cache.Make(ind)
            fmt.Println("选择了DB:",arr[1])    
        }
        if seldb < 0{
            seldb = 0
            cache.Make(0)    
        }
        if arr[0] == "set"{
            if len(arr[0]) == 3{
                cache.Set(arr[1],arr[2])    
            }
            if len(arr[0]) == 4{
                cache.SetEx(arr[1],arr[2],arr[3])
            }
        }
        if arr[0] == "get"{
            fmt.Println(cache.Get(arr[1]))    
        }
        if arr[0] == "del"{
            cache.Del(arr[1])    
        }
	}
}
//创建连接
func (db *Redigo) Make(ind int){
    conn, err := redis.Dial("tcp", "10.102.36.153:26379")
    checkErr(err)
    conn.Do("AUTH", "tttt")
    conn.Do("SELECT", ind)
    db.conn = conn
}
//设置值
func (db Redigo) Set(key,val string){
    db.conn.Do("SET", key, val)
}
//设置带有效期的值
func (db Redigo) SetEx(key,val,time string){
    db.conn.Do("SET",key,val,"EX",time)
}
//取值
func (db Redigo) Get(key string) string{
    val, err := redis.String(db.conn.Do("GET", key))
    if err == nil{
        return val;
    }else{
        return "";    
    }
}
//删除键
func (db Redigo) Del(key string){
    db.conn.Do("DEL",key)
}
//错误检查
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func tip(){
    fmt.Println("指令提示:")
    fmt.Println("db int 选择DB")
    fmt.Println("set key val 设置缓存")
    fmt.Println("set key val ext 设置带过期的缓存")  
    fmt.Println("get key 获取缓存")   
    fmt.Println("del key 删除缓存")   
}