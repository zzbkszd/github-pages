#!/usr/bin/env python
# -*- coding: utf-8 -*-

#函数定义
#python不支持函数的重载，可以设定默认参数
#python支持将函数作为参数传递进另一个函数
def addto (a,b=3):
    if(not isinstance(a,(float,int))):
        raise TypeError("bad operator type")
    return a+b
def nothing ():
    pass

print addto(1,2)

print addto(1)

#lambda表示匿名函数

#装饰器:


def log(func):
    def wapper(*args,**kw):
        print("call for func %s()" % func.__name__)
        return func(*args,**kw)
    return wapper

@log
def now():
    print "20160115"

now()
