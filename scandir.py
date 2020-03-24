#conding=utf8
import os
import re

project_name = 'github.com/hilonfot/network'
localdir = '.'

g = os.walk("./")
with open('gomod_replace','w') as file_object:
    for path,dir_list,file_list in g:
        for dir_name in dir_list:
            if  '.git' not in path and '.idea' not in path and '.git' not in dir_name and '.idea' not in dir_name:
                dir = os.path.join(path, dir_name)[2:]
                file_object.write('%s/%s => %s/%s\r' %(project_name,dir,localdir,dir))

