package utils

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tossp/tsgo/pkg/log"
)

func Zip(dst, src string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = fw.Close()
	}()

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	dir, _ := path.Split(filepath.ToSlash(src))
	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(p string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}
		// 整理打包路径
		fh.Name = strings.TrimPrefix(strings.TrimPrefix(filepath.ToSlash(p), dir), string(filepath.Separator))
		if fi.IsDir() {
			fh.Name += string(filepath.Separator)
		}
		fh.Name = filepath.ToSlash(fh.Name)

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(p)
		if err != nil {
			return
		}

		defer func() {
			_ = fr.Close()
		}()
		// 将打开的文件 Copy 到 w
		_, err = io.Copy(w, fr)
		if err != nil {
			return
		}
		return nil
	})
}

func UnZip(dst, src string) (err error) {
	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
	zr, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	defer func() {
		_ = zr.Close()
	}()

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		if err = unZipLoop(dst, file); err != nil {
			return
		}
	}
	return nil
}
func unZipLoop(dst string, file *zip.File) (err error) {
	p := filepath.Join(dst, file.Name)
	// 如果是目录，就创建目录
	if file.FileInfo().IsDir() {
		if err = os.MkdirAll(p, file.Mode()); err != nil {
			return
		}
		// 因为是目录，跳过当前循环，因为后面都是文件的处理
		return
	}

	// 获取到 Reader
	fr, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = fr.Close()
	}()

	// 创建要写出的文件对应的 Write
	fw, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer func() {
		_ = fw.Close()
	}()
	_, err = io.Copy(fw, fr)
	if err != nil {
		return err
	}
	return
}
