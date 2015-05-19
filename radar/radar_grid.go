package main

import (
	"bufio"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"flag"
	"strings"
	"time"
	
)

type RadarPoint struct {
	Code string
	Lat  float64
	Lon  float64
}

var (
	dburl string = "54.223.146.200:27110"
)

//db.radar.createIndex({"loc":"2dsphere"})
func buildPointDoc(point RadarPoint) bson.M {

	m := bson.M{"loc": bson.M{"type": "Point", "coordinates": []float64{point.Lon, point.Lat}}, "code": point.Code}

	return m
}
func main() {
	flag.Parse();
	rootdir:=flag.String("rootdir","/Users/tiger/radar","rootdir")
	worker:=flag.Int("worker",2,"worker")
	_dburl:=flag.String("db",dburl,"db")
	dburl=*_dburl
	ch := readPoint(*rootdir)
	saveCounter := *worker
	chs := make([]<-chan int, 0, saveCounter)
	for i := 0; i < saveCounter; i++ {
		chs = append(chs, writeToDb(ch))
	}

	cases := make([]reflect.SelectCase, len(chs))
	for i, tch := range chs {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(tch)}
	}
	cmCount := 0
	for {
		_, _, ok := reflect.Select(cases)
		if ok == false {
			cmCount++
		}
		if cmCount == saveCounter {
			break
		}
	}
}

func writeToDb(ch <-chan RadarPoint) <-chan int {
	rch := make(chan int)
	go func() {
		session, err := mgo.Dial(dburl)
		if err != nil {
			panic(err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("radar").C("radar")
		sum:=0
		for {
			point ,ok:= <-ch
			if(ok==true){
				c.Insert(point)
				sum++
				if(sum%500==0){
					fmt.Println("save "+strconv.Itoa(sum)+" success")
				}
			}else{
				break
			}
		}
		fmt.Println("save complete "+strconv.Itoa(sum)+" success")
		time.Sleep(time.Duration(time.Second*2))
		rch <- 1
		close(rch)
	}()
	
	return rch

}
func readPoint(dir string) <-chan RadarPoint {
	ch := make(chan RadarPoint, 10)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("%v\n", err)
		close(ch)
		return ch
	}
	go func() {
		for _, file := range files {
			if file.IsDir() == true {
				continue
			}
			path := dir + string(os.PathSeparator) + file.Name()
			readFile(file.Name(), path, ch)
		}
		close(ch)
	}()

	return ch
}

func readFile(code, filePath string, ch chan RadarPoint) {
	file, err := os.Open(filePath)
	if err != nil {
		println("file open error")
		fmt.Printf("%v\n", err)
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" || len(s) == 0 {
			continue
		}

		ss := strings.Split(s, ",")
		lon, er_lon := strconv.ParseFloat(ss[1], 64)
		lat, er_lat := strconv.ParseFloat(ss[0], 64)
		if er_lon != nil || er_lat != nil {
			fmt.Printf("bad data"+"  "+s )
			continue
		}
		p := RadarPoint{Lon: lon, Lat: lat, Code: code}
		ch <- p
	}

}
