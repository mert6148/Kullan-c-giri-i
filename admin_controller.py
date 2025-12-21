#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Admin Controller - Python3
Gelişmiş admin yönetimi ve kontrolü için controller sınıfı
C/C++ entegrasyonu ile performans optimizasyonu
"""

import sys
import json
import sqlite3
import hashlib
import hmac
import uuid
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Any, Tuple
from functools import wraps
import logging

# C/C++ extension import
try:
    from cpp_admin_wrapper import CppAdminWrapper
    CPP_AVAILABLE = True
except ImportError:
    CPP_AVAILABLE = False
    print("Warning: C/C++ extension not available, using pure Python implementation")
    CppAdminWrapper = None

# Import login system
sys.path.insert(0, str(Path(__file__).parent))
try:
    import print as login_system
except ImportError:
    print("Error: Could not import print module.")
    sys.exit(1)

# Logging setup
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Permission constants
PERM_USER_MANAGE = 'user:manage'
PERM_SYSTEM_CONFIG = 'system:config'
PERM_VIEW_LOGS = 'logs:view'
PERM_MANAGE_ADMINS = 'admin:manage'
PERM_SECURITY = 'security:manage'
PERM_DATABASE = 'database:manage'
PERM_API = 'api:manage'
PERM_ML_CONFIG = 'ml:config'
PERM_OS_CONFIG = 'os:config'

# Role constants
ROLE_SUPER_ADMIN = 'super_admin'
ROLE_ADMIN = 'admin'
ROLE_MODERATOR = 'moderator'
ROLE_VIEWER = 'viewer'

# Session timeout (minutes)
SESSION_TIMEOUT = 30
SESSION_WARNING = 25


class AdminController:
    """Admin Controller - Yönetici işlemleri ve kontrolü"""
    
    def __init__(self, db_path: str = "login_system.db"):
        """Admin Controller başlatma"""
        self.db_path = db_path
        self.audit_log_file = "admin_audit.log"
        self.roles = self._initialize_roles()
        self.sessions = {}
        
        # C/C++ extension kullanımı
        if CPP_AVAILABLE:
            logger.info("Using C/C++ extension for performance")
            self.use_cpp = True
        else:
            logger.info("Using pure Python implementation")
            self.use_cpp = False
        
        # Veritabanı bağlantısı
        self._init_admin_db()
    
    def _initialize_roles(self) -> Dict[str, List[str]]:
        """Rol ve izin eşleştirmelerini başlat"""
        return {
            ROLE_SUPER_ADMIN: [
                PERM_USER_MANAGE, PERM_SYSTEM_CONFIG, PERM_VIEW_LOGS,
                PERM_MANAGE_ADMINS, PERM_SECURITY, PERM_DATABASE,
                PERM_API, PERM_ML_CONFIG, PERM_OS_CONFIG
            ],
            ROLE_ADMIN: [
                PERM_USER_MANAGE, PERM_VIEW_LOGS, PERM_SECURITY,
                PERM_API, PERM_ML_CONFIG
            ],
            ROLE_MODERATOR: [
                PERM_USER_MANAGE, PERM_VIEW_LOGS
            ],
            ROLE_VIEWER: [
                PERM_VIEW_LOGS
            ]
        }
    
    def _init_admin_db(self):
        """Admin veritabanı tablolarını oluştur"""
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        
        # Admin sessions table
        c.execute("""
            CREATE TABLE IF NOT EXISTS admin_sessions (
                id TEXT PRIMARY KEY,
                username TEXT NOT NULL,
                role TEXT NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                expires_at TIMESTAMP NOT NULL,
                ip_address TEXT,
                user_agent TEXT,
                is_active BOOLEAN DEFAULT 1
            )
        """)
        
        # Admin audit log table
        c.execute("""
            CREATE TABLE IF NOT EXISTS admin_audit_log (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                admin_username TEXT NOT NULL,
                action TEXT NOT NULL,
                resource TEXT,
                details TEXT,
                ip_address TEXT,
                timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                status TEXT DEFAULT 'success'
            )
        """)
        
        # Admin permissions table
        c.execute("""
            CREATE TABLE IF NOT EXISTS admin_permissions (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT UNIQUE NOT NULL,
                role TEXT NOT NULL,
                permissions TEXT,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        """)
        
        conn.commit()
        conn.close()
    
    def login(self, username: str, password: str, role: str = ROLE_ADMIN) -> Tuple[bool, Optional[str], Optional[str]]:
        """
        Admin girişi
        Returns: (success, session_id, error_message)
        """
        try:
            # Kullanıcı doğrulama
            conn = login_system.get_db_connection()
            c = conn.cursor()
            c.execute("SELECT id, salt, hash FROM users WHERE username = ?", (username,))
            user = c.fetchone()
            conn.close()
            
            if not user:
                self._log_audit(username, "login_attempt", "failed", "User not found")
                return False, None, "Kullanıcı bulunamadı"
            
            user_id, salt, stored_hash = user
            
            # Şifre doğrulama
            if not login_system.verify_password(password, salt, stored_hash):
                self._log_audit(username, "login_attempt", "failed", "Invalid password")
                return False, None, "Geçersiz şifre"
            
            # Rol kontrolü
            if not self._check_role_permission(username, role):
                self._log_audit(username, "login_attempt", "failed", f"Invalid role: {role}")
                return False, None, "Geçersiz rol"
            
            # Session oluştur
            session_id = str(uuid.uuid4())
            expires_at = datetime.now() + timedelta(minutes=SESSION_TIMEOUT)
            
            conn = sqlite3.connect(self.db_path)
            c = conn.cursor()
            c.execute("""
                INSERT INTO admin_sessions 
                (id, username, role, expires_at, ip_address, user_agent)
                VALUES (?, ?, ?, ?, ?, ?)
            """, (
                session_id, username, role, expires_at.isoformat(),
                self._get_client_ip(), self._get_user_agent()
            ))
            conn.commit()
            conn.close()
            
            self.sessions[session_id] = {
                'username': username,
                'role': role,
                'expires_at': expires_at,
                'created_at': datetime.now()
            }
            
            self._log_audit(username, "login", "success", f"Role: {role}")
            logger.info(f"Admin login successful: {username} ({role})")
            
            return True, session_id, None
            
        except Exception as e:
            logger.error(f"Login error: {str(e)}")
            return False, None, f"Hata: {str(e)}"
    
    def logout(self, session_id: str) -> bool:
        """Admin çıkışı"""
        try:
            if session_id in self.sessions:
                username = self.sessions[session_id]['username']
                del self.sessions[session_id]
                
                conn = sqlite3.connect(self.db_path)
                c = conn.cursor()
                c.execute("UPDATE admin_sessions SET is_active = 0 WHERE id = ?", (session_id,))
                conn.commit()
                conn.close()
                
                self._log_audit(username, "logout", "success")
                logger.info(f"Admin logout: {username}")
                return True
            return False
        except Exception as e:
            logger.error(f"Logout error: {str(e)}")
            return False
    
    def validate_session(self, session_id: str) -> Tuple[bool, Optional[Dict]]:
        """
        Session doğrulama
        Returns: (is_valid, session_data)
        """
        try:
            if session_id in self.sessions:
                session = self.sessions[session_id]
                if datetime.now() < session['expires_at']:
                    return True, session
                else:
                    # Session süresi dolmuş
                    del self.sessions[session_id]
                    return False, None
            
            # Veritabanından kontrol et
            conn = sqlite3.connect(self.db_path)
            c = conn.cursor()
            c.execute("""
                SELECT username, role, expires_at FROM admin_sessions 
                WHERE id = ? AND is_active = 1
            """, (session_id,))
            result = c.fetchone()
            conn.close()
            
            if result:
                username, role, expires_at_str = result
                expires_at = datetime.fromisoformat(expires_at_str)
                
                if datetime.now() < expires_at:
                    session_data = {
                        'username': username,
                        'role': role,
                        'expires_at': expires_at
                    }
                    self.sessions[session_id] = session_data
                    return True, session_data
                else:
                    # Session süresi dolmuş, pasif yap
                    conn = sqlite3.connect(self.db_path)
                    c = conn.cursor()
                    c.execute("UPDATE admin_sessions SET is_active = 0 WHERE id = ?", (session_id,))
                    conn.commit()
                    conn.close()
            
            return False, None
            
        except Exception as e:
            logger.error(f"Session validation error: {str(e)}")
            return False, None
    
    def has_permission(self, session_id: str, permission: str) -> bool:
        """İzin kontrolü"""
        is_valid, session = self.validate_session(session_id)
        if not is_valid or not session:
            return False
        
        role = session.get('role', ROLE_VIEWER)
        role_permissions = self.roles.get(role, [])
        
        return permission in role_permissions
    
    def require_permission(self, session_id: str, permission: str) -> bool:
        """İzin gerektiren işlemler için kontrol"""
        if not self.has_permission(session_id, permission):
            self._log_audit(
                self.sessions.get(session_id, {}).get('username', 'unknown'),
                "permission_denied",
                "failed",
                f"Permission required: {permission}"
            )
            return False
        return True
    
    def get_users(self, session_id: str, limit: int = 100) -> Tuple[bool, Optional[List[Dict]], Optional[str]]:
        """Kullanıcı listesi alma"""
        if not self.require_permission(session_id, PERM_USER_MANAGE):
            return False, None, "İzin yetersiz"
        
        try:
            users = login_system.list_users()
            limited_users = users[:limit]
            
            self._log_audit(
                self.sessions[session_id]['username'],
                "get_users",
                "success",
                f"Count: {len(limited_users)}"
            )
            
            return True, limited_users, None
        except Exception as e:
            logger.error(f"Get users error: {str(e)}")
            return False, None, str(e)
    
    def create_user(self, session_id: str, username: str, password: str, 
                   full_name: str = None, email: str = None) -> Tuple[bool, Optional[str]]:
        """Kullanıcı oluşturma"""
        if not self.require_permission(session_id, PERM_USER_MANAGE):
            return False, "İzin yetersiz"
        
        try:
            success = login_system.create_user(username, password, full_name)
            if success:
                # Email ekle
                if email:
                    login_system.set_user_attribute(username, "email", email)
                
                self._log_audit(
                    self.sessions[session_id]['username'],
                    "create_user",
                    "success",
                    f"Username: {username}"
                )
                return True, None
            else:
                return False, "Kullanıcı zaten mevcut"
        except Exception as e:
            logger.error(f"Create user error: {str(e)}")
            return False, str(e)
    
    def delete_user(self, session_id: str, username: str) -> Tuple[bool, Optional[str]]:
        """Kullanıcı silme"""
        if not self.require_permission(session_id, PERM_USER_MANAGE):
            return False, "İzin yetersiz"
        
        try:
            success = login_system.delete_user(username)
            if success:
                self._log_audit(
                    self.sessions[session_id]['username'],
                    "delete_user",
                    "success",
                    f"Username: {username}"
                )
                return True, None
            else:
                return False, "Kullanıcı bulunamadı"
        except Exception as e:
            logger.error(f"Delete user error: {str(e)}")
            return False, str(e)
    
    def get_audit_logs(self, session_id: str, limit: int = 50) -> Tuple[bool, Optional[List[Dict]], Optional[str]]:
        """Audit logları alma"""
        if not self.require_permission(session_id, PERM_VIEW_LOGS):
            return False, None, "İzin yetersiz"
        
        try:
            conn = sqlite3.connect(self.db_path)
            c = conn.cursor()
            c.execute("""
                SELECT admin_username, action, resource, details, ip_address, 
                       timestamp, status
                FROM admin_audit_log
                ORDER BY timestamp DESC
                LIMIT ?
            """, (limit,))
            
            logs = []
            for row in c.fetchall():
                logs.append({
                    'admin_username': row[0],
                    'action': row[1],
                    'resource': row[2],
                    'details': row[3],
                    'ip_address': row[4],
                    'timestamp': row[5],
                    'status': row[6]
                })
            
            conn.close()
            return True, logs, None
        except Exception as e:
            logger.error(f"Get audit logs error: {str(e)}")
            return False, None, str(e)
    
    def get_system_stats(self, session_id: str) -> Tuple[bool, Optional[Dict], Optional[str]]:
        """Sistem istatistikleri"""
        if not self.has_permission(session_id, PERM_VIEW_LOGS):
            return False, None, "İzin yetersiz"
        
        try:
            # C/C++ extension kullanarak performans optimizasyonu
            if CPP_AVAILABLE and self.use_cpp and CppAdminWrapper:
                stats = CppAdminWrapper.get_system_stats()
            else:
                stats = self._get_system_stats_python()
            
            return True, stats, None
        except Exception as e:
            logger.error(f"Get system stats error: {str(e)}")
            return False, None, str(e)
    
    def _get_system_stats_python(self) -> Dict:
        """Python implementasyonu ile sistem istatistikleri"""
        conn = login_system.get_db_connection()
        c = conn.cursor()
        
        # Kullanıcı sayısı
        c.execute("SELECT COUNT(*) FROM users")
        user_count = c.fetchone()[0]
        
        # Aktif session sayısı
        c.execute("SELECT COUNT(*) FROM sessions WHERE logout_ts IS NULL")
        active_sessions = c.fetchone()[0]
        
        # Admin session sayısı
        conn_admin = sqlite3.connect(self.db_path)
        c_admin = conn_admin.cursor()
        c_admin.execute("SELECT COUNT(*) FROM admin_sessions WHERE is_active = 1")
        admin_sessions = c_admin.fetchone()[0]
        conn_admin.close()
        
        conn.close()
        
        return {
            'total_users': user_count,
            'active_sessions': active_sessions,
            'admin_sessions': admin_sessions,
            'timestamp': datetime.now().isoformat()
        }
    
    def _check_role_permission(self, username: str, role: str) -> bool:
        """Rol izin kontrolü"""
        # Veritabanından rol kontrolü
        conn = sqlite3.connect(self.db_path)
        c = conn.cursor()
        c.execute("SELECT role FROM admin_permissions WHERE username = ?", (username,))
        result = c.fetchone()
        conn.close()
        
        if result:
            db_role = result[0]
            # Rol hiyerarşisi kontrolü
            role_hierarchy = {
                ROLE_SUPER_ADMIN: 4,
                ROLE_ADMIN: 3,
                ROLE_MODERATOR: 2,
                ROLE_VIEWER: 1
            }
            
            requested_level = role_hierarchy.get(role, 0)
            user_level = role_hierarchy.get(db_role, 0)
            
            return user_level >= requested_level
        
        # İlk admin oluşturma
        if role == ROLE_SUPER_ADMIN:
            conn = sqlite3.connect(self.db_path)
            c = conn.cursor()
            c.execute("""
                INSERT OR IGNORE INTO admin_permissions (username, role, permissions)
                VALUES (?, ?, ?)
            """, (username, role, json.dumps(self.roles.get(role, []))))
            conn.commit()
            conn.close()
            return True
        
        return False
    
    def _log_audit(self, username: str, action: str, status: str, 
                  details: str = None, resource: str = None):
        """Audit log kaydı"""
        try:
            conn = sqlite3.connect(self.db_path)
            c = conn.cursor()
            c.execute("""
                INSERT INTO admin_audit_log 
                (admin_username, action, resource, details, ip_address, status)
                VALUES (?, ?, ?, ?, ?, ?)
            """, (
                username, action, resource, details,
                self._get_client_ip(), status
            ))
            conn.commit()
            conn.close()
            
            # Dosyaya da yaz (yedek)
            log_entry = {
                'timestamp': datetime.now().isoformat(),
                'username': username,
                'action': action,
                'status': status,
                'details': details,
                'resource': resource,
                'ip': self._get_client_ip()
            }
            
            with open(self.audit_log_file, 'a', encoding='utf-8') as f:
                f.write(json.dumps(log_entry, ensure_ascii=False) + '\n')
                
        except Exception as e:
            logger.error(f"Audit log error: {str(e)}")
    
    def _get_client_ip(self) -> str:
        """İstemci IP adresini al"""
        # Flask request context'inden alınabilir
        try:
            from flask import request
            return request.remote_addr if request else "127.0.0.1"
        except:
            return "127.0.0.1"
    
    def _get_user_agent(self) -> str:
        """User agent string'i al"""
        try:
            from flask import request
            if request and request.user_agent:
                return request.user_agent.string[:200]
        except:
            pass
        return "Unknown"


# Global admin controller instance
_admin_controller = None

def get_admin_controller() -> AdminController:
    """Global admin controller instance'ı al"""
    global _admin_controller
    if _admin_controller is None:
        _admin_controller = AdminController()
    return _admin_controller

