#!/usr/bin/env bash
RED=`tput setaf 1`
GREEN=`tput setaf 2`
YELLOW=`tput setaf 3`
BLUE=`tput setaf 4`
PURPLE=`tput setaf 5`
CYAN=`tput setaf 6`
BRIGHT=`tput setaf 7`

BRED=`tput bold && tput setaf 1`
BGREEN=`tput bold && tput setaf 2`
BYELLOW=`tput bold && tput setaf 3`
BBLUE=`tput bold && tput setaf 4`
BPURPLE=`tput bold && tput setaf 5`
BCYAN=`tput bold && tput setaf 6`
BBRIGHT=`tput bold && tput setaf 7`

RESET=`tput sgr0`
exitfn () {
    trap SIGINT              
    echo; echo 'Interrupted by user!!'
    exit                     
}
depth=1
dry_run=0
folder="${!#}"
paramnum=$#
email=""
upload='default'
optspec=":h-:"
other_args=''


trap "exitfn" INT

local_projects() {
  number_of_repos=0
  concatenated_repos=()
  for REPO_PATH in `find $folder -maxdepth $depth  -type d  -wholename "*/.git" -exec dirname {} \;`; do
    concatenated_repos+=("$REPO_PATH")
    #total_repos=$((total_repos+1))
    ((number_of_repos+=1))
    REPO_NAME=`echo $REPO_PATH | sed -e 's/\/$//g' | rev | cut -d'/' -f 1 | rev `  >&2;
    
    if [ ! $dry_run -eq 0 ];then 
      echo "     ${BGREEN}✔️${RESET}  ${BBRIGHT}$REPO_NAME${RESET}"
    fi
  done
    
  if [ $number_of_repos -eq 0 ]; then
    echo "${BRED}✖${RESET} We couldn't find any GIT repos under path ${CYAN}$folder${RESET}"
    echo " "  
    return
  fi 
  echo "Found ${GREEN}$number_of_repos${RESET} repos under ${CYAN}$folder${RESET}"
  echo " "
  if [ ! $dry_run -eq 0 ];then 
      echo " "
      echo "  You specified the ${BBRIGHT}--dry${RESET} flag."
      return
  fi
  #repostring=`printf "%s➕" "${concatenated_repos[@]}" | sed -e 's/➕$//g' `
  repostring=`printf "%s|,|" "${concatenated_repos[@]}" | sed -e 's/|,|$//g' `

  python src/main.py "$repostring" $other_args ;\
  wait ;\
  return
}

while getopts "$optspec" optchar; do
     case "${optchar}" in
        -)
            case "${OPTARG}" in
                dry)
                    dry_run=1
                    ;;
                depth=*)
                    val=${OPTARG#*=}
                    opt=${OPTARG%=$val}
                    depth=$val
                    ;;
                *)
                    other_args="$other_args --$OPTARG"
                    ;;
            esac;;
        h)
            echo -e ""
            echo -e "collect your repos info using: "
            echo -e ""
            echo -e " ${GREEN}./run.sh${RESET} [${BBRIGHT}--dry${RESET}] [${BBRIGHT}--depth=2${RESET}] [${BBRIGHT}--email=user@domain.com${RESET}] [${BBRIGHT}--skip_upload${RESET}] ${RESET} ${BCYAN}<folder>${RESET} "
            echo -e ""
            echo -e "Examples:"
            echo -e "${GREEN}./run.sh help${RESET}                               display this help message";echo
            echo -e "${GREEN}./run.sh${RESET} ${BCYAN}~/my_repo${RESET}                   parse info of ${BCYAN}my_repo${RESET}"
            echo -e "${GREEN}./run.sh${RESET} ${BBRIGHT}--dry${RESET}  ${BCYAN}~/my_repo${RESET}            just display repos that would be examined"
            echo -e "${GREEN}./run.sh${RESET} ${BBRIGHT}--skip_upload${RESET}  ${BCYAN}~/my_repo${RESET}    skip auto upload prompt"
            echo -e "${GREEN}./run.sh${RESET} --email=${BBRIGHT}me@mail.com${RESET} ${BCYAN}~/repo${RESET}  preselect ${BCYAN}me@mail.com${RESET} in author list"
            
            echo; echo -e "${GREEN}./run.sh${RESET} ${BBRIGHT}--depth=2${RESET}  ${BCYAN}~/projects${RESET}        sarch recursively N levels from ${BCYAN}~/projects${RESET}"
            echo -e "               (this option replaces the ${BRIGHT}--output${RESET} parameter with folder names)"
            echo -e ""
            echo -e "Combine them all (options available in ${GREEN}main.py${RESET} will be forwared untouched)";echo
            echo -e "${GREEN}./run.sh ${RESET} --depth=2 --parse_libraries --email=me@mail.com --skip_upload  ${BCYAN}~/projects${RESET}"

            echo ""
            exit 0
            ;;
       
            
      
        *)
            if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
                echo "Non-option argument: '-${OPTARG}'" >&2
            fi
            ;;
    esac
done 

  echo "Search repos on ${BPURPLE}$folder${RESET}, ${BCYAN}$depth${RESET} folders deep "
  #$other_args"
  local_projects
  echo "Finished"
  echo " "
trap SIGINT
exit 0