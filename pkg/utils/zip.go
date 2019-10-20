package utils

import (
    "archive/zip"
    "io/ioutil"
    "log"
    "os"
    "strings"
)

func DeCompress(tarFile, dest string) error {
    if strings.HasSuffix(tarFile, ".zip") {
        return zipDeCompress(tarFile, dest)
    }
    return nil
}

func zipDeCompress(zipFile, dest string) error {
    or, err := zip.OpenReader(zipFile)
    if err != nil {
        return err
    }
    defer func() {
        _ = or.Close()
    }()
    log.Print(" 压缩文件", zipFile, " 解压到", dest)
    for _, item := range or.File {
        if item.FileInfo().IsDir() {
            _ = os.Mkdir(dest+item.Name, 0777)
            continue
        }
        rc, _ := item.Open()
        dst, _ := createFile(dest + item.Name)
        payload, err := ioutil.ReadAll(rc)
        n, err := dst.Write(payload)
        if err != nil {
            log.Print(dest + item.Name)
            log.Print(err)
        } else {
            log.Print(n/1024, "kb", "  ", dest+item.Name)
        }
    }

    return nil
}

func createFile(name string) (*os.File, error) {
    err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
    if err != nil {
        return nil, err
    }
    return os.Create(name)
}
