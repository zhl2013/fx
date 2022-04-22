package reducer

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/antonmedv/fx/pkg/json"
	. "github.com/antonmedv/fx/pkg/theme"
)

func GenerateCode(args []string) string {
	lang, ok := os.LookupEnv("FX_LANG")
	if !ok {
		lang = "node"
	}
	switch lang {
	case "node":
		return nodejs(args)
	case "python", "python3":
		return python(args)
	case "ruby":
		return ruby(args)
	default:
		panic("unknown lang")
	}
}

func Reduce(object interface{}, args []string, theme Theme) {
	var cmd *exec.Cmd
	lang, ok := os.LookupEnv("FX_LANG")
	if !ok {
		lang = "node"
	}
	switch lang {
	case "node":
		cmd = CreateNodejs(args)
	case "python", "python3":
		cmd = CreatePython(lang, args)
	case "ruby":
		cmd = CreateRuby(args)
	default:
		panic("unknown lang")
	}

	cmd.Stdin = strings.NewReader(Stringify(object))
	output, err := cmd.CombinedOutput()
	if err == nil {
		dec := json.NewDecoder(bytes.NewReader(output))
		dec.UseNumber()
		jsonObject, err := Parse(dec)
		if err != nil {
			fmt.Print(string(output))
			return
		}
		if str, ok := jsonObject.(string); ok {
			fmt.Println(str)
		} else {
			fmt.Println(PrettyPrint(jsonObject, 1, theme))
		}
		if dec.InputOffset() < int64(len(output)) {
			fmt.Print(string(output[dec.InputOffset():]))
		}
	} else {
		exitCode := 1
		status, ok := err.(*exec.ExitError)
		if ok {
			exitCode = status.ExitCode()
		} else {
			fmt.Println(err.Error())
		}
		fmt.Print(string(output))
		os.Exit(exitCode)
	}
}

func trace(args []string, i int) (pre, post, pointer string) {
	pre = strings.Join(args[:i], " ")
	if len(pre) > 20 {
		pre = "..." + pre[len(pre)-20:]
	}
	post = strings.Join(args[i+1:], " ")
	if len(post) > 20 {
		post = post[:20] + "..."
	}
	pointer = fmt.Sprintf(
		"%v %v %v",
		strings.Repeat(" ", len(pre)),
		strings.Repeat("^", len(args[i])),
		strings.Repeat(" ", len(post)),
	)
	return
}
