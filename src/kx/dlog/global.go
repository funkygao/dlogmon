package dlog

var (
    workerConstructors = map[string]WorkerConstructor{
        "amf": NewAmfWorker}

    amfLineValidatorRegexes = [...][]string{
        {"AMF_SLOW", "PHP.CDlog"}, // must exists
        {"Q=DLog.log"}}            // must not exists
)
