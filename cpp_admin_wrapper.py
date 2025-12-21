#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
C/C++ Admin Extension Wrapper
Python wrapper for C/C++ extension module
"""

import sys
from typing import Dict, Optional, Tuple

# Try to import C/C++ extension
try:
    import cpp_admin_extension
    CPP_AVAILABLE = True
except ImportError:
    CPP_AVAILABLE = False
    print("Warning: C/C++ extension not available, using Python fallback")


class CppAdminWrapper:
    """C/C++ extension wrapper class"""
    
    @staticmethod
    def get_system_stats() -> Dict:
        """Get system statistics using C/C++ extension"""
        if CPP_AVAILABLE:
            try:
                stats = cpp_admin_extension.get_system_stats()
                # Convert to standard dict format
                return {
                    'total_users': stats.get('total_users', 0),
                    'active_sessions': stats.get('active_sessions', 0),
                    'admin_sessions': stats.get('admin_sessions', 0),
                    'cpu_usage': stats.get('cpu_usage', 0.0),
                    'memory_usage': stats.get('memory_usage', 0),
                    'timestamp': stats.get('timestamp', 0)
                }
            except Exception as e:
                print(f"C/C++ extension error: {e}")
                return CppAdminWrapper._get_stats_fallback()
        return CppAdminWrapper._get_stats_fallback()
    
    @staticmethod
    def fast_hash_password(password: str, salt: str) -> str:
        """Fast password hashing using C/C++"""
        if CPP_AVAILABLE:
            try:
                return cpp_admin_extension.fast_hash_password(password, salt)
            except Exception as e:
                print(f"C/C++ extension error: {e}")
                return CppAdminWrapper._hash_fallback(password, salt)
        return CppAdminWrapper._hash_fallback(password, salt)
    
    @staticmethod
    def validate_input_fast(input_str: str, min_len: int = 3, max_len: int = 100) -> Tuple[bool, str]:
        """Fast input validation using C/C++"""
        if CPP_AVAILABLE:
            try:
                result = cpp_admin_extension.validate_input_fast(input_str, min_len, max_len)
                if isinstance(result, tuple) and len(result) == 2:
                    return result[0], result[1]
                return result, ""
            except Exception as e:
                print(f"C/C++ extension error: {e}")
                return CppAdminWrapper._validate_fallback(input_str, min_len, max_len)
        return CppAdminWrapper._validate_fallback(input_str, min_len, max_len)
    
    @staticmethod
    def fast_string_compare(str1: str, str2: str) -> bool:
        """Timing-safe string comparison using C/C++"""
        if CPP_AVAILABLE:
            try:
                return cpp_admin_extension.fast_string_compare(str1, str2)
            except Exception as e:
                print(f"C/C++ extension error: {e}")
                return CppAdminWrapper._compare_fallback(str1, str2)
        return CppAdminWrapper._compare_fallback(str1, str2)
    
    # Fallback methods (pure Python)
    @staticmethod
    def _get_stats_fallback() -> Dict:
        """Fallback system stats"""
        return {
            'total_users': 0,
            'active_sessions': 0,
            'admin_sessions': 0,
            'cpu_usage': 0.0,
            'memory_usage': 0,
            'timestamp': 0
        }
    
    @staticmethod
    def _hash_fallback(password: str, salt: str) -> str:
        """Fallback password hashing"""
        import hashlib
        combined = salt + password
        return hashlib.sha256(combined.encode()).hexdigest()
    
    @staticmethod
    def _validate_fallback(input_str: str, min_len: int, max_len: int) -> Tuple[bool, str]:
        """Fallback input validation"""
        if not input_str:
            return False, "Input is empty"
        
        length = len(input_str)
        if length < min_len or length > max_len:
            return False, f"Length must be between {min_len} and {max_len}"
        
        dangerous_chars = ["'", '"', ";", "--", "/*", "*/"]
        for char in dangerous_chars:
            if char in input_str:
                return False, f"Dangerous character detected: {char}"
        
        return True, ""
    
    @staticmethod
    def _compare_fallback(str1: str, str2: str) -> bool:
        """Fallback timing-safe comparison"""
        import hmac
        return hmac.compare_digest(str1, str2)


# Update admin_controller.py to use wrapper
def update_admin_controller_with_cpp():
    """Update admin controller to use C/C++ wrapper"""
    # This will be called during import
    pass

