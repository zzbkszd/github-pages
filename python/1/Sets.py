#!/usr/bin/env python
# -*- coding: utf-8 -*-

#python常用集合

#list 可变

l_list = ["item1","item2","item3"]

print(l_list)

print(l_list[1])

print(l_list[-1])

l_list.append("item4")

print(l_list)

l_list.insert(2,"item5");

print(l_list)

item = l_list.pop()

print(l_list,item)

#tuple 不可变

l_tuple = ("item1",2,item)

print l_tuple

l_tuple_single = (1,)

print l_tuple_single


#切片
print l_list[1:3]


#流程控制

age = 211

if age==20:
    print age
elif age>=70:
    print 71
else :
    print 21

for name in l_list :
    print name

for num in range(1,10):
    print num

#生成器
l = [x*x for x in range(1,10)]
g = (x*x for x in range(1,10))

for num in g:
    print num

#字典和集合

s_dic = {"name":"ay","age":20}
print s_dic

print s_dic["name"]

s_set = set(l_list)

print s_set






