#!/usr/bin/env python
#encoding=utf-8
'''dlog sql extractor and normalizer'''

import re
import sys

sql_find = re.compile(r'sqls?\^~?\d*!?.*')
sql_replace = re.compile(r'sqls?\^~?\d*!?')

MAX_SQL_DISPLAY_LEN = 110

def normalize_sql(sql):
    sql = re.sub(r"'(.*?)'", 'S', sql)
    sql = re.sub(r"\d{1,}", 'N', sql)
    return sql 

def is_skipped_url(url):
    if 'kaixin002' in url:
        return None
    if '/admin/' in url:
        # 后台的不计算在内
        return None

def extract_line(line):
    '''从日志的一行记录里提取结构化信息'''
    items = line.split(None, 10)
    time_span = items[0].replace('>', '')
    if len(items) < 10:
        return None

    body = items[10]
    ctx = body[2:body.find('}')]
    if not ctx:
        return None

    ctx = dict(map(lambda r: r.split('^'), ctx.split('; ')))
    time = float(items[6][2:])
    if time == 0.0:
        return None

    service = items[9][2:]
    method, url, rid = ctx['CALLER'].split('+')
    if method == 'CLI' or is_skipped_url(url):
        return None

    sql = None
    isSquery = 'sQuery' in service
    isMquery = 'mQuery' in service
    if isSquery or isMquery:
        # mQuery not handled
        header = 'sQ:' if isSquery else 'mQ:' 
        sql = sql_find.findall(line)[0].split('}')[0] 
        sql = header + sql_replace.sub('', sql) 
        if len(sql) > MAX_SQL_DISPLAY_LEN: 
            sql = sql[:MAX_SQL_DISPLAY_LEN] + '...' + str(len(sql) - MAX_SQL_DISPLAY_LEN)

	return url, rid, service, 1000.0 * time, sql, time_span

if __name__ == '__main__':
    while True:
        line = sys.stdin.readline()
        if not line:
            break

        parsed_result = extract_line(line)
        if parsed_result is None:
            print parsed_result
            sys.stdout.flush()
            continue

        url, rid, service, time, sql, time_span = parsed_result
        print url, rid, service, time, normalize_sql(sql), time_span
        sys.stdout.flush()
