#!/bin/bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $(basename "$0") [OPTIONS]

OPTION                DESCRIPTION                               REQUIRED
--help
-d, --dbname          database to dump                          V
-h, --host            database server host or socket directory  V
-p, --port            database server port number               V
-U, --username        connect as specified database user        V
-O, --output          output path                               V
-n, --schema          dump only schemas matching pattern
-N, --exclude-schema  do not dump any schemas matching pattern
-v, --verbose
-j, --jobs

EOF
}

main() {
  local dbname=''
  local host=''
  local port=''
  local username=''
  local output=''
  local schema=()
  local exclude_schema=()
  local verbose="true"
  local jobs="2"

  getopt_short_opts='d:h:p:U:O:n:N:vj:'
  getopt_long_opts='dbname:,host:,port:,username:,output:,schema:,exclude-schema:,verbose,jobs:,help'
  VALID_ARGS=$(getopt -o "${getopt_short_opts}" --long "${getopt_long_opts}" -- "$@")

  # shellcheck disable=SC2181
  if [ $? != 0 ]; then
    echo "error parsing options"
    usage
    exit 1
  fi

  eval set -- "$VALID_ARGS"
  while true; do
    case "$1" in
    -d | --dbname)
      dbname="${2}"
      shift 2
      ;;
    -h | --host)
      host="${2}"
      shift 2
      ;;
    -p | --port)
      port="${2}"
      shift 2
      ;;
    -U | --username)
      username="${2}"
      shift 2
      ;;
    -O | --output)
      output="${2}"
      shift 2
      ;;
    -n | --schema)
      schema+=("${2}")
      shift 2
      ;;
    -N | --exclude-schema)
      exclude_schema+=("${2}")
      shift 2
      ;;
    -v | --verbose)
      verbose=true
      shift
      ;;
    -j | --jobs)
      jobs="${2}"
      shift 2
      ;;
    --help)
      usage
      exit 0
      ;;
    --)
      shift
      break
      ;;
    *)
      echo "unexpected argument ${1}"
      usage
      exit 1
      ;;
    esac
  done

  # check remaining
  shift $((OPTIND - 1))
  remaining_args="${*}"
  if [ -n "${remaining_args}" ]; then
    printf "\n[error]: remaining args are not allowed: ${remaining_args[*]}\n\n"
    usage
    exit 1
  fi

  # check that required parameters were set
  local req_parameters=('dbname' 'host' 'port' 'username' 'output')
  for req_param in "${req_parameters[@]}"; do
    if [ -z "${!req_param:-}" ]; then
      printf "\n[error]: required parameter is not set: ${req_param}\n\n"
      usage
      exit 1
    fi
  done

  # debug variables
  echo "dbname=${dbname}"
  echo "host=${host}"
  echo "port=${port}"
  echo "username=${username}"
  echo "output=${output}"
  echo "schema=${schema[*]}"
  echo "exclude_schema=${exclude_schema[*]}"
  echo "verbose=${verbose}"
  echo "jobs=${jobs}"

}

main "${@}"
