"""
Ingestion engine for the Kaggle sensor dataset.
Streams rows as JSON to the Go API bridge (http://localhost:8080/api/sensor).
"""

import asyncio
import csv
import json
import logging
import time
from pathlib import Path

import aiohttp

# Configure logging
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

# Bridge endpoint
SENSOR_URL = "http://localhost:8080/api/sensor"

# Path to the dataset (adjust if needed)
DATASET_PATH = Path(__file__).parents[3] / "data" / "raw" / "shipment-sensor-dataset.csv"

# Concurrency limit (number of simultaneous POST requests)
CONCURRENCY = 10

async def send_sensor_data(session: aiohttp.ClientSession, row: dict) -> dict:
    """Send a single sensor row to the bridge."""
    payload = {
        "sensor_id": row.get("Sensor_ID", "unknown"),
        "temp_c": float(row.get("Temp_Max_C", 0.0)),
        "humidity": float(row.get("Humidity_Pct", 0.0)),
        "timestamp": row.get("Timestamp", ""),
    }
    try:
        async with session.post(SENSOR_URL, json=payload, timeout=5) as resp:
            if resp.status == 200:
                return {"success": True, "sensor_id": payload["sensor_id"]}
            else:
                text = await resp.text()
                return {"success": False, "sensor_id": payload["sensor_id"], "error": text}
    except Exception as e:
        return {"success": False, "sensor_id": payload["sensor_id"], "error": str(e)}

async def worker(session: aiohttp.ClientSession, queue: asyncio.Queue, results: list):
    """Worker that pulls rows from the queue and sends them."""
    while True:
        row = await queue.get()
        if row is None:  # poison pill
            queue.task_done()
            break
        result = await send_sensor_data(session, row)
        results.append(result)
        if not result["success"]:
            logger.warning(f"Failed for {result['sensor_id']}: {result.get('error')}")
        queue.task_done()

async def main():
    # Read CSV rows into a queue
    queue = asyncio.Queue()
    rows = []
    with open(DATASET_PATH, newline='', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            rows.append(row)
    logger.info(f"Loaded {len(rows)} sensor records.")

    # Enqueue all rows
    for row in rows:
        await queue.put(row)
    # Add poison pills for workers
    for _ in range(CONCURRENCY):
        await queue.put(None)

    results = []
    async with aiohttp.ClientSession() as session:
        workers = [asyncio.create_task(worker(session, queue, results)) for _ in range(CONCURRENCY)]
        await queue.join()  # wait until all tasks are done
        for w in workers:
            w.cancel()  # cancel workers (they are idle after poison pills)

    success_count = sum(1 for r in results if r["success"])
    logger.info(f"Sent {len(results)} records. Success: {success_count}, Failed: {len(results)-success_count}")

if __name__ == "__main__":
    start = time.time()
    asyncio.run(main())
    logger.info(f"Total time: {time.time()-start:.2f}s")