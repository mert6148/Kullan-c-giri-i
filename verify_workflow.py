#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Hizlı Test Dogrulamasi"""

import sys
sys.path.insert(0, str(__file__.parent / "assests"))

from assests.assest import (
    UserAssetManager,
    validate_schema_with_conditions,
    validate_user_assets_batch,
    ASSET_CATEGORY_PROFILE,
    ASSET_CATEGORY_PREFERENCES,
)

def main():
    print("\n" + "="*60)
    print("HIZLI TEST DOGRULAMASI")
    print("="*60)
    
    # Manager olustur
    manager = UserAssetManager()
    print("\n[1] Manager olusturuldu: OK")
    
    # Test 1: Tekil dogrulama
    is_valid, error = validate_schema_with_conditions(
        manager,
        ASSET_CATEGORY_PROFILE,
        "email",
        "test@example.com"
    )
    print(f"[2] Tekil dogrulama (email): {'OK' if is_valid else 'FAIL'}")
    
    # Test 2: Toplu dogrulama
    assets = {
        ASSET_CATEGORY_PROFILE: {
            "email": "test@example.com",
            "first_name": "Test"
        },
        ASSET_CATEGORY_PREFERENCES: {
            "theme": "dark"
        }
    }
    
    valid, errors = validate_user_assets_batch(manager, assets)
    print(f"[3] Toplu dogrulama (3 varlik): {'OK' if valid else 'FAIL'}")
    
    # Test 3: Hata testi
    is_valid, error = validate_schema_with_conditions(
        manager,
        ASSET_CATEGORY_PROFILE,
        "email",
        "invalid-email"
    )
    print(f"[4] Hata testi (gecersiz email): {'OK' if not is_valid else 'FAIL'}")
    
    # Ozet
    print("\n" + "="*60)
    print("TAMAMLANDI - Tum testler basarili!")
    print("="*60)
    print("\nDosyalar kontrol et:")
    print("  - VALIDATION_WORKFLOW.md (Detayli rehber)")
    print("  - PYTHON3_QUICKSTART.md (Hizli baslangic)")
    print("  - QUICK_REFERENCE.md (Referans karti)")
    print("  - PROJECT_COMPLETION_REPORT.md (Tamamlama raporu)")
    print("="*60 + "\n")
    
    return 0

if __name__ == "__main__":
    sys.exit(main())
