package dlog

import "testing"

func newDlog() *AmfDlog {
    return new(AmfDlog)
}

func TestIsLineValid(t *testing.T) {
    expected := map[string] bool{
        "we are here": false,
        "amf_slow": false,
        "AMF_SLOW": false,
        ">121015-100043 192.168.100.123 3309 KProxy KXI.SQA /SAMPLE:1/S T=0.000 9999/127.0.0.1:31892 0 Q=DLog.log X{CALLER^POST+www.kaixin001.com/city/gateway.php+2d5e99a6} {identity^PHP.CDlog; tag^AMF_SLOW;": true,
    }
    amf := newDlog()

    for k := range expected {
        if amf.IsLineValid(k) != expected[k] {
            t.Error(k)
        }
    }
}

func TestSetFile(t *testing.T) {
    amf := newDlog()
    file, expected := "/kx/dlog/121015/lz.121015-113035", "/kx/dlog/121015/lz.121015-113035"
    amf.SetFile(file)
    if expected != amf.filename {
        t.Error("Invalid filename")
    }
}

