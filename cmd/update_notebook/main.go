package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var JupytextPipes = []string{
	"--pipe",
	"/workspace/.venv/bin/isort - --treat-comment-as-code \"# %%\" --float-to-top",
	"--pipe",
	"/workspace/.venv/bin/black -",
}

const JupytextPath = "/workspace/.venv/bin/jupytext"

func ExecJupytext(format, filepath string) error {
	cmd := exec.Command(JupytextPath, JupytextPipes...)
	cmd.Args = append(cmd.Args, "--set-formats", format)
	cmd.Args = append(cmd.Args, filepath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func DeleteJupytextMetadata(filepath string) error {
	log.Printf("delete jupytext notebook metadata: %s", filepath)
	re := regexp.MustCompile(`^\s*.*\-\-\-\s(.*\s)+\-\-\-\s+`)
	buffer, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	match := re.FindIndex(buffer)

	switch len(match) {
	case 0:
		log.Printf("not matched")
		return nil
	case 2:
		return os.WriteFile(filepath, buffer[match[1]:], 0664)

	default:
		return errors.New("unknown format: invalid submatch")
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("invalid args: %s", strings.Join(os.Args, " "))
	}

	targetFilepath := os.Args[1]
	targetFiledir := filepath.Dir(targetFilepath)
	targetFilename := filepath.Base(targetFilepath)
	targetFileext := filepath.Ext(targetFilepath)
	targetBasename := func() string {
		switch targetFileext {
		case ".py":
			re := regexp.MustCompile("^(.+)\\.py$")
			founds := re.FindAllStringSubmatch(targetFilename, -1)
			if len(founds) != 1 {
				log.Fatalf("cannot parse filename: %s", targetFilename)
			}
			return founds[0][1]

		case ".ipynb":
			re := regexp.MustCompile("^@(.+)\\.ipynb$")
			founds := re.FindAllStringSubmatch(targetFilename, -1)
			if len(founds) != 1 {
				log.Fatalf("cannot parse filename: %s", targetFilename)
			}
			return founds[0][1]

		default:
			log.Fatalf("unknown fileext: %s", targetFileext)
			return ""
		}
	}()
	switch targetBasename {
	case "__init__":
		format := "@__init__/ipynb,README/md:markdown,__init__/py:percent"
		destFilepathes := map[string]string{
			".py": filepath.Join(targetFiledir, targetBasename+".py"),
			".md": filepath.Join(targetFiledir, "README.md"),
		}

		if err := ExecJupytext(format, targetFilepath); err != nil {
			log.Fatal(err)
		}

		switch targetFileext {
		case ".ipynb":
			if err := ExecJupytext(format, destFilepathes[".py"]); err != nil {
				log.Fatal(err)
			}
		case ".py":
		default:
			log.Fatalf("unknown extension: %s", targetFileext)
		}

		for ext, destFilepath := range destFilepathes {
			if ext != ".ipynb" {
				if err := DeleteJupytextMetadata(destFilepath); err != nil {
					log.Fatal(err)
				}
			}
		}

	default:
		format := "@/ipynb,docs//md:markdown,py:percent"
		destFilepathes := map[string]string{
			".py": filepath.Join(targetFiledir, targetBasename+".py"),
			".md": filepath.Join(targetFiledir, "docs", targetBasename+".md"),
		}
		if err := ExecJupytext(format, targetFilepath); err != nil {
			log.Fatal(err)
		}

		switch targetFileext {
		case ".ipynb":
			if err := ExecJupytext(format, destFilepathes[".py"]); err != nil {
				log.Fatal(err)
			}
		case ".py":
		default:
			log.Fatalf("unknown extension: %s", targetFileext)
		}

		for ext, destFilepath := range destFilepathes {
			if ext != ".ipynb" {
				if err := DeleteJupytextMetadata(destFilepath); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
