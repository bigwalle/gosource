# tar 简介
    Unix和类Unix系统上的压缩打包工具，可以将多个文件合并为一个文件，打包后的文件后缀亦为“tar”。tar文件格式已经成为POSIX标准，最初是POSIX.1-1988，当前是POSIX.1-2001。本程序最初的设计目的是将文件备份到磁带上（tape archive），因而得名tar。
    
    常用的tar是自由软件基金会开发的GNU版，稳定版本是1.28，发布于2014年7月27日
    同时，它有多个压缩率不同的版本，如tar.xz和tar.gz，前者的压缩率更高，但可能有兼容性问题。
    
    ## 再看 go中的 archive/tar
    
    原文：Package tar implements access to tar archives.
    译文：tar 包实现了对 tar 文件格式读写的操作。
    
    
    原文:Tape archives (tar) are a file format for storing a sequence of files that can be read and written in a streaming manner. This package aims to cover most variations of the format, including those produced by GNU and BSD tar tools.
     
    译文:

## 例如 
```golang
package main

import (
  "archive/tar"
  "bytes"
  "fmt"
  "io"
  "log"
  "os"
)

func main(){
  // Create and add some files to the archive.
  //创建并添加一些文件到压缩文件中
  var buf bytes.Buffer
  tw:=tar.NewWriter(&buf)
  var files = []struct{
      Name ,Body string
  }{
    {"readme.txt", "This archive contains some text files."},
    {"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
    {"todo.txt", "Get animal handling license."},
  }
    
    //Mode 表示分配给文件的权限   
    //0600表示分配给文件的权限
    // 一共四位数，第一位数表示gid/uid一般不用，剩下三位分别表示owner,group,other的权限
    // 每个数可以转换为三位二进制数，分别表示rwx(读，写，执行)权限，为1表示有权限，0无权限
    // 如6是上面第二个数，可以表示为二进制数110,表示owner有读，写权限，无执行权限
  for _,file :=range files {
       hdr:=tar.Header{
         Name:       file.Name,
         Size:       int64(len(file.Body)),
         Mode:       0600, //表示分配给文件的权限
       }
       if err:=tw.WriteHeader(&hdr);err!=nil{
         log.Fatal(err)
       }
       if _,err:=tw.Write([]byte(file.Body));err!=nil{
         log.Fatal(err)
       }
  }
  if err:=tw.Close();err!=nil{
    log.Fatal(err)
  }

  // Open and iterate through the files in the archive.
  tr := tar.NewReader(&buf)
  for {
    hdr, err := tr.Next()
    if err == io.EOF {
      break // End of archive
    }
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("Contents of %s:\n", hdr.Name)
    if _, err := io.Copy(os.Stdout, tr); err != nil {
      log.Fatal(err)
    }
    fmt.Println()
  }


}

```
