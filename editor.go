package noteofcli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
)

var wsre = regexp.MustCompile("\\s")

func Edit(editor, text string) ([]byte, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) == 0 || editor == "" {
		return ioutil.ReadAll(os.Stdin)
	} else if editor != "" {
		return ExecEditor(editor, text)
	}

	return []byte{}, fmt.Errorf("editor error")
}

func ExecEditor(editor, text string) ([]byte, error) {
	parts := wsre.Split(editor, -1)

	tmpfile, err := ioutil.TempFile("", "post")
	tmpfile.WriteString(text)
	tmpPath := tmpfile.Name()
	tmpfile.Close()
	if err != nil {
		return []byte{}, err
	}

	args := []string{}
	if len(parts) > 1 {
		args = append(args, parts[1:]...)
	}

	args = append(args, tmpPath)

	log.Println(parts[0], args)
	cmd := exec.Command(parts[0], args...)
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println()

	body, err := ioutil.ReadFile(tmpPath)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}
