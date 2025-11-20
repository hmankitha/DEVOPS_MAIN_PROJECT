from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.database import get_db
from app.models import Inventory, Product
from app.schemas import InventoryResponse, InventoryUpdate
from app.middleware.auth import get_current_user

router = APIRouter()

@router.get("/{product_id}", response_model=InventoryResponse)
async def get_inventory(product_id: str, db: Session = Depends(get_db)):
    """Get inventory for a product"""
    inventory = db.query(Inventory).filter(Inventory.product_id == product_id).first()
    if not inventory:
        raise HTTPException(status_code=404, detail="Inventory not found")
    return inventory

@router.put("/{product_id}", response_model=InventoryResponse)
async def update_inventory(
    product_id: str,
    inventory_update: InventoryUpdate,
    db: Session = Depends(get_db),
    current_user: dict = Depends(get_current_user)
):
    """Update inventory (Admin only)"""
    if current_user.get("role") != "admin":
        raise HTTPException(status_code=403, detail="Admin access required")
    
    db_inventory = db.query(Inventory).filter(Inventory.product_id == product_id).first()
    if not db_inventory:
        raise HTTPException(status_code=404, detail="Inventory not found")
    
    update_data = inventory_update.dict(exclude_unset=True)
    for key, value in update_data.items():
        setattr(db_inventory, key, value)
    
    # Calculate available quantity
    db_inventory.available_quantity = db_inventory.quantity - db_inventory.reserved_quantity
    
    db.commit()
    db.refresh(db_inventory)
    return db_inventory

@router.post("/{product_id}/reserve")
async def reserve_inventory(
    product_id: str,
    quantity: int,
    db: Session = Depends(get_db),
    current_user: dict = Depends(get_current_user)
):
    """Reserve inventory for an order"""
    inventory = db.query(Inventory).filter(Inventory.product_id == product_id).first()
    if not inventory:
        raise HTTPException(status_code=404, detail="Inventory not found")
    
    if inventory.available_quantity < quantity:
        raise HTTPException(status_code=400, detail="Insufficient inventory")
    
    inventory.reserved_quantity += quantity
    inventory.available_quantity = inventory.quantity - inventory.reserved_quantity
    db.commit()
    
    return {"success": True, "message": "Inventory reserved successfully"}

@router.post("/{product_id}/release")
async def release_inventory(
    product_id: str,
    quantity: int,
    db: Session = Depends(get_db),
    current_user: dict = Depends(get_current_user)
):
    """Release reserved inventory"""
    inventory = db.query(Inventory).filter(Inventory.product_id == product_id).first()
    if not inventory:
        raise HTTPException(status_code=404, detail="Inventory not found")
    
    inventory.reserved_quantity -= quantity
    inventory.available_quantity = inventory.quantity - inventory.reserved_quantity
    db.commit()
    
    return {"success": True, "message": "Inventory released successfully"}
