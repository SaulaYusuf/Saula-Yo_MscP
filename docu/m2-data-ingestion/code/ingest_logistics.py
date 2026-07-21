"""
Ingest smart_logistics_dataset.csv into the Master contract (RecordHandover)
via the Go bridge at /api/logistics.
"""

import asyncio
import csv
import logging
import time
from pathlib import Path

import aiohttp

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

LOGISTICS_URL = "http://localhost:8080/api/logistics"
DATASET_PATH = Path(__file__).parents[3] / "data" / "raw" / "smart_logistics_dataset.csv"
CONCURRENCY = 10

async def send_logistics(session, row):
    payload = {
        "shipment_id": row.get("Asset_ID", "unknown"),
        "origin": row.get("Latitude", ""),          # or use a meaningful origin
        "destination": row.get("Longitude", ""),
        "status": row.get("Shipment_Status", "UNKNOWN"),
        "timestamp": row.get("Timestamp", ""),
    }
    try:
        async with session.post(LOGISTICS_URL, json=payload, timeout=5) as resp:
            if resp.status == 200:
                return {"success": True, "shipment": payload["shipment_id"]}
            else:
                text = await resp.text()
                return {"success": False, "shipment": payload["shipment_id"], "error": text}
    except Exception as e:
        return {"success": False, "shipment": payload["shipment_id"], "error": str(e)}

async def worker(session, queue, results):
    while True:
        row = await queue.get()
        if row is None:
            queue.task_done()
            break
        result = await send_logistics(session, row)
        results.append(result)
        if not result["success"]:
            logger.warning(f"Failed for {result['shipment']}: {result.get('error')}")
        queue.task_done()

async def main():
    queue = asyncio.Queue()
    rows = []
    with open(DATASET_PATH, newline='', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            rows.append(row)
    logger.info(f"Loaded {len(rows)} logistics records.")

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