package net

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Cache ...
type Cache struct {
	Tmp string
}

var cache *Cache

// NewCache ...
func NewCache(tmp string) *Cache {
	if cache != nil {
		return cache
	}

	var path_ string

	if filepath.IsAbs(tmp) {
		path_ = tmp
	} else {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		path_ = path.Join(path.Dir(dir), tmp)
	}

	_ = os.MkdirAll(tmp, os.ModePerm)
	cache = &Cache{Tmp: path_}

	cache.prepareUnicodePath()

	return cache
}

//从unicode获得cache文件句柄
func (c *Cache) ReaderUnicode(base_url string, wd string) (reader io.ReadCloser, e error) {
	wd_str := c.unicodePath(wd)

	stat, e := os.Stat(wd_str)

	if (e == nil && stat.Size() != 0) || !os.IsNotExist(e) {
		open, e := os.Open(wd_str)
		if e != nil {
			panic(e)
		}
		return open, nil
	}

	closer, err := UnicodeRequest(base_url, wd)
	if err != nil {
		Log.Fatal(err)
	}

	return c.CacheUnicode(closer, wd)
}

//从URL获得cache文件句柄
func (c *Cache) Reader(req_url string) (reader io.ReadCloser, e error) {
	req_url_path := strings.Split(req_url, "://")[1]

	file_path := path.Join(path.Dir(c.Tmp), req_url_path)

	stat, e := os.Stat(file_path)

	if (e == nil && stat.Size() != 0) || !os.IsNotExist(e) {
		open, e := os.Open(file_path)
		if e != nil {
			panic(e)
		}
		return open, nil
	}

	closer, err := GetRequest(req_url)
	if err != nil {
		Log.Fatal(err)
	}

	Prepare(file_path)

	return c.Cache(closer, req_url)
}

//从URL和http响应内容句柄获得cache文件句柄
func (c *Cache) Cache(closer io.ReadCloser, req_url string) (io.ReadCloser, error) {
	req_url_path := strings.Split(req_url, "://")[1]

	file_path := path.Join(path.Dir(c.Tmp), req_url_path)

	stat, e := os.Stat(file_path)

	if (e == nil && stat.Size() != 0) || !os.IsNotExist(e) {
		return nil, os.ErrExist
	}
	file, e := os.OpenFile(file_path, os.O_TRUNC|os.O_CREATE|os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if e != nil {
		panic(e)
	}
	defer file.Close()
	_, e = io.Copy(file, closer)
	if e != nil {
		panic(e)
	}
	open, e := os.Open(file_path)
	if e != nil {
		panic(e)
	}

	return open, nil
}

//从Unicode和http响应内容句柄获得cache文件句柄
func (c *Cache) CacheUnicode(closer io.ReadCloser, wd string) (io.ReadCloser, error) {
	wd_str := c.unicodePath(wd)

	stat, e := os.Stat(wd_str)

	if (e == nil && stat.Size() != 0) || !os.IsNotExist(e) {
		return nil, os.ErrExist
	}
	file, e := os.OpenFile(wd_str, os.O_TRUNC|os.O_CREATE|os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if e != nil {
		panic(e)
	}
	defer file.Close()
	_, e = io.Copy(file, closer)
	if e != nil {
		panic(e)
	}
	open, e := os.Open(wd_str)
	if e != nil {
		panic(e)
	}

	return open, nil
}

//准备Unicode的cache路径
func (c *Cache) prepareUnicodePath() {
	dir := path.Join(path.Dir(c.Tmp), "unicode")
	dir_stat, _ := os.Stat(dir)

	if dir_stat != nil && !dir_stat.IsDir() {
		panic("cache dir conflict to a file")
	}

	_ = os.MkdirAll(dir, os.ModePerm)
}

//返回unicode数串对应页面的缓存文件路径
func (c *Cache) unicodePath(wd string) string {
	if len(wd) > 6 {
		panic("输入unicode数串长度不对")
	}

	dir := path.Join(path.Dir(c.Tmp), "unicode")

	result := path.Join(dir, wd+".htm")

	return result
}

// URL存为cache文件
func (c *Cache) Get(req_url string) (e error) {
	req_url_path := strings.Split(req_url, "://")[1]

	file_path := path.Join(path.Dir(c.Tmp), req_url_path)

	stat, e := os.Stat(file_path)

	if (e == nil && stat.Size() != 0) || !os.IsNotExist(e) {
		return os.ErrExist
	}

	closer, err := GetRequest(req_url)
	if err != nil {
		Log.Fatal(err)
	}

	Prepare(file_path)

	file, e := os.OpenFile(file_path, os.O_TRUNC|os.O_CREATE|os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if e != nil {
		return e
	}
	written, e := io.Copy(file, closer)
	if e != nil {
		return e
	}
	//ignore written
	_ = written
	closer.Close()
	return nil
}

//为目标路径创建文件夹，返回目标路径的绝对路径
func Prepare(path string) string {
	s, e := filepath.Abs(path)
	if e != nil {
		panic(e)
	}
	dir, _ := filepath.Split(s)
	_ = os.MkdirAll(dir, os.ModePerm)
	return s
}

// URL对应的缓存文件转存成指定文件
func (c *Cache) Save(req_url string, to string) (written int64, e error) {
	req_url_path := strings.Split(req_url, "://")[1]

	file_path := path.Join(path.Dir(c.Tmp), req_url_path)
	info, e := os.Stat(file_path)
	if e != nil && os.IsNotExist(e) {
		panic(errors.Wrap(e, "cache get error"))
	}
	if info.IsDir() {
		panic("cache get a dir")
	}

	abs_to := Prepare(to)

	file, e := os.Open(file_path)
	if e != nil {
		panic(e)
	}

	pj := path.Join(abs_to)

	openFile, e := os.OpenFile(pj, os.O_TRUNC|os.O_CREATE|os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if e != nil {
		panic(e)
	}
	return io.Copy(openFile, file)
}
