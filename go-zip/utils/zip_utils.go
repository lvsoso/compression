package utils

import (
	"archive/zip"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	kzip "github.com/klauspost/compress/zip"
)

func Pack(srcDir string, zipFileName string) error {
	os.RemoveAll(zipFileName)
	zipfile, _ := os.Create(zipFileName)
	defer zipfile.Close()
	zw := zip.NewWriter(zipfile)
	defer zw.Close()

	walkFunc := func(path string, info os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		if path == srcDir {
			return nil
		}

		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))
		if info.IsDir() {
			fh.Name += "/"
		}

		fh.Method = zip.Deflate

		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		if !fh.Mode().IsRegular() {
			return nil
		}

		rc, err := os.Open(path)
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		fmt.Println(fh.Name, " completed !")
		return nil
	}

	err := filepath.Walk(srcDir, walkFunc)
	return err
}

func Unpack(dstPath string, zipFileName string) error {
	archive, err := zip.OpenReader(zipFileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dstPath, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dstPath)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}

		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func KPack(srcDir string, zipFileName string) error {
	os.RemoveAll(zipFileName)
	zipfile, _ := os.Create(zipFileName)
	defer zipfile.Close()
	zw := kzip.NewWriter(zipfile)
	zw.RegisterCompressor(kzip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		// return flate.NewWriter(out, flate.BestSpeed)
		return flate.NewWriter(out, flate.BestCompression)
		//return flate.NewWriter(out, flate.DefaultCompression)
	})
	defer zw.Close()

	walkFunc := func(path string, info os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		if path == srcDir {
			return nil
		}

		fh, err := kzip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))
		if info.IsDir() {
			fh.Name += "/"
		}

		fh.Method = kzip.Deflate

		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		if !fh.Mode().IsRegular() {
			return nil
		}

		rc, err := os.Open(path)
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}

		err = zw.Flush()
		if err != nil {
			return err
		}

		fmt.Println(fh.Name, " completed !")
		return nil
	}

	err := filepath.Walk(srcDir, walkFunc)
	return err
}

func KUnpack(dstPath string, zipFileName string) error {
	archive, err := kzip.OpenReader(zipFileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dstPath, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dstPath)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}

		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func KZipToZip(srcZipPath string, dstZipPath string, ignore map[string]bool) error {
	r, err := kzip.OpenReader(srcZipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	f, err := os.Create(dstZipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := kzip.NewWriter(f)
	zw.RegisterCompressor(kzip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestSpeed)
	})
	defer zw.Close()

	for _, f := range r.File {
		fmt.Printf("Contents of %s: size: %d \n", f.Name, f.FileInfo().Size()/1024/1024)
		if _, ok := ignore[f.Name]; ok {
			fmt.Println("ignore file : ", f.Name)
			continue
		}

		rc, err := f.OpenRaw()
		if err != nil {
			return err
		}

		h := &kzip.FileHeader{
			Name:               f.Name,
			Method:             f.Method,
			Flags:              f.Flags,
			CRC32:              f.CRC32,
			CompressedSize64:   f.CompressedSize64,
			UncompressedSize64: f.UncompressedSize64,
		}
		w, err := zw.CreateRaw(h)
		if err != nil {
			return err
		}

		_, err = io.CopyN(w, rc, int64(f.CompressedSize64))
		if err != nil {
			return err
		}
	}

	tfName := "readme.txt"
	tfContent := "This archive contains some text files."
	tf, err := zw.Create(tfName)
	if err != nil {
		return err
	}
	_, err = tf.Write([]byte(tfContent))
	if err != nil {
		log.Fatal(err)
	}

	err = zw.SetComment("test comment")
	return err
}
