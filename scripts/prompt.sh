#!/bin/bash

PROMPT_LENGTH=${PROMPT_LENGTH:=80}
PREFIX="----"

prompt() {
  local message="$1"
  local total_length=$PROMPT_LENGTH
  local line_char='-'
  local line_length=$((total_length - ${#message}))
  local line=$(printf "%*s" $line_length | tr ' ' $line_char)

  printf "%s %s %s\n" "$PREFIX" "$message" "$line"
}

end() {
  prompt "Done"
  echo
}
