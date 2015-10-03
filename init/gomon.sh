#!/bin/bash

usage()
{
    echo "Usage: $0 {start|stop|restart|status}"
    exit 1
}

[ $# -gt 0 ] || usage

ACTION=$1

GOMON_WD="/opt/self-monitoring/"

GOMON_BINARY="gomon"
GOMON_CONFIG="config.json"

GOMON_PID="/var/lock/gomon.pid"
GOMON_LOG="/var/log/gomon.log"

export GOMON_WD

RUN_CMD="$GOMON_WD$GOMON_BINARY --config=$GOMON_WD$GOMON_CONFIG"

case "$ACTION" in
    start)
        echo "Starting: "

        if [ -f $GOMON_PID ]
        then
            echo "Already running!"
            exit 1
        fi

        cd $BAMBOO_INSTALL
        nohup $RUN_CMD > $GOMON_LOG 2>&1 & echo $! > $GOMON_PID

        echo "Started at `date`"
        echo "Running pid="`cat $GOMON_PID`
        ;;

    stop)
        PID=`cat $GOMON_PID 2>/dev/null`
        echo "Shutting down: $PID"
        kill $PID 2>/dev/null
        sleep 2
        kill -9 $PID 2>/dev/null
        rm -f $GOMON_PID
        echo "Stopped at `date`"
        ;;

    restart)
        $0 stop $*
        sleep 5
        $0 start $*
        ;;

    status)
        if [ -f $GOMON_PID ]
        then
            echo "Running pid="`cat $GOMON_PID`
            exit 0
        else
            echo "Currently not running."
        fi
        exit 1
        ;;

    *)
        usage
        ;;
esac

exit 0