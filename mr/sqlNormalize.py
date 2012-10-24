#!/usr/bin/env python
#encoding=utf-8
'''dlog sql extractor and normalizer'''

import re

sql_find = re.compile(r'sqls?\^~?\d*!?.*')
sql_replace = re.compile(r'sqls?\^~?\d*!?')

MAX_SQL_DISPLAY_LEN = 110

def normalize_sql(sql):
    sql = re.sub(r"'(.*?)'", 'S', sql)
    sql = re.sub(r"\d{1,}", 'N', sql)
    return sql 

def skipped_url(url):
    if 'kaixin002' in url:
        return None
    if '/admin/' in url:
        # 后台的不计算在内
        return None

def extract_line(line):
    '''从日志的一行记录里提取结构化信息'''
    items = line.split(None, 10)
    time_span = items[0].replace('>', '')
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
    if method == 'CLI' or skipped_url(url):
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
    line = '>121024-163518 192.168.100.123 3309 KProxy KXI.SQA /SAMPLE:1/A T=0.001 9999/127.0.0.1:42162 377 Q=DBMan#app.sQuery X{CALLER^POST+www.kaixin001.com/city/gateway.php+18021f84; MASTER^~F} {kind^s_user_city_gray_inspect; hintId^88546966; convert^~T; sql^~!select orderid,fuid,binspect,extdata,ctime from s_user_city_gray_inspect where uid = 88546966  ~} A=0 {converted^~T; affectedRowNumber^0; fields^[orderid; fuid; binspect; extdata; ctime]; rows^[]}'
    url, rid, service, time, sql, time_span = extract_line(line)
    print url, rid, service, time, sql, time_span
    print normalize_sql(sql)
