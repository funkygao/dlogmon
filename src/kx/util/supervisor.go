// TODO not finished yet
package util

// The interface for recoverable types.
type Recoverable interface {
    Recover(Recoverable, Any)
}

type cryForHelpMsg struct {
    r   Recoverable
    i   Any
}

// The supervisor itself.
type Supervisor struct {
    supervisor *Supervisor
    chHelp     chan *cryForHelpMsg
}

func NewSupervisor(s *Supervisor) *Supervisor {
    this := &Supervisor{s, make(chan *cryForHelpMsg)}

    go s.backend()
    return this
}

func (this *Supervisor) Help(r Recoverable, i Any) {
    this.chHelp <- &cryForHelpMsg{r, i}
}

func (this *Supervisor) Recover(r Recoverable, err Any) {
}

func (this *Supervisor) backend() {
    defer func() {
        if i := recover(); i != nil {
            if this.supervisor != nil {
                this.supervisor.Help(this, i)
            }
        }
    }()
}
