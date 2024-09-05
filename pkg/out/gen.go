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
	for k, v := range o.Opts {
		varname := getVariableNameFromKey(k)
		if v.DefaultValue != "" {
			res += fmt.Sprintf("  local %s=\"%s\"\n", varname, v.DefaultValue)
		} else {
			// if default value is not specified, set it as empty string
			res += fmt.Sprintf("  local %s=''\n", varname)
		}
	}
	return res
}

// Gen short option. If 'Short' field is specified, it'll be used,
// first letter of option instead
func getOneShort(k string, v cfg.Opt) string {
	shortOpts := ""
	if v.Short != "" {
		shortOpts += v.Short
	} else {
		shortOpts += string(k[0])
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
	for k, v := range o.Opts {
		shortOpts += getOneShort(k, v)
		longOpts += k

		if !v.Flag {
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
	for k, v := range o.Opts {
		if v.Optional {
			continue
		}
		varname := getVariableNameFromKey(k)

		res += fmt.Sprintf("  if [ -z \"${%s}\" ]; then\n", varname)
		res += fmt.Sprintf("    printf \"\\n[error] required arg is empty: %s\\n\\n\"\n", k)
		res += fmt.Sprintf("    usage\n")
		res += fmt.Sprintf("    exit 1\n")
		res += fmt.Sprintf("  fi\n")
	}
	return res + "\n"
}

func getMaxPadding(o cfg.Opts) int {
	max := 0
	for k, v := range o.Opts {
		s := fmt.Sprintf("-%s, --%s\n", getOneShort(k, v), k)
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
	for k, v := range o.Opts {
		sh := getOneShort(k, v)

		if v.Desc != "" {
			pad := getPadding(fmt.Sprintf("-%s, --%s\n", sh, k), maxPad)
			optsDesc += fmt.Sprintf("  -%s, --%s %s %s\n", sh, k, pad, v.Desc)
		} else {
			optsDesc += fmt.Sprintf("  -%s, --%s\n", sh, k)
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

	for k, v := range opts.Opts {
		varname := getVariableNameFromKey(k)

		oneOpt := ""
		if v.Flag {
			oneOpt += fmt.Sprintf("    -%s | --%s)\n", getOneShort(k, v), k)
			oneOpt += fmt.Sprintf("      %s=true\n", varname)
			oneOpt += "      shift\n"
			oneOpt += "      ;;\n"
			res += oneOpt
		} else {
			oneOpt += fmt.Sprintf("    -%s | --%s)\n", getOneShort(k, v), k)
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
	for k := range opts.Opts {
		varname := getVariableNameFromKey(k)
		res += fmt.Sprintf(`  echo "%s=${%s}"`+"\n", varname, varname)
	}

	res += "}" + "\n\n"
	res += `main "${@}"` + "\n\n"

	return res
}
