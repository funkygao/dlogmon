=========================
dlog monitor
=========================

:Author: Gao Peng <funky.gao@gmail.com>
:Description: kaixin dlog monitor framework, simulation of hadoop map reduce framework
:Revision: $Id$

.. contents:: Table Of Contents
.. section-numbering::


Introduction
============
This is a single-node map reduce framework implementation. It is widely used in kaixin's
dlog analysis and reporting system.


Package with potential usage
============================

consisent hash
--------------
/Users/gaopeng/github/dlogmon/src/github.com/stathat/consistent

msgpack
-------
/Users/gaopeng/github/dlogmon/src/github.com/msgpack


Todo
====

- why first strings.Contains(100.123) will be slower to 280000 per sec

  enhance performance of strings.Contains

- signal skipped not work

- chain of workers and reducers

- Job abstraction

  it has member Option

- MapReduce

  Input -> Map -> Shuffle -> Reduce -> Output

- select count(*) as count, service from dlog.20121212 where time>2000 group by service order by count desc

- if group, each group is a dedicated reducer and goroutine
