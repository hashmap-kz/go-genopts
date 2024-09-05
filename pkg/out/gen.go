package out

import (
	"fmt"
	"strings"

	"github.com/hashmap-kz/go-genopts/pkg/cfg"
)

// Declare local variables, set empty values:
// local myvar=""
func genLocals(o cfg.Opts) string {
	res := ""
	for _, k := range o.Opts {
		varname := getVariableNameFromKey(k.Name)
		if k.DefaultValue != "" {
			res += fmt.Sprintf("  local %s=\"%s\"\n", varname, k.DefaultValue)
		} else {
			if k.Type == cfg.OptTypeList {
				res += fmt.Sprintf("  local %s=()\n", varname)
			} else if k.Type == cfg.OptTypeBool {
				res += fmt.Sprintf("  local %s='false'\n", varname)
			} else {
				// if default value is not specified, set it as empty string
				res += fmt.Sprintf("  local %s=''\n", varname)
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
	return strings.ReplaceAll(k, "-", "_")
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
			res += fmt.Sprintf("  if [ -z \"${%s[*]}\" ]; then\n", varname)
		} else {
			res += fmt.Sprintf("  if [ -z \"${%s}\" ]; then\n", varname)
		}

		res += fmt.Sprintf("    printf \"\\n[error] required arg is empty: %s\\n\\n\"\n", k.Name)
		res += "    usage\n"
		res += "    exit 1\n"
		res += "  fi\n"
	}
	return res + "\n"
}

func getMaxPadding(o cfg.Opts) int {
	max := 0
	for _, k := range o.Opts {
		s := fmt.Sprintf("-%s, --%s\n", getOneShort(k), k.Name)
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

func getPadding(what string, max int) string {
	diff := max - len(what) + 1
	return strings.Repeat(" ", diff)
}

func genUsage(o cfg.Opts) string {
	maxPad := getMaxPadding(o)

	optsDesc := "usage() {\n"
	optsDesc += "	cat <<EOF\n"
	optsDesc += `Usage: $(basename "$0") [OPTION]` + "\n\n"
	optsDesc += "Options:\n"
	for _, k := range o.Opts {
		sh := getOneShort(k)

		if k.Desc != "" {
			pad := getPadding(fmt.Sprintf("-%s, --%s\n", sh, k.Name), maxPad)
			optsDesc += fmt.Sprintf("  -%s, --%s %s %s\n", sh, k.Name, pad, k.Desc)
		} else {
			optsDesc += fmt.Sprintf("  -%s, --%s\n", sh, k.Name)
		}
	}

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
	validArgs := fmt.Sprintf(`  VALID_ARGS=$(getopt -o %s --long %s -- "$@")`, shorts, longs)
	res += fmt.Sprintln(validArgs)

	hdr := `
	# shellcheck disable=SC2181
	if [ $? != 0 ]; then
		printf "error parsing options"
		usage
		exit 1
	fi

	eval set -- "$VALID_ARGS"
	while true; do
		case "$1" in
		`
	res += hdr + "\n"

	for _, k := range opts.Opts {
		varname := getVariableNameFromKey(k.Name)

		oneOpt := ""
		if k.Type == cfg.OptTypeBool {
			oneOpt += fmt.Sprintf("    -%s | --%s)\n", getOneShort(k), k.Name)
			oneOpt += fmt.Sprintf("      %s=true\n", varname)
			oneOpt += "      shift\n"
			oneOpt += "      ;;\n"
			res += oneOpt
		} else if k.Type == cfg.OptTypeList {
			oneOpt += fmt.Sprintf("    -%s | --%s)\n", getOneShort(k), k.Name)
			oneOpt += fmt.Sprintf(`      %s+=("${2}")`+"\n", varname)
			oneOpt += "      shift 2\n"
			oneOpt += "      ;;\n"
			res += oneOpt
		} else {
			oneOpt += fmt.Sprintf("    -%s | --%s)\n", getOneShort(k), k.Name)
			oneOpt += fmt.Sprintf(`      %s="${2}"`+"\n", varname)
			oneOpt += "      shift 2\n"
			oneOpt += "      ;;\n"
			res += oneOpt
		}
	}

	ftr := `
		--)
			shift
			break
			;;
		*)
			printf "unexpected argument ${1}"
			usage
			exit 1
			;;
		esac
	done

	# check remaining
	shift $((OPTIND - 1))
	remaining_args="${*}"
	if [ -n "${remaining_args}" ]; then
		echo "remaining args are not allowed: ${remaining_args[*]}"
		usage
		exit 1
	fi
		`
	res += ftr + "\n"

	res += "  # set checks\n"
	res += genChecks(opts)

	res += "  # debug variables\n"
	for _, k := range opts.Opts {
		varname := getVariableNameFromKey(k.Name)
		if k.Type == cfg.OptTypeList {
			res += fmt.Sprintf(`  echo "%s=${%s[*]}"`+"\n", varname, varname)
		} else {
			res += fmt.Sprintf(`  echo "%s=${%s}"`+"\n", varname, varname)
		}
	}

	res += "}" + "\n\n"
	res += `main "${@}"` + "\n\n"

	return res
}
