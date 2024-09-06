#!/bin/bash
set -euo pipefail

usage() {
	cat <<EOF
Usage: $(basename "$0") [OPTIONS]

OPTION                DESCRIPTION                             
--help                                                        
-d, --dbname          database to dump                        
-h, --host            database server host or socket directory
-p, --port            database server port number             
-U, --username        connect as specified database user      
-O, --output          output path                             
-n, --schema          dump only schemas matching pattern      
-N, --exclude-schema  do not dump any schemas matching pattern
-v, --verbose                                                 

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

  VALID_ARGS=$(getopt -o d:h:p:U:O:n:N:v --long dbname:,host:,port:,username:,output:,schema:,exclude-schema:,verbose,help -- "$@")
  # shellcheck disable=SC2181
  if [ $? != 0 ]; then
    echo "error parsing options: $?"
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
    echo "remaining args are not allowed: ${remaining_args[*]}"
    usage
    exit 1
  fi
  # set checks
  if [ -z "${dbname}" ]; then
    printf "\n[error] required arg is empty: dbname\n\n"
    usage
    exit 1
  fi
  if [ -z "${host}" ]; then
    printf "\n[error] required arg is empty: host\n\n"
    usage
    exit 1
  fi
  if [ -z "${port}" ]; then
    printf "\n[error] required arg is empty: port\n\n"
    usage
    exit 1
  fi
  if [ -z "${username}" ]; then
    printf "\n[error] required arg is empty: username\n\n"
    usage
    exit 1
  fi
  if [ -z "${output}" ]; then
    printf "\n[error] required arg is empty: output\n\n"
    usage
    exit 1
  fi

  # debug variables
  echo "dbname=${dbname}"
  echo "host=${host}"
  echo "port=${port}"
  echo "username=${username}"
  echo "output=${output}"
  echo "schema=${schema[*]}"
  echo "exclude_schema=${exclude_schema[*]}"
  echo "verbose=${verbose}"

}

main "${@}"


