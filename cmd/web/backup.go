package main

import (
	"fmt"
	"os"
	"zombiezen.com/go/sqlite"
)

// backup creates a backup of the database
func backup() error {
	// Open database connections.
	src, err := sqlite.OpenConn(os.Getenv("SRC_DB"))
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := sqlite.OpenConn(os.Getenv("DST_DB"))
	if err != nil {
		return err
	}
	defer func() {
		if err := dst.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create Backup object.
	backup, err := sqlite.NewBackup(dst, "main", src, "main")
	if err != nil {
		return err
	}

	// Perform online backup/copy.
	_, err1 := backup.Step(-1)
	err2 := backup.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// Backup_Step shows how to use Step
// to prevent holding a read lock on the source database
// during the entire copy.
//func Backup_Step() {
//	// Open database connections.
//	src, err := sqlite.OpenConn(os.Args[1])
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//	defer src.Close()
//	dst, err := sqlite.OpenConn(os.Args[2])
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//	defer func() {
//		if err := dst.Close(); err != nil {
//			fmt.Fprintln(os.Stderr, err)
//		}
//	}()
//
//	// Create Backup object.
//	backup, err := sqlite.NewBackup(dst, "main", src, "main")
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//	defer func() {
//		if err := backup.Close(); err != nil {
//			fmt.Fprintln(os.Stderr, err)
//		}
//	}()
//
//	// Each iteration of this loop copies 5 database pages,
//	// waiting 250ms between iterations.
//	for {
//		more, err := backup.Step(5)
//		if !more {
//			if err != nil {
//				fmt.Fprintln(os.Stderr, err)
//				os.Exit(1)
//			}
//			break
//		}
//		time.Sleep(250 * time.Millisecond)
//	}
//}
