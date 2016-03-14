#!/usr/bin/env python
# -*- coding: utf-8 -*-
#python 数据类型

#整数的两种表示方法
n_integer = 1234
n_oxInteger = 0xa1b2c3

print n_integer,n_oxInteger

#浮点数，就这么一个，不分精度
n_float = 3.14

#字符串，支持常见转义
c_string = "Python"
#忽视转义的字符串
c_string_r = r"\n\r\t"
#多行字符串
c_string_multi = ''' this is line 1
this is line 2
this is line 3'''

print "c_string",c_string
print "c_string_r",c_string_r
print "c_string_multi",c_string_multi

#使用占位符进行输出字符串

print r"%d+%d=%d % (2,3,2+3)" 
print "%d+%d=%d" % (2,3,2+3)

#编码问题
print u"中文".encode("GB2312")

b_boolean = True

print "Boolean is True or False, 注意第一个字母大写"

print "使用 and,or,not 进行运算"

n_none = None

print "None is null"
