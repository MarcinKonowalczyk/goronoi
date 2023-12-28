# https://github.com/MarcinKonowalczyk/run_sh
# Bash script run by a keyboard shortcut, called with the current file path $1
# This is intended as an example, but also contains a bunch of useful path partitions
# Feel free to delete everything in here and make it do whatever you want.

echo "Hello from run script! ^_^"

_VERSION="0.2" # Version of this script

# The directory of the main project from which this script is running
# https://stackoverflow.com/a/246128/2531987
ROOT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
ROOT_FOLDER="${ROOT_FOLDER%/*}" # Strip .vscode folder
PROJECT_NAME="${ROOT_FOLDER##*/}" # Project name

FULL_FILE_PATH="$1"
_RELATIVE_FILE_PATH="${FULL_FILE_PATH##*$ROOT_FOLDER/}" # Relative path of the current file

# Split the relative file path into an array
RELATIVE_PATH_PARTS=(${_RELATIVE_FILE_PATH//\// })
DEPTH=${#RELATIVE_PATH_PARTS[@]}
DEPTH=$((DEPTH-1))

# Couple of useful variables
FILENAME="${RELATIVE_PATH_PARTS[$DEPTH]}"

# If the file has an extension, get it otherwise set it to empty string
EXTENSION="" && [[ "$FILENAME" == *.* ]] && EXTENSION="${FILENAME##*.}"

########################################

GREEN='\033[0;32m'; YELLOW='\033[0;33m'; RED='\033[0;31m'; NC='\033[0m'; PURPLE='\033[0;34m'; DARK_GRAY='\033[1;30m';

function logo() {
    TEXT=(
        " ______   __  __   __   __ " "    " "    ______   __  __   "
        "/\\  == \\ /\\ \\/\\ \\ /\\ \"-.\\ \\ " "    " "  /\\  ___\\ /\\ \\_\\ \\  "
        "\\ \\  __< \\ \\ \\_\\ \\\\\\ \\ \\-.  \\ " "  __" " \\ \\___  \\\\\\ \\  __ \\ "
        " \\ \\_\\ \\_\\\\\\ \\_____\\\\\\ \\_\\\\\\\"\\_\\ " "/\\_\\\\" " \\/\\_____\\\\\\ \\_\\ \\_\\\\"
        "  \\/_/ /_/ \\/_____/ \\/_/ \\/_/ " "\\/_/" "  \\/_____/ \\/_/\\/_/"
    )
    echo -e "$PURPLE${TEXT[0]}$DARK_GRAY${TEXT[1]}$PURPLE${TEXT[2]}$NC"
    echo -e "$PURPLE${TEXT[3]}$DARK_GRAY${TEXT[4]}$PURPLE${TEXT[5]}$NC"
    echo -e "$PURPLE${TEXT[6]}$DARK_GRAY${TEXT[7]}$PURPLE${TEXT[8]}$NC"
    echo -e "$PURPLE${TEXT[9]}$DARK_GRAY${TEXT[10]}$PURPLE${TEXT[11]}$NC"
    echo -e "$PURPLE${TEXT[12]}$DARK_GRAY${TEXT[13]}$PURPLE${TEXT[14]} ${DARK_GRAY}v${_VERSION}$NC"
    echo ""
}

function info() {
    echo -e "PROJECT_NAME        : $GREEN${PROJECT_NAME}$NC  # project name (name of the project folder)"
    echo -e "RELATIVE_PATH_PARTS : $GREEN${RELATIVE_PATH_PARTS[@]}$NC  # relative path of the current file split into an array"
    echo -e "DEPTH               : $GREEN${DEPTH}$NC  # depth of the current file (number of folders deep)"
    echo -e "FILENAME            : $GREEN${FILENAME}$NC  # just the filename (equivalent to RELATIVE_PATH_PARTS[DEPTH])"
    echo -e "EXTENSION           : $GREEN${EXTENSION}$NC  # just the extension of the current file"
    echo -e "ROOT_FOLDER         : $GREEN${ROOT_FOLDER}$NC  # full path to the root folder of the project"
    echo -e "FULL_FILE_PATH      : $GREEN${FULL_FILE_PATH}$NC  # full path of the current file"
}


# VERBOSE=true
VERBOSE=false
[ "${RELATIVE_PATH_PARTS[0]}" = ".vscode" ] && [ ${RELATIVE_PATH_PARTS[$DEPTH]} = "run.sh" ] && [ $DEPTH -eq 1 ] && VERBOSE=true
if $VERBOSE; then logo; info; fi

########################################

go run .