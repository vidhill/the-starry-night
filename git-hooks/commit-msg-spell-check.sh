#!/bin/sh

# slightly modified version from example at https://phrye.com/tools/git-spell-check/ 

if ! type aspell > /dev/null; then
  echo "Warning: aspell is not installed"
  exit
fi

words="$(grep -v '^#' "$1" \
         | aspell --camel-case -d en_GB "--personal=$HOME/.config/git/.dict.en.pws" list)"
if [ -n "$words" ]; then
  cmd="git"
  msg="$1"
  vcsh_repo="$VCSH_REPO_NAME"
  if [ -z "$VCSH_REPO_NAME" ]; then
    vcsh_repo="$VCSH_DIRECTORY"
  fi
  if [ -n "$vcsh_repo" ]; then
    if [ "$vcsh_repo" = past ]; then
      exit
    fi
    cmd="vcsh $vcsh_repo"
    msg="$HOME/$1"
  fi
  echo "Error: Spellcheck failure in commit message"
  echo "- the following word(s) are misspelled"
  echo "  $words" | sed 's/^/  /'
  echo ""
  echo "To add words to personal dictionary run: aspell-add-word.sh WORD"
  exit 1
fi

self=commit-msg
if [ -f "$GIT_DIR/hooks/$self" ]; then
  exec "$GIT_DIR/hooks/$self" "$@"
fi


