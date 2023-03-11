#!/usr/bin/env python3

import sys
import sqlite3
import json
import struct
import socket


db = sqlite3.connect('/home/nc/tabs.db')
db.execute('''
CREATE TABLE IF NOT EXISTS tabs (client_id TEXT PRIMARY KEY, tabs JSON);
''')

def read_message():
    rawLength = sys.stdin.buffer.read(4)
    if len(rawLength) == 0:
        sys.exit(0)
    messageLength = struct.unpack('@I', rawLength)[0]
    message = sys.stdin.buffer.read(messageLength).decode('utf-8')
    return json.loads(message)

tabs = read_message()
with open('/home/nc/tmp', 'w') as out:
    try:
        db.execute('insert or replace into tabs (client_id, tabs) values (?, ?)', (socket.gethostname(), json.dumps(tabs)))
        db.commit()
    except Exception as e:
        import traceback
        print('error', file=out)
        traceback.print_exception(e, file=out)
        raise
