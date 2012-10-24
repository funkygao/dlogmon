=========================
dlog monitor
=========================

:Author: Gao Peng <funky.gao@gmail.com>
:Description: kaixin dlog monitor framework
:Revision: $Id$

.. contents:: Table Of Contents
.. section-numbering::


Project
============

dependency
----------

::

    make dep


compile
-------

::

    export GOPATH=$GOPATH:this_dir
    
    go install kx/mapreduce
    go install kx/dlogmon
    go test kx/mapreduce
    
    cd src/kx/mapreduce
    go test mapreduce_test.go
    
    go clean


Todo
====

- why first strings.Contains(100.123) will be slower to 280000 per sec

  enhance performance of strings.Contains

- distributed

- chunk size

  control cpu load

- signal skipped not work

- timespan in db

- external map reduce stream

  ExtractLineInfo deprecated

- chain of workers and reducers
