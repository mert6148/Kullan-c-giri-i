#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
C/C++ Extension Setup Script
Python3 için C/C++ extension modülünü derlemek için setup scripti
"""

from setuptools import setup, Extension
from pybind11 import get_cmake_dir
import pybind11
import sys

# C/C++ extension modülü tanımı
cpp_admin_extension = Extension(
    'cpp_admin_extension',
    sources=['cpp_admin_extension.cpp'],
    include_dirs=[],
    language='c++',
    extra_compile_args=['-std=c++11', '-O3'] if sys.platform != 'win32' else ['/std:c++11', '/O2'],
    extra_link_args=[]
)

setup(
    name='cpp_admin_extension',
    version='1.0.0',
    description='C/C++ Admin Extension for Python3 - Performance optimization',
    ext_modules=[cpp_admin_extension],
    zip_safe=False,
    python_requires='>=3.8',
)

