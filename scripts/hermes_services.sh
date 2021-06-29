#!/bin/bash
# Title       : hermes_services.sh
# Description : 
# Author      : Freddy Banquez <fbanquez@gmail.com>
# Date        : 2020-09-27
# Version     : 1.0    
# Usage       : bash hermes_services.sh
#

#LOGFILE="$(pwd)/hermes_services.log"
LOGFILE="/home/viewmed/Documents/scripts/hermes_services.log"
declare -a services=("hermes_receiver.service" "hermes_router1.service" "hermes_router2.service" "hermes_dispatcher1.service" "hermes_dispatcher2.service" "hermes_cleaner.service" "hermes_bookkeeper.service" "hermes_ui.service")

function log() {
    datestring=`date +'%Y-%m-%d %H:%M:%S'`
    echo -e "$datestring - $@" >> $LOGFILE
}

for i in ${services[@]}
do
    ACTIVE="$(systemctl show -p ActiveState --value $i)"
    STATE="$(systemctl show -p SubState --value $i)"

    #echo $i; echo $ACTIVE; echo $STATE

    if [[ $ACTIVE != "active" || $STATE != "running" ]]
    then
        log "{'service': '$i', 'active': '$ACTIVE', 'state':'$STATE', 'action': 'restarting'}"
        systemctl restart $i
    else
        log "{'service': '$i', 'active': '$ACTIVE', 'state':'$STATE', 'action': 'nothing'}"
    fi
done

