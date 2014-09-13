# -*- coding: utf-8 -*-

from gevent import socket
import struct
import gevent

def create_connection(address, timeout=None, **ssl_args):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sock.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 0)
    
    if timeout:
        sock.settimeout(timeout)
    if ssl_args:
        from gevent.ssl import wrap_socket
        sock = wrap_socket(sock, **ssl_args)
        
    host = address[0]
    port = int(address[1]) 
    sock.connect((host, port))
    
    return sock


if __name__ == "__main__":
    import struct
    sock = create_connection(("127.0.0.1", 7005))
    while True:
        sock.sendall(struct.pack(">I5s", 5, "hello"))
        print repr(sock.recv(10))
        gevent.sleep(3)
