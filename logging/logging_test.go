package logging_test

import (
	"kickcore/logging"
	"os"
	"testing"
)

func TestLevelSystem(t *testing.T) {
	l, _ := logging.NewLogger(logging.LEVEL_WARN, nil)

	if n := l.Log(logging.LEVEL_ERROR, "LEVEL_ERROR"); n == 0 {
		t.Fatal("LEVEL_ERROR")
	}

	if n := l.Log(logging.LEVEL_WARN, "LEVEL_WARN"); n == 0 {
		t.Fatal("LEVEL_WARN")
	}

	if n := l.Log(logging.LEVEL_INFO, "LEVEL_INFO"); n != 0 {
		t.Fatal("LEVEL_INFO")
	}
}

func TestFilename(t *testing.T) {
	tempdir := t.TempDir()

	l, _ := logging.NewLogger(logging.LEVEL_INFO, &logging.Config{Filename: tempdir + "/" + "filename_test.log"})

	l.Log(logging.LEVEL_WARN, "TestFilename")

	b, err := os.ReadFile(tempdir + "/" + "filename_test.log")
	if err != nil {
		t.Fatal(err)
	}

	// To Check Append Mode
	count := len(b)
	if count == 0 {
		t.Fatal("logging didn't write to file.")
	}

	l, _ = logging.NewLogger(logging.LEVEL_INFO, &logging.Config{
		Filename: tempdir + "/" + "filename_test.log", Append: true,
	})

	l.Log(logging.LEVEL_WARN, "TestAppend")

	b, err = os.ReadFile(tempdir + "/" + "filename_test.log")
	if err != nil {
		t.Fatal(err)
	}

	if len(b) <= count {
		t.Fatal("logging didn't append to file.")
	}
}

func TestFileObject(t *testing.T) {
	tempdir := t.TempDir()
	file, err := os.OpenFile(tempdir+"/"+"filename_test.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	l, _ := logging.NewLogger(logging.LEVEL_INFO, &logging.Config{FileObject: file})
	l.Log(logging.LEVEL_WARN, "TestFilename")

	file.Close()

	b, err := os.ReadFile(tempdir + "/" + "filename_test.log")
	if err != nil {
		t.Fatal(err)
	}

	// To Check Append Mode
	count := len(b)
	if count == 0 {
		t.Fatal("logging didn't write to file.")
	}
}
