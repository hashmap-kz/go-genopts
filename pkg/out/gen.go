package out

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashmap-kz/go-genopts/pkg/cfg"
	"github.com/hashmap-kz/go-genopts/pkg/util"
	"github.com/hashmap-kz/go-texttable/pkg/table"
)

func f(pad int, format string, a ...any) string {
	ws := strings.Repeat(" ", pad)
	return ws + fmt.Sprintf(format, a...) + "\n"
}

func p(pad int, arg string) string {
	ws := strings.Repeat(" ", pad)
	return ws + arg + "\n"
}

// Declare local variables, set empty values:
// local myvar=""
// local myarr=()
// local myopt=false
func genLocals(o cfg.Opts) string {
	res := ""
	for _, k := range o.Opts {
		varname := getVariableNameFromKey(k.Name)
		if k.DefaultValue != "" {
			res += f(2, "local %s=\"%s\"", varname, k.DefaultValue)
		} else {
			if k.Type == cfg.OptTypeList {
				res += f(2, "local %s=()", varname)
			} else if k.Type == cfg.OptTypeBool {
				res += f(2, "local %s='false'", varname)
			} else {
				// if default value is not specified, set it as empty string
				res += f(2, "local %s=''", varname)
			}
		}
	}
	return res
}

// Gen short option. If 'Short' field is specified, it'll be used,
// first letter of option instead
func getOneShort(v cfg.Opt) string {
	shortOpts := ""
	if v.Short != "" {
		shortOpts += v.Short
	} else {
		shortOpts += string(v.Name[0])
	}
	return shortOpts
}

func getVariableNameFromKey(k string) string {
	ident := strings.ReplaceAll(k, "-", "_")
	if !util.NameIsValidIdentifier(ident) {
		log.Fatalf("expect identifier, got: %s", k)
	}
	return ident
}

// Generate opt keys suitable for this format:
// VALID_ARGS=$(getopt -o h:p:u:d: --long host:,port:,username:,dest: -- "$@")
func genOpts(o cfg.Opts) (string, string) {
	shortOpts := ""
	longOpts := ""
	i := 0
	for _, k := range o.Opts {
		shortOpts += getOneShort(k)
		longOpts += k.Name

		if k.Type != cfg.OptTypeBool {
			shortOpts += ":"
			longOpts += ":"
		}

		if i+1 < len(o.Opts) {
			longOpts += ","
		}
		i++
	}

	// note: special handling for '--help'
	if len(longOpts) > 0 {
		longOpts += ",help"
	} else {
		longOpts += "help"
	}

	return shortOpts, longOpts
}

func genChecks(o cfg.Opts) string {
	res := ""
	for _, k := range o.Opts {
		if k.Optional {
			continue
		}
		varname := getVariableNameFromKey(k.Name)

		if k.Type == cfg.OptTypeList {
			res += f(2, `if [ -z "${%s[*]}" ]; then`, varname)
		} else {
			res += f(2, `if [ -z "${%s}" ]; then`, varname)
		}

		res += f(4, `printf "\n[error] required arg is empty: %s\n\n"`, k.Name)
		res += p(4, "usage")
		res += p(4, "exit 1")
		res += p(2, "fi")
	}
	return res + "\n"
}

func genDebugVarsEcho(o cfg.Opts) string {
	res := ""
	for _, k := range o.Opts {
		varname := getVariableNameFromKey(k.Name)
		if k.Type == cfg.OptTypeList {
			res += f(2, `echo "%s=${%s[*]}"`, varname, varname)
		} else {
			res += f(2, `echo "%s=${%s}"`, varname, varname)
		}
	}
	return res + "\n"
}

func req(k cfg.Opt) string {
	if k.Optional {
		return ""
	}
	return "V"
}

func genUsage(o cfg.Opts) string {

	optsDesc := "usage() {\n"
	optsDesc += "	cat <<EOF\n"
	optsDesc += `Usage: $(basename "$0") [OPTIONS]` + "\n\n"

	// pretty print options in a table-based style
	tbl := table.NewTextTable()
	tbl.DefineColumn("OPTION", table.LEFT, table.LEFT)
	tbl.DefineColumn("DESCRIPTION", table.LEFT, table.LEFT)
	tbl.DefineColumn("REQUIRED", table.LEFT, table.LEFT)

	// note: special handling for '--help'
	tbl.InsertAll("--help")
	tbl.EndRow()

	for _, k := range o.Opts {
		sh := getOneShort(k)

		if k.Desc != "" {
			tbl.InsertAll(fmt.Sprintf("-%s, --%s", sh, k.Name), k.Desc, req(k))
			tbl.EndRow()
		} else {
			tbl.InsertAll(fmt.Sprintf("-%s, --%s", sh, k.Name), req(k))
			tbl.EndRow()
		}
	}

	optsDesc += tbl.Print() + "\n"
	optsDesc += "EOF\n"
	optsDesc += "}\n\n"
	return optsDesc
}

func GenOpts(opts cfg.Opts) string {

	res := "#!/bin/bash\n"
	res += "set -euo pipefail\n\n"

	res += genUsage(opts)
	res += "main() {\n"

	// declare empty local vars
	res += fmt.Sprintln(genLocals(opts))

	// declare options list
	shorts, longs := genOpts(opts)
	res += f(2, "getopt_short_opts='%s'", shorts)
	res += f(2, "getopt_long_opts='%s'", longs)
	res += p(2, `VALID_ARGS=$(getopt -o "${getopt_short_opts}" --long "${getopt_long_opts}" -- "$@")`)
	res += "\n"

	// options header (always the same)
	res += p(2, `# shellcheck disable=SC2181`)
	res += p(2, `if [ $? != 0 ]; then`)
	res += p(2, `  echo "error parsing options: $?"`)
	res += p(2, `  usage`)
	res += p(2, `  exit 1`)
	res += p(2, `fi`)
	res += p(2, ``)
	res += p(2, `eval set -- "$VALID_ARGS"`)
	res += p(2, `while true; do`)
	res += p(2, `  case "$1" in`)

	// options cases
	for _, k := range opts.Opts {
		varname := getVariableNameFromKey(k.Name)

		oneOpt := ""
		if k.Type == cfg.OptTypeBool {
			oneOpt += f(4, "-%s | --%s)", getOneShort(k), k.Name)
			oneOpt += f(6, "%s=true", varname)
			oneOpt += p(6, "shift")
			oneOpt += p(6, ";;")
			res += oneOpt
		} else if k.Type == cfg.OptTypeList {
			oneOpt += f(4, "-%s | --%s)", getOneShort(k), k.Name)
			oneOpt += f(6, `%s+=("${2}")`, varname)
			oneOpt += p(6, "shift 2")
			oneOpt += p(6, ";;")
			res += oneOpt
		} else {
			oneOpt += f(4, "-%s | --%s)", getOneShort(k), k.Name)
			oneOpt += f(6, `%s="${2}"`, varname)
			oneOpt += p(6, "shift 2")
			oneOpt += p(6, ";;")
			res += oneOpt
		}
	}

	// note: special handling for '--help'
	// always add help (as a long option)
	res += p(4, "--help)")
	res += p(6, "usage")
	res += p(6, "exit 0")
	res += p(6, ";;")

	// options footer (always the same)
	res += p(2, `  --)`)
	res += p(2, `    shift`)
	res += p(2, `    break`)
	res += p(2, `    ;;`)
	res += p(2, `  *)`)
	res += p(2, `    echo "unexpected argument ${1}"`)
	res += p(2, `    usage`)
	res += p(2, `    exit 1`)
	res += p(2, `    ;;`)
	res += p(2, `  esac`)
	res += p(2, `done`)
	res += p(2, ``)
	res += p(2, `# check remaining`)
	res += p(2, `shift $((OPTIND - 1))`)
	res += p(2, `remaining_args="${*}"`)
	res += p(2, `if [ -n "${remaining_args}" ]; then`)
	res += p(2, `  echo "remaining args are not allowed: ${remaining_args[*]}"`)
	res += p(2, `  usage`)
	res += p(2, `  exit 1`)
	res += p(2, `fi`)
	res += "\n"

	// generate checks for required variables
	res += p(2, "# set checks")
	res += genChecks(opts)

	// generate simple echo output of each parameter, for debug and testing
	res += p(2, "# debug variables")
	res += genDebugVarsEcho(opts)

	res += "}" + "\n\n"
	res += `main "${@}"` + "\n\n"

	return res
}
