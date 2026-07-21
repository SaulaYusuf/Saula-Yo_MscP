import asyncio
import csv
import logging
import time
from pathlib import Path
import aiohttp

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

METADATA_URL = "http://localhost:8080/api/metadata"
DATASET_PATH = Path(__file__).parents[3] / "data" / "raw" / "bdt_mba_supplychain_dataset_2024.csv"
CONCURRENCY = 10

async def send_metadata(session, row):
    payload = {
        "asset_id": row.get("Asset_ID", ""),
        "location": row.get("Location", ""),
        "temperature": float(row.get("Temperature", 0)),
        "vibration": float(row.get("Vibration", 0)),
        "last_maintenance": row.get("Last_Maintenance", ""),
        "condition_score": float(row.get("Condition_Score", 0)),
        "resource_utilization": float(row.get("Resource_Utilization", 0)),
        "delivery_efficiency": float(row.get("Delivery_Efficiency", 0)),
        "downtime_hours": float(row.get("Downtime_Hours", 0)),
        "inventory_level": row.get("Inventory_Level", ""),
        "logistics_cost": float(row.get("Logistics_Cost", 0)),
        "timestamp": row.get("Timestamp", ""),
        "supply_chain_efficiency_label": int(row.get("SupplyChain_Efficiency_Label", 0)),
    }
    try:
        async with session.post(METADATA_URL, json=payload, timeout=5) as resp:
            if resp.status == 200:
                return {"success": True, "asset": payload["asset_id"]}
            else:
                text = await resp.text()
                return {"success": False, "asset": payload["asset_id"], "error": text}
    except Exception as e:
        return {"success": False, "asset": payload["asset_id"], "error": str(e)}

async def worker(session, queue, results):
    while True:
        row = await queue.get()
        if row is None:
            queue.task_done()
            break
        result = await send_metadata(session, row)
        results.append(result)
        if not result["success"]:
            logger.warning(f"Failed for {result['asset']}: {result.get('error')}")
        queue.task_done()

async def main():
    queue = asyncio.Queue()
    rows = []
    with open(DATASET_PATH, newline='', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            rows.append(row)
    logger.info(f"Loaded {len(rows)} metadata records.")

    for row in rows:
        await queue.put(row)
    for _ in range(CONCURRENCY):
        await queue.put(None)

    results = []
    async with aiohttp.ClientSession() as session:
        workers = [asyncio.create_task(worker(session, queue, results)) for _ in range(CONCURRENCY)]
        await queue.join()
        for w in workers:
            w.cancel()

    success = sum(1 for r in results if r["success"])
    logger.info(f"Sent {len(results)} records. Success: {success}, Failed: {len(results)-success}")
    logger.info(f"Total time: {time.time()-start:.2f}s")

if __name__ == "__main__":
    start = time.time()
    asyncio.run(main())