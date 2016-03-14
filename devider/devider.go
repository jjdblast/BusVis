package main

import(
	"os"
	"fmt"
	"io"
	"strings"
	"bufio"

	"github.com/yibowang/BusVis/readline"
)

type FileBuff struct{
	file *os.File
	writer *bufio.Writer
}
func devide(file io.Reader,date string,base string){
	filemap := make(map[string] FileBuff)

	countline := 0
	for {
		line := readline.ReadLine(file)
		if len(line)==0 {
			break
		}
		countline++
		fmt.Printf("%d\r",countline)
		if strings.HasPrefix(line[3],date)&& strings.HasPrefix(line[4],date) {
			fb,find := filemap[line[7]]
			if !find {
				ofile,err := os.Create(base+"/"+line[7]+".csv")
				writer := bufio.NewWriter(ofile)
				if err != nil {
					panic(err)
				}
				fb =  FileBuff{file:ofile,writer:writer}
				filemap[line[7]] = fb
			}
			for i,v := range line {
				if i != 0 {
					fb.writer.WriteString(",")
				}
				fb.writer.WriteString("\"")
				fb.writer.WriteString(v)
				fb.writer.WriteString("\"")
			}
			fb.writer.WriteString("\n")
		}
	}
	defer func(){
		for _,f := range filemap {
			f.writer.Flush()
			f.file.Close()
		}
	}()
}
func devideFile(filestr string,date string,base string){
	file,err := os.Open(filestr)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	defer file.Close()
	devide(r,date,base)
}
func main() {
	if len(os.Args) != 3{
		fmt.Println("format: devider filename date\nexample: devider 20150807.csv 20150807")
		return
	}
	base := os.Args[2]+"/"+"linedata"
	os.MkdirAll(base, 0666)
	devideFile(os.Args[1],os.Args[2],base)
}
