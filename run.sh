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


cr_green=`tput setaf 73`
cr_blue=`tput setaf 25`

RESET=`tput sgr0`
exitfn () {
    trap SIGINT              
    echo; echo 'Interrupted by user!!'
    exit                     
}

trap "exitfn" INT

print_logo() {
  echo
  echo -e "${cr_blue}                                          ':'   ${cr_green}                                ${RESET}"
  echo -e "${cr_blue}                                       .+ymo     ${cr_green}     .+/:.                     ${RESET}"
  echo -e "${cr_blue}                                    -ohmmms'     ${cr_green}    '++++++/-'                 ${RESET}"
  echo -e "${cr_blue}                                ':sdmmmmmh'     ${cr_green}    '/+++++++++/-.              ${RESET}"
  echo -e "${cr_blue}                             '/ymmmmmmmmd.     ${cr_green}     /++++++++++++++:-'          ${RESET}"
  echo -e "${cr_blue}                          .+hmmmmmmmmmmm-       ${cr_green}   :++++++++++++++++++/-'       ${RESET}"
  echo -e "${cr_blue}                       -odmmmmmmmmmmmmm/        ${cr_green}   -:+++++++++++++++++++++:.    ${RESET}"
  echo -e "${cr_blue}                   ':sdmmmmmmmmmmmmmmm+         ${cr_green}      .-/+++++++++++++++++++/'  ${RESET}"
  echo -e "${cr_blue}                ./ymmmmmmmmmmmmmmmmmmo         ${cr_green}          '-:+++++++++++++++++-  ${RESET}"
  echo -e "${cr_blue}             -+hmmmmmmmmmmmmmmmmmds/'         ${cr_green}              /++++++++++++++++-  ${RESET}"
  echo -e "${cr_blue}         ':odmmmmmmmmmmmmmmmmmho-           ${cr_green}            '.:++++++++++++++++++.  ${RESET}"
  echo -e "${cr_blue}      '/smmmmmmmmmmmmmmmmmmy/.                ${cr_green}       '-/++++++++++++++++++/-'   ${RESET}"
  echo -e "${cr_blue}    /ymmmmmmmmmmmmmmmmmds:'                ${cr_green}      '.:+++++++++++++++++++:-'      ${RESET}"
  echo -e "${cr_blue}   ommmmmmmmmmmmmmmmd+-                   ${cr_green}    .-/++++++++++++++++++/:.          ${RESET}"
  echo -e "${cr_blue}   ymmmmmmmmmmmmmmmm+                    ${cr_green}  .:+++++++++++++++++++/-'             ${RESET}"
  echo -e "${cr_blue}   ymmmmmmmmmmmmmmmmmds/'              ${cr_green}   :+++++++++++++++++/:.'                ${RESET}"
  echo -e "${cr_blue}   .hmmmmmmmmmmmmmmmmmmmmho:'         ${cr_green}   -+++++++++++++++++:                    ${RESET}"
  echo -e "${cr_blue}     .+ymmmmmmmmmmmmmmmmmmmmms       ${cr_green}   -+++++++++++++++++++:'                  ${RESET}"
  echo -e "${cr_blue}        ':sdmmmmmmmmmmmmmmmmd-      ${cr_green}   .+++++++++++++++++++++/'                 ${RESET}"
  echo -e "${cr_blue}            .+ymmmmmmmmmmmmm:       ${cr_green}  .++++++++/-/+++++++++++++-                ${RESET}"
  echo -e "${cr_blue}               ':ohmmmmmmmm/      ${cr_green}   '++++++:.   ':+++++++++++++:'              ${RESET}"
  echo -e "${cr_blue}                   ./smmmmo      ${cr_green}   '/++/-'        .+++++++++++++/'             ${RESET}"
  echo -e "${cr_blue}                       -++      ${cr_green}    /:.'            ':++++++++++++/             ${RESET}"
  echo
}
parameter_definition() {
  echo -e " ${BBRIGHT} $1 ${RESET}    $2  $3"
}

print_help() {
  echo -e ""
  echo -e "collect your repos info using: "
  echo -e ""
  echo -e " ${GREEN}./run.sh${RESET} [--help | options]  ${BCYAN}<folder>${RESET} "

  echo -e "Parameters:"
  echo -e "   ${BCYAN}folder${RESET}            local folder of the repo to analyze ${BBRIGHT}required${RESET}"
  
  echo -e "Options:";
  
  parameter_definition "--help                " "display this help message";echo
  parameter_definition "--output=<output_file>" "json output file, defaults to ${BBRIGHT}repo_data.json${RESET}"
  parameter_definition "--dry                 " "display repos found in ${BCYAN}<folder>${RESET}, then exit"
  parameter_definition "--email=<email>       " "preselect ${BBRIGHT}<email>${RESET} in author list"
  parameter_definition "--skip_upload         " "skip upload prompt after collecting repo info"
  parameter_definition "--parse_libraries     " "besides file extension, infer language from commit contents"
  parameter_definition "--depth=<number>      " "search recursively <number> levels from ${BCYAN}<folder>${RESET}"
  
  echo;echo -e " ${BBRIGHT} *${RESET} (if more than one repo is found, results go to ${BBRIGHT}<output>/<repo_folder>.json${BBRIGHT}"
  
  echo;echo -e "Combine them all (other options available in ${GREEN}main.py${RESET} will be forwared untouched). Eg:";echo
  echo -e "${GREEN}./run.sh ${RESET} --output=huge_repo.json --depth=2 --email=me@mail.com --skip_upload  ${BCYAN}~/projects${RESET}"
  echo -e "${GREEN}./run.sh ${RESET} --output=result_folder --parse_libraries  ${BCYAN}~/projects${RESET}"

  echo ""
}



local_projects() {
  number_of_repos=0
  for REPO_PATH in `find $folder -maxdepth $depth  -type d  -wholename "*/.git" -exec dirname {} \;`; do
    concatenated_repos+=("$REPO_PATH")
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
  
  python src/main.py --depth=$depth $other_args "$folder" ;\
  wait ;\
  return
 
}
depth=1
dry_run=0
folder="${!#}"
paramnum=$#
email=""
upload='default'
optspec=":h-:"
other_args=''
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

                help)
                    print_help
                    exit 0
                    ;;
                logo)
                    print_logo
                    exit 0
                    ;;
                *)
                    other_args="$other_args --$OPTARG"
                    ;;
            esac;;
        h)
            print_help
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