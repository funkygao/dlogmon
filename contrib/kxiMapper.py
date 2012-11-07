#!/usr/bin/env python
#encoding=utf-8
'''dlog sql extractor and normalizer'''

import re
import sys
import json

sql_find = re.compile(r'sqls?\^~?\d*!?.*')
sql_replace = re.compile(r'sqls?\^~?\d*!?')

MAX_SQL_DISPLAY_LEN = 110

# when the line is invalid, how do we send back to dlogmon?
INVALID_FEEDBACK = "None"

def is_valid(line):
    if '/SAMPLE:1/A' not in line:
        return False
    return True

def normalize_sql(sql):
    sql = re.sub(r"'(.*?)'", 'S', sql)
    sql = re.sub(r"\d{1,}", 'N', sql)
    return sql 

def is_skipped_url(url):
    if 'kaixin002' in url:
        return True
    if '/admin/' in url:
        # 后台的不计算在内
        return True
    return False

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
    method, url, rid = ctx['CALLER'].split('+', 3)
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

def feedback_invalid():
    print 
    sys.stdout.flush()

def feedback_json(url, rid, service, time, sql, time_span):
    if sql is None:
        sql = ""
    else:
        sql = normalize_sql(sql)

    obj = {"u": url, "i": rid, "s": service, "t": time, "q": sql}
    try:
        encoded = json.dumps(obj)
    except:
        # TODO gb18030
        feedback_invalid()
        return

    print encoded
    sys.stdout.flush()

if __name__ == '__main__':
    while True:
        line = sys.stdin.readline()
        if not line:
            break

        if not is_valid(line):
            feedback_invalid()
            continue

        parsed_result = extract_line(line)
        if parsed_result is None:
            feedback_invalid()
            continue

        url, rid, service, time, sql, time_span = parsed_result
        feedback_json(url, rid, service, time, sql, time_span)
