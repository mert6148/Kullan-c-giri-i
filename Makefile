# Makefile for C/C++ Admin Extension
# Python3 için C/C++ extension derleme

PYTHON = python3
CC = gcc
CXX = g++
CFLAGS = -Wall -O3 -fPIC
CXXFLAGS = -Wall -O3 -fPIC -std=c++11
LDFLAGS = -shared

# Windows için
ifeq ($(OS),Windows_NT)
    EXT = pyd
    CXXFLAGS += /std:c++11 /O2
    LDFLAGS = /DLL
else
    EXT = so
endif

# Source files
CPP_SOURCE = cpp_admin_extension.cpp
TARGET = cpp_admin_extension.$(EXT)

# Default target
all: $(TARGET)

# Build using setuptools (recommended)
build-setuptools:
	$(PYTHON) setup_cpp_extension.py build_ext --inplace

# Manual build
$(TARGET): $(CPP_SOURCE)
	$(CXX) $(CXXFLAGS) $(LDFLAGS) -o $(TARGET) $(CPP_SOURCE) $(shell $(PYTHON)-config --includes) $(shell $(PYTHON)-config --ldflags)

# Clean build artifacts
clean:
	rm -f $(TARGET)
	rm -rf build/
	rm -rf *.egg-info/
	find . -type d -name __pycache__ -exec rm -r {} +
	find . -type f -name "*.pyc" -delete

# Install extension
install: build-setuptools
	$(PYTHON) setup_cpp_extension.py install

# Test extension
test:
	$(PYTHON) -c "import cpp_admin_extension; print('C/C++ extension loaded successfully')"

# Help
help:
	@echo "Available targets:"
	@echo "  all              - Build C/C++ extension"
	@echo "  build-setuptools - Build using setuptools (recommended)"
	@echo "  clean            - Remove build artifacts"
	@echo "  install          - Install extension"
	@echo "  test             - Test extension import"
	@echo "  help             - Show this help message"

.PHONY: all build-setuptools clean install test help

