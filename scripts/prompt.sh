#!/bin/bash

PROMPT_LENGTH=80
LINE=" --------- "

prompt() {
  local message="$1"
  local total_length=PROMPT_LENGTH
  local line_char='-'
  local line_length=$(( (total_length - ${#message}) / 2 ))
  local line=$(printf "%*s" $line_length | tr ' ' $line_char)

  # Print the generated prompt
  printf "%s %s %s\n" "$line" "$message" "$line"
}
