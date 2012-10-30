package dlog

import (
    "github.com/bmizerany/assert"
    "testing"
)

func newWorker() IWorker {
    option := new(Option)
    m := NewManager(option)
    return NewAmfWorker(m, "", "", 1)
}

func TestIsLineValid(t *testing.T) {
    expected := map[string]bool{
        "we are here": false,
        "amf_slow":    false,
        "AMF_SLOW":    false,
        `>121015-180201 192.168.100.123 10282 KP:PHP.CDlog AMF_SLOW POST+www.kaixin001.com/city/gateway.php+30600bfc {'calltime':2043,'classname':'CCityConfig','method':'callfunc','args':['47116815_1226_47116815_1350293555_9af3436e3a716f7afc298bb77ece48fe','16616590',"CCityConfig.callfunc","getConfig","68680510"]}`: true,
    }
    amf := newWorker()

    for k := range expected {
        assert.Equal(t, expected[k], amf.IsLineValid(k))
    }
}
