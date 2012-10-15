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

- stream

  integrate with python, php

- profiler

- ReadLines should be move to dlog.go
